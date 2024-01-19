package tui

import (
	"fmt"
	"time"

	"atheris/requests"

	"github.com/charmbracelet/bubbles/list"
)

// Item represents a request on the list
type item struct {
	id        string
	createdAt time.Time
	name      string
	method    string
	path      string
}

func (i item) Title() string {
	if i.name == "" {
		return i.id
	}
	return i.name
}

func (i item) Description() string {
	return fmt.Sprintf("%s %s at %s", i.method, i.path, i.createdAt)
}

func (i item) FilterValue() string {
	if i.name == "" {
		return i.id
	}
	return i.name
}

func ItemsFromRequests(rqs []requests.Request) []list.Item {
	items := make([]list.Item, len(rqs))
	for i, rq := range rqs {
		items[i] = item{
			id:        rq.ID.String(),
			createdAt: rq.CreatedAt,
			name:      string(rq.Name),
			method:    rq.Method,
			path:      rq.Path,
		}
	}
	return items
}
