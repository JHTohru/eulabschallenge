package product

import (
	"context"

	"github.com/google/uuid"
)

type Getter struct {
	finder Finder
}

func NewGetter(f Finder) *Getter {
	return &Getter{
		finder: f,
	}
}

// Get gets a product record with the given id from the database. If such record
// doesn't exist, Get returns an ErrNotFound error. Get returns any database
// errors that happen.
func (g *Getter) Get(ctx context.Context, id uuid.UUID) (*Product, error) {
	prd, err := g.finder.Find(ctx, id)
	if err != nil {
		return nil, err
	}
	if prd == nil {
		return nil, ErrNotFound
	}
	return prd, err
}
