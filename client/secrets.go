package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/artashesbalabekyan/barbican-sdk-go/xerror"
)

// Create stores the given key in Barbican if and only
// if no entry with the given name exists.
//
// If no such entry exists, Create returns ErrKeyExists.
func (s *Connection) Create(ctx context.Context, name string, value []byte) error {
	const (
		SecretType      = "opaque"
		ContentType     = "application/octet-stream"
		ContentEncoding = "base64"
		Algorithm       = "aes"
		BitLength       = 256
		Mode            = "cbc"
	)
	// Check if key already exists
	if err := s.verifyKeyDoesNotExist(ctx, name); err != xerror.ErrKeyExists {
		return err
	}

	// Create new key
	request, err := json.Marshal(SecretCreateRequest{
		SecretType:             SecretType,
		Name:                   name,
		Payload:                base64.StdEncoding.EncodeToString(value),
		PayloadContentType:     ContentType,
		PayloadContentEncoding: ContentEncoding,
		Algorithm:              Algorithm,
		BitLength:              BitLength,
		Mode:                   Mode,
	})
	if err != nil {
		return err
	}
	s.client.HttpPost(ctx, endpoint(s.config.Endpoint, "/v1/secrets"), request, nil)
	return nil
}

func (s *Connection) GetSecret(ctx context.Context, name string) (*BarbicanSecret, error) {
	url := endpoint(s.config.Endpoint, "/v1/secrets") + "?name=" + name
	resp, err := s.client.HttpGet(ctx, url, nil)
	if err != nil {
		return nil, err
	}

	var response BarbicanSecretsResponse

	if err := json.Unmarshal(resp, &response); err != nil {
		return nil, fmt.Errorf("barbican: failed to fetch '%s': failed to parse key metadata: %v", name, err)
	}

	if len(response.Secrets) == 0 {
		return nil, xerror.ErrKeyNotFound
	}

	return &response.Secrets[0], nil
}

func (s *Connection) GetSecretWithPayload(ctx context.Context, name string) (*BarbicanSecretWithPayload, error) {
	secret, err := s.GetSecret(ctx, name)
	if err != nil {
		return nil, err
	}

	// now we can get the secret payload
	url := endpoint(secret.SecretRef, "/payload")
	payload, err := s.client.HttpGet(ctx, url, nil)
	if err != nil {
		return nil, err
	}

	secPayload := &BarbicanSecretWithPayload{
		Secret:  *secret,
		Payload: payload,
	}

	return secPayload, nil
}

// Delete deletes the key associated with the given name
// from Barbican. It may not return an error if no
// entry for the given name exists.
func (s *Connection) DeleteSecret(ctx context.Context, name string) error {
	secret, err := s.GetSecret(ctx, name)
	if err != nil {
		return err
	}

	// Now, we can delete the key using its UUID.
	url := endpoint(secret.SecretRef)
	_, err = s.client.HttpDelete(ctx, url, nil, nil)
	return err
}

// List returns a new Iterator over the Barbican.
//
// The returned iterator may or may not reflect any
// concurrent changes to the Barbican - i.e.
// creates or deletes. Further, it does not provide any
// ordering guarantees.
func (s *Connection) ListSecrets(ctx context.Context) (*iterator, error) {
	var cancel context.CancelCauseFunc
	ctx, cancel = context.WithCancelCause(ctx)
	values := make(chan string, 10)

	go func() {
		defer close(values)

		var next string
		const limit = 200 // We limit a listing page to 200. This an arbitrary but reasonable value.
		for {
			reqURL := endpoint(s.config.Endpoint, "/v1/secrets") + "?sort=name:asc&limit=" + fmt.Sprint(limit)
			if next != "" {
				reqURL = next
			}

			resp, err := s.client.HttpGet(ctx, reqURL, nil)
			if err != nil {
				cancel(fmt.Errorf("barbican: failed to list keys: %v", err))
			}

			var keys BarbicanSecretsResponse
			if err := json.Unmarshal(resp, &keys); err != nil {
				cancel(fmt.Errorf("barbican: failed to list keys: failed to parse server response: %v", err))
				break
			}
			if len(keys.Secrets) == 0 {
				break
			}
			for _, k := range keys.Secrets {
				select {
				case values <- k.Name:
				case <-ctx.Done():
					return
				}
			}
			next = keys.Next
			if next == "" {
				break
			}
		}
	}()
	return &iterator{
		ch:     values,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

// Checks if a key already exists or not, if so returns ErrKeyExists
func (s *Connection) verifyKeyDoesNotExist(ctx context.Context, name string) error {
	_, err := s.GetSecret(ctx, name)
	if err != nil {
		return xerror.ErrKeyExists
	}
	return nil
}
