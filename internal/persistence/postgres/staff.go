package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/miprokop/fication/internal/models"
	"github.com/uptrace/bun"
)

type StaffRepo struct {
	DB  *bun.DB
	ctx context.Context
}

func (s *StaffRepo) SaveFile(ctx context.Context, image models.StaffImage) error {
	_, err := s.DB.NewInsert().Model(image).Exec(ctx)
	return err
}

func (s *StaffRepo) CreateStaffUser(ctx context.Context, staff *models.StaffSignUp) (uuid.UUID, error) {
	tx, err := s.DB.DB.Begin()
	if err != nil {
		return [16]byte{}, err
	}
	_, err = s.DB.NewInsert().Model(staff).Exec(ctx)
	if err != nil {
		tx.Rollback()
		return [16]byte{}, fmt.Errorf("can not create staff; err: %s", err)
	}
	if staff.CurrentImage != "" {
		var image = &models.StaffImage{
			UserID:    staff.ID,
			ImagePath: staff.CurrentImage,
		}
		_, err = s.DB.NewInsert().Model(image).Exec(ctx)
		if err != nil {
			tx.Rollback()
			return uuid.UUID{}, fmt.Errorf(`can not insert image in staff creation; err: %s`, err)
		}
	}
	return staff.ID, tx.Commit()
}

func (s *StaffRepo) GetStaffAuth(ctx context.Context, email, password string) (*models.Staff, error) {
	var staff = new(models.Staff)

	err := s.DB.NewSelect().Model(staff).Where("email = ?", email).Relation("Position").Where("password = ?", password).Scan(ctx)
	return staff, err
}

func (s *StaffRepo) GetStaff(ctx context.Context, id uuid.UUID) (*models.Staff, error) {
	var staff = new(models.Staff)

	err := s.DB.NewSelect().Model(staff).
		Relation("Position").
		Relation("Team").
		Relation("Organization").
		Relation("Images").
		Where("staff.id = ?", id).
		Scan(ctx)

	return staff, err
}

func (s *StaffRepo) GetStaffByEvent(ctx context.Context, eventID uuid.UUID) ([]*models.Staff, error) {
	var staff = make([]*models.Staff, 0)
	err := s.DB.NewSelect().
		Model(&staff).
		Join("JOIN staff_events ON staff.id = staff_events.user_id").
		Where("staff_events.event_id = ?", eventID).
		Scan(ctx)

	return staff, err
}

func (s *StaffRepo) GetStaffByStep(ctx context.Context, stepID uuid.UUID) ([]*models.Staff, error) {
	var staff = make([]*models.Staff, 0)
	err := s.DB.NewSelect().
		Model(&staff).
		Join("JOIN staff_step ON staff.id = staff_step.user_id").
		Where("staff_step.step_id = ?", stepID).
		Scan(ctx)

	return staff, err
}

func (s *StaffRepo) GetStaffByOrganization(ctx context.Context, organizationName uuid.UUID) ([]models.Staff, error) {
	//TODO implement me
	panic("implement me")
}

func (s *StaffRepo) DeleteStaff(ctx context.Context, id uuid.UUID) error {
	_, err := s.DB.NewDelete().Model(&models.Staff{}).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *StaffRepo) UpdateStaff(ctx context.Context, staff *models.Staff) error {
	_, err := s.DB.NewUpdate().OmitZero().Model(staff).Where("id = ?", staff.ID).Exec(ctx)
	return err
}

func (s *StaffRepo) SetStaffRole(ctx context.Context, role models.StaffRole) error {
	//TODO implement me
	panic("implement me")
}

func (s *StaffRepo) GetInvites(ctx context.Context, id uuid.UUID) ([]models.StaffEvents, error) {
	//TODO implement me
	panic("implement me")
}

func (s *StaffRepo) GetStaffPrizes(ctx context.Context, id uuid.UUID) ([]models.Prize, error) {
	//TODO implement me
	panic("implement me")
}

func (s *StaffRepo) GetRole(ctx context.Context, id uuid.UUID) (*models.Position, error) {
	//TODO implement me
	panic("implement me")
}

func (s *StaffRepo) GetAllPositions(ctx context.Context, orgID uuid.UUID) ([]models.Position, error) {
	//TODO implement me
	panic("implement me")
}

func (s *StaffRepo) CreatePosition(ctx context.Context, position *models.Position) error {
	//TODO implement me
	panic("implement me")
}

func (s *StaffRepo) UpdatePosition(ctx context.Context, position *models.Position) error {
	//TODO implement me
	panic("implement me")
}

func (s *StaffRepo) DeletePosition(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (s *StaffRepo) AssignPosition(ctx context.Context, userID, positionID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (s *StaffRepo) GrantPermission(ctx context.Context, granterID, positionID uuid.UUID, perm models.Permission) error {
	//TODO implement me
	panic("implement me")
}

func (s *StaffRepo) RevokePermission(ctx context.Context, roleID uuid.UUID, perm models.Permission) error {
	//TODO implement me
	panic("implement me")
}

func NewStaffRepo(ctx context.Context, DB *bun.DB) *StaffRepo {
	return &StaffRepo{DB: DB, ctx: ctx}
}
