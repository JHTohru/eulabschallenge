package product

import "context"

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

func (l *Lister) List(ctx context.Context, pageSize, pageNumber int) ([]*Product, int, error) {
	if pageSize <= 0 {
		return nil, 0, ErrInvalidPageSize
	}
	if pageNumber <= 0 {
		return nil, 0, ErrInvalidPageNumber
	}

	total, err := l.counter.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return nil, 0, nil
	}

	limit, offset := pageSize, pageSize*pageNumber
	prds, err := l.fetcher.Fetch(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return prds, total, nil
}
