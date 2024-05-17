package product

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type inserterStub struct {
	err error
}

func (i *inserterStub) Insert(_ context.Context, _ *Product) error {
	return i.err
}

type finderStub struct {
	prd *Product
	err error
}

func (f *finderStub) Find(_ context.Context, _ uuid.UUID) (*Product, error) {
	return f.prd, f.err
}

type fetcherStub struct {
	prds []*Product
	err  error
}

func (f *fetcherStub) Fetch(_ context.Context, _, _ int) ([]*Product, error) {
	return f.prds, f.err
}

type counterStub struct {
	total int
	err   error
}

func (c *counterStub) Count(_ context.Context) (int, error) {
	return c.total, c.err
}

type saverStub struct {
	err error
}

func (s *saverStub) Save(_ context.Context, _ *Product) error {
	return s.err
}

type removerStub struct {
	err error
}

func (r *removerStub) Remove(_ context.Context, _ uuid.UUID) error {
	return r.err
}

func mustParseTime(value string) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05", value)
	if err != nil {
		panic(err)
	}
	return t
}

var errFake = fmt.Errorf("fake error")
