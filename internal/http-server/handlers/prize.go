package handlers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/miprokop/fication/internal/models"
	"net/http"
	"net/url"
	"time"
)

func (h *Handler) CreatePrize(c *gin.Context) {
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

	if !staff.HasPermission(models.PrizeCreate) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}
	var prize *models.Prize

	if err := c.Bind(&prize); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not get input model in creating prize: %s", err).Error())
		return
	}

	if !models.OneOf(prize.PrizeStatus) {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("incorrent prize status: %s, want one of: %s, %s, %s, %s",
			prize.PrizeStatus, models.Legendary, models.Mith, models.Common, models.Rare))
		return
	}
	if prize.PrizeType == models.Image || prize.PrizeType == models.Medal {
		_, err = url.ParseRequestURI(prize.Data)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("can not url: %s", prize.Data))
			return
		}
	}

	prize.ID = uuid.New()
	prize.CreatedBy = userID.(uuid.UUID)
	prize.CurrentCount = prize.Count
	if prize.CreationDate == "" {
		prize.CreationDate = time.Now().Format(time.RFC3339)
	}
	err = h.Service.Prize.CreatePrize(ctx, prize)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not create model: %s", err).Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"created": prize.ID,
	})
}

func (h *Handler) GetPrize(c *gin.Context) {
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

	if !staff.HasPermission(models.PrizeGetByID) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in prize getting: %s", err).Error())
		return
	}

	prize, err := h.Service.Prize.GetPrize(ctx, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not create model: %s", err).Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"prize": prize,
	})
}

func (h *Handler) GetPrizesByType(c *gin.Context) {
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

	if !staff.HasPermission(models.PrizeGetByID) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}

	prizeType, err := models.NewPrizeType(c.Param("type"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in prize getting: %s", err).Error())
		return
	}

	prizes, err := h.Service.Prize.GetPrizesByType(ctx, prizeType)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not create model: %s", err).Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"prizes": prizes,
	})
}

func (h *Handler) GetPrizes(c *gin.Context) {
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

	if !staff.HasPermission(models.PrizeGetAll) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}

	prizes, err := h.Service.Prize.GetPrizes(ctx, userID.(uuid.UUID))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not create model: %s", err).Error())
		return
	}
	if len(prizes) == 0 {
		c.JSON(http.StatusOK, "no prizes created by this user")
	} else {
		c.JSON(http.StatusOK, map[string]interface{}{
			"prizes": prizes,
		})
	}
}

func (h *Handler) UpdatePrize(c *gin.Context) {
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

	if !staff.HasPermission(models.PrizeUpdate) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in updating prize: %s", err).Error())
		return
	}

	var prize *models.Prize

	if err := c.Bind(&prize); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not get input update model in updating prize: %s", err).Error())
		return
	}
	if prize.PrizeType == models.Image || prize.PrizeType == models.Medal {
		if prize.Data != "" {
			_, err = url.ParseRequestURI(prize.Data)
			if err != nil {
				newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("can not url: %s", prize.Data))
				return
			}
		}
	}

	prize.ID = id

	err = h.Service.Prize.UpdatePrize(ctx, prize)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not update model in updating prize: %s", err).Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"updated": true,
	})
}

func (h *Handler) GivePrize(c *gin.Context) {
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

	if !staff.HasPermission(models.PrizeGive) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in updating prize: %s", err).Error())
		return
	}

	var staffID *models.StaffID

	if err := c.Bind(&staffID); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not get input update model in updating prize: %s", err).Error())
		return
	}
	err = h.Service.Prize.GivePrize(ctx, staffID.StaffID, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not update model in updating prize: %s", err).Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"updated": true,
	})
}
