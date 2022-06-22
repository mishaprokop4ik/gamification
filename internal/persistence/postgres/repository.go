package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/miprokop/fication/internal/models"
	log "github.com/sirupsen/logrus"
)

type Repository struct {
	Staff        Staff
	Organization Organization
	Team         Team
	Prize        Prize
	Step         Step
	Event        Event
}

func NewRepository(db *Postgres) *Repository {
	ctx := context.Background()
	db.DB.RegisterModel((*models.OrganizationsTypes)(nil))
	db.DB.RegisterModel((*models.StaffEvents)(nil))
	exists, err := db.DB.NewSelect().Model(&models.DefaultOrganization).Where("name = 'default'").Exists(ctx)
	if err != nil {
		panic(err)
	}
	if !exists {
		models.DefaultOrganization.ID = uuid.New()
		_, _ = db.DB.NewInsert().Model(&models.DefaultOrganization).Exec(ctx)
	} else {
		_ = db.DB.NewSelect().Model(&models.DefaultOrganization).Where("name = 'default'").Scan(ctx)
	}
	exists, err = db.DB.NewSelect().Model(&models.DefaultPosition).Where("name = ?", models.DefaultPosition.Name).Exists(ctx)
	if err != nil {
		panic(err)
	}
	if !exists {
		models.DefaultPosition.ID = uuid.New()
		_, err := db.DB.NewInsert().Model(&models.DefaultPosition).Exec(ctx)
		if err != nil {
			log.Println(err)
		}
		_, err = db.DB.NewInsert().Model(&models.DefaultPosition.Permissions).Exec(ctx)
		if err != nil {
			log.Println(err)
		}
	} else {
		err = db.DB.NewSelect().Model(&models.DefaultPosition).Where("name = 'none'").Scan(ctx)
		if err != nil {
			log.Println(err)
		}
		for i := range models.DefaultPosition.Permissions {
			models.DefaultPosition.Permissions[i].PositionID = models.DefaultPosition.ID
		}
	}
	exists, err = db.DB.NewSelect().Model(&models.DefaultTeam).Where("name = ?", models.DefaultTeam.Name).Exists(ctx)
	if err != nil {
		panic(err)
	}
	if !exists {
		models.DefaultTeam.ID = uuid.New()
		models.DefaultTeam.OrganizationID = models.DefaultOrganization.ID
		_, err := db.DB.NewInsert().Model(&models.DefaultTeam).Exec(ctx)
		if err != nil {
			log.Println(err)
		}
	} else {
		err = db.DB.NewSelect().Model(&models.DefaultTeam).Where("name = 'none'").Scan(ctx)
		if err != nil {
			log.Println(err)
		}
	}
	return &Repository{
		Staff:        NewStaffRepo(ctx, db.DB),
		Organization: NewOrganizationRepo(ctx, db.DB),
		Team:         NewTeamRepo(ctx, db.DB),
		Prize:        NewPrizeRepo(ctx, db.DB),
		Step:         NewStepRepo(ctx, db.DB),
		Event:        NewEventRepo(ctx, db.DB),
	}
}

type StaffAuth interface {
	CreateStaffUser(ctx context.Context, staff *models.StaffSignUp) (uuid.UUID, error)
	GetStaffAuth(ctx context.Context, email, password string) (*models.Staff, error)
}

type Staff interface {
	StaffAuth
	GetStaff(ctx context.Context, id uuid.UUID) (*models.Staff, error)
	GetStaffByEvent(ctx context.Context, eventID uuid.UUID) ([]*models.Staff, error)
	GetStaffByStep(ctx context.Context, stepID uuid.UUID) ([]*models.Staff, error)
	GetStaffByOrganization(ctx context.Context, organizationName uuid.UUID) ([]models.Staff, error)
	DeleteStaff(ctx context.Context, id uuid.UUID) error
	UpdateStaff(ctx context.Context, staff *models.Staff) error
	SetStaffRole(ctx context.Context, role models.StaffRole) error
	GetInvites(ctx context.Context, id uuid.UUID) ([]models.StaffEvents, error)
	GetStaffPrizes(ctx context.Context, id uuid.UUID) ([]models.Prize, error)
	SaveFile(ctx context.Context, image models.StaffImage) error
	GetRole(ctx context.Context, id uuid.UUID) (*models.Position, error)
	GetAllPositions(ctx context.Context, orgID uuid.UUID) ([]models.Position, error)
	CreatePosition(ctx context.Context, position *models.Position) error
	UpdatePosition(ctx context.Context, position *models.Position) error
	DeletePosition(ctx context.Context, id uuid.UUID) error
	AssignPosition(ctx context.Context, userID, positionID uuid.UUID) error
	GrantPermission(ctx context.Context, granterID, positionID uuid.UUID, perm models.Permission) error
	RevokePermission(ctx context.Context, roleID uuid.UUID, perm models.Permission) error
}

type Organization interface {
	GetOrganizations(ctx context.Context) ([]*models.Organization, error)
	GetOrganization(ctx context.Context, id uuid.UUID) (*models.Organization, error)
	CreateOrganization(ctx context.Context, org *models.Organization) error
	UpdateOrganization(ctx context.Context, org *models.Organization) error
	AddUsersToOrg(ctx context.Context, staff *models.Staff) error
	DeleteOrganization(ctx context.Context, id uuid.UUID) error
	GetOrganizationEvents(ctx context.Context, id uuid.UUID) ([]*models.Event, error)
	GetOrganizationStaff(ctx context.Context, orgID uuid.UUID) ([]models.Staff, error)

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
	UpdateTeam(ctx context.Context, team *models.Team) error
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
