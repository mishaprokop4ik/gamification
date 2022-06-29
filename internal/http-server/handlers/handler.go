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
		auth.POST("/sign-up", h.signUp) // done
		auth.POST("/sign-in", h.signIn) // done
	}
	api := router.Group("/api", h.identity) // done
	{
		organization := api.Group("/org")
		{
			organization.GET("/", h.GetAllOrganizations)               // get all done
			organization.GET("/:id", h.GetOrganization)                // get by id done
			organization.PUT("/:id", h.UpdateOrganization)             // update org / done
			organization.PUT("/staff/:id", h.AddStaffToOrganization)   // to check done
			organization.DELETE("/:id", h.DeleteOrganization)          // delete done
			organization.GET("/staff/:id", h.GetStaffByOrganizationID) // to check done
			organization.GET("/event/:id", h.GetOrganizationEvents)    // to check
			organization.POST("/", h.CreateOrganization)               // get all events in org  done

			organizationType := organization.Group("/type")
			{
				organizationType.POST("/", h.CreateOrganizationType)      //create new
				organizationType.GET("/", h.GetOrganizationTypes)         // get all done
				organizationType.GET("/:id", h.GetOrganizationTypeByID)   // get by id done
				organizationType.PUT("/:id", h.UpdateOrganizationType)    // update done
				organizationType.DELETE("/:id", h.DeleteOrganizationType) // delete done
			}
		}

		team := api.Group("/team")
		{
			team.POST("/", h.CreateTeam)                     // create new done
			team.GET("/org/:id", h.GetTeamsByOrganizationID) // get all by company done
			team.GET("/event/:id", h.GetTeamsByEventID)      // done
			team.GET("/:id", h.GetTeamByID)                  // get by id done
			team.PUT("/:id", h.UpdateTeam)                   // update done
			team.DELETE("/:id", h.DeleteTeamByID)            // delete done
		}

		user := api.Group("/user")
		{
			user.GET("/event/:id", h.GetAllUsersInEvent) // all user in event with their status done
			user.GET("/step/:id", h.GetAllUsersInStep)   // all user in step with their status to check
			user.GET("/:id", h.GetStaffByID)             // get by id done
			user.PUT("/:id", h.UpdateStaffByID)          // update done
			user.DELETE("/:id", h.DeleteStaff)           // delete done
			user.POST("/", h.CreateStaff)                // done
			user.GET("/prizes/:id", h.GetStaffPrizes)    // get all prizes // to check
			user.GET("/invites", h.GetStaffInvites)      // get users invites done
			user.PUT("/photo", h.UploadImage)            // done
			user.GET("/image/:id", h.GetImage)           // done

			position := user.Group("/position")
			{
				position.PUT("/:id", h.UpdatePosition)               // update done
				position.PUT("/perm/:id", h.RemovePermissions)       // remove done
				position.PUT("/give/:id", h.GivePosition)            // give done
				position.PUT("/take/:id", h.TakePosition)            // give done
				position.POST("/", h.CreatePosition)                 // create done
				position.DELETE("/:id", h.DeletePosition)            // delete done
				position.GET("/org/:id", h.GetOrganizationPositions) // view comp position done
				position.GET("/:id", h.GetPosition)                  // get position by id done
			}
		}

		event := api.Group("/event")
		{
			event.POST("/", h.CreateEvent)                    // create new event with steps done
			event.POST("/invite/:id", h.AssignStaffToEvent)   // create new event with steps done
			event.POST("/invitation/:id", h.AnswerInvitation) // done
			event.GET("/invitation/", h.GetInvitations)       // done
			event.GET("/:id", h.GetEventByID)                 // get event by id with all data done
			event.GET("/staff/:role", h.GetUserEvents)        // all user's event done
			event.GET("/team/:id", h.GetTeamEvents)           // all user's event done
			event.PUT("/:id", h.UpdateEvent)                  // update event done
			event.GET("/score/:id", h.GetStaffScore)          // get score to check

			step := event.Group("/step")
			{
				step.PUT("/:id", h.UpdateStep)           // update step in event
				step.GET("/:id", h.GetStep)              // get all step info
				step.DELETE("/:id", h.DeleteStep)        // get all step info
				step.GET("/steps/:id", h.GetSteps)       // get all steps in event
				step.POST("/", h.CreateStep)             // create new step
				step.GET("/prizes/:id", h.GetStepPrizes) // get step prizes
				step.PUT("/status/:id", h.PassStaff)     // staff pass step
				step.PUT("/assign/:id", h.AssignStaff)   // add staff to step
			}
		}
		prize := api.Group("/prize")
		{
			prize.POST("/", h.CreatePrize)              // create new event with steps
			prize.GET("/:id", h.GetPrize)               // get event by id with all data
			prize.GET("/", h.GetPrizes)                 // all user's prizes
			prize.GET("/user/:type", h.GetPrizesByType) // all user's prizes
			prize.PUT("/:id", h.UpdatePrize)            // update event
			prize.POST("/give/:id", h.GivePrize)        // give prize to user
		}
	}

	return router
}
