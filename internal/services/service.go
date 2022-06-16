package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/miprokop/fication/internal/models"
)

type Servicer interface {
	Auth
	Staff
	Organization
	Team
	Prize
	Step
	Event
}

type Auth interface {
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Staff interface {
	CreateStaffUser(ctx context.Context, staff *models.Staff) error
	GetStaff(ctx context.Context, id uuid.UUID) (*models.Staff, error)
	GetStaffByEvent(ctx context.Context, eventID uuid.UUID) ([]*models.Staff, error)
	GetStaffByStep(ctx context.Context, stepID uuid.UUID) ([]*models.Staff, error)
	GetStaffByOrganization(ctx context.Context, organizationName uuid.UUID) ([]models.Staff, error)
	DeleteStaff(ctx context.Context, id uuid.UUID) error
	UpdateStaff(ctx context.Context, staff *models.Staff) error
	SetStaffRole(ctx context.Context, role models.StaffRole) error
	GetInvites(ctx context.Context, id uuid.UUID) ([]models.Invite, error)
	GetStaffPrizes(ctx context.Context, id uuid.UUID) ([]models.Prize, error)

	GetRole(ctx context.Context, id uuid.UUID) (*models.Role, error)
	GetAllRoles(ctx context.Context, orgID uuid.UUID) ([]models.Role, error)
	CreateRole(ctx context.Context, role *models.Role) error
	UpdateRole(ctx context.Context, role *models.Role) error
	DeleteRole(ctx context.Context, id uuid.UUID) error
	AssignRole(ctx context.Context, userID, roleID uuid.UUID) error
	GrantPermission(ctx context.Context, granterID, roleID uuid.UUID, perm models.Permission) error
	RevokePermission(ctx context.Context, roleID uuid.UUID, perm models.Permission) error
}

type Organization interface {
	GetOrganization(ctx context.Context, id uuid.UUID) (models.Organization, error)
	CreateOrganization(ctx context.Context, org *models.Organization) error
	UpdateOrganization(ctx context.Context, org *models.Organization) error
	AddUsersToOrg(ctx context.Context, orgID uuid.UUID, userIDs []uuid.UUID) error
	DeleteOrganization(ctx context.Context, id uuid.UUID) error
	GetOrganizationEvents(ctx context.Context, id uuid.UUID) []*models.Event
	GetOrganizationStaff(ctx context.Context, orgID uuid.UUID) ([]models.StaffInfo, error)

	CreateOrganizationType(ctx context.Context, orgType *models.OrganizationType) error
	GetOrganizationTypeByID(ctx context.Context, id uuid.UUID) (*models.OrganizationType, error)
	GetOrganizationTypes(ctx context.Context) ([]*models.OrganizationType, error)
	UpdateOrganizationType(ctx context.Context, orgType *models.OrganizationType) error
	DeleteOrganizationTypeByID(ctx context.Context, id uuid.UUID) error
}

type Team interface {
	CreateTeam(ctx context.Context, team *models.Team) error
	GetTeamsByOrganizationID(ctx context.Context, id uuid.UUID) ([]*models.Team, error)
	GetTeamsByEvent(ctx context.Context, eventID uuid.UUID) ([]*models.Team, error)
	GetTeamByID(ctx context.Context, id uuid.UUID) (*models.Team, error)
	UpdateTeam(ctx context.Context, orgType *models.Team) error
	DeleteTeam(ctx context.Context, id uuid.UUID) error
}

type Prize interface {
	CreatePrize(ctx context.Context, prize *models.Prize) error
	GetPrize(ctx context.Context, id uuid.UUID) (*models.Prize, error)
	GetPrizes(ctx context.Context, userID uuid.UUID) ([]*models.Prize, error)
	GetAllPrizes(ctx context.Context) ([]*models.Prize, error)
	DeletePrize(ctx context.Context, id uuid.UUID) error
	GivePrize(ctx context.Context, userID, prizeID uuid.UUID) error
	UpdatePrize(ctx context.Context, prize *models.Prize) error
}

type Step interface {
	CreateStep(ctx context.Context, step *models.Step) error
	GetStep(ctx context.Context, id uuid.UUID) (*models.Step, error)
	GetSteps(ctx context.Context, stepID uuid.UUID) ([]*models.Step, error)
	DeleteStep(ctx context.Context, id uuid.UUID) error
	UpdateStep(ctx context.Context, step *models.Step) error
}

type Event interface {
	CreateEvent(ctx context.Context, step *models.Event) error
	GetEvent(ctx context.Context, id uuid.UUID) (*models.Event, error)
	GetEventsByOrgID(ctx context.Context, orgID uuid.UUID) ([]*models.Step, error)
	GetEventsByCommandID(ctx context.Context, commandID uuid.UUID) ([]*models.Step, error)
	DeleteEvent(ctx context.Context, id uuid.UUID) error
	UpdateEvent(ctx context.Context, step *models.Event) error
	GetStaffsEvents(ctx context.Context, id uuid.UUID) ([]*models.Event, error)
}

type Repo interface {
	Staff
	Organization
	Team
	Prize
	Step
	Event
}
