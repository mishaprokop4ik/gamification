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

type Step struct {
	bun.BaseModel `bun:"table:step,alias:step"`

	ID           uuid.UUID  `json:"id"`
	EventID      uuid.UUID  `json:"event_id"`
	Name         string     `json:"name"`
	CreationDate string     `json:"creation_date"`
	EndDate      string     `json:"end_date"`
	Task         string     `json:"task"`
	Status       StepStatus `json:"status"`
	Description  string     `json:"description"`
}
