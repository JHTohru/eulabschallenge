package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProductSaver(t *testing.T) {
	t.Parallel()
	db, dbName := newTmpDB(t)
	defer db.Close()
	defer dropDB(dbName)

	prd := randomProduct()

	insertProducts(t, db, prd)

	prd.Name = randomString(20)
	prd.Description = randomString(100)
	prd.UpdatedAt = time.Now()
	ps := NewProductSaver(db)
	ctx := context.Background()

	err := ps.Update(ctx, prd)

	assert.Nil(t, err)
	assert.True(t, productExists(t, db, prd))

	if err := db.Close(); err != nil {
		t.Fatal(err)
	}

	err = ps.Update(ctx, prd)

	assertErrIsDatabaseIsClosed(t, err)
}
