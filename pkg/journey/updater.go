package journey

import "context"

type updater struct {
	repo Repository
}

func NewUpdater(repo Repository) Updater {
	return &updater{repo: repo}
}

func (s *updater) UpdateJourney(ctx context.Context, id string, j *Journey) error {
	/* TODO: Fixme
	if err := j.Validate(); err != nil {
		return err
	}
	*/

	return s.repo.Update(ctx, id, j)
}
