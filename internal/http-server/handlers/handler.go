package handlers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/miprokop/fication/internal/services"
	"net/http"
	"time"

	_ "github.com/miprokop/fication/docs"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

type Handler struct {
	Service *services.Service
}

func NewHandler(services *services.Service) *Handler {
	return &Handler{Service: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodPost, http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}
	api := router.Group("/api", h.identity)
	{
		organization := api.Group("/org")
		{
			organization.GET("/", h.GetAllOrganizations)
			organization.GET("/:id", h.GetOrganization)
			organization.PUT("/:id", h.UpdateOrganization)
			organization.PUT("/staff/:id", h.AddStaffToOrganization)
			organization.DELETE("/:id", h.DeleteOrganization)
			organization.GET("/staff/:id", h.GetStaffByOrganizationID)
			organization.GET("/event/:id", h.GetOrganizationEvents)
			organization.POST("/", h.CreateOrganization)

			organizationType := organization.Group("/type")
			{
				organizationType.POST("/", h.CreateOrganizationType)
				organizationType.GET("/", h.GetOrganizationTypes)
				organizationType.GET("/:id", h.GetOrganizationTypeByID)
				organizationType.PUT("/:id", h.UpdateOrganizationType)
				organizationType.DELETE("/:id", h.DeleteOrganizationType)
			}
		}

		team := api.Group("/team")
		{
			team.POST("/", h.CreateTeam)
			team.GET("/org/:id", h.GetTeamsByOrganizationID)
			team.GET("/event/:id", h.GetTeamsByEventID)
			team.GET("/:id", h.GetTeamByID)
			team.PUT("/:id", h.UpdateTeam)
			team.DELETE("/:id", h.DeleteTeamByID)
		}

		user := api.Group("/user")
		{
			user.GET("/event/:id", h.GetAllUsersInEvent)
			user.GET("/step/:id", h.GetAllUsersInStep)
			user.GET("/:id", h.GetStaffByID)
			user.PUT("/:id", h.UpdateStaffByID)
			user.DELETE("/:id", h.DeleteStaff)
			user.POST("/", h.CreateStaff)
			user.GET("/prizes/:id", h.GetStaffPrizes)
			user.GET("/invites", h.GetStaffInvites)
			user.PUT("/photo", h.UploadImage)
			user.GET("/image/:id", h.GetImage)

			position := user.Group("/position")
			{
				position.PUT("/:id", h.UpdatePosition)
				position.PUT("/perm/:id", h.RemovePermissions)
				position.PUT("/give/:id", h.GivePosition)
				position.PUT("/take/:id", h.TakePosition)
				position.POST("/", h.CreatePosition)
				position.DELETE("/:id", h.DeletePosition)
				position.GET("/org/:id", h.GetOrganizationPositions)
				position.GET("/:id", h.GetPosition)
			}
		}

		event := api.Group("/event")
		{
			event.POST("/", h.CreateEvent) // ads
			event.POST("/invite/:id", h.AssignStaffToEvent)
			event.POST("/invitation/:id", h.AnswerInvitation) // ads
			event.GET("/invitation/", h.GetInvitations)       // ads
			event.GET("/:id", h.GetEventByID)                 // ads
			event.GET("/staff/:role", h.GetUserEvents)        // ads
			event.GET("/team/:id", h.GetTeamEvents)           // ads
			event.PUT("/:id", h.UpdateEvent)
			event.GET("/score/:id", h.GetStaffScore) // ads
			event.DELETE("/remove/:id", h.RemoveStaffFromEvent)
			event.DELETE("/:id", h.DeleteEvent)

			step := event.Group("/step")
			{
				step.PUT("/:id", h.UpdateStep)
				step.GET("/:id", h.GetStep)
				step.DELETE("/:id", h.DeleteStep)
				step.GET("/steps/:id", h.GetSteps)
				step.POST("/", h.CreateStep)
				step.GET("/prizes/:id", h.GetStepPrizes)
				step.PUT("/status/:id", h.PassStaff)
				step.PUT("/assign/:id", h.AssignStaff)
			}
		}
		prize := api.Group("/prize")
		{
			prize.POST("/", h.CreatePrize)
			prize.GET("/:id", h.GetPrize)
			prize.GET("/", h.GetPrizes)
			prize.GET("/user/:type", h.GetPrizesByType)
			prize.PUT("/:id", h.UpdatePrize)
			prize.POST("/give/:id", h.GivePrize)
		}
	}

	return router
}
