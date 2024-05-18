package product

import (
	"context"

	"github.com/google/uuid"
)

type Inserter interface {
	Insert(ctx context.Context, prd *Product) error
}

type Finder interface {
	Find(ctx context.Context, id uuid.UUID) (*Product, error)
}

type Fetcher interface {
	Fetch(ctx context.Context, limit, offset int) ([]*Product, error)
}

type Counter interface {
	Count(ctx context.Context) (int, error)
}

type Saver interface {
	Save(ctx context.Context, prd *Product) error
}

type Remover interface {
	Remove(ctx context.Context, id uuid.UUID) error
}
