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
			user.PUT("/:id", h.UpdateStaffByID)       // update done
			user.DELETE("/:id", h.DeleteStaff)        // delete done
			user.POST("/", h.CreateStaff)             // done
			user.GET("/prizes/:id", h.GetStaffPrizes) // get all prizes
			user.GET("/invites", h.GetStaffInvites)   // get users invites done
			user.PUT("/photo", h.UploadImage)         // done
			user.GET("/image/:id", h.GetImage)        // done

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
			event.GET("/score/:id", h.GetStaffScore)          // get score done

			step := event.Group("/step")
			{
				step.PUT("/:id", h.UpdateStep)           // update step in event done
				step.GET("/:id", h.GetStep)              // get all step info done
				step.DELETE("/:id", h.DeleteStep)        // get all step info done
				step.GET("/steps/:id", h.GetSteps)       // get all steps in event done
				step.POST("/", h.CreateStep)             // create new step done
				step.GET("/prizes/:id", h.GetStepPrizes) // get step prizes done
				step.PUT("/status/:id", h.PassStaff)     // staff pass step done
				step.PUT("/assign/:id", h.AssignStaff)   // add staff to step done
			}
		}
		prize := api.Group("/prize")
		{
			prize.POST("/", h.CreatePrize)              // create new event with steps done
			prize.GET("/:id", h.GetPrize)               // get event by id with all data done
			prize.GET("/", h.GetPrizes)                 // all user's prizes done
			prize.GET("/user/:type", h.GetPrizesByType) // all user's prizes done
			prize.PUT("/:id", h.UpdatePrize)            // update prize done
			prize.POST("/give/:id", h.GivePrize)        // give prize to user done
		}
	}

	return router
}
