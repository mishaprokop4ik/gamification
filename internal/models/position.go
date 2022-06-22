package models

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var Developer = Position{
	CompanyID:    DefaultOrganization.ID,
	Organization: Organization{},
	Name:         "",
	Permissions:  nil,
}

var DefaultPosition = Position{
	CompanyID: DefaultOrganization.ID,
	Name:      "none",
	Permissions: []*Permission{
		{
			Permission: EventDelete,
		},
		{
			Permission: EventCreate,
		},
		{
			Permission: EventGetAll,
		},
		{
			Permission: EventGetByID,
		},
		{
			Permission: EventUpdate,
		},
	},
}

type Position struct {
	bun.BaseModel `bun:"table:position,alias:position"`
	ID            uuid.UUID     `json:"id" bun:",pk"`
	CompanyID     uuid.UUID     `json:"company_id"`
	Organization  Organization  `bun:"rel:belongs-to,join:company_id=id"`
	Name          string        `json:"name"`
	Permissions   []*Permission `bun:"rel:has-many,join:id=position_id"`
}

func (p *Position) HasPermission(perm PermissionName) bool {
	for _, permission := range p.Permissions {
		if permission.Permission == perm {
			return true
		}
	}
	return false
}
