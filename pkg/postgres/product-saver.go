package postgres

import (
	"context"
	"database/sql"

	"github.com/JHTohru/eulabschallenge/pkg/product"
)

type ProductSaver struct {
	db *sql.DB
}

func NewProductSaver(db *sql.DB) *ProductSaver {
	return &ProductSaver{
		db: db,
	}
}

func (ps *ProductSaver) Update(ctx context.Context, prd *product.Product) error {
	query := "UPDATE product SET name = $1, description = $2, created_at = $3, " +
		"updated_at = $4 WHERE id = $5"
	_, err := ps.db.ExecContext(ctx, query, prd.Name, prd.Description,
		prd.CreatedAt.UTC(), prd.UpdatedAt.UTC(), prd.ID)
	return err
}
