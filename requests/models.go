package requests

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// TODO: Move the panics

// NullableString is a string that can be null
type NullableString string

// Scan checks if the value is nil and sets the string to empty if it is
// Otherwise it just uses the value as a string
func (ns *NullableString) Scan(value any) error {
	if value == nil {
		*ns = ""
		return nil
	}
	v, ok := value.(string)
	if !ok {
		panic(fmt.Sprintf("unsupported type: %T", v))
	}
	*ns = NullableString(v)
	return nil
}

// APIData represents the data needed to create a new API
type APIData struct {
	Prefix string
	URL    string
}

// API represents an existing API that will be mapped from the prefix
type API struct {
	ID        uuid.UUID
	Timestamp time.Time
	Prefix    string
	URL       string
}

// RequestData represents the data needed to create a new request
type RequestData struct {
	Method   string
	Prefix   string
	Path     string
	Response ResponseData
}

// Request represents a made request
type Request struct {
	ID        uuid.UUID
	CreatedAt time.Time `db:"created_at"`
	Name      NullableString
	Method    string
	Path      string
	Response  Response
}

// ResponseData represents the data needed to create a new response
type ResponseData struct {
	Status string
	Body   []byte
}

// Response represents a response to a request
type Response struct {
	Status  string
	Content map[string]any
}

// Scan unmarshals the response body into a map
func (c *Response) Scan(value any) error {
	v, ok := value.([]byte)
	if !ok {
		panic(fmt.Sprintf("unsupported type: %T", v))
	}
	return json.Unmarshal(v, &c)
}
