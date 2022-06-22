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
	OrganizationID uuid.UUID      `json:"organization_id"`
	StaffEvents    []*StaffEvents `json:"staff" bun:"m2m:staff_events,join:Event=Staff"`
	Steps          []*Step        `json:"steps" bun:"rel:has-many,join:id=event_id"`
}
