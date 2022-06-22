package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/miprokop/fication/internal/models"
	"github.com/miprokop/fication/internal/persistence/postgres"
)

type StepService struct {
	repo postgres.Step
	ctx  context.Context
}

func (s StepService) CreateStep(ctx context.Context, step *models.Step) error {
	//TODO implement me
	panic("implement me")
}

func (s StepService) GetStep(ctx context.Context, id uuid.UUID) (*models.Step, error) {
	//TODO implement me
	panic("implement me")
}

func (s StepService) GetSteps(ctx context.Context, stepID uuid.UUID) ([]*models.Step, error) {
	//TODO implement me
	panic("implement me")
}

func (s StepService) DeleteStep(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (s StepService) UpdateStep(ctx context.Context, step *models.Step) error {
	//TODO implement me
	panic("implement me")
}

func NewStepService(ctx context.Context, repo postgres.Step) *StepService {
	return &StepService{repo: repo, ctx: ctx}
}
