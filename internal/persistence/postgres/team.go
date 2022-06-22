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
	err := t.DB.NewSelect().Model(&staffEvents).Where("event_id = ?", eventID).Scan(ctx)
	if err != nil {
		return nil, err
	}

	var teams []*models.Team
	var teamsSet = make(map[string]*models.Team)
	for i := 0; i < len(staffEvents); i++ {
		if _, ok := teamsSet[staffEvents[i].Staff.Team.Name]; !ok {
			teamsSet[staffEvents[i].Staff.Team.Name] = staffEvents[i].Staff.Team
			teams = append(teams, staffEvents[i].Staff.Team)
		}
	}
	return teams, err
}

func (t *TeamRepo) GetTeamByID(ctx context.Context, id uuid.UUID) (*models.Team, error) {
	var team = new(models.Team)
	err := t.DB.NewSelect().
		Model(team).
		Relation("Organization").
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
