package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/miprokop/fication/internal/models"
	"github.com/miprokop/fication/internal/persistence/postgres"
)

type TeamService struct {
	repo postgres.Team
	ctx  context.Context
}

func (t *TeamService) GetTeamByName(ctx context.Context, orgID uuid.UUID, name string) (*models.Team, error) {
	return t.repo.GetTeamByName(ctx, orgID, name)
}

func (t *TeamService) CreateTeam(ctx context.Context, team *models.Team) error {
	return t.repo.CreateTeam(ctx, team)
}

func (t *TeamService) GetTeamsByOrganizationID(ctx context.Context, id uuid.UUID) ([]*models.Team, error) {
	return t.repo.GetTeamsByOrganizationID(ctx, id)
}

func (t *TeamService) GetTeamsByEvent(ctx context.Context, eventID uuid.UUID) ([]*models.Team, error) {
	return t.repo.GetTeamsByEvent(ctx, eventID)
}

func (t *TeamService) GetTeamByID(ctx context.Context, id uuid.UUID) (*models.Team, error) {
	return t.repo.GetTeamByID(ctx, id)
}

func (t *TeamService) UpdateTeam(ctx context.Context, orgType *models.Team) error {
	return t.repo.UpdateTeam(ctx, orgType)
}

func (t *TeamService) DeleteTeam(ctx context.Context, id uuid.UUID) error {
	return t.repo.DeleteTeam(ctx, id)
}

func NewTeamService(ctx context.Context, repo postgres.Team) *TeamService {
	return &TeamService{repo: repo, ctx: ctx}
}
