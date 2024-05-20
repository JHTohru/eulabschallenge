package product

import (
	"context"

	"github.com/google/uuid"
)

type Inserter interface {
	// Insert inserts a product record into a database.
	Insert(ctx context.Context, prd *Product) error
}

type Finder interface {
	// Find finds a single product record from a database. If there is no
	// produc record with the given id, Find returns a nil result and nil error.
	Find(ctx context.Context, id uuid.UUID) (*Product, error)
}

type Fetcher interface {
	// Fetch fetches from the database a slice of products sorted by their
	// creation time. The returned slice has at most the size of limit.
	// The first offset product records in the database are ignored.
	Fetch(ctx context.Context, limit, offset int) ([]*Product, error)
}

type Counter interface {
	// Count counts the total number of product records in the database.
	Count(ctx context.Context) (int, error)
}

type Saver interface {
	// Save saves a product record in the database.
	Save(ctx context.Context, prd *Product) error
}

type Remover interface {
	// Remove removes a product record from the database.
	Remove(ctx context.Context, id uuid.UUID) error
}
