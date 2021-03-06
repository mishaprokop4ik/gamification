package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/miprokop/fication/internal/models"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"path/filepath"
)

func (h *Handler) GetStaffInvites(c *gin.Context) {
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
	if !staff.HasOneOfPermissions(models.StaffGetInvites, models.StaffGetSelfInvites) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}

	invites, err := h.Service.Staff.GetInvites(ctx, userID.(uuid.UUID))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"invites": invites,
	})
}

func (h *Handler) GetStaffPrizes(c *gin.Context) {
	ctx := context.Background()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in users in events: %s", err).Error())
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
	staff, err := h.Service.Staff.GetStaff(ctx, id)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not get staff by id: %s", err).Error())
		return
	}
	if id == userID.(uuid.UUID) {
		if !staff.HasPermission(models.StaffSelfGet) {
			newErrorResponse(c, http.StatusForbidden,
				"no access to this action")
			return
		}
	} else {
		if !staff.HasPermission(models.PrizeStaffAll) {
			newErrorResponse(c, http.StatusForbidden,
				"no access to this action")
			return
		}
	}
	prizes, err := h.Service.Staff.GetStaffPrizes(ctx, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"prizes": prizes,
	})
}

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
	userID, ok := c.Get("userID")
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError,
			"there is no userID in context")
		return
	}
	_, ok = userID.(uuid.UUID)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError,
			"there is incorrect userID in context")
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in getting staff by id: %s", err).Error())
		return
	}

	staff, err := h.Service.Staff.GetStaff(ctx, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError,
			fmt.Errorf("can not get staff by id: %s", err).Error())
		return
	}

	if id == userID.(uuid.UUID) {
		if !staff.HasPermission(models.StaffSelfGet) {
			newErrorResponse(c, http.StatusForbidden,
				"no access to this action")
			return
		}
	} else {
		if !staff.HasPermission(models.PrizeStaffAll) {
			newErrorResponse(c, http.StatusForbidden,
				"no access to this action")
			return
		}
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"staff": staff,
	})
}

func (h *Handler) UpdateStaffByID(c *gin.Context) {
	ctx := context.Background()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in updating staff: %s", err).Error())
		return
	}
	userID, ok := c.Get("userID")
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError,
			"there is no userID in context")
		return
	}
	staff, err := h.Service.Staff.GetStaff(ctx, id)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not get staff by id: %s", err).Error())
		return
	}
	_, ok = userID.(uuid.UUID)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError,
			"there is incorrect userID in context")
		return
	}
	if id == userID.(uuid.UUID) {
		if !staff.HasPermission(models.StaffSelfUpdate) {
			newErrorResponse(c, http.StatusForbidden,
				"no access to this action")
			return
		}
	} else {
		if !staff.HasPermission(models.StaffUpdate) {
			newErrorResponse(c, http.StatusForbidden,
				"no access to this action")
			return
		}
	}
	var input *models.StaffSignUp
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if !staff.Sex.IsCorrect(string(staff.Sex)) {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("incorrect sex input: %s; want: %s, %s", staff.Sex,
			models.Male, models.Female))
		return
	}
	if input.Email == viper.GetString("admin.email") {
		newErrorResponse(c, http.StatusForbidden,
			fmt.Sprintf("use this email is forbidden"))
		return
	}
	if !(string(input.TextColor) == "") && !input.TextColor.IsHex() ||
		!(string(input.BackgroundColor) == "") && !input.BackgroundColor.IsHex() {
		newErrorResponse(c, http.StatusBadRequest,
			fmt.Sprintf("incorrect color format: %s, %s, want: #000000",
				input.TextColor, "input.BackgroundColor"))
		return
	}
	staffUpdate := &models.Staff{
		ID:              id,
		FirstName:       input.FirstName,
		LastName:        input.LastName,
		Email:           input.Email,
		Password:        input.Password,
		Sex:             input.Sex,
		AdditionalInfo:  input.AdditionalInfo,
		TeamID:          input.TeamID,
		PositionID:      input.PositionID,
		OrganizationID:  input.OrganizationID,
		TextColor:       input.TextColor,
		BackgroundColor: input.BackgroundColor,
	}

	err = h.Service.Staff.UpdateStaff(ctx, staffUpdate)
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
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not parse input id in deleting staff: %s", err).Error())
		return
	}
	userID, ok := c.Get("userID")
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError,
			"there is no userID in context")
		return
	}
	staff, err := h.Service.Staff.GetStaff(ctx, id)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("can not get staff by id: %s", err).Error())
		return
	}
	_, ok = userID.(uuid.UUID)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError,
			"there is incorrect userID in context")
		return
	}
	if id == userID.(uuid.UUID) {
		if !staff.HasPermission(models.StaffSelfDelete) {
			newErrorResponse(c, http.StatusForbidden,
				"no access to this action")
			return
		}
	} else {
		if !staff.HasPermission(models.StaffDelete) {
			newErrorResponse(c, http.StatusForbidden,
				"no access to this action")
			return
		}
	}
	err = h.Service.Staff.DeleteStaff(ctx, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError,
			fmt.Errorf("can not delete staff by id: %s", err).Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"deleted": true,
	})
}

func (h *Handler) GetImage(c *gin.Context) {
	id, ok := c.Get("userID")
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id is not found")
		return
	}
	fileName := c.Param("id")
	endpointFile := fmt.Sprintf("%s/%s/%s", imagePath, id, fileName)
	c.File(endpointFile)
}

func (h *Handler) CreateStaff(c *gin.Context) {
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

	if !staff.HasPermission(models.StaffCreate) {
		newErrorResponse(c, http.StatusForbidden,
			"no access to this action")
		return
	}
	var input *models.StaffSignUp

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not create model: %s", err).Error())
		return
	}
	if !staff.Sex.IsCorrect(string(staff.Sex)) {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("incorrect sex input: %s; want: %s, %s", staff.Sex,
			models.Male, models.Female))
		return
	}
	if input.Email == viper.GetString("admin.email") {
		newErrorResponse(c, http.StatusForbidden,
			fmt.Sprintf("use this email is forbidden"))
		return
	}
	if !(string(input.TextColor) == "") && !input.TextColor.IsHex() ||
		!(string(input.BackgroundColor) == "") && !input.BackgroundColor.IsHex() {
		newErrorResponse(c, http.StatusBadRequest,
			fmt.Sprintf("incorrect color format: %s, %s, want: #000000",
				input.TextColor, "input.BackgroundColor"))
		return
	}

	input.ID = uuid.New()
	input.TextColor = "#000000"
	input.BackgroundColor = "#fffff"
	err = h.Service.Staff.CreateStaffUser(c.Request.Context(), input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) UploadImage(c *gin.Context) {
	id, ok := c.Get("userID")
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError,
			"there is no userID in context")
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if filepath.Ext(file.Filename) != ".png" {
		newErrorResponse(c, http.StatusBadRequest,
			fmt.Sprintf("this format is unsupported: %s; want: png", filepath.Ext(file.Filename)))
		return
	}
	dir := fmt.Sprintf("%s/%s", imagePath, id)
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
	dst := fmt.Sprintf("%s/%s/%s", imagePath, id, file.Filename)
	err = c.SaveUploadedFile(file, dst)
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
