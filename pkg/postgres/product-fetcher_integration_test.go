package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductFetcher(t *testing.T) {
	t.Parallel()

	db, dbName := newTmpDB(t)
	defer db.Close()
	defer dropDB(dbName)

	prds := randomProducts(10)

	insertProducts(t, db, prds...)

	prdsWant := prds[2:5]
	pf := NewProductFetcher(db)
	ctx := context.Background()

	prdsGot, err := pf.Fetch(ctx, 3, 2)

	assert.Nil(t, err)
	assertProductListsAreEqual(t, prdsGot, prdsWant)

	if err := db.Close(); err != nil {
		t.Fatal(err)
	}

	prdsGot, err = pf.Fetch(ctx, 3, 2)

	assertErrIsDatabaseIsClosed(t, err)
	assert.Nil(t, prdsGot)
}
