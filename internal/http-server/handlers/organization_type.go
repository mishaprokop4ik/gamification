package handlers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/miprokop/fication/internal/models"
	"net/http"
)

func (h *Handler) CreateOrganizationType(c *gin.Context) {
	ctx := context.Background()

	var orgType *models.OrganizationType

	if err := c.Bind(&orgType); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not get input model in creating org type: %s", err).Error())
		return
	}

	orgType.ID = uuid.New()

	err := h.Service.Organization.CreateOrganizationType(ctx, orgType)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not create model org type: %s", err).Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"created": orgType.ID,
	})
}

func (h *Handler) GetOrganizationTypes(c *gin.Context) {
	ctx := context.Background()
	organizationTypes, err := h.Service.Organization.GetOrganizationTypes(ctx)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"organization_types": organizationTypes,
	})
}

func (h *Handler) GetOrganizationTypeByID(c *gin.Context) {
	ctx := context.Background()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in getting org type by id: %s", err).Error())
		return
	}
	organization, err := h.Service.Organization.GetOrganizationTypeByID(ctx, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"organization_type": organization,
	})
}

func (h *Handler) DeleteOrganizationType(c *gin.Context) {
	ctx := context.Background()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in deleting org: %s", err).Error())
		return
	}
	err = h.Service.Organization.DeleteOrganizationTypeByID(ctx, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"deleted": true,
	})
}

func (h *Handler) UpdateOrganizationType(c *gin.Context) {
	ctx := context.Background()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in updating org: %s", err).Error())
		return
	}

	var orgType *models.OrganizationType

	if err := c.Bind(&orgType); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not get input update model in updating org: %s", err).Error())
		return
	}

	orgType.ID = id

	err = h.Service.Organization.UpdateOrganizationType(ctx, orgType)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("canupdate model in updating org: %s", err).Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"updated": true,
	})
}
