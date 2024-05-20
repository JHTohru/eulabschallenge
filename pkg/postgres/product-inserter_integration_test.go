package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductInserter(t *testing.T) {
	t.Parallel()

	db, dbName := newTmpDB(t)
	defer dropDB(t, dbName)
	defer db.Close()

	pi := NewProductInserter(db)
	ctx := context.Background()
	prd := randomProduct()

	err := pi.Insert(ctx, prd)

	assert.Nil(t, err)
	assert.True(t, productExists(t, db, prd))

	if err := db.Close(); err != nil {
		t.Fatal(err)
	}

	err = pi.Insert(ctx, prd)

	assertErrIsDatabaseIsClosed(t, err)
}
