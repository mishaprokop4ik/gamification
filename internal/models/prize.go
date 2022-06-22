package models

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type PrizeType string

const (
	Image      PrizeType = "image"
	Medal      PrizeType = "medal"
	Background PrizeType = "background"
	Text       PrizeType = "text"
)

type PrizeStatus string

const (
	Common    PrizeStatus = "common"
	Rare      PrizeStatus = "rare"
	Mith      PrizeStatus = "mith"
	Legendary PrizeStatus = "legendary"
)

type Prize struct {
	bun.BaseModel `bun:"table:prize,alias:prize"`

	ID           uuid.UUID   `json:"id"`
	StepID       uuid.UUID   `json:"step_id"`
	Step         *Step       `json:"step" bun:"rel:belongs-to,join:step_id=id"`
	Name         string      `json:"name"`
	CreationDate string      `json:"creation_date"`
	PrizeType    PrizeType   `json:"prize_type"`
	PrizeStatus  PrizeStatus `json:"prize_status"`
	CreatedBy    uuid.UUID   `json:"created_by"`
	Staff        *Staff      `json:"staff" bun:"rel:belongs-to,join:created_by=id"`
	Count        uint        `json:"count"`
	CurrentCount uint        `json:"current_count"`
	FileURL      string      `json:"file_url"`
	Description  string      `json:"description"`
}
