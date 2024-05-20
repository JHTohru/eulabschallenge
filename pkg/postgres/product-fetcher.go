package postgres

import (
	"context"
	"database/sql"

	"github.com/JHTohru/eulabschallenge/pkg/product"
)

type ProductFetcher struct {
	db *sql.DB
}

func NewProductFetcher(db *sql.DB) *ProductFetcher {
	return &ProductFetcher{
		db: db,
	}
}

func (pf *ProductFetcher) Fetch(ctx context.Context, limit, offset int) ([]*product.Product, error) {
	query := "SELECT id, name, description, created_at, updated_at " +
		"FROM product ORDER BY created_at ASC " +
		"LIMIT $1 OFFSET $2"
	rows, err := pf.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var prds []*product.Product
	for rows.Next() {
		prd := &product.Product{}
		err := rows.Scan(&prd.ID, &prd.Name, &prd.Description, &prd.CreatedAt, &prd.UpdatedAt)
		if err != nil {
			return nil, err
		}
		prds = append(prds, prd)
	}
	return prds, nil
}
