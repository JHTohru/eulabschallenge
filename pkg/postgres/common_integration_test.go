package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"

	"github.com/JHTohru/eulabschallenge/pkg/product"
)

var (
	mainDB *sql.DB
	once   sync.Once
)

func connectToDB(dbName string) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=localhost user=postgres password=password "+
		"dbname=%s port=5432 sslmode=disable", dbName)
	return sql.Open("postgres", dsn)
}

func connectToMainDB() {
	db, err := connectToDB("postgres")
	if err != nil {
		panic(err)
	}
	mainDB = db
}

func newTmpDB(t *testing.T) (*sql.DB, string) {
	t.Helper()
	once.Do(connectToMainDB)
	// creates a new database
	dbName := fmt.Sprintf("tmpdb_%s", uuid.New())
	dbName = strings.ReplaceAll(dbName, "-", "")
	query := fmt.Sprintf("CREATE DATABASE %s", dbName)
	if _, err := mainDB.Exec(query); err != nil {
		t.Fatal(err)
	}
	db, err := connectToDB(dbName)
	if err != nil {
		t.Fatal(err)
	}
	// run migrations files on the new database
	ctx := context.Background()
	migrationsDir := "./infra/migrations"
	if err := goose.UpContext(ctx, db, migrationsDir); err != nil {
		t.Fatal(err)
	}
	return db, dbName
}

func dropDB(t *testing.T, dbName string) {
	t.Helper()

	query := fmt.Sprintf("DROP DATABASE %s WITH (FORCE)", dbName)
	_, err := mainDB.Exec(query)
	if err != nil {
		t.Fatal(err)
	}
}

// randomLetters builds a string of a len size containing random lowercase
// ascii letters.
func randomLetters(len int) string {
	b := make([]byte, len)
	for i := range len {
		b[i] = byte('a' + rand.Intn(25))
	}
	return string(b)
}

// randomTime generates a time not before than from, and at maximum max
// duration after it.
func randomTime(from time.Time, max time.Duration) time.Time {
	n := rand.Int63n(int64(max))
	return from.Add(time.Duration(n))
}

func randomProduct() *product.Product {
	createdAt := randomTime(time.Now(), 24*time.Hour)
	updatedAt := randomTime(createdAt, 24*time.Hour)
	return &product.Product{
		ID:          uuid.New(),
		Name:        randomLetters(20),
		Description: randomLetters(100),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

// randomProducts generate a list of random products of len size sorted by
// their CreateAt field in an ascendent fashion.
func randomProducts(len int) []*product.Product {
	prds := make([]*product.Product, len)
	for i := range len {
		prds[i] = randomProduct()
	}
	sort.Slice(prds, func(i, j int) bool {
		return prds[i].CreatedAt.Before(prds[j].CreatedAt)
	})
	return prds
}

func productExists(t *testing.T, db *sql.DB, prd *product.Product) bool {
	t.Helper()

	query := "SELECT EXISTS(SELECT 1 FROM product WHERE " +
		"id = $1 AND name = $2 AND description = $3 AND " +
		"created_at = $4 AND updated_at = $5)"
	var exists bool
	err := db.QueryRow(query, prd.ID, prd.Name, prd.Description,
		prd.CreatedAt.UTC(),
		prd.UpdatedAt.UTC()).Scan(&exists)
	if err != nil {
		t.Fatal(err)
	}
	return exists
}

func productExistsByID(t *testing.T, db *sql.DB, id uuid.UUID) bool {
	t.Helper()

	query := "SELECT EXISTS(SELECT 1 FROM product WHERE id = $1)"
	var exists bool
	err := db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		t.Fatal(err)
	}
	return exists
}

func insertProducts(t *testing.T, db *sql.DB, prds ...*product.Product) {
	t.Helper()

	for _, prd := range prds {
		query := "INSERT INTO product (id, name, description, created_at, " +
			"updated_at) VALUES ($1, $2, $3, $4, $5)"
		_, err := db.Query(query, prd.ID, prd.Name, prd.Description,
			prd.CreatedAt.UTC(), prd.UpdatedAt.UTC())
		if err != nil {
			t.Fatal(err)
		}
	}
}

func assertErrIsDatabaseIsClosed(t *testing.T, err error) bool {
	t.Helper()

	return assert.ErrorContains(t, err, "sql: database is closed")
}

func assertProductsAreEqual(t *testing.T, got, want *product.Product) bool {
	t.Helper()

	return assert.Conditionf(t, func() bool {
		ok1 := got.ID == want.ID
		ok2 := got.Name == want.Name
		ok3 := got.Description == want.Description

		gotCreatedAt := got.CreatedAt.Round(time.Microsecond)
		wantCreatedAt := want.CreatedAt.Round(time.Microsecond)
		ok4 := gotCreatedAt.Equal(wantCreatedAt)

		gotUpdatedAt := got.UpdatedAt.Round(time.Microsecond)
		wantUpdatedAt := want.UpdatedAt.Round(time.Microsecond)
		ok5 := gotUpdatedAt.Equal(wantUpdatedAt)

		return ok1 && ok2 && ok3 && ok4 && ok5
	}, "Products not equal:\nexpected: %v\nactual: %v", want, got)
}

func assertProductListsAreEqual(t *testing.T, got, want []*product.Product) bool {
	t.Helper()

	var ok bool
	if assert.Equal(t, len(got), len(want)) {
		ok = true
		for i := range got {
			if !assertProductsAreEqual(t, got[i], want[i]) && ok {
				ok = false
			}
		}
	}
	return ok
}
