// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/SSK-TBD/elastic/v7/config"
)

const (
	// Version is the current version of Elastic.
	Version = "7.0.32"

	// DefaultURL is the default endpoint of Elasticsearch on the local machine.
	// It is used e.g. when initializing a new Client without a specific URL.
	DefaultURL = "http://127.0.0.1:9200"

	// DefaultScheme is the default protocol scheme to use when sniffing
	// the Elasticsearch cluster.
	DefaultScheme = "http"

	// DefaultHealthcheckEnabled specifies if healthchecks are enabled by default.
	DefaultHealthcheckEnabled = true

	// DefaultHealthcheckTimeoutStartup is the time the healthcheck waits
	// for a response from Elasticsearch on startup, i.e. when creating a
	// client. After the client is started, a shorter timeout is commonly used
	// (its default is specified in DefaultHealthcheckTimeout).
	DefaultHealthcheckTimeoutStartup = 5 * time.Second

	// DefaultHealthcheckTimeout specifies the time a running client waits for
	// a response from Elasticsearch. Notice that the healthcheck timeout
	// when a client is created is larger by default (see DefaultHealthcheckTimeoutStartup).
	DefaultHealthcheckTimeout = 1 * time.Second

	// DefaultHealthcheckInterval is the default interval between
	// two health checks of the nodes in the cluster.
	DefaultHealthcheckInterval = 60 * time.Second

	// DefaultSnifferEnabled specifies if the sniffer is enabled by default.
	DefaultSnifferEnabled = true

	// DefaultSnifferInterval is the interval between two sniffing procedures,
	// i.e. the lookup of all nodes in the cluster and their addition/removal
	// from the list of actual connections.
	DefaultSnifferInterval = 15 * time.Minute

	// DefaultSnifferTimeoutStartup is the default timeout for the sniffing
	// process that is initiated while creating a new client. For subsequent
	// sniffing processes, DefaultSnifferTimeout is used (by default).
	DefaultSnifferTimeoutStartup = 5 * time.Second

	// DefaultSnifferTimeout is the default timeout after which the
	// sniffing process times out. Notice that for the initial sniffing
	// process, DefaultSnifferTimeoutStartup is used.
	DefaultSnifferTimeout = 2 * time.Second

	// DefaultSendGetBodyAs is the HTTP method to use when elastic is sending
	// a GET request with a body.
	DefaultSendGetBodyAs = "GET"

	// DefaultGzipEnabled specifies if gzip compression is enabled by default.
	DefaultGzipEnabled = false

	// off is used to disable timeouts.
	off = -1 * time.Second
)

var (
	// nilByte is used in JSON marshal/unmarshal
	nilByte = []byte("null")

	// ErrNoClient is raised when no Elasticsearch node is available.
	ErrNoClient = errors.New("no Elasticsearch node available")

	// ErrRetry is raised when a request cannot be executed after the configured
	// number of retries.
	ErrRetry = errors.New("cannot connect after several retries")

	// ErrTimeout is raised when a request timed out, e.g. when WaitForStatus
	// didn't return in time.
	ErrTimeout = errors.New("timeout")

	// noRetries is a retrier that does not retry.
	noRetries = NewStopRetrier()

	// noDeprecationLog is a no-op for logging deprecations.
	noDeprecationLog = func(*http.Request, *http.Response) {}
)

// Doer is an interface to perform HTTP requests.
// It can be used for mocking.
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// ClientOptionFunc is a function that configures a Client.
// It is used in NewClient.
type ClientOptionFunc func(*Client) error

// Client is an Elasticsearch client. Create one by calling NewClient.
type Client struct {
	c Doer // e.g. a net/*http.Client to use for requests

	connsMu sync.RWMutex // connsMu guards the next block
	conns   []*conn      // all connections
	cindex  int          // index into conns

	mu                        sync.RWMutex // guards the next block
	urls                      []string     // set of URLs passed initially to the client
	running                   bool         // true if the client's background processes are running
	errorlog                  Logger       // error log for critical messages
	infolog                   Logger       // information log for e.g. response times
	tracelog                  Logger       // trace log for debugging
	deprecationlog            func(*http.Request, *http.Response)
	scheme                    string          // http or https
	healthcheckEnabled        bool            // healthchecks enabled or disabled
	healthcheckTimeoutStartup time.Duration   // time the healthcheck waits for a response from Elasticsearch on startup
	healthcheckTimeout        time.Duration   // time the healthcheck waits for a response from Elasticsearch
	healthcheckInterval       time.Duration   // interval between healthchecks
	healthcheckStop           chan bool       // notify healthchecker to stop, and notify back
	snifferEnabled            bool            // sniffer enabled or disabled
	snifferTimeoutStartup     time.Duration   // time the sniffer waits for a response from nodes info API on startup
	snifferTimeout            time.Duration   // time the sniffer waits for a response from nodes info API
	snifferInterval           time.Duration   // interval between sniffing
	snifferStop               chan bool       // notify sniffer to stop, and notify back
	decoder                   Decoder         // used to decode data sent from Elasticsearch
	basicAuthUsername         string          // username for HTTP Basic Auth
	basicAuthPassword         string          // password for HTTP Basic Auth
	sendGetBodyAs             string          // override for when sending a GET with a body
	gzipEnabled               bool            // gzip compression enabled or disabled (default)
	requiredPlugins           []string        // list of required plugins
	retrier                   Retrier         // strategy for retries
	retryStatusCodes          []int           // HTTP status codes where to retry automatically (with retrier)
	headers                   http.Header     // a list of default headers to add to each request
}

// NewClient creates a new client to work with Elasticsearch.
//
// NewClient, by default, is meant to be long-lived and shared across
// your application. If you need a short-lived client, e.g. for request-scope,
// consider using NewSimpleClient instead.
//
// The caller can configure the new client by passing configuration options
// to the func.
//
// Example:
//
//	client, err := elastic.NewClient(
//	  elastic.SetURL("http://127.0.0.1:9200", "http://127.0.0.1:9201"),
//	  elastic.SetBasicAuth("user", "secret"))
//
// If no URL is configured, Elastic uses DefaultURL by default.
//
// If the sniffer is enabled (the default), the new client then sniffes
// the cluster via the Nodes Info API
// (see https://www.elastic.co/guide/en/elasticsearch/reference/7.0/cluster-nodes-info.html#cluster-nodes-info).
// It uses the URLs specified by the caller. The caller is responsible
// to only pass a list of URLs of nodes that belong to the same cluster.
// This sniffing process is run on startup and periodically.
// Use SnifferInterval to set the interval between two sniffs (default is
// 15 minutes). In other words: By default, the client will find new nodes
// in the cluster and remove those that are no longer available every
// 15 minutes. Disable the sniffer by passing SetSniff(false) to NewClient.
//
// The list of nodes found in the sniffing process will be used to make
// connections to the REST API of Elasticsearch. These nodes are also
// periodically checked in a shorter time frame. This process is called
// a health check. By default, a health check is done every 60 seconds.
// You can set a shorter or longer interval by SetHealthcheckInterval.
// Disabling health checks is not recommended, but can be done by
// SetHealthcheck(false).
//
// Connections are automatically marked as dead or healthy while
// making requests to Elasticsearch. When a request fails, Elastic will
// call into the Retry strategy which can be specified with SetRetry.
// The Retry strategy is also responsible for handling backoff i.e. the time
// to wait before starting the next request. There are various standard
// backoff implementations, e.g. ExponentialBackoff or SimpleBackoff.
// Retries are disabled by default.
//
// If no HttpClient is configured, then http.DefaultClient is used.
// You can use your own http.Client with some http.Transport for
// advanced scenarios.
//
// An error is also returned when some configuration option is invalid or
// the new client cannot sniff the cluster (if enabled).
func NewClient(options ...ClientOptionFunc) (*Client, error) {
	return DialContext(context.Background(), options...)
}

// NewClientFromConfig initializes a client from a configuration.
func NewClientFromConfig(cfg *config.Config) (*Client, error) {
	options, err := configToOptions(cfg)
	if err != nil {
		return nil, err
	}
	return DialContext(context.Background(), options...)
}

// NewSimpleClient creates a new short-lived Client that can be used in
// use cases where you need e.g. one client per request.
//
// While NewClient by default sets up e.g. periodic health checks
// and sniffing for new nodes in separate goroutines, NewSimpleClient does
// not and is meant as a simple replacement where you don't need all the
// heavy lifting of NewClient.
//
// NewSimpleClient does the following by default: First, all health checks
// are disabled, including timeouts and periodic checks. Second, sniffing
// is disabled, including timeouts and periodic checks. The number of retries
// is set to 1. NewSimpleClient also does not start any goroutines.
//
// Notice that you can still override settings by passing additional options,
// just like with NewClient.
func NewSimpleClient(options ...ClientOptionFunc) (*Client, error) {
	c := &Client{
		c:                         http.DefaultClient,
		conns:                     make([]*conn, 0),
		cindex:                    -1,
		scheme:                    DefaultScheme,
		decoder:                   &DefaultDecoder{},
		healthcheckEnabled:        false,
		healthcheckTimeoutStartup: off,
		healthcheckTimeout:        off,
		healthcheckInterval:       off,
		healthcheckStop:           make(chan bool),
		snifferEnabled:            false,
		snifferTimeoutStartup:     off,
		snifferTimeout:            off,
		snifferInterval:           off,
		snifferStop:               make(chan bool),
		sendGetBodyAs:             DefaultSendGetBodyAs,
		gzipEnabled:               DefaultGzipEnabled,
		retrier:                   noRetries, // no retries by default
		retryStatusCodes:          nil,       // no automatic retries for specific HTTP status codes
		deprecationlog:            noDeprecationLog,
	}

	// Run the options on it
	for _, option := range options {
		if err := option(c); err != nil {
			return nil, err
		}
	}

	// Use a default URL and normalize them
	if len(c.urls) == 0 {
		c.urls = []string{DefaultURL}
	}
	c.urls = canonicalize(c.urls...)

	// If the URLs have auth info, use them here as an alternative to SetBasicAuth
	if c.basicAuthUsername == "" && c.basicAuthPassword == "" {
		for _, urlStr := range c.urls {
			u, err := url.Parse(urlStr)
			if err == nil && u.User != nil {
				c.basicAuthUsername = u.User.Username()
				c.basicAuthPassword, _ = u.User.Password()
				break
			}
		}
	}

	for _, url := range c.urls {
		c.conns = append(c.conns, newConn(url, url))
	}

	// Ensure that we have at least one connection available
	if err := c.mustActiveConn(); err != nil {
		return nil, err
	}

	c.mu.Lock()
	c.running = true
	c.mu.Unlock()

	return c, nil
}

// Dial will call DialContext with a background context.
func Dial(options ...ClientOptionFunc) (*Client, error) {
	return DialContext(context.Background(), options...)
}

// DialContext will connect to Elasticsearch, just like NewClient does.
//
// The context is honoured in terms of e.g. cancellation.
func DialContext(ctx context.Context, options ...ClientOptionFunc) (*Client, error) {
	// Set up the client
	c := &Client{
		c:                         http.DefaultClient,
		conns:                     make([]*conn, 0),
		cindex:                    -1,
		scheme:                    DefaultScheme,
		decoder:                   &DefaultDecoder{},
		healthcheckEnabled:        DefaultHealthcheckEnabled,
		healthcheckTimeoutStartup: DefaultHealthcheckTimeoutStartup,
		healthcheckTimeout:        DefaultHealthcheckTimeout,
		healthcheckInterval:       DefaultHealthcheckInterval,
		healthcheckStop:           make(chan bool),
		snifferEnabled:            DefaultSnifferEnabled,
		snifferTimeoutStartup:     DefaultSnifferTimeoutStartup,
		snifferTimeout:            DefaultSnifferTimeout,
		snifferInterval:           DefaultSnifferInterval,
		snifferStop:               make(chan bool),
		sendGetBodyAs:             DefaultSendGetBodyAs,
		gzipEnabled:               DefaultGzipEnabled,
		retrier:                   noRetries, // no retries by default
		retryStatusCodes:          nil,       // no automatic retries for specific HTTP status codes
		deprecationlog:            noDeprecationLog,
	}

	// Run the options on it
	for _, option := range options {
		if err := option(c); err != nil {
			return nil, err
		}
	}

	// Use a default URL and normalize them
	if len(c.urls) == 0 {
		c.urls = []string{DefaultURL}
	}
	c.urls = canonicalize(c.urls...)

	// If the URLs have auth info, use them here as an alternative to SetBasicAuth
	if c.basicAuthUsername == "" && c.basicAuthPassword == "" {
		for _, urlStr := range c.urls {
			u, err := url.Parse(urlStr)
			if err == nil && u.User != nil {
				c.basicAuthUsername = u.User.Username()
				c.basicAuthPassword, _ = u.User.Password()
				break
			}
		}
	}

	// Check if we can make a request to any of the specified URLs
	if c.healthcheckEnabled {
		if err := c.startupHealthcheck(ctx, c.healthcheckTimeoutStartup); err != nil {
			return nil, err
		}
	}

	
	for _, url := range c.urls {
		c.conns = append(c.conns, newConn(url, url))
	}

	if c.healthcheckEnabled {
		// Perform an initial health check
		c.healthcheck(ctx, c.healthcheckTimeoutStartup, true)
	}
	// Ensure that we have at least one connection available
	if err := c.mustActiveConn(); err != nil {
		return nil, err
	}

	if c.healthcheckEnabled {
		go c.healthchecker() // start goroutine periodically ping all nodes of the cluster
	}

	c.mu.Lock()
	c.running = true
	c.mu.Unlock()

	return c, nil
}

// DialWithConfig will use the configuration settings parsed from config package
// to connect to Elasticsearch.
//
// The context is honoured in terms of e.g. cancellation.
func DialWithConfig(ctx context.Context, cfg *config.Config) (*Client, error) {
	options, err := configToOptions(cfg)
	if err != nil {
		return nil, err
	}
	return DialContext(ctx, options...)
}

func configToOptions(cfg *config.Config) ([]ClientOptionFunc, error) {
	var options []ClientOptionFunc
	if cfg != nil {
		if cfg.URL != "" {
			options = append(options, SetURL(cfg.URL))
		}
		if cfg.Errorlog != "" {
			f, err := os.OpenFile(cfg.Errorlog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return nil, errors.Wrap(err, "unable to initialize error log")
			}
			l := log.New(f, "", 0)
			options = append(options, SetErrorLog(l))
		}
		if cfg.Tracelog != "" {
			f, err := os.OpenFile(cfg.Tracelog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return nil, errors.Wrap(err, "unable to initialize trace log")
			}
			l := log.New(f, "", 0)
			options = append(options, SetTraceLog(l))
		}
		if cfg.Infolog != "" {
			f, err := os.OpenFile(cfg.Infolog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return nil, errors.Wrap(err, "unable to initialize info log")
			}
			l := log.New(f, "", 0)
			options = append(options, SetInfoLog(l))
		}
		if cfg.Username != "" || cfg.Password != "" {
			options = append(options, SetBasicAuth(cfg.Username, cfg.Password))
		}
		if cfg.Sniff != nil {
			options = append(options, SetSniff(*cfg.Sniff))
		}
		if cfg.Healthcheck != nil {
			options = append(options, SetHealthcheck(*cfg.Healthcheck))
		}
	}
	return options, nil
}

// SetHttpClient can be used to specify the http.Client to use when making
// HTTP requests to Elasticsearch.
func SetHttpClient(httpClient Doer) ClientOptionFunc {
	return func(c *Client) error {
		if httpClient != nil {
			c.c = httpClient
		} else {
			c.c = http.DefaultClient
		}
		return nil
	}
}

// SetBasicAuth can be used to specify the HTTP Basic Auth credentials to
// use when making HTTP requests to Elasticsearch.
func SetBasicAuth(username, password string) ClientOptionFunc {
	return func(c *Client) error {
		c.basicAuthUsername = username
		c.basicAuthPassword = password
		return nil
	}
}

// SetURL defines the URL endpoints of the Elasticsearch nodes. Notice that
// when sniffing is enabled, these URLs are used to initially sniff the
// cluster on startup.
func SetURL(urls ...string) ClientOptionFunc {
	return func(c *Client) error {
		switch len(urls) {
		case 0:
			c.urls = []string{DefaultURL}
		default:
			c.urls = urls
		}
		// Check URLs
		for _, urlStr := range c.urls {
			if _, err := url.Parse(urlStr); err != nil {
				return err
			}
		}
		return nil
	}
}

// SetScheme sets the HTTP scheme to look for when sniffing (http or https).
// This is http by default.
func SetScheme(scheme string) ClientOptionFunc {
	return func(c *Client) error {
		c.scheme = scheme
		return nil
	}
}

// SetSniff enables or disables the sniffer (enabled by default).
func SetSniff(enabled bool) ClientOptionFunc {
	return func(c *Client) error {
		c.snifferEnabled = enabled
		return nil
	}
}

// SetSnifferTimeoutStartup sets the timeout for the sniffer that is used
// when creating a new client. The default is 5 seconds. Notice that the
// timeout being used for subsequent sniffing processes is set with
// SetSnifferTimeout.
func SetSnifferTimeoutStartup(timeout time.Duration) ClientOptionFunc {
	return func(c *Client) error {
		c.snifferTimeoutStartup = timeout
		return nil
	}
}

// SetSnifferTimeout sets the timeout for the sniffer that finds the
// nodes in a cluster. The default is 2 seconds. Notice that the timeout
// used when creating a new client on startup is usually greater and can
// be set with SetSnifferTimeoutStartup.
func SetSnifferTimeout(timeout time.Duration) ClientOptionFunc {
	return func(c *Client) error {
		c.snifferTimeout = timeout
		return nil
	}
}

// SetSnifferInterval sets the interval between two sniffing processes.
// The default interval is 15 minutes.
func SetSnifferInterval(interval time.Duration) ClientOptionFunc {
	return func(c *Client) error {
		c.snifferInterval = interval
		return nil
	}
}

// SetHealthcheck enables or disables healthchecks (enabled by default).
func SetHealthcheck(enabled bool) ClientOptionFunc {
	return func(c *Client) error {
		c.healthcheckEnabled = enabled
		return nil
	}
}

// SetHealthcheckTimeoutStartup sets the timeout for the initial health check.
// The default timeout is 5 seconds (see DefaultHealthcheckTimeoutStartup).
// Notice that timeouts for subsequent health checks can be modified with
// SetHealthcheckTimeout.
func SetHealthcheckTimeoutStartup(timeout time.Duration) ClientOptionFunc {
	return func(c *Client) error {
		c.healthcheckTimeoutStartup = timeout
		return nil
	}
}

// SetHealthcheckTimeout sets the timeout for periodic health checks.
// The default timeout is 1 second (see DefaultHealthcheckTimeout).
// Notice that a different (usually larger) timeout is used for the initial
// healthcheck, which is initiated while creating a new client.
// The startup timeout can be modified with SetHealthcheckTimeoutStartup.
func SetHealthcheckTimeout(timeout time.Duration) ClientOptionFunc {
	return func(c *Client) error {
		c.healthcheckTimeout = timeout
		return nil
	}
}

// SetHealthcheckInterval sets the interval between two health checks.
// The default interval is 60 seconds.
func SetHealthcheckInterval(interval time.Duration) ClientOptionFunc {
	return func(c *Client) error {
		c.healthcheckInterval = interval
		return nil
	}
}

// SetMaxRetries sets the maximum number of retries before giving up when
// performing a HTTP request to Elasticsearch.
//
// Deprecated: Replace with a Retry implementation.
func SetMaxRetries(maxRetries int) ClientOptionFunc {
	return func(c *Client) error {
		if maxRetries < 0 {
			return errors.New("MaxRetries must be greater than or equal to 0")
		} else if maxRetries == 0 {
			c.retrier = noRetries
		} else {
			// Create a Retrier that will wait for 100ms (+/- jitter) between requests.
			// This resembles the old behavior with maxRetries.
			ticks := make([]int, maxRetries)
			for i := 0; i < len(ticks); i++ {
				ticks[i] = 100
			}
			backoff := NewSimpleBackoff(ticks...)
			c.retrier = NewBackoffRetrier(backoff)
		}
		return nil
	}
}

// SetGzip enables or disables gzip compression (disabled by default).
func SetGzip(enabled bool) ClientOptionFunc {
	return func(c *Client) error {
		c.gzipEnabled = enabled
		return nil
	}
}

// SetDecoder sets the Decoder to use when decoding data from Elasticsearch.
// DefaultDecoder is used by default.
func SetDecoder(decoder Decoder) ClientOptionFunc {
	return func(c *Client) error {
		if decoder != nil {
			c.decoder = decoder
		} else {
			c.decoder = &DefaultDecoder{}
		}
		return nil
	}
}

// SetRequiredPlugins can be used to indicate that some plugins are required
// before a Client will be created.
func SetRequiredPlugins(plugins ...string) ClientOptionFunc {
	return func(c *Client) error {
		if c.requiredPlugins == nil {
			c.requiredPlugins = make([]string, 0)
		}
		c.requiredPlugins = append(c.requiredPlugins, plugins...)
		return nil
	}
}

// SetErrorLog sets the logger for critical messages like nodes joining
// or leaving the cluster or failing requests. It is nil by default.
func SetErrorLog(logger Logger) ClientOptionFunc {
	return func(c *Client) error {
		c.errorlog = logger
		return nil
	}
}

// SetInfoLog sets the logger for informational messages, e.g. requests
// and their response times. It is nil by default.
func SetInfoLog(logger Logger) ClientOptionFunc {
	return func(c *Client) error {
		c.infolog = logger
		return nil
	}
}

// SetTraceLog specifies the log.Logger to use for output of HTTP requests
// and responses which is helpful during debugging. It is nil by default.
func SetTraceLog(logger Logger) ClientOptionFunc {
	return func(c *Client) error {
		c.tracelog = logger
		return nil
	}
}

// SetSendGetBodyAs specifies the HTTP method to use when sending a GET request
// with a body. It is GET by default.
func SetSendGetBodyAs(httpMethod string) ClientOptionFunc {
	return func(c *Client) error {
		c.sendGetBodyAs = httpMethod
		return nil
	}
}

// SetRetrier specifies the retry strategy that handles errors during
// HTTP request/response with Elasticsearch.
func SetRetrier(retrier Retrier) ClientOptionFunc {
	return func(c *Client) error {
		if retrier == nil {
			retrier = noRetries // no retries by default
		}
		c.retrier = retrier
		return nil
	}
}

// SetRetryStatusCodes specifies the HTTP status codes where the client
// will retry automatically. Notice that retries call the specified retrier,
// so calling SetRetryStatusCodes without setting a Retrier won't do anything
// for retries.
func SetRetryStatusCodes(statusCodes ...int) ClientOptionFunc {
	return func(c *Client) error {
		c.retryStatusCodes = statusCodes
		return nil
	}
}

// SetHeaders adds a list of default HTTP headers that will be added to
// each requests executed by PerformRequest.
func SetHeaders(headers http.Header) ClientOptionFunc {
	return func(c *Client) error {
		c.headers = headers
		return nil
	}
}

// String returns a string representation of the client status.
func (c *Client) String() string {
	c.connsMu.Lock()
	conns := c.conns
	c.connsMu.Unlock()

	var buf bytes.Buffer
	for i, conn := range conns {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(conn.String())
	}
	return buf.String()
}

// IsRunning returns true if the background processes of the client are
// running, false otherwise.
func (c *Client) IsRunning() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.running
}

// Start starts the background processes like sniffing the cluster and
// periodic health checks. You don't need to run Start when creating a
// client with NewClient; the background processes are run by default.
//
// If the background processes are already running, this is a no-op.
func (c *Client) Start() {
	c.mu.RLock()
	if c.running {
		c.mu.RUnlock()
		return
	}
	c.mu.RUnlock()

	if c.healthcheckEnabled {
		go c.healthchecker()
	}

	c.mu.Lock()
	c.running = true
	c.mu.Unlock()

	c.infof("elastic: client started")
}

// Stop stops the background processes that the client is running,
// i.e. sniffing the cluster periodically and running health checks
// on the nodes.
//
// If the background processes are not running, this is a no-op.
func (c *Client) Stop() {
	c.mu.RLock()
	if !c.running {
		c.mu.RUnlock()
		return
	}
	c.mu.RUnlock()

	if c.healthcheckEnabled {
		c.healthcheckStop <- true
		<-c.healthcheckStop
	}

	if c.snifferEnabled {
		c.snifferStop <- true
		<-c.snifferStop
	}

	c.mu.Lock()
	c.running = false
	c.mu.Unlock()

	c.infof("elastic: client stopped")
}

// errorf logs to the error log.
func (c *Client) errorf(format string, args ...interface{}) {
	if c.errorlog != nil {
		c.errorlog.Printf(format, args...)
	}
}

// infof logs informational messages.
func (c *Client) infof(format string, args ...interface{}) {
	if c.infolog != nil {
		c.infolog.Printf(format, args...)
	}
}

// tracef logs to the trace log.
func (c *Client) tracef(format string, args ...interface{}) {
	if c.tracelog != nil {
		c.tracelog.Printf(format, args...)
	}
}

// dumpRequest dumps the given HTTP request to the trace log.
func (c *Client) dumpRequest(r *http.Request) {
	if c.tracelog != nil {
		out, err := httputil.DumpRequestOut(r, true)
		if err == nil {
			c.tracef("%s\n", string(out))
		}
	}
}

// dumpResponse dumps the given HTTP response to the trace log.
func (c *Client) dumpResponse(resp *http.Response) {
	if c.tracelog != nil {
		out, err := httputil.DumpResponse(resp, true)
		if err == nil {
			c.tracef("%s\n", string(out))
		}
	}
}

// extractHostname returns the URL from the http.publish_address setting.
func (c *Client) extractHostname(scheme, address string) string {
	var (
		host string
		port string

		addrs = strings.Split(address, "/")
		ports = strings.Split(address, ":")
	)

	if len(addrs) > 1 {
		host = addrs[0]
	} else {
		host = strings.Split(addrs[0], ":")[0]
	}
	port = ports[len(ports)-1]

	return fmt.Sprintf("%s://%s:%s", scheme, host, port)
}

// updateConns updates the clients' connections with new information
// gather by a sniff operation.
func (c *Client) updateConns(conns []*conn) {
	c.connsMu.Lock()

	// Build up new connections:
	// If we find an existing connection, use that (including no. of failures etc.).
	// If we find a new connection, add it.
	var newConns []*conn
	for _, conn := range conns {
		var found bool
		for _, oldConn := range c.conns {
			// Notice that e.g. in a Kubernetes cluster the NodeID might be
			// stable while the URL has changed.
			if oldConn.NodeID() == conn.NodeID() && oldConn.URL() == conn.URL() {
				// Take over the old connection
				newConns = append(newConns, oldConn)
				found = true
				break
			}
		}
		if !found {
			// New connection didn't exist, so add it to our list of new conns.
			c.infof("elastic: %s joined the cluster", conn.URL())
			newConns = append(newConns, conn)
		}
	}

	c.conns = newConns
	c.cindex = -1
	c.connsMu.Unlock()
}

// healthchecker periodically runs healthcheck.
func (c *Client) healthchecker() {
	c.mu.RLock()
	timeout := c.healthcheckTimeout
	interval := c.healthcheckInterval
	c.mu.RUnlock()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-c.healthcheckStop:
			// we are asked to stop, so we signal back that we're stopping now
			c.healthcheckStop <- true
			return
		case <-ticker.C:
			c.healthcheck(context.Background(), timeout, false)
		}
	}
}

// healthcheck does a health check on all nodes in the cluster. Depending on
// the node state, it marks connections as dead, sets them alive etc.
// If healthchecks are disabled and force is false, this is a no-op.
// The timeout specifies how long to wait for a response from Elasticsearch.
func (c *Client) healthcheck(parentCtx context.Context, timeout time.Duration, force bool) {
	c.mu.RLock()
	if !c.healthcheckEnabled && !force {
		c.mu.RUnlock()
		return
	}
	headers := c.headers
	basicAuth := c.basicAuthUsername != "" || c.basicAuthPassword != ""
	basicAuthUsername := c.basicAuthUsername
	basicAuthPassword := c.basicAuthPassword
	c.mu.RUnlock()

	c.connsMu.RLock()
	conns := c.conns
	c.connsMu.RUnlock()

	for _, conn := range conns {
		// Run the HEAD request against ES with a timeout
		ctx, cancel := context.WithTimeout(parentCtx, timeout)
		defer cancel()

		// Goroutine executes the HTTP request, returns an error and sets status
		var status int
		errc := make(chan error, 1)
		go func(url string) {
			req, err := NewRequest("HEAD", url)
			if err != nil {
				errc <- err
				return
			}
			if basicAuth {
				req.SetBasicAuth(basicAuthUsername, basicAuthPassword)
			}
			if len(headers) > 0 {
				for key, values := range headers {
					for _, v := range values {
						req.Header.Add(key, v)
					}
				}
			}
			if req.Header.Get("User-Agent") == "" {
				req.Header.Add("User-Agent", "elastic/"+Version+" ("+runtime.GOOS+"-"+runtime.GOARCH+")")
			}
			res, err := c.c.Do((*http.Request)(req).WithContext(ctx))
			if res != nil {
				status = res.StatusCode
				if res.Body != nil {
					res.Body.Close()
				}
			}
			errc <- err
		}(conn.URL())

		// Wait for the Goroutine (or its timeout)
		select {
		case <-ctx.Done(): // timeout
			c.errorf("elastic: %s is dead", conn.URL())
			conn.MarkAsDead()
		case err := <-errc:
			if err != nil {
				c.errorf("elastic: %s is dead", conn.URL())
				conn.MarkAsDead()
				break
			}
			if status >= 200 && status < 300 {
				conn.MarkAsAlive()
			} else {
				conn.MarkAsDead()
				c.errorf("elastic: %s is dead [status=%d]", conn.URL(), status)
			}
		}
	}
}

// startupHealthcheck is used at startup to check if the server is available
// at all.
func (c *Client) startupHealthcheck(parentCtx context.Context, timeout time.Duration) error {
	c.mu.Lock()
	urls := c.urls
	headers := c.headers
	basicAuth := c.basicAuthUsername != "" || c.basicAuthPassword != ""
	basicAuthUsername := c.basicAuthUsername
	basicAuthPassword := c.basicAuthPassword
	c.mu.Unlock()

	// If we don't get a connection after "timeout", we bail.
	var lastErr error
	start := time.Now()
	done := false
	for !done {
		for _, url := range urls {
			req, err := http.NewRequest("HEAD", url, nil)
			if err != nil {
				return err
			}
			if basicAuth {
				req.SetBasicAuth(basicAuthUsername, basicAuthPassword)
			}
			if len(headers) > 0 {
				for key, values := range headers {
					for _, v := range values {
						req.Header.Add(key, v)
					}
				}
			}
			ctx, cancel := context.WithTimeout(parentCtx, timeout)
			defer cancel()
			req = req.WithContext(ctx)
			res, err := c.c.Do(req)
			if err != nil {
				lastErr = err
			} else if res.StatusCode >= 200 && res.StatusCode < 300 {
				return nil
			} else if res.StatusCode == http.StatusUnauthorized {
				lastErr = &Error{Status: res.StatusCode}
			}
		}
		select {
		case <-parentCtx.Done():
			lastErr = parentCtx.Err()
			done = true
		case <-time.After(1 * time.Second):
			if time.Since(start) > timeout {
				done = true
			}
		}
	}
	if lastErr != nil {
		if IsContextErr(lastErr) || IsUnauthorized(lastErr) {
			return lastErr
		}
		return errors.Wrapf(ErrNoClient, "health check timeout: %v", lastErr)
	}
	return errors.Wrap(ErrNoClient, "health check timeout")
}

// next returns the next available connection, or ErrNoClient.
func (c *Client) next() (*conn, error) {
	// We do round-robin here.
	// TODO(oe) This should be a pluggable strategy, like the Selector in the official clients.
	c.connsMu.Lock()
	defer c.connsMu.Unlock()

	i := 0
	numConns := len(c.conns)
	for {
		i++
		if i > numConns {
			break // we visited all conns: they all seem to be dead
		}
		c.cindex++
		if c.cindex >= numConns {
			c.cindex = 0
		}
		conn := c.conns[c.cindex]
		if !conn.IsDead() {
			return conn, nil
		}
	}

	// We have a deadlock here: All nodes are marked as dead.
	// If sniffing is disabled, connections will never be marked alive again.
	// So we are marking them as alive--if sniffing is disabled.
	// They'll then be picked up in the next call to PerformRequest.
	if !c.snifferEnabled {
		c.errorf("elastic: all %d nodes marked as dead; resurrecting them to prevent deadlock", len(c.conns))
		for _, conn := range c.conns {
			conn.MarkAsAlive()
		}
	}

	// We tried hard, but there is no node available
	return nil, errors.Wrap(ErrNoClient, "no available connection")
}

// mustActiveConn returns nil if there is an active connection,
// otherwise ErrNoClient is returned.
func (c *Client) mustActiveConn() error {
	c.connsMu.Lock()
	defer c.connsMu.Unlock()

	for _, c := range c.conns {
		if !c.IsDead() {
			return nil
		}
	}
	return errors.Wrap(ErrNoClient, "no active connection found")
}

// -- PerformRequest --

// PerformRequestOptions must be passed into PerformRequest.
type PerformRequestOptions struct {
	Method           string
	Path             string
	Params           url.Values
	Body             interface{}
	ContentType      string
	IgnoreErrors     []int
	Retrier          Retrier
	RetryStatusCodes []int
	Headers          http.Header
	MaxResponseSize  int64
	Stream           bool
}

// PerformRequest does a HTTP request to Elasticsearch.
// It returns a response (which might be nil) and an error on failure.
//
// Optionally, a list of HTTP error codes to ignore can be passed.
// This is necessary for services that expect e.g. HTTP status 404 as a
// valid outcome (Exists, IndicesExists, IndicesTypeExists).
//
// If Stream is set, the returned BodyReader field must be closed, even
// if PerformRequest returns an error.
func (c *Client) PerformRequest(ctx context.Context, opt PerformRequestOptions) (*Response, error) {
	start := time.Now().UTC()

	c.mu.RLock()
	timeout := c.healthcheckTimeout
	basicAuth := c.basicAuthUsername != "" || c.basicAuthPassword != ""
	basicAuthUsername := c.basicAuthUsername
	basicAuthPassword := c.basicAuthPassword
	sendGetBodyAs := c.sendGetBodyAs
	gzipEnabled := c.gzipEnabled
	healthcheckEnabled := c.healthcheckEnabled
	retrier := c.retrier
	if opt.Retrier != nil {
		retrier = opt.Retrier
	}
	retryStatusCodes := c.retryStatusCodes
	if opt.RetryStatusCodes != nil {
		retryStatusCodes = opt.RetryStatusCodes
	}
	defaultHeaders := c.headers
	c.mu.RUnlock()

	// retry returns true if statusCode indicates the request is to be retried
	retry := func(statusCode int) bool {
		for _, code := range retryStatusCodes {
			if code == statusCode {
				return true
			}
		}
		return false
	}

	var err error
	var conn *conn
	var req *Request
	var resp *Response
	var retried bool
	var n int

	// Change method if sendGetBodyAs is specified.
	if opt.Method == "GET" && opt.Body != nil && sendGetBodyAs != "GET" {
		opt.Method = sendGetBodyAs
	}

	for {
		pathWithParams := opt.Path
		if len(opt.Params) > 0 {
			pathWithParams += "?" + opt.Params.Encode()
		}

		// Get a connection
		conn, err = c.next()
		if errors.Cause(err) == ErrNoClient {
			n++
			if !retried {
				// Force a healtcheck as all connections seem to be dead.
				c.healthcheck(ctx, timeout, false)
				if healthcheckEnabled {
					retried = true
					continue
				}
			}
			wait, ok, rerr := retrier.Retry(ctx, n, nil, nil, err)
			if rerr != nil {
				return nil, rerr
			}
			if !ok {
				return nil, err
			}
			retried = true
			time.Sleep(wait)
			continue // try again
		}
		if err != nil {
			c.errorf("elastic: cannot get connection from pool")
			return nil, err
		}

		req, err = NewRequest(opt.Method, conn.URL()+pathWithParams)
		if err != nil {
			c.errorf("elastic: cannot create request for %s %s: %v", strings.ToUpper(opt.Method), conn.URL()+pathWithParams, err)
			return nil, err
		}
		if basicAuth {
			req.SetBasicAuth(basicAuthUsername, basicAuthPassword)
		}
		if opt.ContentType != "" {
			req.Header.Set("Content-Type", opt.ContentType)
		}
		for key, value := range opt.Headers {
			for _, v := range value {
				req.Header.Add(key, v)
			}
		}
		if len(defaultHeaders) > 0 {
			for key, value := range defaultHeaders {
				for _, v := range value {
					req.Header.Add(key, v)
				}
			}
		}
		if req.Header.Get("User-Agent") == "" {
			req.Header.Set("User-Agent", "elastic/"+Version+" ("+runtime.GOOS+"-"+runtime.GOARCH+")")
		}

		// Set body
		if opt.Body != nil {
			err = req.SetBody(opt.Body, gzipEnabled)
			if err != nil {
				c.errorf("elastic: couldn't set body %+v for request: %v", opt.Body, err)
				return nil, err
			}
		}

		// Tracing
		c.dumpRequest((*http.Request)(req))

		// Get response
		res, err := c.c.Do((*http.Request)(req).WithContext(ctx))
		if IsContextErr(err) {
			// Proceed, but don't mark the node as dead
			return nil, err
		}
		if err != nil {
			n++
			wait, ok, rerr := retrier.Retry(ctx, n, (*http.Request)(req), res, err)
			if rerr != nil {
				c.errorf("elastic: %s is dead", conn.URL())
				conn.MarkAsDead()
				return nil, rerr
			}
			if !ok {
				c.errorf("elastic: %s is dead", conn.URL())
				conn.MarkAsDead()
				return nil, err
			}
			retried = true
			time.Sleep(wait)
			continue // try again
		}
		if retry(res.StatusCode) {
			n++
			wait, ok, rerr := retrier.Retry(ctx, n, (*http.Request)(req), res, err)
			if rerr != nil {
				c.errorf("elastic: %s is dead", conn.URL())
				conn.MarkAsDead()
				return nil, rerr
			}
			if ok {
				// retry
				retried = true
				time.Sleep(wait)
				continue // try again
			}
		}

		if !opt.Stream {
			defer res.Body.Close()
		}

		// Tracing
		c.dumpResponse(res)

		// Log deprecation warnings as errors
		if len(res.Header["Warning"]) > 0 {
			c.deprecationlog((*http.Request)(req), res)
			for _, warning := range res.Header["Warning"] {
				c.errorf("Deprecation warning: %s", warning)
			}
		}

		// Check for errors
		if err := checkResponse((*http.Request)(req), res, opt.IgnoreErrors...); err != nil {
			// No retry if request succeeded
			// We still try to return a response.
			resp, _ = c.newResponse(res, opt.MaxResponseSize, opt.Stream)
			return resp, err
		}

		// We successfully made a request with this connection
		conn.MarkAsHealthy()

		resp, err = c.newResponse(res, opt.MaxResponseSize, opt.Stream)
		if err != nil {
			return nil, err
		}

		break
	}

	duration := time.Now().UTC().Sub(start)
	c.infof("%s %s [status:%d, request:%.3fs]",
		strings.ToUpper(opt.Method),
		req.URL.Redacted(),
		resp.StatusCode,
		float64(int64(duration/time.Millisecond))/1000)

	return resp, nil
}

// TermVectors returns information and statistics on terms in the fields
// of a particular document.
func (c *Client) TermVectors(index string) *TermvectorsService {
	builder := NewTermvectorsService(c)
	builder = builder.Index(index)
	return builder
}

// MultiTermVectors returns information and statistics on terms in the fields
// of multiple documents.
func (c *Client) MultiTermVectors() *MultiTermvectorService {
	return NewMultiTermvectorService(c)
}

// -- Search APIs --

// Search is the entry point for searches.
func (c *Client) Search(indices ...string) *SearchService {
	return NewSearchService(c).Index(indices...)
}

// MultiSearch is the entry point for multi searches.
func (c *Client) MultiSearch() *MultiSearchService {
	return NewMultiSearchService(c)
}

// TODO Search Template
// TODO Search Exists API

// Validate allows a user to validate a potentially expensive query without executing it.
func (c *Client) Validate(indices ...string) *ValidateService {
	return NewValidateService(c).Index(indices...)
}

// SearchShards returns statistical information about nodes and shards.
func (c *Client) SearchShards(indices ...string) *SearchShardsService {
	return NewSearchShardsService(c).Index(indices...)
}

// Scroll through documents. Use this to efficiently scroll through results
// while returning the results to a client.
func (c *Client) Scroll(indices ...string) *ScrollService {
	return NewScrollService(c).Index(indices...)
}

// OpenPointInTime opens a new Point in Time.
func (c *Client) OpenPointInTime(indices ...string) *OpenPointInTimeService {
	return NewOpenPointInTimeService(c).Index(indices...)
}

// ClosePointInTime closes an existing Point in Time.
func (c *Client) ClosePointInTime(id string) *ClosePointInTimeService {
	return NewClosePointInTimeService(c).ID(id)
}

// -- Scripting APIs --

// GetScript reads a stored script in Elasticsearch.
// Use PutScript for storing a script.
func (c *Client) GetScript() *GetScriptService {
	return NewGetScriptService(c)
}

// PutScript allows saving a stored script in Elasticsearch.
func (c *Client) PutScript() *PutScriptService {
	return NewPutScriptService(c)
}

// DeleteScript allows removing a stored script from Elasticsearch.
func (c *Client) DeleteScript() *DeleteScriptService {
	return NewDeleteScriptService(c)
}
