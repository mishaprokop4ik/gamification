package handlers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/miprokop/fication/internal/models"
	"net/http"
)

func (h *Handler) CreateTeam(c *gin.Context) {
	ctx := context.Background()

	var team *models.Team

	if err := c.Bind(&team); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not get input model in creating org: %s", err).Error())
		return
	}

	team.ID = uuid.New()

	err := h.Service.Team.CreateTeam(ctx, team)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not create model: %s", err).Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"created": team.ID,
	})
}

func (h *Handler) GetTeamsByOrganizationID(c *gin.Context) {
	ctx := context.Background()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in teams by org: %s", err).Error())
		return
	}

	teams, err := h.Service.Team.GetTeamsByOrganizationID(ctx, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not create model: %s", err).Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"teams": teams,
	})
}

func (h *Handler) GetTeamsByEventID(c *gin.Context) {
	ctx := context.Background()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in teams by org: %s", err).Error())
		return
	}

	teams, err := h.Service.Team.GetTeamsByEvent(ctx, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not create model: %s", err).Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"teams": teams,
	})
}

func (h *Handler) GetTeamByID(c *gin.Context) {
	ctx := context.Background()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in teams by org: %s", err).Error())
		return
	}

	team, err := h.Service.Team.GetTeamByID(ctx, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not create model: %s", err).Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"team": team,
	})
}

func (h *Handler) UpdateTeam(c *gin.Context) {
	ctx := context.Background()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in updating team: %s", err).Error())
		return
	}

	var team *models.Team

	if err := c.Bind(&team); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not get input update model in updating team: %s", err).Error())
		return
	}

	team.ID = id

	err = h.Service.Team.UpdateTeam(ctx, team)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not update model in updating team: %s", err).Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"updated": true,
	})
}

func (h *Handler) DeleteTeamByID(c *gin.Context) {
	ctx := context.Background()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in teams deletion: %s", err).Error())
		return
	}

	err = h.Service.Team.DeleteTeam(ctx, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not delete team: %s", err).Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"deleted": true,
	})
}
