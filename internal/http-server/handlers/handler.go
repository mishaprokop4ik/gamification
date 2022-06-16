package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/miprokop/fication/internal/services"
)

type Handler struct {
	Service *services.Servicer
}

func NewHandler(services *services.Servicer) *Handler {
	return &Handler{Service: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api", h.userIdentity)
	{
		organization := api.Group("/org")
		{
			organization.GET("/")          // get all
			organization.GET("/:id")       // get by id
			organization.PUT("/:id")       // update org
			organization.PUT("/staff/:id") // add staff to org
			organization.DELETE("/:id")    // delete
			organization.GET("/event/:id") // get all events in org
		}

		organizationType := api.Group("/org-type")
		{
			organizationType.POST("/")      //create new
			organizationType.GET("/")       // get all
			organizationType.GET("/:id")    // get by id
			organizationType.PUT("/:id")    // update
			organizationType.DELETE("/:id") // delete
		}

		team := api.Group("/team")
		{
			team.POST("/")         //create new
			team.GET("/org/:id")   // get all by company
			team.GET("/event/:id") // get all teams by event
			team.GET("/:id")       // get by id
			team.PUT("/:id")       // update
			team.DELETE("/:id")    // delete
		}

		user := api.Group("/user")
		{
			user.POST("/")         //create new // TODO STOP THERE
			user.GET("/org")       // get all by org
			user.GET("/event/:id") // all user in event with their status
			user.GET("/step/:id")  // all user in event with their status
			user.GET("/")          // get by id
			user.PUT("/")          // update
			user.DELETE("/")       // delete
			user.GET("/prizes")    // get all prizes
			user.GET("/invites")   // get users invites

			role := user.Group("/role")
			{
				role.PUT("/:id")
				role.PUT("/give/:id")
				role.POST("/")
				role.DELETE("/:id")
				role.GET("/:comp")
				role.GET("/:id")
			}
		}

		event := api.Group("/event")
		{
			event.POST("/")         // create new event with steps
			event.GET("/:id")       // get event by id with all data
			event.GET("/")          // all user's event
			event.PUT("/:id")       // update event
			event.GET("/score/:id") // get score

			step := event.Group("/step")
			{
				step.PUT("/:id")       // update step in event
				step.GET("/:id")       // get all step info
				step.GET("/steps/:id") // get all steps in event
				step.POST("/")         // create new step
				step.GET("/users/:id") // get all users in step
				step.GET("/prizes")
			}
		}
		prize := api.Group("/prize")
		{
			prize.POST("/")         // create new event with steps
			prize.GET("/:id")       // get event by id with all data
			prize.GET("/")          // all user's prizes
			prize.PUT("/:id")       // update event
			prize.POST("/give/:id") // give prize to user
		}
	}

	return router
}
