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
