package handlers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/miprokop/fication/internal/models"
	"net/http"
)

// UpdatePosition
// @Summary Update position by id
// @Security ApiKeyAuth
// @Tags positions
// @Description Update position by id
// @ID update-position
// @Accept  json
// @Produce  json
// @Param input body updatePosition true "position type"
// @Success 200 {object} boolean
// @Failure 400,403 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/user/position/:id [put]
func (h *Handler) UpdatePosition(c *gin.Context) {
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

	if !staff.HasPermission(models.PositionUpdate) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in deleting org: %s", err).Error())
		return
	}
	var position *models.Position

	if err := c.Bind(&position); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not get input model in creating position: %s", err).Error())
		return
	}
	position.ID = id
	if len(position.Permissions) != 0 {
		for i := range position.Permissions {
			position.Permissions[i].PositionID = position.ID
			position.Permissions[i].GrantedBy = userID.(uuid.UUID)
		}
	}
	err = h.Service.Staff.UpdatePosition(ctx, position)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not create model: %s", err).Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"updated": true,
	})
}

// RemovePermissions
// @Summary Remove Permissions From Position
// @Security ApiKeyAuth
// @Tags positions
// @Description Remove Permissions From Position
// @ID remove-permissions
// @Accept  json
// @Produce  json
// @Param input body permissions true "position id with permissions"
// @Success 200 {object} boolean
// @Failure 400,403 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/user/perm/:id [put]
func (h *Handler) RemovePermissions(c *gin.Context) {
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

	if !staff.HasPermission(models.PositionUpdate) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in deleting org: %s", err).Error())
		return
	}
	var permissions *models.Permissions

	if err := c.Bind(&permissions); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not get input model in creating position: %s", err).Error())
		return
	}
	permissions.PositionID = id
	for i := range permissions.Permissions {
		permissions.Permissions[i].PositionID = id
	}
	err = h.Service.Staff.RemovePermissionsFromPosition(ctx, *permissions)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not create model: %s", err).Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"removed": true,
	})
}

func (h *Handler) DeletePosition(c *gin.Context) {
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

	if !staff.HasPermission(models.PositionDelete) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in deleting org: %s", err).Error())
		return
	}
	err = h.Service.Staff.DeletePosition(ctx, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"deleted": true,
	})
}

// CreatePosition
// @Summary Create Position type
// @Security ApiKeyAuth
// @Tags positions
// @Description create models.Position
// @ID create-position
// @Accept  json
// @Produce  json
// @Param input body positionCreate true "position info"
// @Success 200 {string} uuid
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/user/position/ [post]
func (h *Handler) CreatePosition(c *gin.Context) {
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

	if !staff.HasPermission(models.PositionCreate) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}
	var position *models.Position

	if err := c.Bind(&position); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not get input model in creating position: %s", err).Error())
		return
	}

	position.ID = uuid.New()
	for i := range position.Permissions {
		position.Permissions[i].PositionID = position.ID
		position.Permissions[i].GrantedBy = staff.ID
	}
	err = h.Service.Staff.CreatePosition(ctx, position)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not create model: %s", err).Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"created": position.ID,
	})
}

func (h *Handler) TakePosition(c *gin.Context) {
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

	if !staff.HasPermission(models.PositionGive) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}
	var staffIDs []*models.StaffID

	if err := c.Bind(&staffIDs); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not get input model in creating position: %s", err).Error())
		return
	}

	for _, staffID := range staffIDs {
		err = h.Service.Staff.RemoveFromPosition(ctx, staffID.StaffID)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not create model: %s", err).Error())
			return
		}
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"removed": true,
	})
}

func (h *Handler) GivePosition(c *gin.Context) {
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

	if !staff.HasPermission(models.PositionGive) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}
	var staffIDs []*models.StaffID

	if err := c.Bind(&staffIDs); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not get input model in creating position: %s", err).Error())
		return
	}

	positionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in users in events: %s", err).Error())
		return
	}
	for _, staffID := range staffIDs {
		err = h.Service.Staff.AssignPosition(ctx, staffID.StaffID, positionID)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not create model: %s", err).Error())
			return
		}
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"gave": true,
	})
}

func (h *Handler) GetOrganizationPositions(c *gin.Context) {
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

	if !staff.HasPermission(models.PositionGetAll) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in getting org: %s", err).Error())
		return
	}
	positions, err := h.Service.Staff.GetAllPositions(ctx, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"positions": positions,
	})
}

func (h *Handler) GetPosition(c *gin.Context) {
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

	if !staff.HasPermission(models.PositionGetByID) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in getting org: %s", err).Error())
		return
	}
	position, err := h.Service.Staff.GetPosition(ctx, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"position": position,
	})
}
