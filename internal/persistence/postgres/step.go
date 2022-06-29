package postgres

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/miprokop/fication/internal/models"
	"github.com/uptrace/bun"
	"time"
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
	_, err := s.DB.NewInsert().Model(&staff).Exec(ctx)
	return err
}

func (s *StepRepo) PassStaff(ctx context.Context, staff models.StepStaff) error {
	_, err := s.DB.NewUpdate().OmitZero().
		Model(&staff).
		Where("staff_id = ?", staff.StaffID).
		Where("step_id = ?", staff.StepID).
		Exec(ctx)
	return err
}

func (s *StepRepo) CreateStep(ctx context.Context, step *models.Step) error {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	_, err = tx.NewInsert().
		Model(step).
		ExcludeColumn("end_date").
		ExcludeColumn("creation_date").
		Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}
	creationTime, _ := time.Parse(time.RFC3339, step.CreationDate)
	_, err = tx.NewUpdate().Model(step).OmitZero().Where("id = ?", step.ID).Set("creation_date = ?", creationTime).Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}
	endTime, _ := time.Parse(time.RFC3339, step.EndDate)
	_, err = tx.NewUpdate().Model(step).OmitZero().Where("id = ?", step.ID).Set("end_date = ?", endTime).Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}
	if len(step.Prizes) != 0 {
		_, err = tx.NewInsert().Model(&step.Prizes).ExcludeColumn("creation_date").Exec(ctx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	if len(step.Images) != 0 {
		_, err = tx.NewInsert().Model(&step.Images).Exec(ctx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (s *StepRepo) GetStep(ctx context.Context, id uuid.UUID) (*models.Step, error) {
	step := new(models.Step)
	err := s.DB.NewSelect().
		Model(step).
		Relation("Prizes").
		Relation("Images").
		Relation("ActiveStaff").
		Where("id = ?", id).
		Scan(ctx)
	return step, err
}

func (s *StepRepo) GetSteps(ctx context.Context, eventID uuid.UUID) ([]*models.Step, error) {
	steps := new([]*models.Step)
	err := s.DB.NewSelect().Model(steps).Relation("Images").Where("event_id = ?", eventID).Scan(ctx)
	return *steps, err
}

func (s *StepRepo) DeleteStep(ctx context.Context, id uuid.UUID) error {
	_, err := s.DB.NewDelete().Model(&models.Step{}).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *StepRepo) UpdateStep(ctx context.Context, step *models.Step) error {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	_, err = tx.NewUpdate().OmitZero().Model(step).Where("id = ?", step.ID).Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}
	if len(step.Prizes) != 0 {
		for i := range step.Prizes {
			_, err = tx.NewUpdate().OmitZero().Model(step.Prizes[i]).Where("id = ?", step.Prizes[i].ID).Exec(ctx)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	if len(step.Images) != 0 {
		_, err = tx.NewInsert().Model(&step.Images).Exec(ctx)
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
