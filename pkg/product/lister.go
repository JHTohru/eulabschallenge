package product

import (
	"context"
	"math"
)

type Lister struct {
	counter Counter
	fetcher Fetcher
}

func NewLister(c Counter, f Fetcher) *Lister {
	return &Lister{
		counter: c,
		fetcher: f,
	}
}

// List lists the existing product records sorted by their creation times.
// The existing product records are divided in sets called pages. Each page has
// pageSize records. List returns the pageNumberth subset with at most pageSize
// records and the total of pages with such size. List returns an
// ErrInvalidPageSize error for a pageSize less than one. It also returns an
// ErrInvalidPageNumber error for a pageNumber less than one. List returns any
// database errors that happen.
func (l *Lister) List(ctx context.Context, pageSize, pageNumber int) ([]*Product, int, error) {
	if pageSize < 1 {
		return nil, 0, ErrInvalidPageSize
	}
	if pageNumber < 1 {
		return nil, 0, ErrInvalidPageNumber
	}

	total, err := l.counter.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	limit, offset := pageSize, (pageNumber-1)*pageSize
	prds, err := l.fetcher.Fetch(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return prds, totalPages, nil
}
