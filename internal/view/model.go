// view package provides simple way to track views on a given URL.
// It also implements simple querying API.
package view

import (
	"context"
	"time"
)

// View represents entity model of a single view.
type View struct {
	ID        uint
	URL       string
	CreatedAt time.Time
}

// ViewCollection represents a collection of View domain entities.
type ViewCollection []View

// Filter holds parameters available for filtering Views.
type Filter struct {
	URL    string
	After  time.Time
	Before time.Time
}

// Repository defines a storage API for View entity.
type Repository interface {
	Create(context.Context, View) (View, error)
	Filter(context.Context, Filter) (ViewCollection, error)
}
