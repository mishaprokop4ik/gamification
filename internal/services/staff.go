package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/miprokop/fication/internal/models"
	"github.com/miprokop/fication/internal/persistence/postgres"
)

type StaffService struct {
	repo postgres.Staff
	ctx  context.Context
}

func (s *StaffService) RemovePermissionsFromPosition(ctx context.Context, permissions models.Permissions) error {
	for _, p := range permissions.Permissions {
		p := p
		return s.repo.RemovePermissionsFromPosition(ctx, &p)
	}
	return nil
}

func (s *StaffService) RemoveFromPosition(ctx context.Context, userID uuid.UUID) error {
	staff := models.Staff{ID: userID}
	return s.repo.RemoveFromPosition(ctx, &staff)
}

func (s *StaffService) GetDefaultPosition(ctx context.Context, orgID uuid.UUID) (models.Position, error) {
	return s.repo.GetDefaultPosition(ctx, orgID)
}

func (s *StaffService) UploadImage(ctx context.Context, image models.StaffImage) error {
	return s.repo.SaveFile(ctx, image)
}

func (s *StaffService) CreateStaffUser(ctx context.Context, staff *models.StaffSignUp) error {
	if !staff.Sex.IsCorrect(string(staff.Sex)) {
		return fmt.Errorf("incorrect sex input: %s; want: %s, %s", staff.Sex,
			models.Male, models.Female)
	}
	staff.Password = generatePasswordHash(staff.Password)
	_, err := s.repo.CreateStaffUser(ctx, staff)
	return err
}

func (s *StaffService) GetStaff(ctx context.Context, id uuid.UUID) (*models.Staff, error) {
	return s.repo.GetStaff(ctx, id)
}

func (s *StaffService) GetStaffByEvent(ctx context.Context, eventID uuid.UUID) ([]*models.Staff, error) {
	return s.repo.GetStaffByEvent(ctx, eventID)
}

func (s *StaffService) GetStaffByStep(ctx context.Context, stepID uuid.UUID) ([]*models.Staff, error) {
	return s.repo.GetStaffByStep(ctx, stepID)
}

func (s *StaffService) GetStaffByOrganization(ctx context.Context, organizationName uuid.UUID) ([]models.Staff, error) {
	//TODO implement me
	panic("implement me")
}

func (s *StaffService) DeleteStaff(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteStaff(ctx, id)
}

func (s *StaffService) UpdateStaff(ctx context.Context, staff *models.Staff) error {
	return s.repo.UpdateStaff(ctx, staff)
}

func (s *StaffService) SetStaffRole(ctx context.Context, role models.StaffRole) error {
	//TODO implement me
	panic("implement me")
}

func (s *StaffService) GetInvites(ctx context.Context, id uuid.UUID) ([]models.StaffEvents, error) {
	return s.repo.GetInvites(ctx, id)
}

func (s *StaffService) GetStaffPrizes(ctx context.Context, id uuid.UUID) ([]models.Prize, error) {
	return s.repo.GetStaffPrizes(ctx, id)
}

func (s *StaffService) GetPosition(ctx context.Context, id uuid.UUID) (*models.Position, error) {
	return s.repo.GetRole(ctx, id)
}

func (s *StaffService) GetAllPositions(ctx context.Context, orgID uuid.UUID) ([]models.Position, error) {
	return s.repo.GetAllPositions(ctx, orgID)
}

func (s *StaffService) CreatePosition(ctx context.Context, position *models.Position) error {
	return s.repo.CreatePosition(ctx, position)
}

func (s *StaffService) UpdatePosition(ctx context.Context, position *models.Position) error {
	return s.repo.UpdatePosition(ctx, position)
}

func (s *StaffService) DeletePosition(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeletePosition(ctx, id)
}

func (s *StaffService) AssignPosition(ctx context.Context, userID, positionID uuid.UUID) error {
	staff := models.Staff{
		ID:         userID,
		PositionID: positionID,
	}
	return s.repo.AssignPosition(ctx, &staff)
}

func (s *StaffService) GrantPermission(ctx context.Context, granterID, positionID uuid.UUID, perm models.Permission) error {
	//TODO implement me
	panic("implement me")
}

func (s *StaffService) RevokePermission(ctx context.Context, positionID uuid.UUID, perm models.Permission) error {
	//TODO implement me
	panic("implement me")
}

func NewStaffService(ctx context.Context, repo postgres.Staff) *StaffService {
	return &StaffService{repo: repo, ctx: ctx}
}
