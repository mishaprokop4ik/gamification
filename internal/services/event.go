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

func (e *EventService) GetInvites(ctx context.Context, staffID uuid.UUID) ([]*models.StaffEvents, error) {
	return e.repo.GetInvites(ctx, staffID)
}

func (e *EventService) GetStaffScore(ctx context.Context, eventID, staffID uuid.UUID) (models.StaffScore, error) {
	return e.repo.GetStaffScore(ctx, eventID, staffID)
}

func (e *EventService) AnswerInvitation(ctx context.Context, events models.StaffEvents) error {
	return e.repo.AnswerInvitation(ctx, events)
}

func (e *EventService) AssignStaff(ctx context.Context, events []models.StaffEvents, eventID uuid.UUID) error {
	for _, event := range events {
		event.ID = uuid.New()
		event.EventID = eventID
		if event.StaffRole == "" {
			event.StaffRole = models.Default
		}
		event.Status = models.InProgress
		err := e.repo.AssignStaff(ctx, event)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *EventService) CreateEvent(ctx context.Context, event *models.Event) error {
	return e.repo.CreateEvent(ctx, event)
}

func (e *EventService) GetEvent(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	return e.repo.GetEvent(ctx, id)
}

func (e *EventService) GetEventsByTeamID(ctx context.Context, orgID uuid.UUID) ([]*models.Event, error) {
	return e.repo.GetEventsByTeamID(ctx, orgID)
}

func (e *EventService) GetEventsByCommandID(ctx context.Context, commandID uuid.UUID) ([]*models.Event, error) {
	return e.repo.GetEventsByCommandID(ctx, commandID)
}

func (e *EventService) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	return e.repo.DeleteEvent(ctx, id)
}

func (e *EventService) UpdateEvent(ctx context.Context, event *models.Event) error {
	return e.repo.UpdateEvent(ctx, event)
}

func (e *EventService) GetStaffsEvents(ctx context.Context, id uuid.UUID, role string) ([]*models.Event, error) {
	return e.repo.GetStaffsEvents(ctx, id, role)
}

func NewEventService(ctx context.Context, repo postgres.Event) *EventService {
	return &EventService{repo: repo, ctx: ctx}
}
