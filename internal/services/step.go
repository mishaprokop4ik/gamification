package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/miprokop/fication/internal/models"
	"github.com/miprokop/fication/internal/persistence/postgres"
	log "github.com/sirupsen/logrus"
	"time"
)

type StepService struct {
	repo postgres.Step
	ctx  context.Context
}

func (s *StepService) GetStepPrizes(ctx context.Context, id uuid.UUID) ([]*models.Prize, error) {
	return s.repo.GetStepPrizes(ctx, id)
}

func (s *StepService) CreateStep(ctx context.Context, step *models.Step,
	creationTime, endTime time.Time) error {
	if creationTime.Round(10*time.Minute) != time.Now().Round(10*time.Minute) {
		s.createByTime(ctx, step, creationTime, endTime)
	} else {
		steps, err := s.repo.GetSteps(ctx, step.EventID)
		if err != nil {
			return err
		}
		step.Level = uint(len(steps) + 1)
		err = s.repo.CreateStep(ctx, step)
		if err != nil {
			return err
		}
	}

	s.updateByTime(ctx, step, creationTime, endTime)
	return nil
}

func (s *StepService) GetStep(ctx context.Context, id uuid.UUID) (*models.Step, error) {
	return s.repo.GetStep(ctx, id)
}

func (s *StepService) GetSteps(ctx context.Context, eventID uuid.UUID) ([]*models.Step, error) {
	return s.repo.GetSteps(ctx, eventID)
}

func (s *StepService) DeleteStep(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteStep(ctx, id)
}

func (s *StepService) AssignStaff(ctx context.Context, staffID, stepID uuid.UUID) error {
	staffStep := models.StepStaff{
		ID:             uuid.New(),
		StepID:         stepID,
		StaffID:        staffID,
		Accomplishment: models.Process,
		StartDate:      time.Now().Format(time.RFC3339),
	}
	return s.repo.AssignStaff(ctx, staffStep)
}

func (s *StepService) PassStaff(ctx context.Context, id uuid.UUID, status models.StepStatus) error {
	staffStep := models.StepStaff{
		ID:             id,
		Accomplishment: status,
	}
	return s.repo.PassStaff(ctx, staffStep)
}

func (s *StepService) UpdateStep(ctx context.Context, step *models.Step) error {
	var (
		endTime    time.Time
		createTime time.Time
	)

	var err error
	var toUpdate bool

	oldStep, err := s.repo.GetStep(ctx, step.ID)
	if err != nil {
		return err
	}
	createTime, err = time.Parse(time.RFC3339, oldStep.CreationDate)
	if err != nil {
		return err
	}

	if step.EndDate != "" {
		endTime, err = time.Parse(time.RFC3339, step.EndDate)
		if err != nil {
			return err
		}
		toUpdate = true
	}
	if toUpdate {
		if !endTime.After(createTime) {
			return fmt.Errorf("incorrent end time and creation time: %s, %s", endTime,
				createTime)
		}
		s.updateByTime(ctx, step, createTime, endTime)
	}

	step.Status = models.Changed
	err = s.repo.UpdateStep(ctx, step)

	return err
}

func (s *StepService) updateByTime(ctx context.Context, step *models.Step,
	creationTime, endTime time.Time) {
	time.AfterFunc(endTime.Sub(creationTime), func() {
		updateStep := step
		updateStep.Status = models.Finished
		err := s.repo.UpdateStep(ctx, updateStep)
		if err != nil {
			log.Println(err)
			return
		}
	})
}

func (s *StepService) createByTime(ctx context.Context, step *models.Step,
	creationTime, endTime time.Time) {
	time.AfterFunc(endTime.Sub(creationTime), func() {
		err := s.repo.CreateStep(ctx, step)
		if err != nil {
			log.Println(err)
			return
		}
	})
}

func NewStepService(ctx context.Context, repo postgres.Step) *StepService {
	return &StepService{repo: repo, ctx: ctx}
}
