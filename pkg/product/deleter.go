package product

import (
	"context"

	"github.com/google/uuid"
)

type Deleter struct {
	finder  Finder
	remover Remover
}

func NewDeleter(f Finder, r Remover) *Deleter {
	return &Deleter{
		finder:  f,
		remover: r,
	}
}

// Delete deletes a product record with the given id from the database and
// returns it. If such record doesn't exist, Delete returns an ErrNotFound
// error. Delete returns any database errors that happen.
func (d *Deleter) Delete(ctx context.Context, id uuid.UUID) (*Product, error) {
	prd, err := d.finder.Find(ctx, id)
	if err != nil {
		return nil, err
	}
	if prd == nil {
		return nil, ErrNotFound
	}

	if err := d.remover.Remove(ctx, id); err != nil {
		return nil, err
	}

	return prd, nil
}
