package requests

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// TODO: Move the panics

type NullableString string

func (ns *NullableString) Scan(value any) error {
	fmt.Println("scanning", value)
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

type APIData struct {
	Prefix string
	URL    string
	Secret bool
}

type API struct {
	ID        uuid.UUID
	Timestamp time.Time
	Prefix    string
	URL       string
	Secret    bool
}

type RequestData struct {
	Method   string
	Prefix   string
	Path     string
	Response ResponseData
}

type Request struct {
	ID        uuid.UUID
	CreatedAt time.Time `db:"created_at"`
	Name      NullableString
	Method    string
	Path      string
	Response  Response
}

type ResponseData struct {
	Status string
	Body   []byte
}

type Response struct {
	Status  string
	Content map[string]any
}

func (c *Response) Scan(value any) error {
	v, ok := value.([]byte)
	if !ok {
		panic(fmt.Sprintf("unsupported type: %T", v))
	}
	return json.Unmarshal(v, &c)
}
