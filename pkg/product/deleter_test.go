package product

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

type removerStub struct {
	err error
}

func (r *removerStub) Remove(_ context.Context, _ uuid.UUID) error {
	return r.err
}

func TestDeleter(t *testing.T) {
	id := uuid.New()
	prd := &Product{
		ID:          id,
		Name:        "lorem ipsum",
		Description: "dolor sit amet",
		CreatedAt:   mustParseTime("2024-03-12 12:44:22"),
		UpdatedAt:   mustParseTime("2024-03-12 12:44:22"),
	}
	tests := map[string]struct {
		id        uuid.UUID
		findPrd   *Product
		findErr   error
		removeErr error
		prdWant   *Product
		errWant   error
	}{
		"it must err for a failing find": {
			id:        id,
			findPrd:   nil,
			findErr:   errFake,
			removeErr: nil,
			prdWant:   nil,
			errWant:   errFake,
		},
		"it must err for a missing product": {
			id:        id,
			findPrd:   nil,
			findErr:   nil,
			removeErr: nil,
			prdWant:   nil,
			errWant:   ErrNotFound,
		},
		"it must err for a failing remove": {
			id:        id,
			findPrd:   prd,
			findErr:   nil,
			removeErr: errFake,
			prdWant:   nil,
			errWant:   errFake,
		},
		"happy path": {
			id:        id,
			findPrd:   prd,
			findErr:   nil,
			removeErr: nil,
			prdWant:   prd,
			errWant:   nil,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			f := &finderStub{test.findPrd, test.findErr}
			r := &removerStub{test.removeErr}
			d := &Deleter{f, r}

			prd, err := d.Delete(context.Background(), test.id)

			if !reflect.DeepEqual(prd, test.prdWant) {
				t.Errorf("unexpected product; got %v, want %v", prd, test.prdWant)
			}

			if !errors.Is(err, test.errWant) {
				t.Errorf("unexpected error; got %v, want %v", err, test.errWant)
			}
		})
	}
}
