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
					models.DefaultPosition,
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
	events, err := h.Service.Organization.GetOrganizationEvents(ctx, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not get org events: %s", err).Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"events": events,
	})
}

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

	type users struct {
		Staff []*models.StaffInsertion `json:"users"`
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
