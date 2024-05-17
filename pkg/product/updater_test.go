package product

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestUpdater_Update(t *testing.T) {
	now := time.Now()
	id := uuid.New()
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
			now:     now,
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
			now:     now,
			findPrd: nil,
			findErr: errFake,
			saveErr: nil,
			prdWant: nil,
			errWant: errFake,
		},
		"it must err for a missing product": {
			id: id,
			in: &Input{
				Name:        "consectetur",
				Description: "adipiscing elit",
			},
			now:     now,
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
			now: now,
			findPrd: &Product{
				ID:          id,
				Name:        "lorem ipsum",
				Description: "dolor sit amet",
				CreatedAt:   mustParseTime("2024-03-12 12:44:22"),
				UpdatedAt:   mustParseTime("2024-03-12 12:44:22"),
			},
			findErr: nil,
			saveErr: errFake,
			prdWant: nil,
			errWant: errFake,
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
				CreatedAt:   mustParseTime("2024-03-12 12:44:22"),
				UpdatedAt:   mustParseTime("2024-03-12 12:44:22"),
			},
			findErr: nil,
			saveErr: nil,
			prdWant: &Product{
				ID:          id,
				Name:        "consectetur",
				Description: "adipiscing elit",
				CreatedAt:   mustParseTime("2024-03-12 12:44:22"),
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
