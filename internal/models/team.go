package models

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

const DefaultTeamName = "none"

var DefaultTeam = Team{
	Name:        DefaultTeamName,
	Description: "",
}

type Team struct {
	bun.BaseModel  `bun:"table:team,alias:team"`
	ID             uuid.UUID    `json:"id" bun:",pk"`
	Name           string       `json:"name"`
	Description    string       `json:"description"`
	OrganizationID uuid.UUID    `json:"organization_id"`
	Organization   Organization `json:"organization" bun:"rel:belongs-to,join:organization_id=id"`
	Staff          []*Staff     `json:"staff" bun:"rel:has-many,join:id=team_id"`
}
