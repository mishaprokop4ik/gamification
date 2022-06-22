package models

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"strings"
)

type StaffRole string

type Sex string

const (
	Male   Sex = "male"
	Female Sex = "female"
)

func (s *Sex) IsCorrect(input string) bool {
	return input == string(Male) || input == string(Female)
}

type HexColor string

func (h *HexColor) IsHex() bool {
	return len(*h) == 7 && strings.HasPrefix(string(*h), "#")
}

type Status string

type StaffEvents struct {
	bun.BaseModel `bun:"table:staff_events,alias:s_e"`

	ID        uuid.UUID `bun:",pk"`
	StaffID   uuid.UUID `json:"staff_id" bun:"user_id"`
	Staff     *Staff    `bun:"rel:belongs-to,join:user_id=id"`
	EventID   uuid.UUID `json:"event_id" bun:"event_id"`
	Event     *Event    `bun:"rel:belongs-to,join:event_id=id"`
	Status    Status    `json:"status"`
	StaffRole string    `json:"staff_role" bun:"user_role"`
}

type StaffLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type StaffImage struct {
	bun.BaseModel `bun:"table:staff_image,alias:staff_image"`

	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"userID"`
	Staff     Staff     `bun:"rel:belongs-to,join:user_id=id" json:"staff"`
	ImagePath string    `json:"image_path"`
}

type Staff struct {
	bun.BaseModel `bun:"table:staff,alias:staff"`

	ID              uuid.UUID      `json:"id" bun:",pk"`
	FirstName       string         `json:"first_name"`
	LastName        string         `json:"last_name"`
	Email           string         `json:"email"`
	Password        string         `json:"password"`
	Sex             Sex            `json:"sex"`
	AdditionalInfo  string         `json:"additional_info"`
	TeamID          uuid.UUID      `json:"team_id"`
	Team            *Team          `json:"team" bun:"rel:belongs-to,join:team_id=id"`
	PositionID      uuid.UUID      `json:"position_id"`
	Position        *Position      `json:"position" bun:"rel:belongs-to,join:position_id=id"`
	OrganizationID  uuid.UUID      `json:"organization_id" bun:"company_id"`
	Organization    *Organization  `json:"-" bun:"rel:belongs-to,join:company_id=id"`
	TextColor       HexColor       `json:"text_color"`
	BackgroundColor HexColor       `json:"background_color"`
	Events          []*StaffEvents `json:"events" bun:"m2m:staff_events,join:Staff=Event"`
	CurrentImage    string         `json:"current_image" bun:"-"`
	Images          []*StaffImage  `json:"images" bun:"rel:has-many,join:id=user_id"`
}

type StaffSignUp struct {
	bun.BaseModel `bun:"table:staff,alias:staff"`

	ID              uuid.UUID `pg:",pk"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	Email           string    `json:"email"`
	Password        string    `json:"password"`
	Sex             Sex       `json:"sex"`
	AdditionalInfo  string    `json:"additional_info"`
	TeamID          uuid.UUID `json:"team_id"`
	PositionID      uuid.UUID `json:"position_id"`
	OrganizationID  uuid.UUID `json:"organization_id" bun:"company_id"`
	TextColor       HexColor  `json:"text_color"`
	BackgroundColor HexColor  `json:"background_color"`
	CurrentImage    string    `json:"current_image" bun:"-"`
}

type PermissionName string

const (
	EventCreate  PermissionName = "event-create"
	EventUpdate  PermissionName = "event-update"
	EventDelete  PermissionName = "event-delete"
	EventGetByID PermissionName = "event-get-by-id"
	EventGetAll  PermissionName = "event-get-all"
)

type Permission struct {
	PositionID uuid.UUID      `json:"position_id"`
	Permission PermissionName `json:"permission"`
	GrantedBy  uuid.UUID      `json:"granted_by"`
	Position   *Position      `bun:"rel:belongs-to,join:position_id=id"`
}

func (u *Staff) HasPermission(perm PermissionName) bool {
	if u.Position == nil {
		return false
	}

	return u.Position.HasPermission(perm)
}

func (u *Staff) HasOneOfPermissions(perms ...PermissionName) bool {
	for _, perm := range perms {
		if u.HasPermission(perm) {
			return true
		}
	}

	return false
}

type StaffInfo struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	BackgroundColor string `json:"background_color"`
	TextColor       string `json:"text_color"`
	CompanyName     string `json:"company_name"`
	Position        string `json:"position"`
}

type InviteStatus string

const (
	Accepted   InviteStatus = "accepted"
	InProgress InviteStatus = "none"
	Declared   InviteStatus = "declared"
)

type EventStaffRole string

const (
	Admin   EventStaffRole = "admin"
	Default EventStaffRole = "default"
	Creator EventStaffRole = "creator"
)
