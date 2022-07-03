package handlers

import (
	"context"
	"encoding/json"
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
	"strings"
)

const (
	authHeader = "Authorization"
	imagePath  = "upload/files/staff"
)

func (h *Handler) identity(c *gin.Context) {
	header := c.GetHeader(authHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}
	temp := headerParts[1]
	userID, err := h.Service.Auth.ParseToken(temp)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth token")
		return
	}
	c.Set("userID", userID)
}

// @Summary SignUp
// @Tags auth
// @Description create models.Staff
// @Description if no organization_id in input
// @Description push staff do default org
// @Description do the same stuff with position_id and team_id if organization_id is empty
// @Description also in registering you can pass image
// @ID create-staff
// @Accept  json
// @Produce  json
// @Param input body models.Staff true "staff account info"
// @Success 200
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	ctx := context.Background()
	reqData := []byte(c.PostForm("json"))
	var input models.StaffSignUp

	if err := json.Unmarshal(reqData, &input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	id := uuid.New()
	file, _ := c.FormFile("file")
	if file != nil && file.Filename != "" {
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
		err := c.SaveUploadedFile(file, dst)
		if err != nil {
			log.Error(err)
			_ = os.Remove(dst)
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		input.CurrentImage = dst
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
	if input.TeamID == (uuid.UUID{}) && input.OrganizationID != (uuid.UUID{}) {
		defaultTeam, err := h.Service.Team.GetTeamByName(ctx, input.OrganizationID, models.DefaultTeamName)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not create defaul team: %s", err).Error())
			return
		}
		input.TeamID = defaultTeam.ID
	}

	if input.PositionID == (uuid.UUID{}) && input.OrganizationID != (uuid.UUID{}) {
		positions, err := h.Service.Staff.GetAllPositions(ctx, input.OrganizationID)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can not get positions: %s", err).Error())
			return
		}
		for _, p := range positions {
			if p.Name == "none" {
				input.PositionID = p.ID
			}
		}
	}

	if input.OrganizationID == (uuid.UUID{}) {
		input.OrganizationID = models.DefaultOrganization.ID
		input.PositionID = models.DefaultPosition.ID
		input.TeamID = models.DefaultTeam.ID
	}

	input.ID = id
	input.TextColor = "#000000"
	input.BackgroundColor = "#fffff"
	err := h.Service.Staff.CreateStaffUser(c.Request.Context(), &input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusCreated)
}

// @Summary signIn
// @Tags auth
// @Description sign in staff to get token
// @Description token is used in authorization
// @ID sign-in-staff
// @Accept  json
// @Produce  json
// @Param input body models.StaffLogin true "staff account log in info"
// @Success 200 {string} token
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var input models.StaffLogin
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, id, orgID, err := h.Service.Auth.GenerateToken(input.Email, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token":           token,
		"id":              id,
		"organization_id": orgID,
	})
}
