package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type ProductRemover struct {
	db *sql.DB
}

func NewProductRemover(db *sql.DB) *ProductRemover {
	return &ProductRemover{
		db: db,
	}
}

func (pr *ProductRemover) Remove(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM product WHERE id = $1"
	_, err := pr.db.QueryContext(ctx, query, id)
	return err
}
