package product

import (
	"context"
	"fmt"
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

func mustParseTime(value string) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05", value)
	if err != nil {
		panic(err)
	}
	return t
}

var errFake = fmt.Errorf("fake error")
