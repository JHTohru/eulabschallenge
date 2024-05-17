package product

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGetter_Get(t *testing.T) {
	id := uuid.New()
	fakeErr := fmt.Errorf("fake error")
	now := time.Now()
	prd := &Product{
		ID:          id,
		Name:        "lorem ipsum",
		Description: "dolor sit amet",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	tests := map[string]struct {
		id      uuid.UUID
		findPrd *Product
		findErr error
		prdWant *Product
		errWant error
	}{
		"it must err for a failing find": {
			id:      id,
			findPrd: nil,
			findErr: fakeErr,
			prdWant: nil,
			errWant: fakeErr,
		},
		"it must err for a missing product": {
			id:      id,
			findPrd: nil,
			findErr: nil,
			prdWant: nil,
			errWant: ErrNotFound,
		},
		"happy path": {
			id:      id,
			findPrd: prd,
			findErr: nil,
			prdWant: prd,
			errWant: nil,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			f := &finderStub{test.findPrd, test.findErr}
			g := &Getter{finder: f}

			prd, err := g.Get(context.Background(), test.id)

			if !reflect.DeepEqual(prd, test.prdWant) {
				t.Errorf("unexpected product; got %v, want %v", prd, test.prdWant)
			}
			if !errors.Is(err, test.errWant) {
				t.Errorf("unexpected product; got %v, want %v", err, test.errWant)
			}
		})
	}
}
