package handlers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/miprokop/fication/internal/models"
	"net/http"
	"time"
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

	staffRole := c.Param("role")

	events, err := h.Service.Event.GetStaffsEvents(ctx, id.(uuid.UUID), staffRole)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError,
			fmt.Errorf("can not get events by staff role: %s", err).Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"events": events,
	})
}

func (h *Handler) CreateEvent(c *gin.Context) {
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

	if !staff.HasPermission(models.EventCreate) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}
	var event *models.Event

	if err := c.Bind(&event); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not get input model in creating org: %s", err).Error())
		return
	}
	if event.CreationDate == "" {
		event.CreationDate = time.Now().Format(time.RFC3339)
	}
	if event.OrganizationID == (uuid.UUID{}) {
		event.OrganizationID = models.DefaultOrganization.ID
	}
	if event.EventType == "" {
		event.EventType = "public"
	}
	if event.EventStatus == "" {
		event.EventStatus = "process"
	}
	event.ID = uuid.New()
	for i := range event.StaffEvents {
		event.StaffEvents[i].EventID = event.ID
	}
	event.StaffEvents = append(event.StaffEvents, &models.StaffEvents{
		ID:        uuid.New(),
		StaffID:   userID.(uuid.UUID),
		EventID:   event.ID,
		Status:    "accepted",
		StaffRole: models.Creator,
	})
	err = h.Service.Event.CreateEvent(ctx, event)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not create event: %s", err).Error())
		return
	}
	if len(event.Steps) != 0 {
		for i, step := range event.Steps {
			step := *step
			step.ID = uuid.New()
			step.Status = models.Process
			step.EventID = event.ID
			step.Level = uint(i + 1)
			if step.CreationDate == "" {
				step.CreationDate = time.Now().Format(time.RFC3339)
			}
			if step.Status == "" {
				step.Status = models.Process
			}
			creationTime, err := time.Parse(time.RFC3339, step.CreationDate)
			if err != nil {
				newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("incorrent creation time: %s", err).Error())
				return
			}

			endTime, err := time.Parse(time.RFC3339, step.EndDate)
			if err != nil {
				newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("incorrent end time: %s", err).Error())
				return
			}
			if !endTime.After(creationTime) {
				newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("incorrent end time and creation time: %s, %s", endTime,
					creationTime).Error())
				return
			}

			err = h.Service.Step.CreateStep(ctx, &step, creationTime, endTime)
			if err != nil {
				newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not create model: %s", err).Error())
				return
			}
		}
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"created": event.ID,
	})
}

func (h *Handler) AssignStaffToEvent(c *gin.Context) {
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

	if !staff.HasOneOfPermissions(models.EventCreate, models.EventUpdate, models.StaffSelfUpdate) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input: %s", err).Error())
		return
	}

	var staffEvents []models.StaffEvents

	if err := c.Bind(&staffEvents); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not get input: %s", err).Error())
		return
	}
	err = h.Service.Event.AssignStaff(ctx, staffEvents, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not create model: %s", err).Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"assigned": true,
	})
}

func (h *Handler) AnswerInvitation(c *gin.Context) {
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

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in answer invites: %s", err).Error())
		return
	}

	var staffEvents models.StaffEvents

	if err := c.Bind(&staffEvents); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not get input model in creating org: %s", err).Error())
		return
	}
	staffEvents.StaffID = userID.(uuid.UUID)
	staffEvents.EventID = id
	err = h.Service.Event.AnswerInvitation(ctx, staffEvents)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not create model: %s", err).Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"answered": true,
	})
}

func (h *Handler) GetInvitations(c *gin.Context) {
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

	invites, err := h.Service.Event.GetInvites(ctx, userID.(uuid.UUID))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not create model: %s", err).Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"invites": invites,
	})
}

func (h *Handler) GetEventByID(c *gin.Context) {
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

	if !staff.HasPermission(models.EventGetByID) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in teams by org: %s", err).Error())
		return
	}

	event, err := h.Service.Event.GetEvent(ctx, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not create model: %s", err).Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"event": event,
	})
}

func (h *Handler) UpdateEvent(c *gin.Context) {
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

	if !staff.HasOneOfPermissions(models.EventUpdate, models.OrganizationUpdate) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in updating team: %s", err).Error())
		return
	}

	var event *models.Event

	if err := c.Bind(&event); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not get input update model in updating team: %s", err).Error())
		return
	}

	event.ID = id

	err = h.Service.Event.UpdateEvent(ctx, event)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not update model in updating team: %s", err).Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"updated": true,
	})
}

func (h *Handler) DeleteEvent(c *gin.Context) {
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

	if !staff.HasOneOfPermissions(models.EventDelete, models.OrganizationDelete) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in teams deletion: %s", err).Error())
		return
	}

	err = h.Service.Event.DeleteEvent(ctx, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not delete team: %s", err).Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"deleted": true,
	})
}

func (h *Handler) GetTeamEvents(c *gin.Context) {
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

	if !staff.HasOneOfPermissions(models.EventGetByID, models.TeamGetByID) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in teams by org: %s", err).Error())
		return
	}

	events, err := h.Service.Event.GetEventsByTeamID(ctx, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not create model: %s", err).Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"events": events,
	})
}

func (h *Handler) GetStaffScore(c *gin.Context) {
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

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in teams by org: %s", err).Error())
		return
	}

	events, err := h.Service.Event.GetStaffScore(ctx, id, userID.(uuid.UUID))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not create model: %s", err).Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"events": events,
	})
}
