package models

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type StepStatus string

const (
	Finished StepStatus = "finished"
	Process  StepStatus = "process"
	Canceled StepStatus = "canceled"
	Changed  StepStatus = "changed"
)

type Accomplishment string

const (
	InProcess    Accomplishment = "process"
	ReadyToCheck Accomplishment = "ready-check"
	Done         Accomplishment = "done"
	Failed       Accomplishment = "failed"
	Cheated      Accomplishment = "cheated"
)

type StepStatusRequest struct {
	StaffID    uuid.UUID      `json:"staff_id"`
	StepStatus Accomplishment `json:"step_status"`
	Score      uint           `json:"score"`
}

type Step struct {
	bun.BaseModel `bun:"table:step,alias:step"`

	ID           uuid.UUID    `json:"id" bun:",pk"`
	EventID      uuid.UUID    `json:"event_id"`
	Name         string       `json:"name"`
	CreationDate string       `json:"creation_date"`
	EndDate      string       `json:"end_date"`
	Prizes       []*Prize     `json:"prizes" bun:"rel:has-many,join:id=step_id"`
	Task         string       `json:"task"`
	MaxScore     uint         `json:"max_score"`
	Level        uint         `json:"level"`
	Status       StepStatus   `json:"status" bun:"step_status"`
	Images       []*StepImage `json:"images" bun:"rel:has-many,join:id=step_id"`
	ActiveStaff  []*Staff     `json:"active_staff" bun:"m2m:staff_step,join:Step=Staff"`
	Description  string       `json:"description"`
}

type StepImage struct {
	bun.BaseModel `bun:"table:step_image,alias:st_img"`

	ID       uuid.UUID `json:"id" bun:",pk"`
	StepID   uuid.UUID `json:"step_id" bun:",pk"`
	Step     *Step     `json:"step" bun:"rel:belongs-to,join:step_id=id"`
	ImageURL string    `json:"image_url"`
}

type StepStaff struct {
	bun.BaseModel `bun:"table:staff_step,alias:staff_step"`

	ID             uuid.UUID      `json:"id" bun:",pk"`
	StepID         uuid.UUID      `json:"step_id"`
	Step           *Step          `json:"step" bun:"rel:belongs-to,join:step_id=id"`
	StaffID        uuid.UUID      `json:"staff_id"`
	Staff          *Staff         `json:"staff" bun:"rel:belongs-to,join:staff_id=id"`
	Accomplishment Accomplishment `json:"accomplishment"`
	Score          uint           `json:"score"`
	StartDate      string         `json:"start_date"`
}
