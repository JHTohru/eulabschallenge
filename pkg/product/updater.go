package product

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Updater struct {
	finder      Finder
	saver       Saver
	currentTime func() time.Time
}

func NewUpdater(f Finder, s Saver) *Updater {
	return &Updater{
		finder:      f,
		saver:       s,
		currentTime: func() time.Time { return time.Now() },
	}
}

func (u *Updater) Update(ctx context.Context, id uuid.UUID, in *Input) (*Product, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	prd, err := u.finder.Find(ctx, id)
	if err != nil {
		return nil, err
	}
	if prd == nil {
		return nil, ErrNotFound
	}

	prd.Name = in.Name
	prd.Description = in.Description
	prd.UpdatedAt = u.currentTime()
	if err := u.saver.Save(ctx, prd); err != nil {
		return nil, err
	}

	return prd, nil
}
