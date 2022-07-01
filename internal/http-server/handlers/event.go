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

// GetUserEvents
// @Summary Get Events where staff takes a part by role
// @Security ApiKeyAuth
// @Tags events
// @Description Get Events where staff takes a part by role
// @ID get-staff-events
// @Accept  json
// @Produce  json
// @Success 200 {object} eventsResponse
// @Failure 400,403 {} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/staff/:role [get]
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

// CreateEvent
// @Summary Create event
// @Security ApiKeyAuth
// @Tags events
// @Description create models.Event
// @Description if no org id in request
// @Description use a default org id
// @Description if no type use public type
// @Description event status by default process
// @Description add a user who created this event to a member of this event
// @Description if some steps there creates it
// @ID create-event
// @Accept  json
// @Produce  json
// @Param input body eventRequest true "event info"
// @Success 200 {string} uuid
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/event/ [post]
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

// AssignStaffToEvent
// @Summary Assign Staff
// @Security ApiKeyAuth
// @Tags events
// @Description Assign staff array
// @ID assign-staff-to-event
// @Accept  json
// @Produce  json
// @Param input body []staffEvents true "staff ids and events"
// @Success 200 {string} uuid
// @Failure 400,403 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/event/invite/:id [post]
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

// AnswerInvitation
// @Summary Answer invite
// @Security ApiKeyAuth
// @Tags events
// @Description Answer invite by ID
// @Description status can be only models.InviteStatus
// @ID answer-invite
// @Accept  json
// @Produce  json
// @Param input body staffEvents true "staff and invites relation"
// @Success 200 {string} uuid
// @Failure 400,403 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/event/invitation/:id [post]
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

// GetInvitations
// @Summary Get all invitations
// @Security ApiKeyAuth
// @Tags events
// @Description Get all invitations by current user
// @ID get-invites
// @Accept  json
// @Produce  json
// @Success 200 {object} []staffEvents
// @Failure 400,403 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/event/invitation/ [get]
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

// GetEventByID
// @Summary Get event by id
// @Security ApiKeyAuth
// @Tags events
// @Description Get all invitations by current user
// @ID get-invite-by-id
// @Accept  json
// @Produce  json
// @Success 200 {object} eventAllData
// @Failure 400,403 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/event/:id [get]
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

// UpdateEvent
// @Summary Update event by id
// @Security ApiKeyAuth
// @Tags events
// @Description Update event by ID
// @ID update-event
// @Accept  json
// @Produce  json
// @Param input body eventRequestUpdate true "event"
// @Success 200 {object} eventAllData
// @Failure 400,403 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/event/:id [put]
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
	for i := range event.StaffEvents {
		event.StaffEvents[i].EventID = id
	}
	err = h.Service.Event.UpdateEvent(ctx, event)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not update model in updating team: %s", err).Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"updated": true,
	})
}

// DeleteEvent
// @Summary Delete Event By ID
// @Security ApiKeyAuth
// @Tags events
// @Description delete event by id
// @Description also delete chained steps
// @ID delete-event
// @Accept  json
// @Produce  json
// @Success 200 {object} boolean
// @Failure 400,403 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/event/:id [delete]
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

// GetTeamEvents
// @Summary Get all events in team
// @Security ApiKeyAuth
// @Tags events
// @Description Get all team's events by ID
// @ID team-events
// @Accept  json
// @Produce  json
// @Success 200 {object} []eventAllData
// @Failure 400,403 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/event/team/:id [get]
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

// RemoveStaffFromEvent
// @Summary Remove staff from event
// @Security ApiKeyAuth
// @Tags events
// @Description Remove staff from event and connected steps
// @ID event-remove-staff
// @Accept  json
// @Param input body models.StaffID true "staffID"
// @Produce  json
// @Success 200 {object} boolean
// @Failure 400,403 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/event/remove/:id [delete]
func (h *Handler) RemoveStaffFromEvent(c *gin.Context) {
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

	if !staff.HasOneOfPermissions(models.EventCreate, models.EventDelete) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in teams by org: %s", err).Error())
		return
	}
	var staffID *models.StaffID

	if err := c.Bind(&staffID); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not get input model in removing staff from event: %s", err).Error())
		return
	}
	event := models.StaffEvents{
		StaffID: staffID.StaffID,
		EventID: id,
	}

	err = h.Service.Event.RemoveStaffFromEvent(ctx, event)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not create model: %s", err).Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"removed": true,
	})
}

// GetStaffScore
// @Summary Get Staff Score In Event
// @Security ApiKeyAuth
// @Tags events
// @Description get staff score in event by ID
// @ID delete-event
// @Accept  json
// @Produce  json
// @Success 200 {object} integer
// @Failure 400,403 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/event/score/:id [get]
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
