package gobarbicansdk

import (
	"context"

	"github.com/artashesbalabekyan/barbican-sdk-go/client"
	"github.com/artashesbalabekyan/barbican-sdk-go/xhttp"
)

func NewConnection(ctx context.Context, config *xhttp.Config) (*client.Connection, error) {
	return client.New(ctx, config)
}
