package product

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCreator_Create(t *testing.T) {
	now := time.Now()
	id := uuid.New()
	tests := map[string]struct {
		in          *Input
		now         time.Time
		id          uuid.UUID
		inserterErr error
		prdWant     *Product
		errWant     error
	}{
		"it must err for an invalid input": {
			in:          nil,
			now:         time.Time{},
			id:          uuid.UUID{},
			inserterErr: nil,
			prdWant:     nil,
			errWant:     ErrNilInput,
		},
		"it must err for a failing insert": {
			in: &Input{
				Name:        "lorem ipsum",
				Description: "dolor sit amet",
			},
			now:         time.Time{},
			id:          uuid.UUID{},
			inserterErr: errFake,
			prdWant:     nil,
			errWant:     errFake,
		},
		"happy path": {
			in: &Input{
				Name:        "lorem ipsum",
				Description: "dolor sit amet",
			},
			now:         now,
			id:          id,
			inserterErr: nil,
			prdWant: &Product{
				ID:          id,
				Name:        "lorem ipsum",
				Description: "dolor sit amet",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			errWant: nil,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			i := &inserterStub{test.inserterErr}
			c := &Creator{
				inserter:    i,
				currentTime: func() time.Time { return test.now },
				newID:       func() uuid.UUID { return test.id },
			}

			prd, err := c.Create(context.Background(), test.in)

			if !reflect.DeepEqual(prd, test.prdWant) {
				t.Errorf("unexpected product; got %v, want %v", prd, test.prdWant)
			}
			if !errors.Is(err, test.errWant) {
				t.Errorf("unexpected error; got %v, want %v", err, test.errWant)
			}
		})
	}
}
