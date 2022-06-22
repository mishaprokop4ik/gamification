package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/miprokop/fication/internal/models"
	"github.com/miprokop/fication/internal/persistence/postgres"
)

type PrizeService struct {
	repo postgres.Prize
	ctx  context.Context
}

func (p PrizeService) CreatePrize(ctx context.Context, prize *models.Prize) error {
	//TODO implement me
	panic("implement me")
}

func (p PrizeService) GetPrize(ctx context.Context, id uuid.UUID) (*models.Prize, error) {
	//TODO implement me
	panic("implement me")
}

func (p PrizeService) GetPrizes(ctx context.Context, userID uuid.UUID) ([]*models.Prize, error) {
	//TODO implement me
	panic("implement me")
}

func (p PrizeService) GetAllPrizes(ctx context.Context) ([]*models.Prize, error) {
	//TODO implement me
	panic("implement me")
}

func (p PrizeService) DeletePrize(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (p PrizeService) GivePrize(ctx context.Context, userID, prizeID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (p PrizeService) UpdatePrize(ctx context.Context, prize *models.Prize) error {
	//TODO implement me
	panic("implement me")
}

func NewPrizeService(ctx context.Context, repo postgres.Prize) *PrizeService {
	return &PrizeService{repo: repo, ctx: ctx}
}
