package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/miprokop/fication/internal/models"
	"github.com/miprokop/fication/internal/persistence/postgres"
)

type OrganizationService struct {
	repo postgres.Organization
	ctx  context.Context
}

func (o *OrganizationService) GetOrganizations(ctx context.Context) ([]*models.Organization, error) {
	return o.repo.GetOrganizations(ctx)
}

func (o *OrganizationService) GetOrganization(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	return o.repo.GetOrganization(ctx, id)
}

func (o *OrganizationService) CreateOrganization(ctx context.Context, org *models.Organization, userID uuid.UUID) error {
	return o.repo.CreateOrganization(ctx, org, userID)
}

func (o *OrganizationService) UpdateOrganization(ctx context.Context, org *models.Organization) error {
	return o.repo.UpdateOrganization(ctx, org)
}

func (o *OrganizationService) AddUsersToOrg(ctx context.Context, orgID uuid.UUID, users []*models.StaffInsertion) error {
	org, err := o.repo.GetOrganization(ctx, orgID)
	if err != nil {
		return err
	}
	for i := 0; i < len(users); i++ {
		var hasPosition bool
		for j := 0; j < len(org.Positions); j++ {
			if org.Positions[j].ID == users[i].PositionID {
				hasPosition = true
			}
		}
		if !hasPosition {
			return fmt.Errorf("no such position in this org")
		}
		var hasTeam bool
		for j := 0; j < len(org.Teams); j++ {
			if org.Teams[j].ID == users[i].TeamID {
				hasTeam = true
			}
		}
		if !hasTeam {
			return fmt.Errorf("no such team in this org")
		}
		staff := &models.Staff{
			ID:             users[i].ID,
			OrganizationID: orgID,
			PositionID:     users[i].PositionID,
			TeamID:         users[i].TeamID,
		}

		if err := o.repo.AddUsersToOrg(ctx, staff); err != nil {
			return err
		}
	}
	return nil
}

func (o *OrganizationService) DeleteOrganization(ctx context.Context, id uuid.UUID) error {
	return o.repo.DeleteOrganization(ctx, id)
}

func (o *OrganizationService) GetOrganizationEvents(ctx context.Context, id uuid.UUID) ([]*models.Event, error) {
	return o.repo.GetOrganizationEvents(ctx, id)
}

func (o *OrganizationService) GetOrganizationStaff(ctx context.Context, orgID uuid.UUID) ([]models.StaffInfo, error) {
	staff, err := o.repo.GetOrganizationStaff(ctx, orgID)
	if err != nil {
		return nil, err
	}

	var responseStaff = make([]models.StaffInfo, len(staff))
	for i := 0; i < len(staff); i++ {
		responseStaff[i] = models.StaffInfo{
			FirstName:       staff[i].FirstName,
			LastName:        staff[i].LastName,
			BackgroundColor: string(staff[i].BackgroundColor),
			TextColor:       string(staff[i].BackgroundColor),
			CompanyName:     staff[i].Organization.Name,
			Position:        staff[i].Position.Name,
		}
	}

	return responseStaff, err
}

func (o *OrganizationService) CreateOrganizationType(ctx context.Context, orgType *models.OrganizationType) error {
	return o.repo.CreateOrganizationType(ctx, orgType)
}

func (o *OrganizationService) GetOrganizationTypeByID(ctx context.Context, id uuid.UUID) (*models.OrganizationType, error) {
	return o.repo.GetOrganizationTypeByID(ctx, id)
}

func (o *OrganizationService) GetOrganizationTypes(ctx context.Context) ([]*models.OrganizationType, error) {
	return o.repo.GetOrganizationTypes(ctx)
}

func (o *OrganizationService) UpdateOrganizationType(ctx context.Context, orgType *models.OrganizationType) error {
	return o.repo.UpdateOrganizationType(ctx, orgType)
}

func (o *OrganizationService) DeleteOrganizationTypeByID(ctx context.Context, id uuid.UUID) error {
	return o.repo.DeleteOrganizationTypeByID(ctx, id)
}

func NewOrganizationService(ctx context.Context, repo postgres.Organization) *OrganizationService {
	return &OrganizationService{repo: repo, ctx: ctx}
}
