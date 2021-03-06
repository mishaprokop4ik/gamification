package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/miprokop/fication/internal/models"
	"github.com/uptrace/bun"
)

type TeamRepo struct {
	DB  *bun.DB
	ctx context.Context
}

func (t *TeamRepo) GetTeamByName(ctx context.Context, orgID uuid.UUID, name string) (*models.Team, error) {
	team := new(models.Team)
	err := t.DB.NewSelect().Model(team).Where("organization_id = ?", orgID).Where("name = ?", name).Scan(ctx)
	return team, err
}

func (t *TeamRepo) CreateTeam(ctx context.Context, team *models.Team) error {
	_, err := t.DB.NewInsert().Model(team).Exec(ctx)
	return err
}

func (t *TeamRepo) GetTeamsByOrganizationID(ctx context.Context, id uuid.UUID) ([]*models.Team, error) {
	var teams []*models.Team
	err := t.DB.NewSelect().Model(&teams).Where("organization_id = ?", id).Scan(ctx)
	return teams, err
}

func (t *TeamRepo) GetTeamsByEvent(ctx context.Context, eventID uuid.UUID) ([]*models.Team, error) {
	var staffEvents []models.StaffEvents
	err := t.DB.NewSelect().
		Model(&staffEvents).
		Relation("Staff").
		Where("staff_events.event_id = ?", eventID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	ids := make([]uuid.UUID, len(staffEvents))
	for i := range staffEvents {
		ids[i] = staffEvents[i].Staff.TeamID
	}
	var teams = make([]*models.Team, 0)
	err = t.DB.NewSelect().
		Model(&teams).
		Where("id IN (?)", bun.In(ids)).
		Scan(ctx)
	return teams, err
}

func (t *TeamRepo) GetTeamByID(ctx context.Context, id uuid.UUID) (*models.Team, error) {
	var team = new(models.Team)
	err := t.DB.NewSelect().
		Model(team).
		Relation("Organization").
		Relation("Staff").
		Where("team.id = ?", id).Scan(ctx)
	return team, err
}

func (t *TeamRepo) UpdateTeam(ctx context.Context, team *models.Team) error {
	_, err := t.DB.NewUpdate().OmitZero().Model(team).Where("id = ?", team.ID).Exec(ctx)
	return err
}

func (t *TeamRepo) DeleteTeam(ctx context.Context, id uuid.UUID) error {
	_, err := t.DB.NewDelete().Model(&models.Team{}).Where("id = ?", id).Exec(ctx)
	return err
}

func NewTeamRepo(ctx context.Context, DB *bun.DB) *TeamRepo {
	return &TeamRepo{DB: DB, ctx: ctx}
}
