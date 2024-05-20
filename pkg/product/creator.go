package product

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Creator struct {
	inserter    Inserter
	currentTime func() time.Time
	newID       func() uuid.UUID
}

func NewCreator(i Inserter) *Creator {
	return &Creator{
		inserter:    i,
		currentTime: time.Now,
		newID:       uuid.New,
	}
}

// Create creates a new product and saves it into the database. Since input is
// validated, Create may return the same errors as *Input.Validate returns.
// Create returns any database failure errors that happen.
func (c *Creator) Create(ctx context.Context, in *Input) (*Product, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	now := c.currentTime()
	prd := &Product{
		ID:          c.newID(),
		Name:        in.Name,
		Description: in.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := c.inserter.Insert(ctx, prd); err != nil {
		return nil, err
	}

	return prd, nil
}
