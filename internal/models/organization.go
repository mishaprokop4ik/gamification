package models

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var DefaultOrganization = Organization{
	Name:       "default",
	WebsiteURL: "https://nure.ua/",
}

type Organization struct {
	bun.BaseModel `bun:"table:organizations,alias:organizations"`

	ID         uuid.UUID          `json:"id" bun:",pk"`
	Name       string             `json:"name"`
	WebsiteURL string             `json:"website_url"`
	Image      string             `json:"image"`
	Types      []OrganizationType `json:"types" bun:"m2m:organizations_types,join:Organization=OrganizationType"`
	Positions  []Position         `json:"positions" bun:"rel:has-many,join:id=company_id"`
	Teams      []Team             `json:"teams" bun:"rel:has-many,join:id=organization_id"`
}

type OrganizationsTypes struct {
	bun.BaseModel    `bun:"table:organizations_types,alias:organizations_types"`
	ID               uuid.UUID        `json:"id" bun:",pk"`
	OrganizationID   uuid.UUID        `json:"organization_id" bun:"org_id"`
	Organization     *Organization    `bun:"rel:belongs-to,join:org_id=id"`
	TypeID           uuid.UUID        `bun:",pk"`
	OrganizationType OrganizationType `json:"type" bun:"rel:belongs-to,join:type_id=id"`
}

type OrganizationType struct {
	bun.BaseModel `bun:"table:org_type,alias:org_type"`
	ID            uuid.UUID      `json:"id" bun:",pk"`
	Name          string         `json:"name"`
	Organization  []Organization `bun:"m2m:organizations_types,join:OrganizationType=Organization"`
}

var DefaultOrganizationType = OrganizationType{
	Name: "none",
}
