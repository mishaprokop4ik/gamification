package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/miprokop/fication/internal/models"
	"github.com/uptrace/bun"
)

type EventRepo struct {
	DB  *bun.DB
	ctx context.Context
}

func (e *EventRepo) RemoveStaffFromEvent(ctx context.Context, events models.StaffEvents) error {
	tx, err := e.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.NewDelete().Model(&events).
		Where("user_id = ?", events.StaffID).
		Where("event_id = ?", events.EventID).
		Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}

	oldEvent, err := e.GetEvent(ctx, events.EventID)
	if err != nil {
		tx.Rollback()
		return err
	}
	ids := make([]uuid.UUID, len(oldEvent.Steps))
	for i := range oldEvent.Steps {
		ids[i] = oldEvent.Steps[i].ID
	}

	_, err = tx.NewDelete().Model(&models.StepStaff{}).
		Where("staff_id = ?", events.StaffID).
		Where("step_id IN (?)", bun.In(ids)).
		Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (e *EventRepo) IsStaffInOrg(ctx context.Context, staffID, teamID uuid.UUID) (bool, error) {
	exists, err := e.DB.NewSelect().Model((*models.Staff)(nil)).Where("organization_id = ?", teamID).Where("staff_id = ?", staffID).Exists(ctx)
	return exists, err
}

func (e *EventRepo) GetStaff(ctx context.Context, id uuid.UUID) (*models.Staff, error) {
	var staff = new(models.Staff)
	err := e.DB.NewSelect().Model(staff).
		Where("staff.id = ?", id).
		Scan(ctx)
	return staff, err
}

func (e *EventRepo) IsStaffInTeam(ctx context.Context, staffID, teamID uuid.UUID) (bool, error) {
	exists, err := e.DB.NewSelect().Model((*models.Staff)(nil)).Where("team_id = ?", teamID).Where("staff_id = ?", staffID).Exists(ctx)
	return exists, err
}

func (e *EventRepo) GetInvites(ctx context.Context, staffID uuid.UUID) ([]*models.StaffEvents, error) {
	invites := new([]*models.StaffEvents)
	err := e.DB.NewSelect().Model(invites).Relation("Event").Where("user_id = ?", staffID).Distinct().Scan(ctx)
	return *invites, err
}

func (e *EventRepo) GetStaffScore(ctx context.Context, eventID, staffID uuid.UUID) (models.StaffScore, error) {
	var score models.StaffScore
	steps := new([]*models.Step)
	err := e.DB.NewSelect().Model(steps).Where("event_id = ?", eventID).Scan(ctx)
	if err != nil {
		return models.StaffScore{}, err
	}
	stepsIDs := make([]uuid.UUID, len(*steps))
	for i := range *steps {
		stepsIDs[i] = (*steps)[i].ID
	}
	stepsStaff := new([]*models.StepStaff)
	err = e.DB.NewSelect().Model(stepsStaff).
		Where("staff_id = ?", staffID).Where("step_id IN (?)", bun.In(stepsIDs)).Scan(ctx)
	if err != nil {
		return models.StaffScore{}, err
	}
	for i := range *stepsStaff {
		score.Score += int((*stepsStaff)[i].Score)
	}

	return score, nil
}

func (e *EventRepo) AnswerInvitation(ctx context.Context, events models.StaffEvents) error {
	_, err := e.DB.NewUpdate().Model(&events).OmitZero().
		Where("event_id = ?", events.EventID).Where("user_id = ?", events.StaffID).Exec(ctx)
	return err
}

func (e *EventRepo) AssignStaff(ctx context.Context, events models.StaffEvents) error {
	exists, err := e.DB.NewSelect().
		Model(&models.StaffEvents{}).
		Where("user_id = ?", events.StaffID).
		Where("event_id = ?", events.EventID).Exists(ctx)
	if err != nil {
		return err
	}
	if !exists {
		_, err = e.DB.NewInsert().Model(&events).Exec(ctx)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("invitation exists")
	}
	return nil
}

func NewEventRepo(ctx context.Context, DB *bun.DB) *EventRepo {
	return &EventRepo{DB: DB, ctx: ctx}
}

func (e *EventRepo) CreateEvent(ctx context.Context, event *models.Event) error {
	tx, err := e.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	_, err = tx.NewInsert().Model(event).Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}
	if len(event.StaffEvents) != 0 {
		_, err = tx.NewInsert().
			Model(&event.StaffEvents).Exec(ctx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (e *EventRepo) GetEvent(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	event := new(models.Event)
	err := e.DB.NewSelect().
		Model(event).
		Relation("Steps").
		Where("event.id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	staffEvents := new([]*models.StaffEvents)
	err = e.DB.NewSelect().
		Model(staffEvents).
		Where("event_id = ?", id).
		Scan(ctx)
	event.StaffEvents = *staffEvents
	return event, err
}

func (e *EventRepo) GetEventsByTeamID(ctx context.Context, teamID uuid.UUID) ([]*models.Event, error) {
	events := new([]*models.Event)
	staff := new([]*models.Staff)
	staffEvents := new([]*models.StaffEvents)
	err := e.DB.NewSelect().Model(staff).Where("team_id = ?", teamID).Scan(ctx)
	if err != nil {
		return nil, err
	}
	ids := make([]uuid.UUID, len(*staff))
	for i := range *staff {
		ids[i] = (*staff)[i].ID
	}
	err = e.DB.NewSelect().Model(staffEvents).Where("user_id IN (?)", bun.In(ids)).Scan(ctx)
	if err != nil {
		return nil, err
	}
	ids = nil
	ids = make([]uuid.UUID, len(*staffEvents))
	for i := range *staffEvents {
		ids[i] = (*staffEvents)[i].EventID
	}
	err = e.DB.NewSelect().Model(events).Where("id IN (?)", bun.In(ids)).Scan(ctx)
	return *events, err
}

func (e *EventRepo) GetEventsByCommandID(ctx context.Context, commandID uuid.UUID) ([]*models.Event, error) {
	staff := new([]*models.Staff)
	err := e.DB.NewSelect().Model(staff).Where("team_id = ?", commandID).Scan(ctx)
	if err != nil {
		return nil, err
	}
	staffIDs := make([]uuid.UUID, len(*staff))
	for i := range *staff {
		staffIDs[i] = (*staff)[i].ID
	}
	events := new([]*models.Event)
	err = e.DB.NewSelect().Model(events).Relation("StaffEvents").
		Where("staff_events.user_id IN (?)", bun.In(staffIDs)).Scan(ctx)
	return *events, nil
}

func (e *EventRepo) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	_, err := e.DB.NewDelete().Model(&models.Event{}).Where("id = ?", id).Exec(ctx)
	return err
}

func (e *EventRepo) UpdateEvent(ctx context.Context, event *models.Event) error {
	tx, err := e.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	_, err = tx.NewUpdate().OmitZero().Model(event).Where("id = ?", event.ID).Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}
	if len(event.StaffEvents) != 0 {
		_, err = tx.NewInsert().
			Model(&event.StaffEvents).Exec(ctx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return nil
}

func (e *EventRepo) GetStaffsEvents(ctx context.Context, id uuid.UUID, role string) ([]*models.Event, error) {
	var staffEvents = new([]*models.StaffEvents)

	err := e.DB.NewSelect().Model(staffEvents).
		Where("user_id = ?", id).
		Where("user_role = ?", role).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	ids := make([]uuid.UUID, len(*staffEvents))
	for i := range *staffEvents {
		ids[i] = (*staffEvents)[i].EventID
	}
	events := new([]*models.Event)
	err = e.DB.NewSelect().Model(events).
		Where("id IN (?)", bun.In(ids)).
		Scan(ctx)
	return *events, err
}
