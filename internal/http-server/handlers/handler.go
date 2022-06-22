package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/miprokop/fication/internal/services"
)

type Handler struct {
	Service *services.Service
}

func NewHandler(services *services.Service) *Handler {
	return &Handler{Service: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}
	api := router.Group("/api", h.identity)
	{
		organization := api.Group("/org")
		{
			organization.GET("/", h.GetAllOrganizations)               // get all
			organization.GET("/:id", h.GetOrganization)                // get by id
			organization.PUT("/:id", h.UpdateOrganization)             // update org /
			organization.PUT("/staff/:id", h.AddStaffToOrganization)   // to check
			organization.DELETE("/:id", h.DeleteOrganization)          // delete
			organization.GET("/staff/:id", h.GetStaffByOrganizationID) // to check
			organization.GET("/event/:id", h.GetOrganizationEvents)    // to check
			organization.POST("/", h.CreateOrganization)               // get all events in org

			organizationType := organization.Group("/type")
			{
				organizationType.POST("/", h.CreateOrganizationType)      //create new
				organizationType.GET("/", h.GetOrganizationTypes)         // get all
				organizationType.GET("/:id", h.GetOrganizationTypeByID)   // get by id
				organizationType.PUT("/:id", h.UpdateOrganizationType)    // update
				organizationType.DELETE("/:id", h.DeleteOrganizationType) // delete
			}
		}

		team := api.Group("/team")
		{
			team.POST("/", h.CreateTeam)                     //create new
			team.GET("/org/:id", h.GetTeamsByOrganizationID) // get all by company
			team.GET("/event/:id", h.GetTeamsByEventID)      // to check
			team.GET("/:id", h.GetTeamByID)                  // get by id
			team.PUT("/:id", h.UpdateTeam)                   // update
			team.DELETE("/:id", h.DeleteTeamByID)            // delete
		}

		user := api.Group("/user")
		{
			user.GET("/event/:id", h.GetAllUsersInEvent) // all user in event with their status
			user.GET("/step/:id", h.GetAllUsersInStep)   // all user in step with their status
			user.GET("/", h.GetStaffByID)                // get by id
			user.PUT("/", h.UpdateStaffByID)             // update
			user.DELETE("/", h.DeleteStaff)              // delete
			//user.GET("/prizes")                          // get all prizes
			//user.GET("/invites")                         // get users invites
			user.PUT("/photo", h.UploadImage)
			position := user.Group("/position")
			{
				position.PUT("/:id")
				position.PUT("/give/:id")
				position.POST("/")
				position.DELETE("/:id")
				position.GET("/:comp")
				position.GET("/:id")
			}
		}

		event := api.Group("/event")
		{
			event.POST("/")         // create new event with steps
			event.GET("/:id")       // get event by id with all data
			event.GET("/")          // all user's event
			event.PUT("/:id")       // update event
			event.GET("/score/:id") // get score
			event.POST("/invite")
			//
			//	step := event.Group("/step")
			//	{
			//		step.PUT("/:id")       // update step in event
			//		step.GET("/:id")       // get all step info
			//		step.GET("/steps/:id") // get all steps in event
			//		step.POST("/")         // create new step
			//		step.GET("/users/:id") // get all users in step
			//		step.GET("/prizes")
			//	}
		}
		//prize := api.Group("/prize")
		//{
		//	prize.POST("/")         // create new event with steps
		//	prize.GET("/:id")       // get event by id with all data
		//	prize.GET("/")          // all user's prizes
		//	prize.PUT("/:id")       // update event
		//	prize.POST("/give/:id") // give prize to user
		//}
	}

	return router
}
