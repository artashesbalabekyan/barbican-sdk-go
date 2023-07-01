package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/artashesbalabekyan/barbican-sdk-go/xhttp"
)

const (
	errValidatePrefix = "invalid barbican config: %s"
)

func New(ctx context.Context, config *xhttp.Config) (Conn, error) {
	return newConnection(ctx, config)
}

func newConnection(ctx context.Context, config *xhttp.Config) (Conn, error) {
	if err := validateConfig(config); err != nil {
		return nil, err
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
	return &Client{
		config: config,
		client: client,
	}, nil
}

// The path elements will not be URL-escaped.
func endpoint(endpoint string, elems ...string) string {
	endpoint = strings.TrimSpace(endpoint)
	url, _ := url.JoinPath(endpoint, elems...)
	return url
}

func validateConfig(config *xhttp.Config) error {
	if config == nil {
		return fmt.Errorf(errValidatePrefix, "config is nil")
	}
	if config.Endpoint == "" {
		return fmt.Errorf(errValidatePrefix, "endpoint is empty")
	}
	if config.Login.AuthUrl == "" {
		return fmt.Errorf(errValidatePrefix, "auth url is empty")
	}
	if config.Login.ProjectName == "" {
		return fmt.Errorf(errValidatePrefix, "project name is empty")
	}
	if config.Login.ProjectDomain == "" {
		return fmt.Errorf(errValidatePrefix, "project domain is empty")
	}
	if config.Login.Username == "" {
		return fmt.Errorf(errValidatePrefix, "username is empty")
	}
	if config.Login.UserDomainName == "" {
		return fmt.Errorf(errValidatePrefix, "user domain name is empty")
	}
	if config.Login.Password == "" {
		return fmt.Errorf(errValidatePrefix, "password is empty")
	}
	return nil
}
