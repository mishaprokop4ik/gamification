package postgres

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/miprokop/fication/internal/models"
	"github.com/uptrace/bun"
	"time"
)

type PrizeRepo struct {
	DB  *bun.DB
	ctx context.Context
}

func (p *PrizeRepo) GetPrizesByType(ctx context.Context, prizeType models.PrizeType) ([]*models.Prize, error) {
	var prizes = new([]*models.Prize)
	err := p.DB.NewSelect().Model(prizes).Where("prize_type = ?", prizeType).Scan(ctx)

	return *prizes, err
}

func (p *PrizeRepo) CreatePrize(ctx context.Context, prize *models.Prize) error {
	tx, err := p.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.NewInsert().Model(prize).ExcludeColumn("creation_date").Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}
	creationTime, err := time.Parse(time.RFC3339, prize.CreationDate)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.NewUpdate().Model(&models.Prize{}).OmitZero().
		Set("creation_date = ?", creationTime).Where("id = ?", prize.ID).Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (p *PrizeRepo) GetPrize(ctx context.Context, id uuid.UUID) (*models.Prize, error) {
	var prize = new(models.Prize)
	err := p.DB.NewSelect().Model(prize).Where("id = ?", id).Scan(ctx)
	return prize, err
}

func (p *PrizeRepo) GetPrizes(ctx context.Context, userID uuid.UUID) ([]*models.Prize, error) {
	var prizes = new([]*models.Prize)
	err := p.DB.NewSelect().Model(prizes).Where("created_by = ?", userID).Scan(ctx)
	return *prizes, err
}

func (p *PrizeRepo) GetAllPrizes(ctx context.Context) ([]*models.Prize, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PrizeRepo) DeletePrize(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (p *PrizeRepo) GivePrize(ctx context.Context, staffPrize *models.StaffPrize) error {
	_, err := p.DB.NewInsert().Model(staffPrize).Exec(ctx)
	return err
}

func (p *PrizeRepo) UpdatePrize(ctx context.Context, prize *models.Prize) error {
	_, err := p.DB.NewUpdate().Model(prize).OmitZero().Where("id = ?", prize.ID).Exec(ctx)
	return err
}

func NewPrizeRepo(ctx context.Context, DB *bun.DB) *PrizeRepo {
	return &PrizeRepo{DB: DB, ctx: ctx}
}
