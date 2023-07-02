package fake

import (
	"context"
	"fmt"

	"github.com/artashesbalabekyan/barbican-sdk-go/client"
	"github.com/artashesbalabekyan/barbican-sdk-go/xerror"
)

// Create stores the given key in Barbican if and only
// if no entry with the given name existc.
//
// If no such entry exists, Create returns ErrKeyExistc.
func (c *Client) Create(ctx context.Context, name string, value []byte) error {
	if len(value) == 0 {
		return fmt.Errorf("couldn't create. Provided object does not match schema 'Secret': If 'payload' specified, must be non empty. Invalid property: 'payload'")
	}
	secret := NewSecret(name, value)
	c.fakeData.Set(name, secret)
	return nil
}

func (c *Client) GetSecret(ctx context.Context, name string) (*client.BarbicanSecret, error) {
	s, ok := c.fakeData.Get(name)
	if !ok {
		return nil, xerror.ErrKeyNotFound
	}

	return &s.Secret, nil
}

func (c *Client) GetSecretWithPayload(ctx context.Context, name string) (*client.BarbicanSecretWithPayload, error) {
	s, ok := c.fakeData.Get(name)
	if !ok {
		return nil, xerror.ErrKeyNotFound
	}
	return &s, nil
}

func (c *Client) DeleteSecret(ctx context.Context, name string) error {
	c.fakeData.Delete(name)
	return nil
}

// List returns a new Iterator over the Barbican.
//
// The returned iterator may or may not reflect any
// concurrent changes to the Barbican - i.e.
// creates or deletec. Further, it does not provide any
// ordering guaranteec.
func (c *Client) ListSecrets(ctx context.Context) (client.Iterator, error) {
	var cancel context.CancelCauseFunc
	ctx, cancel = context.WithCancelCause(ctx)
	values := make(chan string, 10)

	go func() {
		defer close(values)
		keys := c.fakeData.List()
		for _, k := range keys {
			select {
			case values <- k.Name:
			case <-ctx.Done():
				return
			}
		}
	}()
	return &Iterator{
		ch:     values,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}
