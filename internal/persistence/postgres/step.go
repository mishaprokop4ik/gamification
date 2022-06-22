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

func (s StepRepo) CreateStep(ctx context.Context, step *models.Step) error {
	//TODO implement me
	panic("implement me")
}

func (s StepRepo) GetStep(ctx context.Context, id uuid.UUID) (*models.Step, error) {
	//TODO implement me
	panic("implement me")
}

func (s StepRepo) GetSteps(ctx context.Context, stepID uuid.UUID) ([]*models.Step, error) {
	//TODO implement me
	panic("implement me")
}

func (s StepRepo) DeleteStep(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (s StepRepo) UpdateStep(ctx context.Context, step *models.Step) error {
	//TODO implement me
	panic("implement me")
}

func NewStepRepo(ctx context.Context, DB *bun.DB) *StepRepo {
	return &StepRepo{DB: DB, ctx: ctx}
}
