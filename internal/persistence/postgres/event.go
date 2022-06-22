package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/miprokop/fication/internal/models"
	"github.com/uptrace/bun"
)

type EventRepo struct {
	DB  *bun.DB
	ctx context.Context
}

func NewEventRepo(ctx context.Context, DB *bun.DB) *EventRepo {
	return &EventRepo{DB: DB, ctx: ctx}
}

func (e *EventRepo) CreateEvent(ctx context.Context, step *models.Event) error {
	//TODO implement me
	panic("implement me")
}

func (e *EventRepo) GetEvent(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (e *EventRepo) GetEventsByOrgID(ctx context.Context, orgID uuid.UUID) ([]*models.Step, error) {
	//TODO implement me
	panic("implement me")
}

func (e *EventRepo) GetEventsByCommandID(ctx context.Context, commandID uuid.UUID) ([]*models.Step, error) {
	//TODO implement me
	panic("implement me")
}

func (e *EventRepo) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (e *EventRepo) UpdateEvent(ctx context.Context, step *models.Event) error {
	//TODO implement me
	panic("implement me")
}

func (e *EventRepo) GetStaffsEvents(ctx context.Context, id uuid.UUID) ([]*models.Event, error) {
	var events = new([]*models.Event)

	err := e.DB.NewSelect().Model(events).Relation("StaffEvents").Where("s_e.user_id = ?", id).Scan(ctx)
	return *events, err
}
