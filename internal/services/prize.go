package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/miprokop/fication/internal/models"
	"github.com/miprokop/fication/internal/persistence/postgres"
)

type PrizeService struct {
	repo postgres.Prize
	ctx  context.Context
}

func (p *PrizeService) GetPrizesByType(ctx context.Context, prizeType models.PrizeType) ([]*models.Prize, error) {
	return p.repo.GetPrizesByType(ctx, prizeType)
}

func (p *PrizeService) CreatePrize(ctx context.Context, prize *models.Prize) error {
	return p.repo.CreatePrize(ctx, prize)
}

func (p *PrizeService) GetPrize(ctx context.Context, id uuid.UUID) (*models.Prize, error) {
	return p.repo.GetPrize(ctx, id)
}

func (p *PrizeService) GetPrizes(ctx context.Context, userID uuid.UUID) ([]*models.Prize, error) {
	return p.repo.GetPrizes(ctx, userID)
}

func (p *PrizeService) GetAllPrizes(ctx context.Context) ([]*models.Prize, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PrizeService) DeletePrize(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (p *PrizeService) GivePrize(ctx context.Context, userID, prizeID uuid.UUID) error {
	prize, err := p.repo.GetPrize(ctx, prizeID)
	if err != nil {
		return err
	}
	if prize.CurrentCount == 0 {
		return fmt.Errorf("can not give prize to user; current count of prizes is zero")
	}
	prize.CurrentCount -= 1
	err = p.repo.UpdatePrize(ctx, prize)
	if err != nil {
		return err
	}

	staffPrize := &models.StaffPrize{
		ID:      uuid.New(),
		StaffID: userID,
		PrizeID: prizeID,
	}
	return p.repo.GivePrize(ctx, staffPrize)
}

func (p *PrizeService) UpdatePrize(ctx context.Context, prize *models.Prize) error {
	return p.repo.UpdatePrize(ctx, prize)
}

func NewPrizeService(ctx context.Context, repo postgres.Prize) *PrizeService {
	return &PrizeService{repo: repo, ctx: ctx}
}
