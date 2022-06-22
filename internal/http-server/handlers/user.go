package handlers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/miprokop/fication/internal/models"
	"net/http"
	"os"
)

func (h *Handler) GetAllUsersInEvent(c *gin.Context) {
	ctx := context.Background()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in users in events: %s", err).Error())
		return
	}

	staff, err := h.Service.Staff.GetStaffByEvent(ctx, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError,
			fmt.Errorf("can not parse input id in users in events: %s", err).Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"staff": staff,
	})
}

func (h *Handler) GetAllUsersInStep(c *gin.Context) {
	ctx := context.Background()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in users in step: %s", err).Error())
		return
	}

	staff, err := h.Service.Staff.GetStaffByStep(ctx, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError,
			fmt.Errorf("can not parse input id in users in step: %s", err).Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"staff": staff,
	})
}

func (h *Handler) GetStaffByID(c *gin.Context) {
	ctx := context.Background()
	id, ok := c.Get("userID")
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError,
			"there is no userID in context")
		return
	}

	staff, err := h.Service.Staff.GetStaff(ctx, id.(uuid.UUID))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError,
			fmt.Errorf("can not get staff by id: %s", err).Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"staff": staff,
	})
}

func (h *Handler) UpdateStaffByID(c *gin.Context) {
	ctx := context.Background()
	id, ok := c.Get("userID")
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError,
			"there is no userID in context")
		return
	}

	var input *models.Staff
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	input.ID = id.(uuid.UUID)
	err := h.Service.Staff.UpdateStaff(ctx, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError,
			fmt.Errorf("can not update staff by id: %s", err).Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"updated": true,
	})
}

func (h *Handler) DeleteStaff(c *gin.Context) {
	ctx := context.Background()
	id, ok := c.Get("userID")
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError,
			"there is no userID in context")
		return
	}

	err := h.Service.Staff.DeleteStaff(ctx, id.(uuid.UUID))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError,
			fmt.Errorf("can not delete staff by id: %s", err).Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"deleted": true,
	})
}

func (h *Handler) UploadImage(c *gin.Context) {
	id, ok := c.Get("userID")
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError,
			"there is no userID in context")
		return
	}
	file, _ := c.FormFile("file")
	dst := fmt.Sprintf("%s/%s", imagePath, file.Filename)
	err := c.SaveUploadedFile(file, dst)
	if err != nil {
		_ = os.Remove(dst)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var staffImage = models.StaffImage{
		ID:        uuid.New(),
		UserID:    id.(uuid.UUID),
		ImagePath: dst,
	}
	err = h.Service.Staff.UploadImage(c.Request.Context(), staffImage)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusCreated)
}
