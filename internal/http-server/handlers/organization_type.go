package handlers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/miprokop/fication/internal/models"
	"net/http"
)

// CreateOrganizationType
// @Summary Create organization type
// @Security ApiKeyAuth
// @Tags organizations-type
// @Description create models.OrganizationType
// @ID create-organization-type
// @Accept  json
// @Produce  json
// @Param input body organizationType true "organization type info"
// @Success 200 {string} uuid
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/org/type/ [post]
func (h *Handler) CreateOrganizationType(c *gin.Context) {
	ctx := context.Background()
	userID, ok := c.Get("userID")
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError,
			"there is no userID in context")
		return
	}
	_, ok = userID.(uuid.UUID)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError,
			"can not parse user id from context")
		return
	}

	staff, err := h.Service.Staff.GetStaff(ctx, userID.(uuid.UUID))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not get staff by id: %s", err).Error())
		return
	}

	if !staff.HasPermission(models.OrganizationTypeCreate) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}
	var orgType *models.OrganizationType

	if err := c.Bind(&orgType); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not get input model in creating org type: %s", err).Error())
		return
	}

	orgType.ID = uuid.New()

	err = h.Service.Organization.CreateOrganizationType(ctx, orgType)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not create model org type: %s", err).Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"created": orgType.ID,
	})
}

// GetOrganizationTypes
// @Summary Get All Organizations Types In Service
// @Security ApiKeyAuth
// @Tags organizations-type
// @Description Get All Organizations Types In Service
// @ID get-all-organizations-types
// @Accept  json
// @Produce  json
// @Success 200 {object} []organizationType
// @Failure 400,403 {} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/org/type/ [get]
func (h *Handler) GetOrganizationTypes(c *gin.Context) {
	ctx := context.Background()
	userID, ok := c.Get("userID")
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError,
			"there is no userID in context")
		return
	}
	_, ok = userID.(uuid.UUID)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError,
			"can not parse user id from context")
		return
	}

	staff, err := h.Service.Staff.GetStaff(ctx, userID.(uuid.UUID))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not get staff by id: %s", err).Error())
		return
	}

	if !staff.HasPermission(models.OrganizationTypeGetAll) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}
	organizationTypes, err := h.Service.Organization.GetOrganizationTypes(ctx)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"organization_types": organizationTypes,
	})
}

// GetOrganizationTypeByID
// @Summary Get Organization Type By ID
// @Security ApiKeyAuth
// @Tags organizations-type
// @Description Get Organization By ID
// @ID get-organization-type
// @Accept  json
// @Produce  json
// @Success 200 {object} organizationType
// @Failure 400,403 {} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/org/type/:id [get]
func (h *Handler) GetOrganizationTypeByID(c *gin.Context) {
	ctx := context.Background()
	userID, ok := c.Get("userID")
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError,
			"there is no userID in context")
		return
	}
	_, ok = userID.(uuid.UUID)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError,
			"can not parse user id from context")
		return
	}

	staff, err := h.Service.Staff.GetStaff(ctx, userID.(uuid.UUID))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not get staff by id: %s", err).Error())
		return
	}

	if !staff.HasPermission(models.OrganizationTypeGetByID) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}
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

// DeleteOrganizationType
// @Summary Delete Organization Type By ID
// @Security ApiKeyAuth
// @Tags organization-type
// @Description delete organization type by id
// @ID delete-organization
// @Accept  json
// @Produce  json
// @Success 200 {object} boolean
// @Failure 400,403 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/org/type/:id [delete]
func (h *Handler) DeleteOrganizationType(c *gin.Context) {
	ctx := context.Background()
	userID, ok := c.Get("userID")
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError,
			"there is no userID in context")
		return
	}
	_, ok = userID.(uuid.UUID)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError,
			"can not parse user id from context")
		return
	}

	staff, err := h.Service.Staff.GetStaff(ctx, userID.(uuid.UUID))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not get staff by id: %s", err).Error())
		return
	}

	if !staff.HasPermission(models.OrganizationTypeDelete) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}
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

// UpdateOrganizationType
// @Summary Update organization type by id
// @Security ApiKeyAuth
// @Tags organizations
// @Description Update organization type by id
// @ID update-event
// @Accept  json
// @Produce  json
// @Param input body organizationType true "organization type"
// @Success 200 {object} boolean
// @Failure 400,403 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/org/type/:id [put]
func (h *Handler) UpdateOrganizationType(c *gin.Context) {
	ctx := context.Background()
	userID, ok := c.Get("userID")
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError,
			"there is no userID in context")
		return
	}
	_, ok = userID.(uuid.UUID)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError,
			"can not parse user id from context")
		return
	}

	staff, err := h.Service.Staff.GetStaff(ctx, userID.(uuid.UUID))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not get staff by id: %s", err).Error())
		return
	}

	if !staff.HasPermission(models.OrganizationTypeUpdate) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in updating orgtype: %s", err).Error())
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
