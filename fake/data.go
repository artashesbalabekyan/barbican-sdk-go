package fake

import (
	"time"

	"github.com/artashesbalabekyan/barbican-sdk-go/client"
)

const (
	timeFormat = "2006-01-02T15:04:05"
)

func newCreationTime() string {
	return time.Now().Format(timeFormat)
}

func NewFakeData(d map[string][]byte) *FakeData {
	f := &FakeData{
		data: make(map[string]client.BarbicanSecretWithPayload),
	}
	for name, payload := range d {
		f.Set(name, NewSecret(name, payload))
	}
	return f
}

func NewSecret(name string, payload []byte) client.BarbicanSecretWithPayload {
	return client.BarbicanSecretWithPayload{
		Secret: client.BarbicanSecret{
			Algorithm:    "aes",
			BitLength:    256,
			ContentTypes: map[string]string{"default": "application/octet-stream"},
			Created:      newCreationTime(),
			Mode:         "cbc",
			Name:         name,
			SecretType:   "opaque",
			Status:       "ACTIVE",
			Updated:      newCreationTime(),
		},
		Payload: payload,
	}
}

func (f *FakeData) Set(name string, secret client.BarbicanSecretWithPayload) {
	f.Lock()
	f.data[name] = secret
	f.Unlock()
}

func (f *FakeData) Get(name string) (secret client.BarbicanSecretWithPayload, exist bool) {
	f.RLock()
	secret, exist = f.data[name]
	f.RUnlock()
	return
}

func (f *FakeData) Delete(name string) {
	f.Lock()
	delete(f.data, name)
	f.Unlock()
}

func (f *FakeData) List() []client.BarbicanSecret {
	f.RLock()
	list := []client.BarbicanSecret{}
	for _, secret := range f.data {
		list = append(list, secret.Secret)
	}
	f.RUnlock()
	return list
}
