package journey

import "context"

type deleter struct {
	repo Repository
}

func NewDeleter(repo Repository) Deleter {
	return &deleter{repo: repo}
}

func (s *deleter) DeleteJourney(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
