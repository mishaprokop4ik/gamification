package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/miprokop/fication/internal/models"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

type stepShortResponse struct {
	ID           uuid.UUID         `json:"id" bun:",pk"`
	EventID      uuid.UUID         `json:"event_id"`
	Name         string            `json:"name"`
	CreationDate string            `json:"creation_date"`
	EndDate      string            `json:"end_date"`
	Task         string            `json:"task"`
	MaxScore     uint              `json:"max_score"`
	Level        uint              `json:"level"`
	Status       models.StepStatus `json:"status" bun:"step_status"`
	Description  string            `json:"description"`
}

type eventAllData struct {
	ID             uuid.UUID            `json:"id" bun:",pk"`
	Name           string               `json:"name"`
	CreationDate   string               `json:"creation_date"`
	EndDate        string               `json:"end_date"`
	Description    string               `json:"description"`
	ImagePath      string               `json:"image_path"`
	CreatedByID    uuid.UUID            `json:"created_by_id" bun:"created_by"`
	EventStatus    string               `json:"event_status"`
	EventType      string               `json:"event_type"`
	OrganizationID uuid.UUID            `json:"organization_id"`
	StaffEvents    []*staffEvents       `json:"staff"`
	Steps          []*stepShortResponse `json:"steps"`
}

type position struct {
	ID        uuid.UUID `json:"id" bun:",pk"`
	CompanyID uuid.UUID `json:"company_id"`
	Name      string    `json:"name"`
}

type permission struct {
	bun.BaseModel `bun:"table:permissions,alias:permissions"`
	PositionID    uuid.UUID             `json:"position_id"`
	Permission    models.PermissionName `json:"permission"`
	GrantedBy     uuid.UUID             `json:"granted_by"`
}

type positionCreate struct {
	ID          uuid.UUID    `json:"id" bun:",pk"`
	CompanyID   uuid.UUID    `json:"company_id"`
	Name        string       `json:"name"`
	Permissions []permission `json:"permissions"`
}

type organizationType struct {
	bun.BaseModel `bun:"table:org_type,alias:org_type"`
	ID            uuid.UUID `json:"id" bun:",pk"`
	Name          string    `json:"name"`
}

type team struct {
	bun.BaseModel  `bun:"table:team,alias:team"`
	ID             uuid.UUID `json:"id" bun:",pk"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	OrganizationID uuid.UUID `json:"organization_id"`
}

type organizationUpdateResponse struct {
	ID         uuid.UUID `json:"id" bun:",pk"`
	Name       string    `json:"name"`
	WebsiteURL string    `json:"website_url"`
	Image      string    `json:"image"`
}

type users struct {
	Staff []*models.StaffInsertion `json:"users"`
}

type insertInfo struct {
	Inserted      bool `json:"inserted"`
	InsertedCount int  `json:"inserted_count"`
}

type organizationResponse struct {
	ID         uuid.UUID          `json:"id" bun:",pk"`
	Name       string             `json:"name"`
	WebsiteURL string             `json:"website_url"`
	Image      string             `json:"image"`
	Types      []organizationType `json:"types"`
	Positions  []position         `json:"positions"`
	Teams      []team             `json:"teams"`
}

type stepImageRequest struct {
	ID       uuid.UUID `json:"id" bun:",pk"`
	StepID   uuid.UUID `json:"step_id" bun:",pk"`
	ImageURL string    `json:"image_url"`
}

type eventShortData struct {
	ID             uuid.UUID `json:"id" bun:",pk"`
	Name           string    `json:"name"`
	CreationDate   string    `json:"creation_date"`
	EndDate        string    `json:"end_date"`
	Description    string    `json:"description"`
	ImagePath      string    `json:"image_path"`
	CreatedByID    uuid.UUID `json:"created_by_id" bun:"created_by"`
	EventStatus    string    `json:"event_status"`
	EventType      string    `json:"event_type"`
	OrganizationID uuid.UUID `json:"organization_id"`
}

type staffEvents struct {
	ID        uuid.UUID             `bun:",pk" json:"id"`
	StaffID   uuid.UUID             `json:"staff_id" bun:"user_id"`
	EventID   uuid.UUID             `json:"event_id" bun:"event_id"`
	Status    models.InviteStatus   `json:"status"`
	StaffRole models.EventStaffRole `json:"staff_role" bun:"user_role"`
}

type eventsResponse struct {
	Data []eventRequest `json:"events"`
}

type prizeRequest struct {
	ID           uuid.UUID          `json:"id" bun:",pk"`
	StepID       uuid.UUID          `json:"step_id"`
	Name         string             `json:"name"`
	CreationDate string             `json:"creation_date"`
	PrizeType    models.PrizeType   `json:"type"`
	PrizeStatus  models.PrizeStatus `json:"status"`
	CreatedBy    uuid.UUID          `json:"created_by"`
	Count        uint               `json:"count"`
	CurrentCount uint               `json:"current_count" bun:"current_count"`
	Data         string             `json:"data"`
	Description  string             `json:"description"`
}

type stepRequest struct {
	ID           uuid.UUID           `json:"id" bun:",pk"`
	EventID      uuid.UUID           `json:"event_id"`
	Name         string              `json:"name"`
	CreationDate string              `json:"creation_date"`
	EndDate      string              `json:"end_date"`
	Prizes       []*prizeRequest     `json:"prizes" bun:"rel:has-many,join:id=step_id"`
	Task         string              `json:"task"`
	MaxScore     uint                `json:"max_score"`
	Level        uint                `json:"level"`
	Status       models.StepStatus   `json:"status" bun:"step_status"`
	Images       []*stepImageRequest `json:"images" bun:"rel:has-many,join:id=step_id"`
	Description  string              `json:"description"`
}

type permissions struct {
	PositionID  uuid.UUID    `json:"position_id"`
	Permissions []permission `json:"permissions"`
}

type updatePosition struct {
	ID          uuid.UUID    `json:"id" bun:",pk"`
	CompanyID   uuid.UUID    `json:"company_id"`
	Name        string       `json:"name"`
	Permissions []permission `json:"permissions" bun:"rel:has-many,join:id=position_id"`
}

type eventRequestUpdate struct {
	ID             uuid.UUID `json:"id" bun:",pk"`
	Name           string    `json:"name"`
	CreationDate   string    `json:"creation_date"`
	EndDate        string    `json:"end_date"`
	Description    string    `json:"description"`
	ImagePath      string    `json:"image_path"`
	CreatedByID    uuid.UUID `json:"created_by_id" bun:"created_by"`
	EventStatus    string    `json:"event_status"`
	EventType      string    `json:"event_type"`
	OrganizationID uuid.UUID `json:"organization_id"`
}

type eventRequest struct {
	ID             uuid.UUID      `json:"id" bun:",pk"`
	Name           string         `json:"name"`
	CreationDate   string         `json:"creation_date"`
	EndDate        string         `json:"end_date"`
	Description    string         `json:"description"`
	ImagePath      string         `json:"image_path"`
	CreatedByID    uuid.UUID      `json:"created_by_id" bun:"created_by"`
	EventStatus    string         `json:"event_status"`
	EventType      string         `json:"event_type"`
	OrganizationID uuid.UUID      `json:"organization_id"`
	StaffEvents    []*staffEvents `json:"staff" bun:"m2m:staff_events,join:Event=Staff"`
	Steps          []*stepRequest `json:"steps" bun:"rel:has-many,join:id=event_id"`
}

type errorResponse struct {
	Message string `json:"message"`
}

type statusResponse struct {
	Status string `json:"status"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
