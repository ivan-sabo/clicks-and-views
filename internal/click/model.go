// click package provides simple way to track clicks on a given URL.
// It also implements simple querying API.
package click

import (
	"context"
	"time"
)

// Click represents entity model of a single click.
type Click struct {
	ID        uint
	URL       string
	CreatedAt time.Time
}

// ClickCollection represents a collection of Click domain entities.
type ClickCollection []Click

// Filter holds parameters available for filtering Clicks.
type Filter struct {
	URL    string
	After  time.Time
	Before time.Time
}

// Repository defines a storage API for Click entity.
type Repository interface {
	Create(context.Context, Click) (Click, error)
	Filter(context.Context, Filter) (ClickCollection, error)
}
