package engine

import (
	"time"
)

const (
	ConfigPathPrefix = "frpc.http.clients"
	Module           = "http"
	Delimiter        = "."
)

type ClientConfig struct {
	// http or https
	Scheme string `toml:"scheme"`

	// cmlb_id、ip_port or domain name
	Address string `toml:"address"`

	// DisableKeepAlive, if true, disable HTTP keep-alive and
	// will only use the connection to the server for a single
	// HTTP request.
	DisableKeepAlive bool `toml:"disable_keep_alive"`

	// MaxIdleConn controls the maximum number of idle (keep-alive)
	// connections across all hosts. Zero means no limit.
	MaxIdleConn int `toml:"max_idle_conn"`

	// MaxIdlePerHost, if non-zero, controls the maximum idle
	// (keep-alive) connections to keep per-host. If zero,
	// net/http.DefaultMaxIdleConnsPerHost is used.
	MaxIdlePerHost int `toml:"max_idle_per_host"`

	// MaxConnPerHost optionally limits the total number of
	// connections per host, including connections in the dialing,
	// active, and idle states. On limit violation, dials will block.
	//
	// Zero means no limit.
	MaxConnPerHost int `toml:"max_conn_per_host"`

	// RequestTimeout specifies a time limit for requests made by this
	// Client. The timeout includes connection time, any
	// redirects, and reading the response body. The timer remains
	// running after Get, Head, Post, or Do return and will
	// interrupt reading of the Response.Body.
	//
	// A RequestTimeout of zero means no timeout.
	//
	// The Client cancels requests to the underlying Transport
	// as if the Request's Context ended.
	//
	// For compatibility, the Client will also use the deprecated
	// CancelRequest method on Transport if found. New
	// RoundTripper implementations should use the Request's Context
	// for cancellation instead of implementing CancelRequest.
	RequestTimeout time.Duration `toml:"request_timeout"`

	// IdleConnTimeout is the maximum amount of time an idle
	// (keep-alive) connection will remain idle before closing
	// itself.
	// Zero means no limit.
	IdleConnTimeout time.Duration `toml:"idle_conn_timeout"`

	// ResponseHeaderTimeout, if non-zero, specifies the amount of
	// time to wait for a server's response headers after fully
	// writing the request (including its body, if any). This
	// time does not include the time to read the response body.
	ResponseHeaderTimeout time.Duration `toml:"response_header_timeout"`

	// ExpectContinueTimeout, if non-zero, specifies the amount of
	// time to wait for a server's first response headers after fully
	// writing the request headers if the request has an
	// "Expect: 100-continue" header. Zero means no timeout and
	// causes the body to be sent immediately, without
	// waiting for the server to approve.
	// This time does not include the time to send the request header.
	ExpectContinueTimeout time.Duration `toml:"expect_continue_timeout"`

	// WriteBufferSize specifies the size of the write buffer used
	// when writing to the transport.
	// If zero, a default (currently 4KB) is used.
	WriteBufferSize int `toml:"write_buffer_size"`

	// ReadBufferSize specifies the size of the read buffer used
	// when reading from the transport.
	// If zero, a default (currently 4KB) is used.
	ReadBufferSize int `toml:"read_buffer_size"`

	// MaxResponseHeaderBytes specifies a limit on how many
	// response bytes are allowed in the server's response
	// header.
	//
	// Zero means to use a default limit.
	MaxResponseHeaderBytes int64 `toml:"max_response_header_bytes"`

	// TransportTimeout is maximum amount of time a dial will wait for a connect to complete.
	TransportTimeout time.Duration `toml:"transport_timeout"`

	// TransportKeepAlive specifies the interval between keep-alive probes for an active network connection
	TransportKeepAlive time.Duration `toml:"transport_keep_alive"`
}

var defaultClientConfig ClientConfig

func Parse() {
	//var ClientsConfigMap sync.Map
	//if err := defaultClientConfig.UnmarshalPartially(configPrefix, &frpcConfig); err != nil {
	//	log.Error(module, "failed to unmarshal frpc config", log.ErrorField(err))
	//	return errors.Wrap(err, "failed to unmarshal frpc config")
	//}
}
