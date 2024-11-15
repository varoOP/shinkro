package sharedhttp

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"time"
)

var UserAgent string = fmt.Sprintf("shinkro/%v (%v;%v)", runtime.Version(), runtime.GOOS, runtime.GOARCH)

var Transport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second, // default transport value
		KeepAlive: 30 * time.Second, // default transport value
	}).DialContext,
	ForceAttemptHTTP2:     true,              // default is true; since HTTP/2 multiplexes a single TCP connection.
	MaxIdleConns:          100,               // default transport value
	MaxIdleConnsPerHost:   10,                // default is 2, so we want to increase the number to use establish more connections.
	IdleConnTimeout:       90 * time.Second,  // default transport value
	ResponseHeaderTimeout: 120 * time.Second, // servers can respond slowly - this should fix some portion of releases getting stuck as pending.
	TLSHandshakeTimeout:   10 * time.Second,  // default transport value
	ExpectContinueTimeout: 1 * time.Second,   // default transport value
	ReadBufferSize:        65536,
	WriteBufferSize:       65536,
	TLSClientConfig: &tls.Config{
		MinVersion: tls.VersionTLS12,
	},
}

var TransportTLSInsecure = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second, // default transport value
		KeepAlive: 30 * time.Second, // default transport value
	}).DialContext,
	ForceAttemptHTTP2:     true,              // default is true; since HTTP/2 multiplexes a single TCP connection.
	MaxIdleConns:          100,               // default transport value
	MaxIdleConnsPerHost:   10,                // default is 2, so we want to increase the number to use establish more connections.
	IdleConnTimeout:       90 * time.Second,  // default transport value
	ResponseHeaderTimeout: 120 * time.Second, // servers can respond slowly - this should fix some portion of releases getting stuck as pending.
	TLSHandshakeTimeout:   10 * time.Second,  // default transport value
	ExpectContinueTimeout: 1 * time.Second,   // default transport value
	ReadBufferSize:        65536,
	WriteBufferSize:       65536,
	TLSClientConfig: &tls.Config{
		InsecureSkipVerify: true,
	},
}
