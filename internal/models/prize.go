package models

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type PrizeType string

const (
	Image      PrizeType = "image"
	Medal      PrizeType = "medal"
	Background PrizeType = "background"
	Text       PrizeType = "text"
)

func NewPrizeType(p string) (PrizeType, error) {
	if p != string(Image) && p != string(Medal) &&
		p != string(Background) && p != string(Text) {
		return "", fmt.Errorf("can not create type, incorrent prize type name: %s; want: %s, %s, %s, %s",
			p, Image, Medal, Background, Text)
	}
	return PrizeType(p), nil
}

type PrizeStatus string

func OneOf(prizeStatus PrizeStatus) bool {
	return prizeStatus == Common || prizeStatus == Rare ||
		prizeStatus == Mith || prizeStatus == Legendary
}

const (
	Common    PrizeStatus = "common"
	Rare      PrizeStatus = "rare"
	Mith      PrizeStatus = "mith"
	Legendary PrizeStatus = "legendary"
)

type Prize struct {
	bun.BaseModel `bun:"table:prize,alias:prize"`

	ID           uuid.UUID     `json:"id" bun:",pk"`
	StepID       uuid.UUID     `json:"step_id"`
	Step         *Step         `json:"step" bun:"rel:belongs-to,join:step_id=id"`
	Name         string        `json:"name"`
	CreationDate string        `json:"creation_date"`
	PrizeType    PrizeType     `json:"type"`
	PrizeStatus  PrizeStatus   `json:"status"`
	CreatedBy    uuid.UUID     `json:"created_by"`
	Staff        *Staff        `json:"staff" bun:"rel:belongs-to,join:created_by=id"`
	Count        uint          `json:"count"`
	CurrentCount uint          `json:"current_count" bun:"current_count"`
	Data         string        `json:"data"`
	Description  string        `json:"description"`
	Prizes       []*StaffPrize `json:"prizes" bun:"m2m:staff_prizes,join:Staff=Prize"`
}

type PrizeRepo struct {
	bun.BaseModel `bun:"table:prize,alias:prize"`

	ID           uuid.UUID   `json:"id" bun:",pk"`
	StepID       uuid.UUID   `json:"step_id"`
	Step         *Step       `json:"step" bun:"rel:belongs-to,join:step_id=id"`
	Name         string      `json:"name"`
	CreationDate time.Time   `json:"creation_date"`
	PrizeType    PrizeType   `json:"type"`
	PrizeStatus  PrizeStatus `json:"status"`
	CreatedBy    uuid.UUID   `json:"created_by"`
	Staff        *Staff      `json:"staff" bun:"rel:belongs-to,join:created_by=id"`
	Count        uint        `json:"count"`
	CurrentCount uint        `json:"current_count"`
	Data         string      `json:"data"`
	Description  string      `json:"description"`
}

type StaffPrize struct {
	bun.BaseModel `bun:"table:staff_prizes,alias:staff_prizes"`

	ID      uuid.UUID `json:"id" bun:",pk"`
	StaffID uuid.UUID `json:"staff_id"`
	Staff   *Staff    `bun:"rel:belongs-to,join:staff_id=id"`
	PrizeID uuid.UUID `json:"prize_id"`
	Prize   *Prize    `bun:"rel:belongs-to,join:prize_id=id"`
}
