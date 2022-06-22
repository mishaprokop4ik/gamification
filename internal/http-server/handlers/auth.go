package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/miprokop/fication/internal/models"
	"net/http"
	"os"
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

func (h *Handler) signUp(c *gin.Context) {
	reqData := []byte(c.PostForm("json"))
	var input models.StaffSignUp

	if err := json.Unmarshal(reqData, &input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	file, _ := c.FormFile("file")
	if file != nil && file.Filename != "" {
		dst := fmt.Sprintf("%s/%s", imagePath, file.Filename)
		err := c.SaveUploadedFile(file, dst)
		if err != nil {
			_ = os.Remove(dst)
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		input.CurrentImage = dst
	}
	if !(string(input.TextColor) == "") && !input.TextColor.IsHex() ||
		!(string(input.BackgroundColor) == "") && !input.BackgroundColor.IsHex() {
		newErrorResponse(c, http.StatusBadRequest,
			fmt.Sprintf("incorrect color format: %s, %s, want: #000000",
				input.TextColor, "input.BackgroundColor"))
		return
	}
	if input.OrganizationID == (uuid.UUID{}) {
		input.OrganizationID = models.DefaultOrganization.ID
		input.PositionID = models.DefaultPosition.ID
		input.TeamID = models.DefaultTeam.ID
	}
	input.ID = uuid.New()
	input.TextColor = "#000000"
	input.BackgroundColor = "#fffff"
	err := h.Service.Staff.CreateStaffUser(c.Request.Context(), &input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) signIn(c *gin.Context) {
	var input models.StaffLogin
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.Service.Auth.GenerateToken(input.Email, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
