package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/miprokop/fication/internal/models"
	"github.com/miprokop/fication/internal/persistence/postgres"
)

type EventService struct {
	repo postgres.Event
	ctx  context.Context
}

func (e *EventService) CreateEvent(ctx context.Context, step *models.Event) error {
	//TODO implement me
	panic("implement me")
}

func (e *EventService) GetEvent(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (e *EventService) GetEventsByOrgID(ctx context.Context, orgID uuid.UUID) ([]*models.Step, error) {
	//TODO implement me
	panic("implement me")
}

func (e *EventService) GetEventsByCommandID(ctx context.Context, commandID uuid.UUID) ([]*models.Step, error) {
	//TODO implement me
	panic("implement me")
}

func (e *EventService) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (e *EventService) UpdateEvent(ctx context.Context, step *models.Event) error {
	//TODO implement me
	panic("implement me")
}

func (e *EventService) GetStaffsEvents(ctx context.Context, id uuid.UUID) ([]*models.Event, error) {
	return e.repo.GetStaffsEvents(ctx, id)
}

func NewEventService(ctx context.Context, repo postgres.Event) *EventService {
	return &EventService{repo: repo, ctx: ctx}
}
