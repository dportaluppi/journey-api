package journey

import "context"

type creator struct {
	repo Repository
}

func NewCreator(repo Repository) Creator {
	return &creator{repo: repo}
}

func (s *creator) CreateJourney(ctx context.Context, j *Journey) (*Journey, error) {
	/* TODO: Fixme
	if err := j.Validate(); err != nil {
		return nil, err
	}
	*/

	return s.repo.Create(ctx, j)
}
