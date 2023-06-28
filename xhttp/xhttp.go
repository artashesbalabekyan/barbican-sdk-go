package xhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"aead.dev/mem"
	"github.com/artashesbalabekyan/barbican-sdk-go/xerror"
)

// authToken is a Barbican authentication token.
// It can be used to authenticate API requests.
type authToken struct {
	Key    string
	Expiry time.Time
}

// client is a Barbican REST API client
// responsible for fetching and renewing
// authentication tokens.
type Client struct {
	http.Client
	config Config

	lock  sync.Mutex
	token authToken
}

// Add auth header to request
func (c *Client) setAuthHeader(ctx context.Context, config Config, h *http.Header) error {
	if c.token.Expiry.Unix() < time.Now().Unix() {
		err := c.Authenticate(ctx, config)
		if err != nil {
			return err
		}
	}
	h.Add("X-Auth-Token", string(c.token.Key))
	return nil
}

// Authenticate tries to obtain a new authentication token
// from the given Barbican endpoint via the given credentials.
//
// Authenticate should be called to obtain the first authentication
// token. This token can then be renewed via RenewApiToken.
func (c *Client) Authenticate(ctx context.Context, config Config) error {
	r := AuthRequest{}
	r.Auth.Identity.Methods = []string{"password"}
	r.Auth.Identity.Password.User.Domain.Name = config.Login.UserDomainName
	r.Auth.Identity.Password.User.Name = config.Login.Username
	r.Auth.Identity.Password.User.Password = config.Login.Password
	r.Auth.Scope.Project.Domain.Name = config.Login.ProjectDomain
	r.Auth.Scope.Project.Name = config.Login.ProjectName

	body, err := json.Marshal(&r)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/v3/auth/tokens", config.Login.AuthUrl)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		err := parseErrorResponse(resp)
		return fmt.Errorf("%s: %v", resp.Status, err)
	}

	const MaxSize = 1 * mem.MiB // An auth. token response should not exceed 1 MiB
	var response AuthResponse
	if err = json.NewDecoder(mem.LimitReader(resp.Body, MaxSize)).Decode(&response); err != nil {
		return err
	}
	if response.Token.ExpiresAt == "" {
		return errors.New("server response does not contain a token expiry")
	}
	expiry, err := time.Parse(time.RFC3339Nano, response.Token.ExpiresAt)
	if err != nil || expiry.Unix() < time.Now().Unix() {
		return errors.New("server response does not contain a valid token expiry")
	}

	token := resp.Header.Get("x-subject-token")
	if token == "" {
		return errors.New("server response does not contain a token header")
	}

	c.lock.Lock()
	c.token = authToken{
		Key:    token,
		Expiry: expiry,
	}
	c.lock.Unlock()
	return nil
}

// parseErrorResponse returns an error containing
// the response status code and response body
// as error message if the response is an error
// response - i.e. status code >= 400.
//
// If the response status code is < 400, e.g. 200 OK,
// parseErrorResponse returns nil and does not attempt
// to read or close the response body.
//
// If resp is an error response, parseErrorResponse reads
// and closes the response body.
func parseErrorResponse(resp *http.Response) error {
	if resp.StatusCode < 400 {
		return nil
	}
	if resp.Body == nil {
		return xerror.NewError(resp.StatusCode, resp.Status)
	}
	defer resp.Body.Close()

	const MaxSize = 1 * mem.MiB
	size := mem.Size(resp.ContentLength)
	if size < 0 || size > MaxSize {
		size = MaxSize
	}

	if contentType := strings.TrimSpace(resp.Header.Get("Content-Type")); strings.HasPrefix(contentType, "application/json") {
		type Response struct {
			Message string `json:"message"`
		}
		var response Response
		if err := json.NewDecoder(mem.LimitReader(resp.Body, size)).Decode(&response); err != nil {
			return err
		}
		return xerror.NewError(resp.StatusCode, response.Message)
	}

	var sb strings.Builder
	if _, err := io.Copy(&sb, mem.LimitReader(resp.Body, size)); err != nil {
		return err
	}
	return xerror.NewError(resp.StatusCode, sb.String())
}

func (c *Client) request(ctx context.Context, method string, address string, payload []byte, params url.Values) ([]byte, error) {

	url, err := url.Parse(address)
	if err != nil {
		return nil, err
	}
	if len(params) > 0 {
		url.RawQuery = params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, url.String(), bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	err = c.setAuthHeader(ctx, c.config, &req.Header)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, xerror.NewError(resp.StatusCode, err.Error())
	}

	return body, nil
}

func (c *Client) HttpGet(ctx context.Context, url string, params url.Values) ([]byte, error) {
	return c.request(ctx, http.MethodGet, url, nil, params)
}

func (c *Client) HttpPost(ctx context.Context, url string, payload []byte, params url.Values) ([]byte, error) {
	return c.request(ctx, http.MethodPost, url, payload, params)
}

func (c *Client) HttpPut(ctx context.Context, url string, payload []byte, params url.Values) ([]byte, error) {
	return c.request(ctx, http.MethodPut, url, payload, params)
}

func (c *Client) HttpDelete(ctx context.Context, url string, payload []byte, params url.Values) ([]byte, error) {
	return c.request(ctx, http.MethodDelete, url, payload, params)
}
