package rpc

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"
)

// ServiceError is an error from server.
type ServiceError string

func (e ServiceError) Error() string {
	return string(e)
}

const (
	// ReaderBufferSize is used for bufio reader.
	ReaderBufferSize = 16 * 1024
	module           = "client.srpc"
	inputModule      = "client.srpc.input"
	heartbeatModule  = "client.srpc.heartbeat"
)

type seqKey struct{}

// Client is the rpc client using http protocol
type Client struct {
	client     *http.Client
	localAddr  net.Addr
	remoteAddr net.Addr
}

func (c *Client) Close() error {
	return nil
}

// LocalAddr implement ConnInfo interface
func (c *Client) LocalAddr() net.Addr {
	if c != nil {
		return c.localAddr
	}
	return nil
}

// RemoteAddr implement ConnInfo interface
func (c *Client) RemoteAddr() net.Addr {
	if c != nil {
		return c.remoteAddr
	}
	return nil
}

func (c *Client) decodeResponse(contentType string,
	body []byte, response interface{}) error {
	if contentType == "application/json" {

	} else if contentType == "application/protobuf" {
		return fmt.Errorf("invalid argument type of response: %T", response)
	}
	return fmt.Errorf("unknown content type: %s", contentType)
}

type Option struct {
	// TLSConfig for tcp and quic
	TLSConfig *tls.Config

	// ConnectTimeout sets timeout for dialing
	ConnectTimeout time.Duration
	// IdleTimeout sets max idle time for underlying net.Conns
	IdleTimeout time.Duration

	// Persistent indicate short or long connection
	Persistent bool
	ConnSize   int
}

func NewDefaultOption() *Option {
	return &Option{
		ConnectTimeout: time.Duration(3) * time.Second,
		IdleTimeout:    time.Duration(180) * time.Second,
	}
}

var ClientMap = make(map[string]map[string]*http.Client)

func newClient(option *Option) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout: option.ConnectTimeout,
			}).DialContext,
			MaxIdleConnsPerHost: option.ConnSize,
			IdleConnTimeout:     option.IdleTimeout,
			DisableKeepAlives:   !option.Persistent,
			TLSClientConfig:     option.TLSConfig,
		},
	}
}

func GetClient(serviceName string, configKey string) (c *http.Client, ok bool) {
	var m map[string]*http.Client
	m, ok = ClientMap[serviceName]
	if ok {
		c, ok = m[configKey]
		return c, ok
	}
	return nil, false
}
