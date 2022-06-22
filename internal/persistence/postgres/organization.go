package postgres

import (
	"context"
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
	err := o.DB.NewSelect().Model(organizations).Scan(ctx)
	return *organizations, err
}

func (o *OrganizationRepo) GetOrganization(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	var organization = new(models.Organization)
	err := o.DB.NewSelect().Model(organization).Where("id = ?", id).Scan(ctx)
	return organization, err
}

func (o *OrganizationRepo) CreateOrganization(ctx context.Context, org *models.Organization) error {
	_, err := o.DB.NewInsert().Model(org).Exec(ctx)
	return err
}

func (o *OrganizationRepo) UpdateOrganization(ctx context.Context, org *models.Organization) error {
	_, err := o.DB.NewUpdate().OmitZero().Model(org).WherePK().Exec(ctx)
	return err
}

func (o *OrganizationRepo) AddUsersToOrg(ctx context.Context, staff *models.Staff) error {
	_, err := o.DB.NewUpdate().Model(staff).Column("company_id").Where("id = ?", staff.ID).Exec(ctx)
	return err
}

func (o *OrganizationRepo) DeleteOrganization(ctx context.Context, id uuid.UUID) error {
	_, err := o.DB.NewDelete().Model(&models.Organization{}).Where("id = ?", id).Exec(ctx)
	return err
}

func (o *OrganizationRepo) GetOrganizationEvents(ctx context.Context, id uuid.UUID) ([]*models.Event, error) {
	var events []*models.Event
	_, err := o.DB.NewSelect().Model(&events).Where("organization_id = ?", id).Exec(ctx)
	return events, err
}

func (o *OrganizationRepo) GetOrganizationStaff(ctx context.Context, orgID uuid.UUID) ([]models.Staff, error) {
	var staff []models.Staff
	_, err := o.DB.NewSelect().Model(&staff).Where("organization_id = ?", orgID).Exec(ctx)
	return staff, err
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
	_, err := o.DB.NewSelect().Model(&organizationTypes).Exec(ctx)
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
