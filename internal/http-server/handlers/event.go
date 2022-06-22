package handlers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/miprokop/fication/internal/models"
	"net/http"
)

func (h *Handler) GetUserEvents(c *gin.Context) {
	ctx := context.Background()
	id, ok := c.Get("userID")
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError,
			"there is no userID in context")
		return
	}

	user, err := h.Service.Staff.GetStaff(ctx, id.(uuid.UUID))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("can not get staff by id; err: %s;", err.Error()))
		return
	}
	if !user.Position.HasPermission(models.EventGetAll) {
		newErrorResponse(c, http.StatusForbidden, "")
		return
	}
	events, err := h.Service.Event.GetStaffsEvents(ctx, id.(uuid.UUID))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError,
			fmt.Errorf("can not get staff by id: %s", err).Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"events": events,
	})
}
