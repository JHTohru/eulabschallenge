package postgres

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProductRemover(t *testing.T) {
	t.Parallel()
	db, dbName := newTmpDB(t)
	defer db.Close()
	defer dropDB(dbName)

	id := uuid.New()
	prd := randomProduct()

	insertProducts(t, db, prd)

	pr := NewProductRemover(db)
	ctx := context.Background()

	err := pr.Remove(ctx, id)

	assert.Nil(t, err)
	assert.False(t, productExistsByID(t, db, id))

	if err := db.Close(); err != nil {
		t.Fatal(err)
	}

	err = pr.Remove(ctx, id)

	assertErrIsDatabaseIsClosed(t, err)
}
