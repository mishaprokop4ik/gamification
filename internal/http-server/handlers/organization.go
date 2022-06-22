package handlers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/miprokop/fication/internal/models"
	"net/http"
	"net/url"
)

func (h *Handler) GetAllOrganizations(c *gin.Context) {
	ctx := context.Background()
	organizations, err := h.Service.Organization.GetOrganizations(ctx)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"organizations": organizations,
	})
}

func (h *Handler) GetStaffByOrganizationID(c *gin.Context) {
	ctx := context.Background()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in getting org: %s", err).Error())
		return
	}
	staff, err := h.Service.Organization.GetOrganizationStaff(ctx, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"staff": staff,
	})
}

func (h *Handler) GetOrganization(c *gin.Context) {
	ctx := context.Background()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in getting org: %s", err).Error())
		return
	}
	organization, err := h.Service.Organization.GetOrganization(ctx, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"organization": organization,
	})
}

func (h *Handler) DeleteOrganization(c *gin.Context) {
	ctx := context.Background()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in deleting org: %s", err).Error())
		return
	}
	err = h.Service.Organization.DeleteOrganization(ctx, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"deleted": true,
	})
}

func (h *Handler) CreateOrganization(c *gin.Context) {
	ctx := context.Background()

	var org *models.Organization

	if err := c.Bind(&org); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not get input model in creating org: %s", err).Error())
		return
	}
	_, err := url.ParseRequestURI(org.WebsiteURL)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("can not url: %s", org.WebsiteURL))
		return
	}
	org.ID = uuid.New()

	err = h.Service.Organization.CreateOrganization(ctx, org)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not create model: %s", err).Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"created": org.ID,
	})
}

func (h *Handler) UpdateOrganization(c *gin.Context) {
	ctx := context.Background()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in updating org: %s", err).Error())
		return
	}

	var org *models.Organization

	if err := c.Bind(&org); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not get input update model in updating org: %s", err).Error())
		return
	}
	if org.WebsiteURL != "" {
		_, err = url.ParseRequestURI(org.WebsiteURL)
		if err != nil {
			newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("can not parse url: %s", org.WebsiteURL))
			return
		}
	}

	if org.Image != "" {
		_, err = url.ParseRequestURI(org.Image)
		if err != nil {
			newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("can not parse url: %s", org.WebsiteURL))
			return
		}
	}
	org.ID = id

	err = h.Service.Organization.UpdateOrganization(ctx, org)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("canupdate model in updating org: %s", err).Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"updated": true,
	})
}

func (h *Handler) GetOrganizationEvents(c *gin.Context) {
	ctx := context.Background()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in getting org events: %s", err).Error())
		return
	}

	events, err := h.Service.Organization.GetOrganizationEvents(ctx, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not modelin getting org events: %s", err).Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"events": events,
	})
}

func (h *Handler) AddStaffToOrganization(c *gin.Context) {
	ctx := context.Background()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in adding staff into org: %s", err).Error())
		return
	}

	type ids struct {
		IDs []uuid.UUID `json:"ids"`
	}

	var staffIDs ids

	if err := c.Bind(&staffIDs); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not bind ids in adding staff into org: %s", err).Error())
		return
	}

	err = h.Service.Organization.AddUsersToOrg(ctx, id, staffIDs.IDs)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not bind staff into org: %s", err).Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"inserted":       true,
		"inserted_count": len(staffIDs.IDs),
	})
}
