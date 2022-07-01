package services

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/miprokop/fication/internal/models"
	"github.com/miprokop/fication/internal/persistence/postgres"
	"github.com/spf13/viper"
	"time"
)

type Service struct {
	Auth         Auth
	Staff        Staff
	Organization Organization
	Team         Team
	Prize        Prize
	Step         Step
	Event        Event
}

type Auth interface {
	GenerateToken(email, password string) (string, uuid.UUID, error)
	ParseToken(accessToken string) (uuid.UUID, error)
}

type Staff interface {
	CreateStaffUser(ctx context.Context, staff *models.StaffSignUp) error
	GetStaff(ctx context.Context, id uuid.UUID) (*models.Staff, error)
	GetStaffByEvent(ctx context.Context, eventID uuid.UUID) ([]*models.Staff, error)
	GetStaffByStep(ctx context.Context, stepID uuid.UUID) ([]*models.Staff, error)
	GetStaffByOrganization(ctx context.Context, organizationName uuid.UUID) ([]models.Staff, error)
	DeleteStaff(ctx context.Context, id uuid.UUID) error
	UpdateStaff(ctx context.Context, staff *models.Staff) error
	SetStaffRole(ctx context.Context, role models.StaffRole) error
	GetInvites(ctx context.Context, id uuid.UUID) ([]models.StaffEvents, error)
	GetStaffPrizes(ctx context.Context, id uuid.UUID) ([]models.Prize, error)
	UploadImage(ctx context.Context, image models.StaffImage) error
	GetPosition(ctx context.Context, id uuid.UUID) (*models.Position, error)
	GetDefaultPosition(ctx context.Context, orgID uuid.UUID) (models.Position, error)
	GetAllPositions(ctx context.Context, orgID uuid.UUID) ([]models.Position, error)
	CreatePosition(ctx context.Context, position *models.Position) error
	RemovePermissionsFromPosition(ctx context.Context, permissions models.Permissions) error
	UpdatePosition(ctx context.Context, position *models.Position) error
	DeletePosition(ctx context.Context, id uuid.UUID) error
	AssignPosition(ctx context.Context, userID, positionID uuid.UUID) error
	RemoveFromPosition(ctx context.Context, userID uuid.UUID) error
	GrantPermission(ctx context.Context, granterID, positionID uuid.UUID, perm models.Permission) error
	RevokePermission(ctx context.Context, positionID uuid.UUID, perm models.Permission) error
}

type Organization interface {
	GetOrganizations(ctx context.Context) ([]*models.Organization, error)
	GetOrganization(ctx context.Context, id uuid.UUID) (*models.Organization, error)
	CreateOrganization(ctx context.Context, org *models.Organization, userID uuid.UUID) error
	UpdateOrganization(ctx context.Context, org *models.Organization) error
	AddUsersToOrg(ctx context.Context, orgID uuid.UUID, users []*models.StaffInsertion) error
	DeleteOrganization(ctx context.Context, id uuid.UUID) error
	GetOrganizationEvents(ctx context.Context, orgID, staffID uuid.UUID) ([]*models.Event, error)
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
	GetTeamByName(ctx context.Context, orgID uuid.UUID, name string) (*models.Team, error)
	GetTeamsByEvent(ctx context.Context, eventID uuid.UUID) ([]*models.Team, error)
	GetTeamByID(ctx context.Context, id uuid.UUID) (*models.Team, error)
	UpdateTeam(ctx context.Context, team *models.Team) error
	DeleteTeam(ctx context.Context, id uuid.UUID) error
}

type Prize interface {
	CreatePrize(ctx context.Context, prize *models.Prize) error
	GetPrize(ctx context.Context, id uuid.UUID) (*models.Prize, error)
	GetPrizes(ctx context.Context, userID uuid.UUID) ([]*models.Prize, error)
	GetAllPrizes(ctx context.Context) ([]*models.Prize, error)
	GetPrizesByType(ctx context.Context, prizeType models.PrizeType) ([]*models.Prize, error)
	DeletePrize(ctx context.Context, id uuid.UUID) error
	GivePrize(ctx context.Context, userID, prizeID uuid.UUID) error
	UpdatePrize(ctx context.Context, prize *models.Prize) error
}

type Step interface {
	CreateStep(ctx context.Context, step *models.Step, creationTime, endTime time.Time) error
	GetStep(ctx context.Context, id uuid.UUID) (*models.Step, error)
	GetSteps(ctx context.Context, eventID uuid.UUID) ([]*models.Step, error)
	DeleteStep(ctx context.Context, id uuid.UUID) error
	GetStepPrizes(ctx context.Context, id uuid.UUID) ([]*models.Prize, error)
	AssignStaff(ctx context.Context, staffID, stepID uuid.UUID) error
	PassStaff(ctx context.Context, stepID, statusID uuid.UUID, status models.Accomplishment,
		score uint) error
	UpdateStep(ctx context.Context, step *models.Step) error
}

type Event interface {
	RemoveStaffFromEvent(ctx context.Context, events models.StaffEvents) error
	GetInvites(ctx context.Context, staffID uuid.UUID) ([]*models.StaffEvents, error)
	CreateEvent(ctx context.Context, event *models.Event) error
	GetEvent(ctx context.Context, id uuid.UUID) (*models.Event, error)
	GetEventsByTeamID(ctx context.Context, orgID uuid.UUID) ([]*models.Event, error)
	AssignStaff(ctx context.Context, events []models.StaffEvents, eventID uuid.UUID) error
	GetEventsByCommandID(ctx context.Context, commandID uuid.UUID) ([]*models.Event, error)
	AnswerInvitation(ctx context.Context, events models.StaffEvents) error
	GetStaffScore(ctx context.Context, eventID, staffID uuid.UUID) (models.StaffScore, error)
	DeleteEvent(ctx context.Context, id uuid.UUID) error
	UpdateEvent(ctx context.Context, step *models.Event) error
	GetStaffsEvents(ctx context.Context, id uuid.UUID, role string) ([]*models.Event, error)
}

func NewService(r *postgres.Repository) *Service {
	ctx := context.Background()

	staff := &models.StaffSignUp{
		ID:              uuid.New(),
		FirstName:       viper.GetString("admin.firstName"),
		LastName:        viper.GetString("admin.LastName"),
		Email:           viper.GetString("admin.email"),
		Password:        generatePasswordHash(viper.GetString("admin.password")),
		Sex:             models.Male,
		AdditionalInfo:  "admin",
		TeamID:          models.DefaultTeam.ID,
		PositionID:      models.AdminPosition.ID,
		OrganizationID:  models.DefaultOrganization.ID,
		TextColor:       models.HexColor(viper.GetString("admin.textColor")),
		BackgroundColor: "#ffffff",
	}

	if _, err := r.Staff.GetStaffAuth(ctx, staff.Email, staff.Password); err != nil && err == sql.ErrNoRows {
		_, err := r.Staff.CreateStaffUser(ctx, staff)
		if err != nil {
			panic(err)
		}
	}
	return &Service{
		Auth:         NewAuthService(ctx, r.Staff),
		Staff:        NewStaffService(ctx, r.Staff),
		Organization: NewOrganizationService(ctx, r.Organization),
		Team:         NewTeamService(ctx, r.Team),
		Prize:        NewPrizeService(ctx, r.Prize),
		Step:         NewStepService(ctx, r.Step),
		Event:        NewEventService(ctx, r.Event),
	}
}
