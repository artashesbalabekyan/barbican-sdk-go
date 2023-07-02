package fake

import (
	"context"

	"github.com/artashesbalabekyan/barbican-sdk-go/client"
)

func New(ctx context.Context, fakeData map[string][]byte) (client.Conn, error) {
	return newConnection(ctx, fakeData)
}

func newConnection(ctx context.Context, fakeData map[string][]byte) (client.Conn, error) {
	return &Client{
		fakeData: NewFakeData(fakeData),
	}, nil
}
