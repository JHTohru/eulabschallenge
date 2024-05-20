package postgres

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/JHTohru/eulabschallenge/pkg/product"
)

type ProductInserter struct {
	db *sql.DB
}

func NewProductInserter(db *sql.DB) *ProductInserter {
	return &ProductInserter{
		db: db,
	}
}

func (pi *ProductInserter) Insert(ctx context.Context, prd *product.Product) error {
	query := "INSERT INTO product (id, name, description, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)"
	_, err := pi.db.ExecContext(ctx, query, prd.ID, prd.Name, prd.Description,
		prd.CreatedAt.UTC(), prd.UpdatedAt.UTC())
	return err
}
