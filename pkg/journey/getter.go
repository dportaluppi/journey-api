package journey

import (
	"context"
	"errors"
	"time"
)

type getter struct {
	repo Repository
}

func NewGetter(repo Repository) Getter {
	return &getter{repo: repo}
}

func (s *getter) GetJourneysByCriteria(ctx context.Context, storefront string, audiences, channels []string, date time.Time, sortBy string) ([]Journey, error) {
	filter := &Filter{
		Storefront: storefront,
		Audiences:  audiences,
		Channels:   channels,
		Date:       date,
	}

	journeys, err := s.repo.GetJourneys(ctx, filter, sortBy)
	if err != nil {
		return nil, err
	}

	if len(journeys) == 0 {
		return nil, errors.New("no journeys found")
	}

	return journeys, nil
}

func (s *getter) GetJourneyByID(ctx context.Context, id string) (*Journey, error) {
	return s.repo.GetByID(ctx, id)
}
