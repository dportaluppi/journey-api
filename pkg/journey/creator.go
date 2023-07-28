package journey

import "context"

type creator struct {
	repo Repository
}

func NewCreator(repo Repository) Creator {
	return &creator{repo: repo}
}

func (s *creator) CreateJourney(ctx context.Context, j *Journey) (*Journey, error) {
	/* TODO: fixme
	if err := j.Validate(); err != nil {
		return nil, err
	}
	*/

	err := s.repo.Create(ctx, j)
	if err != nil {
		return nil, err
	}

	return j, nil
}
