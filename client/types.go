package client

import "github.com/artashesbalabekyan/barbican-sdk-go/xhttp"

type Connection struct {
	config xhttp.Config
	client *xhttp.Client
}

type SecretCreateRequest struct {
	Name                   string `json:"name,omitempty"`                     // (optional) The name of the secret set by the user.
	Expiration             string `json:"expiration,omitempty"`               // (optional) This is a UTC timestamp in ISO 8601 format YYYY-MM-DDTHH:MM:SSZ. If set, the secret will not be available after this time.
	Algorithm              string `json:"algorithm,omitempty"`                // (optional) Metadata provided by a user or system for informational purposes.
	BitLength              int    `json:"bit_length,omitempty"`               // (optional) Metadata provided by a user or system for informational purposes. Value must be greater than zero.
	Mode                   string `json:"mode,omitempty"`                     // (optional) Metadata provided by a user or system for informational purposes.
	Payload                string `json:"payload"`                            // (optional) The secret’s data to be stored. payload_content_type must also be supplied if payload is included.
	PayloadContentType     string `json:"payload_content_type,omitempty"`     // (optional) (required if payload is included) The media type for the content of the payload. For more information see Secret Types
	PayloadContentEncoding string `json:"payload_content_encoding,omitempty"` // (optional) (required if payload is encoded) The encoding used for the payload to be able to include it in the JSON request. Currently only base64 is supported.
	SecretType             string `json:"secret_type,omitempty"`              // (optional) Used to indicate the type of secret being stored. For more information see Secret Types (default: opaque)
}

type BarbicanSecret struct {
	Algorithm    interface{}       `json:"algorithm"`
	BitLength    interface{}       `json:"bit_length"`
	ContentTypes map[string]string `json:"content_types"`
	Created      string            `json:"created"`
	CreatorID    string            `json:"creator_id"`
	Expiration   interface{}       `json:"expiration"`
	Mode         interface{}       `json:"mode"`
	Name         string            `json:"name"`
	SecretRef    string            `json:"secret_ref"`
	SecretType   string            `json:"secret_type"`
	Status       string            `json:"status"`
	Updated      string            `json:"updated"`
}

type BarbicanSecretWithPayload struct {
	Payload []byte         `json:"payload"`
	Secret  BarbicanSecret `json:"secret"`
}

type BarbicanSecretsResponse struct {
	Next     string           `json:"next"`
	Previous string           `json:"previous"`
	Secrets  []BarbicanSecret `json:"secrets"`
	Total    int              `json:"total"`
}
