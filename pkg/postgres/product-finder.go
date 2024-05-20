package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"github.com/JHTohru/eulabschallenge/pkg/product"
)

type ProductFinder struct {
	db *sql.DB
}

func NewProductFinder(db *sql.DB) *ProductFinder {
	return &ProductFinder{
		db: db,
	}
}

func (pg *ProductFinder) Find(ctx context.Context, id uuid.UUID) (*product.Product, error) {
	query := "SELECT id, name, description, created_at, updated_at " +
		"FROM product WHERE id = $1"
	prd := &product.Product{}
	err := pg.db.QueryRowContext(ctx, query, id).Scan(&prd.ID, &prd.Name,
		&prd.Description, &prd.CreatedAt, &prd.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return prd, nil
}
