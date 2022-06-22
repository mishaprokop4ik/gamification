package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/miprokop/fication/internal/models"
	"github.com/uptrace/bun"
)

type PrizeRepo struct {
	DB  *bun.DB
	ctx context.Context
}

func (p PrizeRepo) CreatePrize(ctx context.Context, prize *models.Prize) error {
	//TODO implement me
	panic("implement me")
}

func (p PrizeRepo) GetPrize(ctx context.Context, id uuid.UUID) (*models.Prize, error) {
	//TODO implement me
	panic("implement me")
}

func (p PrizeRepo) GetPrizes(ctx context.Context, userID uuid.UUID) ([]*models.Prize, error) {
	//TODO implement me
	panic("implement me")
}

func (p PrizeRepo) GetAllPrizes(ctx context.Context) ([]*models.Prize, error) {
	//TODO implement me
	panic("implement me")
}

func (p PrizeRepo) DeletePrize(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (p PrizeRepo) GivePrize(ctx context.Context, userID, prizeID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (p PrizeRepo) UpdatePrize(ctx context.Context, prize *models.Prize) error {
	//TODO implement me
	panic("implement me")
}

func NewPrizeRepo(ctx context.Context, DB *bun.DB) *PrizeRepo {
	return &PrizeRepo{DB: DB, ctx: ctx}
}
