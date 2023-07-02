package gobarbicansdk

import (
	"context"

	"github.com/artashesbalabekyan/barbican-sdk-go/client"
	"github.com/artashesbalabekyan/barbican-sdk-go/fake"
	"github.com/artashesbalabekyan/barbican-sdk-go/xhttp"
)

func NewConnection(ctx context.Context, config *xhttp.Config) (client.Conn, error) {
	return client.New(ctx, config)
}

func NewFakeConnection(ctx context.Context, fakeData map[string][]byte) (client.Conn, error) {
	return fake.New(ctx, fakeData)
}
