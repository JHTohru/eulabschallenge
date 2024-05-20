package postgres

import (
	"context"
	"database/sql"
)

type ProductCounter struct {
	db *sql.DB
}

func NewProductCounter(db *sql.DB) *ProductCounter {
	return &ProductCounter{
		db: db,
	}
}

func (pc *ProductCounter) Count(ctx context.Context) (int, error) {
	query := "SELECT COUNT(*) FROM product"
	var total int
	err := pc.db.QueryRowContext(ctx, query).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}
