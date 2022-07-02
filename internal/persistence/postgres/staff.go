package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/miprokop/fication/internal/models"
	"github.com/uptrace/bun"
)

type StaffRepo struct {
	DB  *bun.DB
	ctx context.Context
}

func (s *StaffRepo) RemovePermissionsFromPosition(ctx context.Context, permission *models.Permission) error {
	_, err := s.DB.NewDelete().Model(permission).
		Where("position_id = ?", permission.PositionID).
		Where("permission = ?", permission.Permission).Exec(ctx)
	return err
}

func (s *StaffRepo) RemoveFromPosition(ctx context.Context, staff *models.Staff) error {
	_, err := s.DB.NewUpdate().OmitZero().Model(staff).Where("id = ?", staff.ID).Set("position_id = DEFAULT").Exec(ctx)
	return err
}

func (s *StaffRepo) GetDefaultPosition(ctx context.Context, orgID uuid.UUID) (models.Position, error) {
	p := new(models.Position)
	err := s.DB.NewSelect().Model(p).Where("company_id = ?", orgID).Where("name = ?", models.DefaultPositionName).Scan(ctx)
	return *p, err
}

func (s *StaffRepo) SaveFile(ctx context.Context, image models.StaffImage) error {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	_, err = tx.NewInsert().Model(&image).Exec(ctx)
	if err != nil {
		return err
	}
	_, err = tx.NewUpdate().Model(&models.Staff{}).OmitZero().Set("current_image = ?", image.ImagePath).Where("id = ?", image.UserID).Exec(ctx)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (s *StaffRepo) CreateStaffUser(ctx context.Context, staff *models.StaffSignUp) (uuid.UUID, error) {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return uuid.UUID{}, err
	}
	_, err = tx.NewInsert().Model(staff).Exec(ctx)
	if err != nil {
		tx.Rollback()
		return uuid.UUID{}, fmt.Errorf("can not create staff; err: %s", err)
	}
	if staff.CurrentImage != "" {
		var image = &models.StaffImage{
			UserID:    staff.ID,
			ImagePath: staff.CurrentImage,
		}
		_, err = tx.NewInsert().Model(image).Exec(ctx)
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
	var permissions = new([]*models.Permission)
	err := s.DB.NewSelect().Model(staff).
		Relation("Position").
		Relation("Team").
		Relation("Organization").
		Relation("Images").
		Where("staff.id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	err = s.DB.NewSelect().Model(permissions).
		Where("permissions.position_id = ?", staff.PositionID).
		Scan(ctx)
	staff.Position.Permissions = *permissions
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
		Distinct().
		Join("JOIN staff_step ON staff.id = staff_step.staff_id").
		Where("staff_step.step_id = ?", stepID).
		Scan(ctx)

	return staff, err
}

func (s *StaffRepo) DeleteStaff(ctx context.Context, id uuid.UUID) error {
	_, err := s.DB.NewDelete().Model(&models.Staff{}).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *StaffRepo) UpdateStaff(ctx context.Context, staff *models.Staff) error {
	_, err := s.DB.NewUpdate().OmitZero().Model(staff).WherePK().Exec(ctx)
	return err
}

func (s *StaffRepo) SetStaffRole(ctx context.Context, role models.StaffRole) error {
	//TODO implement me
	panic("implement me")
}

func (s *StaffRepo) GetInvites(ctx context.Context, id uuid.UUID) ([]models.StaffEvents, error) {
	invites := new([]models.StaffEvents)
	err := s.DB.NewSelect().Model(invites).Where("user_id = ?", id).Scan(ctx)
	return *invites, err
}

func (s *StaffRepo) GetStaffPrizes(ctx context.Context, id uuid.UUID) ([]models.Prize, error) {
	var prizes = new([]models.Prize)
	var staffPrizes = new([]models.StaffPrize)
	err := s.DB.NewSelect().Model(staffPrizes).Where("staff_id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	ids := make([]uuid.UUID, len(*staffPrizes))
	for i := range *staffPrizes {
		ids[i] = (*staffPrizes)[i].PrizeID
	}
	err = s.DB.NewSelect().Model(prizes).Where("id IN (?)", bun.In(ids)).Scan(ctx)
	return *prizes, err
}

func (s *StaffRepo) GetRole(ctx context.Context, id uuid.UUID) (*models.Position, error) {
	position := new(models.Position)
	err := s.DB.NewSelect().Model(position).Relation("Permissions").Where("id = ?", id).Scan(ctx)
	return position, err
}

func (s *StaffRepo) GetAllPositions(ctx context.Context, orgID uuid.UUID) ([]models.Position, error) {
	var positions = new([]models.Position)
	err := s.DB.NewSelect().Model(positions).Relation("Permissions").Where("company_id = ?", orgID).Scan(ctx)
	return *positions, err
}

func (s *StaffRepo) CreatePosition(ctx context.Context, position *models.Position) error {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	_, err = tx.NewInsert().Model(position).Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}
	if len(position.Permissions) != 0 {
		_, err = tx.NewInsert().Model(&position.Permissions).Exec(ctx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (s *StaffRepo) UpdatePosition(ctx context.Context, position *models.Position) error {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	if len(position.Permissions) != 0 {
		old := new(models.Position)
		err = tx.NewSelect().Model(old).Relation("Permissions").Where("id = ?", position.ID).Scan(ctx)
		if err != nil {
			return err
		}
		for i := range position.Permissions {
			var exist bool
			for j := range old.Permissions {
				if position.Permissions[i].Permission == old.Permissions[j].Permission {
					exist = true
				}
			}
			if !exist {
				_, err = tx.NewInsert().Model(position.Permissions[i]).Exec(ctx)
				if err != nil {
					tx.Rollback()
					return err
				}
			}
		}
	}
	_, err = tx.NewUpdate().OmitZero().Model(position).Where("id = ?", position.ID).Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}
	return err
}

func (s *StaffRepo) DeletePosition(ctx context.Context, id uuid.UUID) error {
	_, err := s.DB.NewDelete().Model(&models.Position{}).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *StaffRepo) AssignPosition(ctx context.Context, staff *models.Staff) error {
	_, err := s.DB.NewUpdate().Model(staff).OmitZero().Where("id = ?", staff.ID).Exec(ctx)
	return err
}

func NewStaffRepo(ctx context.Context, DB *bun.DB) *StaffRepo {
	return &StaffRepo{DB: DB, ctx: ctx}
}
