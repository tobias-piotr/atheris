package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

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

// TODO: Maybe there would be a better way than 'prefix'
type RequestData struct {
	Prefix   string
	Path     string
	Response Response
}

type Request struct {
	ID        uuid.UUID
	Timestamp time.Time
	Path      string
	Response  Response
}

type Response struct {
	Status  string
	Content map[string]any
}

func main() {
	fmt.Println("Hello, world")
}
