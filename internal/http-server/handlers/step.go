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

func (h *Handler) CreateStep(c *gin.Context) {
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

	if !staff.HasOneOfPermissions(models.StepCreate, models.EventCreate) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}

	var step *models.Step

	if err := c.Bind(&step); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not get input model in creating org: %s", err).Error())
		return
	}

	step.ID = uuid.New()
	step.Status = models.Process
	if step.CreationDate == "" {
		step.CreationDate = time.Now().Format(time.RFC3339)
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
	step.CreationDate = creationTime.String()
	step.EndDate = endTime.String()
	if len(step.Prizes) != 0 {
		for i := range step.Prizes {
			step.Prizes[i].ID = uuid.New()
			step.Prizes[i].StepID = step.ID
			step.Prizes[i].CreatedBy = userID.(uuid.UUID)
		}
	}
	if len(step.Images) != 0 {
		for i := range step.Images {
			step.Images[i].ID = uuid.New()
			step.Images[i].StepID = step.ID
		}
	}
	err = h.Service.Step.CreateStep(ctx, step, creationTime, endTime)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not create step: %s", err).Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"created": step.ID,
	})
}

func (h *Handler) UpdateStep(c *gin.Context) {
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

	if !staff.HasOneOfPermissions(models.StepUpdate, models.EventCreate) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in updating org: %s", err).Error())
		return
	}

	var step *models.Step

	if err := c.Bind(&step); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not get input model in creating org: %s", err).Error())
		return
	}

	step.ID = id

	err = h.Service.Step.UpdateStep(ctx, step)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not create model: %s", err).Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"updated": true,
	})
}

func (h *Handler) GetStep(c *gin.Context) {
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

	if !staff.HasOneOfPermissions(models.StepGetByID, models.StepGetByID) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in updating org: %s", err).Error())
		return
	}

	step, err := h.Service.Step.GetStep(ctx, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not create model: %s", err).Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"step": step,
	})
}

func (h *Handler) DeleteStep(c *gin.Context) {
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

	if !staff.HasOneOfPermissions(models.StepDelete, models.EventDelete) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in updating org: %s", err).Error())
		return
	}

	err = h.Service.Step.DeleteStep(ctx, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not create model: %s", err).Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"deleted": true,
	})
}

func (h *Handler) GetSteps(c *gin.Context) {
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

	if !staff.HasOneOfPermissions(models.StepGetAll, models.EventGetByID, models.EventGetAll) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in updating org: %s", err).Error())
		return
	}

	steps, err := h.Service.Step.GetSteps(ctx, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not create model: %s", err).Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"steps": steps,
	})
}

func (h *Handler) GetStepPrizes(c *gin.Context) {
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

	if !staff.HasOneOfPermissions(models.StepGetAll, models.EventGetByID, models.EventGetAll) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in updating org: %s", err).Error())
		return
	}

	prizes, err := h.Service.Step.GetStepPrizes(ctx, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not create model: %s", err).Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"prizes": prizes,
	})
}

func (h *Handler) PassStaff(c *gin.Context) {
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

	if !staff.HasOneOfPermissions(models.StepGetAll, models.EventGetByID, models.EventGetAll) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in updating org: %s", err).Error())
		return
	}

	var stepStatus *models.StepStatusRequest

	if err := c.Bind(&stepStatus); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not get input model in creating org: %s", err).Error())
		return
	}

	err = h.Service.Step.PassStaff(ctx, id, stepStatus.StaffID, stepStatus.StepStatus, stepStatus.Score)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not create model: %s", err).Error())
		return
	}
	if stepStatus.StepStatus == models.Done {
		step, err := h.Service.Step.GetStep(ctx, id)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not get model: %s", err).Error())
			return
		}
		steps, err := h.Service.Step.GetSteps(ctx, step.EventID)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not get model: %s", err).Error())
			return
		}
		for _, s := range steps {
			if s.Level > step.Level {
				err = h.Service.Step.AssignStaff(ctx, stepStatus.StaffID, s.ID)
				if err != nil {
					newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not create model: %s", err).Error())
					return
				}
				break
			}
		}
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"changed": true,
	})
}

func (h *Handler) AssignStaff(c *gin.Context) {
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

	if !staff.HasOneOfPermissions(models.StepGetAll, models.EventGetByID, models.EventGetAll) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in assign staff to step: %s", err).Error())
		return
	}

	var staffIDs []*models.StaffID

	if err := c.Bind(&staffIDs); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not get input model in assign staff to step: %s", err).Error())
		return
	}
	for _, staffID := range staffIDs {
		err = h.Service.Step.AssignStaff(ctx, staffID.StaffID, id)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not create model: %s", err).Error())
			return
		}
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"assigned": true,
	})
}
