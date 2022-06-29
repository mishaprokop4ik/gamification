package models

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Event struct {
	bun.BaseModel `bun:"table:event,alias:event"`

	ID             uuid.UUID      `json:"id" bun:",pk"`
	Name           string         `json:"name"`
	CreationDate   string         `json:"creation_date"`
	EndDate        string         `json:"end_date"`
	Description    string         `json:"description"`
	ImagePath      string         `json:"image_path"`
	CreatedByID    uuid.UUID      `json:"created_by_id" bun:"created_by"`
	CreatedBy      *Staff         `json:"created_by" bun:"rel:belongs-to,join:created_by=id"`
	EventStatus    string         `json:"event_status"`
	EventType      string         `json:"event_type"`
	OrganizationID uuid.UUID      `json:"organization_id"`
	StaffEvents    []*StaffEvents `json:"staff" bun:"m2m:staff_events,join:Event=Staff"`
	Steps          []*Step        `json:"steps" bun:"rel:has-many,join:id=event_id"`
}

type StaffScore struct {
	Score int `json:"score"`
}
