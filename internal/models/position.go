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

const DefaultPositionName = "none"

var DefaultPosition = Position{
	CompanyID: DefaultOrganization.ID,
	Name:      DefaultPositionName,
	Permissions: []*Permission{
		// org types
		{
			Permission: OrganizationTypeCreate,
		},
		{
			Permission: OrganizationTypeUpdate,
		},
		{
			Permission: OrganizationTypeDelete,
		},
		{
			Permission: OrganizationTypeGetByID,
		},
		{
			Permission: OrganizationTypeGetAll,
		},
		// events
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
		// steps
		{
			Permission: StepDelete,
		},
		{
			Permission: StepCreate,
		},
		{
			Permission: StepGetAll,
		},
		{
			Permission: StepGetByID,
		},
		{
			Permission: StepUpdate,
		},
		// org
		{
			Permission: OrganizationDelete,
		},
		{
			Permission: OrganizationCreate,
		},
		{
			Permission: OrganizationGetAll,
		},
		{
			Permission: OrganizationGetByID,
		},
		{
			Permission: OrganizationUpdate,
		},
		{
			Permission: OrganizationEvents,
		},
		{
			Permission: OrganizationAddStaff,
		},
		// prize
		{
			Permission: PrizeDelete,
		},
		{
			Permission: PrizeCreate,
		},
		{
			Permission: PrizeGive,
		},
		{
			Permission: PrizeGetAll,
		},
		{
			Permission: PrizeGetByID,
		},
		{
			Permission: PrizeUpdate,
		},
		{
			Permission: StaffByOrganizationID,
		},
		// teams
		{
			Permission: TeamDelete,
		},
		{
			Permission: TeamCreate,
		},
		{
			Permission: TeamGetAll,
		},
		{
			Permission: TeamGetByID,
		},
		{
			Permission: TeamUpdate,
		},
		// pos
		{
			Permission: PositionDelete,
		},
		{
			Permission: PositionCreate,
		},
		{
			Permission: PositionGetAll,
		},
		{
			Permission: PositionGive,
		},
		{
			Permission: PositionUpdate,
		},
		{
			Permission: PositionGetByID,
		},
		// staff
		{
			Permission: StaffUpdate,
		},
		{
			Permission: StaffCreate,
		},
		{
			Permission: StaffGetAll,
		},
		{
			Permission: StaffGetByID,
		},
		{
			Permission: StaffSelfUpdate,
		},
		{
			Permission: StaffSelfDelete,
		},
		{
			Permission: StaffSelfGet,
		},
		{
			Permission: StaffDelete,
		},
		{
			Permission: PrizeStaffAll,
		},
	},
}

const (
	DeveloperName = "developer"
	HRName        = "hr"
	PM            = "project manager"
	DM            = "delivery manager"
	QA            = "quality assurance"
)

type Permissions struct {
	PositionID  uuid.UUID    `json:"position_id"`
	Permissions []Permission `json:"permissions"`
}

type Position struct {
	bun.BaseModel `bun:"table:position,alias:position"`
	ID            uuid.UUID     `json:"id" bun:",pk"`
	CompanyID     uuid.UUID     `json:"company_id"`
	Organization  Organization  `json:"-" bun:"rel:belongs-to,join:company_id=id"`
	Name          string        `json:"name"`
	Permissions   []*Permission `json:"permissions" bun:"rel:has-many,join:id=position_id"`
}

var DefaultProgrammingPositions = []Position{
	{
		Name: DeveloperName,
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
				Permission: StaffGetInvites,
			},
			{
				Permission: StaffGetSelfInvites,
			},
			{
				Permission: EventGetByID,
			},
			{
				Permission: EventUpdate,
			},
		},
	},
	{
		Name: HRName,
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
	},
	{
		Name: PM,
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
	},
	{
		Name: QA,
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
	},
	{
		Name: DM,
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
	},
}

func (p *Position) HasPermission(perm PermissionName) bool {
	for _, permission := range p.Permissions {
		if permission.Permission == perm {
			return true
		}
	}
	return false
}
