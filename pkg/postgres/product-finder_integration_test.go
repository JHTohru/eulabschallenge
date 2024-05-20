package postgres

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProductFinder(t *testing.T) {
	t.Parallel()

	db, dbName := newTmpDB(t)
	defer dropDB(t, dbName)
	defer db.Close()

	prdWant := randomProduct()

	insertProducts(t, db, prdWant)

	pf := NewProductFinder(db)
	ctx := context.Background()

	prdGot, err := pf.Find(ctx, prdWant.ID)

	assert.Nil(t, err)
	assertProductsAreEqual(t, prdGot, prdWant)

	prdGot, err = pf.Find(ctx, uuid.New())

	assert.Nil(t, err)
	assert.Nil(t, prdGot)

	if err := db.Close(); err != nil {
		t.Fatal(err)
	}

	prdGot, err = pf.Find(ctx, prdWant.ID)

	assertErrIsDatabaseIsClosed(t, err)
	assert.Nil(t, prdGot)
}
