package postgres

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/miprokop/fication/internal/models"
	"github.com/uptrace/bun"
)

type OrganizationRepo struct {
	DB  *bun.DB
	ctx context.Context
}

func (o *OrganizationRepo) GetOrganizations(ctx context.Context) ([]*models.Organization, error) {
	var organizations = new([]*models.Organization)
	err := o.DB.NewSelect().Model(organizations).
		Relation("Types").Relation("Positions").Relation("Teams").Scan(ctx)
	return *organizations, err
}

func (o *OrganizationRepo) GetOrganization(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	var organization = new(models.Organization)
	err := o.DB.NewSelect().Model(organization).Relation("Types").Relation("Positions").Relation("Teams").Where("id = ?", id).Scan(ctx)
	return organization, err
}

func (o *OrganizationRepo) CreateOrganization(ctx context.Context, org *models.Organization, userID uuid.UUID) error {
	tx, err := o.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	_, err = tx.NewInsert().Model(org).Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}
	for i := range org.Positions {
		org.Positions[i].CompanyID = org.ID
	}
	_, err = tx.NewInsert().Model(&org.Positions).Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}
	for i := range org.Positions {
		for _, p := range org.Positions[i].Permissions {
			p.PositionID = org.Positions[i].ID
			p.GrantedBy = userID
			_, err = tx.NewInsert().Model(p).Exec(ctx)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	for i := range org.Teams {
		org.Teams[i].OrganizationID = org.ID
	}
	_, err = tx.NewInsert().Model(&org.Teams).Exec(ctx)
	var orgTypes = make([]models.OrganizationsTypes, len(org.Types))
	for i := range org.Types {
		orgTypes[i] = models.OrganizationsTypes{
			ID:             uuid.New(),
			OrganizationID: org.ID,
			TypeID:         org.Types[i].ID,
		}
	}
	_, err = tx.NewInsert().Model(&orgTypes).Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (o *OrganizationRepo) UpdateOrganization(ctx context.Context, org *models.Organization) error {
	_, err := o.DB.NewUpdate().OmitZero().Model(org).WherePK().Exec(ctx)
	return err
}

func (o *OrganizationRepo) AddUsersToOrg(ctx context.Context, staff *models.Staff) error {
	_, err := o.DB.NewUpdate().Model(staff).
		Column("company_id").
		Column("team_id").
		Column("position_id").
		Where("id = ?", staff.ID).
		Exec(ctx)
	return err
}

func (o *OrganizationRepo) DeleteOrganization(ctx context.Context, id uuid.UUID) error {
	_, err := o.DB.NewDelete().Model(&models.Organization{}).Where("id = ?", id).Exec(ctx)
	return err
}

func (o *OrganizationRepo) GetOrganizationEvents(ctx context.Context, id uuid.UUID) ([]*models.Event, error) {
	var events []*models.Event
	err := o.DB.NewSelect().Model(&events).Where("organization_id = ?", id).Scan(ctx)
	return events, err
}

func (o *OrganizationRepo) GetOrganizationStaff(ctx context.Context, orgID uuid.UUID) ([]models.Staff, error) {
	staff := new([]models.Staff)
	err := o.DB.NewSelect().Model(staff).Where("staff.company_id = ?", orgID).
		Relation("Position").Relation("Organization").Scan(ctx)
	return *staff, err
}

func (o *OrganizationRepo) CreateOrganizationType(ctx context.Context, orgType *models.OrganizationType) error {
	_, err := o.DB.NewInsert().Model(orgType).Exec(ctx)
	return err
}

func (o *OrganizationRepo) GetOrganizationTypeByID(ctx context.Context, id uuid.UUID) (*models.OrganizationType, error) {
	var organizationType = new(models.OrganizationType)
	err := o.DB.NewSelect().Model(organizationType).Where("id = ?", id).Scan(ctx)
	return organizationType, err
}

func (o *OrganizationRepo) GetOrganizationTypes(ctx context.Context) ([]*models.OrganizationType, error) {
	var organizationTypes []*models.OrganizationType
	err := o.DB.NewSelect().Model(&organizationTypes).Scan(ctx)
	return organizationTypes, err
}

func (o *OrganizationRepo) UpdateOrganizationType(ctx context.Context, orgType *models.OrganizationType) error {
	_, err := o.DB.NewUpdate().OmitZero().Model(orgType).Where("id = ?", orgType.ID).Exec(ctx)
	return err
}

func (o *OrganizationRepo) DeleteOrganizationTypeByID(ctx context.Context, id uuid.UUID) error {
	_, err := o.DB.NewDelete().Model(&models.OrganizationType{}).Where("id = ?", id).Exec(ctx)
	return err
}

func NewOrganizationRepo(ctx context.Context, DB *bun.DB) *OrganizationRepo {
	return &OrganizationRepo{DB: DB, ctx: ctx}
}
