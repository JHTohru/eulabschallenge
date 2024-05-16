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

type finderStub struct {
	prd *Product
	err error
}

func (f *finderStub) Find(_ context.Context, _ uuid.UUID) (*Product, error) {
	return f.prd, f.err
}

type saverStub struct {
	err error
}

func (s *saverStub) Save(_ context.Context, _ *Product) error {
	return s.err
}

func TestUpdater_Update(t *testing.T) {
	now := time.Now()
	id := uuid.New()
	fakeErr := fmt.Errorf("fake error")
	tests := map[string]struct {
		id      uuid.UUID
		in      *Input
		now     time.Time
		findPrd *Product
		findErr error
		saveErr error
		prdWant *Product
		errWant error
	}{
		"it must err for an invalid input": {
			id:      id,
			in:      nil,
			now:     time.Time{},
			findPrd: nil,
			findErr: nil,
			saveErr: nil,
			prdWant: nil,
			errWant: ErrNilInput,
		},
		"it must err for a failing find": {
			id: id,
			in: &Input{
				Name:        "consectetur",
				Description: "adipiscing elit",
			},
			now:     time.Time{},
			findPrd: nil,
			findErr: fakeErr,
			saveErr: nil,
			prdWant: nil,
			errWant: fakeErr,
		},
		"it must err for a missing product": {
			id: id,
			in: &Input{
				Name:        "consectetur",
				Description: "adipiscing elit",
			},
			now:     time.Time{},
			findPrd: nil,
			findErr: nil,
			saveErr: nil,
			prdWant: nil,
			errWant: ErrNotFound,
		},
		"it must err for a failing save": {
			id: id,
			in: &Input{
				Name:        "consectetur",
				Description: "adipiscing elit",
			},
			now: time.Time{},
			findPrd: &Product{
				ID:          id,
				Name:        "lorem ipsum",
				Description: "dolor sit amet",
				CreatedAt:   time.Time{},
				UpdatedAt:   time.Time{},
			},
			findErr: nil,
			saveErr: fakeErr,
			prdWant: nil,
			errWant: fakeErr,
		},
		"happy path": {
			id: id,
			in: &Input{
				Name:        "consectetur",
				Description: "adipiscing elit",
			},
			now: now,
			findPrd: &Product{
				ID:          id,
				Name:        "lorem ipsum",
				Description: "dolor sit amet",
				CreatedAt:   time.Time{},
				UpdatedAt:   time.Time{},
			},
			findErr: nil,
			saveErr: nil,
			prdWant: &Product{
				ID:          id,
				Name:        "consectetur",
				Description: "adipiscing elit",
				CreatedAt:   time.Time{},
				UpdatedAt:   now,
			},
			errWant: nil,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			f := &finderStub{test.findPrd, test.findErr}
			s := &saverStub{test.saveErr}
			u := &Updater{
				finder:      f,
				saver:       s,
				currentTime: func() time.Time { return test.now },
			}

			prd, err := u.Update(context.Background(), test.id, test.in)

			if !reflect.DeepEqual(prd, test.prdWant) {
				t.Errorf("unexpected product; got %v, want %v", prd, test.prdWant)
			}
			if !errors.Is(err, test.errWant) {
				t.Errorf("unexpected erro; got %v, want %v", err, test.errWant)
			}
		})
	}
}
