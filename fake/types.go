package fake

import (
	"sync"

	"github.com/artashesbalabekyan/barbican-sdk-go/client"
)

type Client struct {
	fakeData *FakeData
}

type FakeData struct {
	data map[string]client.BarbicanSecretWithPayload
	sync.RWMutex
}
