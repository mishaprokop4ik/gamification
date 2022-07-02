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

// GetAllOrganizations
// @Summary Get All Organizations In Service
// @Security ApiKeyAuth
// @Tags organizations
// @Description Get All Organizations In Service
// @ID get-all-organizations
// @Accept  json
// @Produce  json
// @Success 200 {object} []organizationResponse
// @Failure 400,403 {} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/org/ [get]
func (h *Handler) GetAllOrganizations(c *gin.Context) {
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

	if !staff.HasPermission(models.OrganizationGetAll) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}

	organizations, err := h.Service.Organization.GetOrganizations(ctx)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"organizations": organizations,
	})
}

// GetStaffByOrganizationID
// @Summary Get All Staff In Organization
// @Security ApiKeyAuth
// @Tags organizations
// @Description Get All Staff In Organization By Organization ID
// @ID get-organization-staff
// @Accept  json
// @Produce  json
// @Success 200 {object} []models.StaffInfo
// @Failure 400,403 {} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/org/staff/:id [get]
func (h *Handler) GetStaffByOrganizationID(c *gin.Context) {
	ctx := context.Background()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in getting org: %s", err).Error())
		return
	}

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

	if !staff.HasPermission(models.StaffByOrganizationID) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}

	staffs, err := h.Service.Organization.GetOrganizationStaff(ctx, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"staff": staffs,
	})
}

// GetOrganization
// @Summary Get Organization By ID
// @Security ApiKeyAuth
// @Tags organizations
// @Description Get Organization By ID
// @ID get-organization
// @Accept  json
// @Produce  json
// @Success 200 {object} organizationResponse
// @Failure 400,403 {} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/org/:id [get]
func (h *Handler) GetOrganization(c *gin.Context) {
	ctx := context.Background()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in getting org: %s", err).Error())
		return
	}

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

	if !staff.HasPermission(models.OrganizationGetByID) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
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

// DeleteOrganization
// @Summary Delete Organization By ID
// @Security ApiKeyAuth
// @Tags organizations
// @Description delete organization by id
// @Description also delete all chained staff, positions, steps, teams
// @ID delete-organization
// @Accept  json
// @Produce  json
// @Success 200 {object} boolean
// @Failure 400,403 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/org/:id [delete]
func (h *Handler) DeleteOrganization(c *gin.Context) {
	ctx := context.Background()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in deleting org: %s", err).Error())
		return
	}

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

	if !staff.HasPermission(models.OrganizationDelete) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
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

// CreateOrganization
// @Summary Create organization
// @Security ApiKeyAuth
// @Tags organizations
// @Description create models.Organization
// @Description create organization
// @ID create-organization
// @Accept  json
// @Produce  json
// @Param input body organizationResponse true "organization info"
// @Success 200 {string} uuid
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/org/ [post]
func (h *Handler) CreateOrganization(c *gin.Context) {
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

	if !staff.HasPermission(models.OrganizationCreate) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}
	var org *models.Organization

	if err := c.Bind(&org); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not get input model in creating org: %s", err).Error())
		return
	}
	_, err = url.ParseRequestURI(org.WebsiteURL)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("can not url: %s", org.WebsiteURL))
		return
	}
	org.ID = uuid.New()
	if org.Types == nil || len(org.Types) == 0 {
		org.Types = []models.OrganizationType{
			models.DefaultOrganizationType,
		}
	}

	for _, orgType := range org.Types {
		if orgType.Name == "none" {
			if len(org.Positions) == 0 {
				org.Positions = []models.Position{
					models.AdminPosition,
				}
				org.Positions[0].ID = uuid.New()
			}
			break
		}
		if orgType.Name == "developer" {
			if org.Positions != nil {
				org.Positions = append(org.Positions, models.DefaultProgrammingPositions...)
			}
		}
	}
	err = h.Service.Organization.CreateOrganization(ctx, org, userID.(uuid.UUID))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not create model: %s", err).Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"created": org.ID,
	})
}

// UpdateOrganization
// @Summary Update organization by id
// @Security ApiKeyAuth
// @Tags organizations
// @Description Update organization by id
// @ID update-event
// @Accept  json
// @Produce  json
// @Param input body organizationUpdateResponse true "organization"
// @Success 200 {object} boolean
// @Failure 400,403 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/org/:id [put]
func (h *Handler) UpdateOrganization(c *gin.Context) {
	ctx := context.Background()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in updating org: %s", err).Error())
		return
	}
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

	if !staff.HasPermission(models.OrganizationUpdate) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
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

// GetOrganizationEvents
// @Summary Get Organization Events By ID
// @Security ApiKeyAuth
// @Tags organizations
// @Description Get Organization Events By ID
// @ID get-organization-events
// @Accept  json
// @Produce  json
// @Success 200 {object} []eventShortData
// @Failure 400,403 {} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/org/event/:id [get]
func (h *Handler) GetOrganizationEvents(c *gin.Context) {
	ctx := context.Background()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in getting org events: %s", err).Error())
		return
	}
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

	if !staff.HasPermission(models.OrganizationEvents) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}
	events, err := h.Service.Organization.GetOrganizationEvents(ctx, id, staff.ID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not get org events: %s", err).Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"events": events,
	})
}

// AddStaffToOrganization
// @Summary Add Staff To Organization
// @Security ApiKeyAuth
// @Tags organizations
// @Description Add Staff To Organization
// @ID add-staff-organization
// @Param input body users true "staff ids"
// @Accept  json
// @Produce  json
// @Success 200 {object} insertInfo
// @Failure 400,403 {} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/org/staff/:id [put]
func (h *Handler) AddStaffToOrganization(c *gin.Context) {
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

	if !staff.HasPermission(models.OrganizationAddStaff) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in adding staff into org: %s", err).Error())
		return
	}

	var requestStaff users

	if err := c.Bind(&requestStaff); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not bind ids in adding staff into org: %s", err).Error())
		return
	}

	err = h.Service.Organization.AddUsersToOrg(ctx, id, requestStaff.Staff)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not bind staff into org: %s", err).Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"inserted":       true,
		"inserted_count": len(requestStaff.Staff),
	})
}
