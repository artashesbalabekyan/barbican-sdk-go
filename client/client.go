package client

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/artashesbalabekyan/barbican-sdk-go/xhttp"
)

func New(ctx context.Context, config *xhttp.Config) (*Connection, error) {
	return newConnection(ctx, config)
}

func newConnection(ctx context.Context, config *xhttp.Config) (*Connection, error) {
	if config.Endpoint == "" {
		return nil, errors.New("barican: endpoint is empty")
	}

	var tlsConfig *tls.Config
	client := &xhttp.Client{
		Client: http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
				Proxy:           http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   10 * time.Second,
					KeepAlive: 10 * time.Second,
					DualStack: true,
				}).DialContext,
				ForceAttemptHTTP2:     true,
				MaxIdleConns:          100,
				IdleConnTimeout:       30 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
		},
	}

	// Authenticate and get token
	if err := client.Authenticate(ctx, *config); err != nil {
		return nil, err
	}
	return &Connection{
		config: *config,
		client: client,
	}, nil
}

// The path elements will not be URL-escaped.
func endpoint(endpoint string, elems ...string) string {
	endpoint = strings.TrimSpace(endpoint)
	url, _ := url.JoinPath(endpoint, elems...)
	return url
}
