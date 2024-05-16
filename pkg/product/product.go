package product

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Inserter interface {
	Insert(ctx context.Context, prd *Product) error
}

type Finder interface {
	Find(ctx context.Context, id uuid.UUID) (*Product, error)
}

type Saver interface {
	Save(ctx context.Context, prd *Product) error
}
