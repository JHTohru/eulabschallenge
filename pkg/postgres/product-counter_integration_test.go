package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductCounter(t *testing.T) {
	t.Parallel()

	db, dbName := newTmpDB(t)
	defer db.Close()
	defer dropDB(dbName)

	prds := randomProducts(10)

	insertProducts(t, db, prds...)

	totalWant := len(prds)
	pc := NewProductCounter(db)
	ctx := context.Background()

	totalGot, err := pc.Count(ctx)

	assert.Nil(t, err)
	assert.Equal(t, totalGot, totalWant)

	if err := db.Close(); err != nil {
		t.Fatal(err)
	}

	totalGot, err = pc.Count(ctx)

	assertErrIsDatabaseIsClosed(t, err)
	assert.Zero(t, totalGot)
}
