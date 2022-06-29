package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/miprokop/fication/internal/models"
	"github.com/uptrace/bun"
)

type StepRepo struct {
	DB  *bun.DB
	ctx context.Context
}

func (s *StepRepo) GetStepPrizes(ctx context.Context, id uuid.UUID) ([]*models.Prize, error) {
	prizes := new([]*models.Prize)
	err := s.DB.NewSelect().Model(prizes).Relation("Step").Where("step.id = ?", id).Scan(ctx)
	return *prizes, err
}

func (s *StepRepo) AssignStaff(ctx context.Context, staff models.StepStaff) error {
	_, err := s.DB.NewInsert().Model(staff).Exec(ctx)
	return err
}

func (s *StepRepo) PassStaff(ctx context.Context, staff models.StepStaff) error {
	_, err := s.DB.NewUpdate().OmitZero().Model(&staff).Exec(ctx)
	return err
}

func (s *StepRepo) CreateStep(ctx context.Context, step *models.Step) error {
	tx, err := s.DB.DB.Begin()
	if err != nil {
		return err
	}
	_, err = s.DB.NewInsert().Model(step).Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}
	if len(step.Prizes) != 0 {
		_, err = s.DB.NewInsert().Model(&step.Prizes).Exec(ctx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	if len(step.Images) != 0 {
		for i := range step.Images {
			step.Images[i].ID = uuid.New()
			step.Images[i].StepID = step.ID
		}
		_, err = s.DB.NewInsert().Model(&step.Images).Exec(ctx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (s *StepRepo) GetStep(ctx context.Context, id uuid.UUID) (*models.Step, error) {
	step := new(models.Step)
	err := s.DB.NewSelect().Model(step).Where("id = ?", id).Scan(ctx)
	return step, err
}

func (s *StepRepo) GetSteps(ctx context.Context, eventID uuid.UUID) ([]*models.Step, error) {
	steps := new([]*models.Step)
	err := s.DB.NewSelect().Model(steps).Where("event_id = ?", eventID).Scan(ctx)
	return *steps, err
}

func (s *StepRepo) DeleteStep(ctx context.Context, id uuid.UUID) error {
	_, err := s.DB.NewDelete().Where("id = ?", id).Exec(ctx)
	return err
}

func (s *StepRepo) UpdateStep(ctx context.Context, step *models.Step) error {
	tx, err := s.DB.DB.Begin()
	if err != nil {
		return err
	}
	_, err = s.DB.NewUpdate().OmitZero().Model(step).Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}
	if len(step.Prizes) != 0 {
		_, err = s.DB.NewUpdate().OmitZero().Model(&step.Prizes).Exec(ctx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func NewStepRepo(ctx context.Context, DB *bun.DB) *StepRepo {
	return &StepRepo{DB: DB, ctx: ctx}
}
