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
	defer dropDB(t, dbName)
	defer db.Close()

	prd := randomProduct()

	insertProducts(t, db, prd)

	prd.Name = randomLetters(20)
	prd.Description = randomLetters(100)
	prd.UpdatedAt = time.Now()
	ps := NewProductSaver(db)
	ctx := context.Background()

	err := ps.Save(ctx, prd)

	assert.Nil(t, err)
	assert.True(t, productExists(t, db, prd))

	if err := db.Close(); err != nil {
		t.Fatal(err)
	}

	err = ps.Save(ctx, prd)

	assertErrIsDatabaseIsClosed(t, err)
}
