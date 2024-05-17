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

type counterStub struct {
	total int
	err   error
}

func (c *counterStub) Count(_ context.Context) (int, error) {
	return c.total, c.err
}

type fetcherStub struct {
	prds []*Product
	err  error
}

func (f *fetcherStub) Fetch(_ context.Context, _, _ int) ([]*Product, error) {
	return f.prds, f.err
}

func mustParseTime(value string) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05", value)
	if err != nil {
		panic(err)
	}
	return t
}

func TestLister_List(t *testing.T) {
	fakeErr := fmt.Errorf("fake error")
	prds := []*Product{
		{
			ID:          uuid.New(),
			Name:        "lorem ipsum",
			Description: "dolor sit amet",
			CreatedAt:   mustParseTime("2024-03-12 12:44:22"),
			UpdatedAt:   mustParseTime("2024-03-12 12:44:22"),
		},
		{
			ID:          uuid.New(),
			Name:        "consectetur",
			Description: "adipiscing elit",
			CreatedAt:   mustParseTime("2024-04-11 13:21:30"),
			UpdatedAt:   mustParseTime("2024-04-11 14:01:59"),
		},
		{
			ID:          uuid.New(),
			Name:        "donec non",
			Description: "convallis nulla",
			CreatedAt:   mustParseTime("2024-04-15 12:04:51"),
			UpdatedAt:   mustParseTime("2024-05-03 23:11:19"),
		},
	}
	tests := map[string]struct {
		pageSize   int
		pageNumber int
		countTotal int
		countErr   error
		fetchPrds  []*Product
		fetchErr   error
		prdsWant   []*Product
		totalWant  int
		errWant    error
	}{
		"it must err for an invalid page size": {
			pageSize:   0,
			pageNumber: 1,
			countTotal: 0,
			countErr:   nil,
			fetchPrds:  nil,
			fetchErr:   nil,
			prdsWant:   nil,
			totalWant:  0,
			errWant:    ErrInvalidPageSize,
		},
		"it must err for an invalid page number": {
			pageSize:   3,
			pageNumber: 0,
			countTotal: 0,
			countErr:   nil,
			fetchPrds:  nil,
			fetchErr:   nil,
			prdsWant:   nil,
			totalWant:  0,
			errWant:    ErrInvalidPageNumber,
		},
		"it must err for a failing count": {
			pageSize:   3,
			pageNumber: 1,
			countTotal: 0,
			countErr:   fakeErr,
			fetchPrds:  nil,
			fetchErr:   nil,
			prdsWant:   nil,
			totalWant:  0,
			errWant:    fakeErr,
		},
		"it must err for a failing fetch": {
			pageSize:   3,
			pageNumber: 1,
			countTotal: 100,
			countErr:   nil,
			fetchPrds:  nil,
			fetchErr:   fakeErr,
			prdsWant:   nil,
			totalWant:  0,
			errWant:    fakeErr,
		},
		"happy path": {
			pageSize:   3,
			pageNumber: 1,
			countTotal: 100,
			countErr:   nil,
			fetchPrds:  prds,
			fetchErr:   nil,
			prdsWant:   prds,
			totalWant:  100,
			errWant:    nil,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			c := &counterStub{test.countTotal, test.countErr}
			f := &fetcherStub{test.fetchPrds, test.fetchErr}
			l := &Lister{c, f}

			prds, total, err := l.List(context.Background(), test.pageSize, test.pageNumber)

			if !reflect.DeepEqual(prds, test.prdsWant) {
				t.Errorf("unexpected products; want %v, got %v", prds, test.prdsWant)
			}

			if total != test.totalWant {
				t.Errorf("unexpected total: want %v, got %v", total, test.totalWant)
			}

			if !errors.Is(err, test.errWant) {
				t.Errorf("unexpected error: want %v, got %v", err, test.errWant)
			}
		})
	}
}
