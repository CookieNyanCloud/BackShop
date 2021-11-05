package v1

import (
	"errors"
	"github.com/cookienyancloud/back/internal/domain"
	"github.com/cookienyancloud/back/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type createEventInput struct {
	Time        time.Time     `json:"time" db:"time"`
	Description string        `json:"description" db:"description"`
	MapFile     string        `json:"mapfile" db:"mapfile"`
	Zones       []domain.Zone `json:"zones" db:"zones"`
}

func (h *Handler) initAdminRoutes(api *gin.RouterGroup) {
	admins := api.Group("/admin")
	{
		admins.POST("/sign-in", h.adminSignIn)
		admins.POST("/auth/refresh", h.adminRefresh)
		authenticated := admins.Group("/", h.adminIdentity)
		{
			authenticated.POST("/create-event", h.adminCreateEvent)
			//authenticated.DELETE("/:id", h.adminDeleteEvent)
			//authenticated.PUT("/:id", h.adminChangeEvent)
		}
	}
}

func (h *Handler) adminSignIn(c *gin.Context) {
	var inp signInInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}
	res, err := h.services.Admins.SignIn(c.Request.Context(), service.UserSignInInput{
		Email:    inp.Email,
		Password: inp.Password,
	})
	if err != nil {
		if errors.Is(err, errUserNotFound) {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}

func (h *Handler) adminRefresh(c *gin.Context) {
	var inp refreshInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}
	res, err := h.services.Admins.RefreshTokens(c.Request.Context(), inp.Authorization)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}

func (h *Handler) adminCreateEvent(c *gin.Context) {
	var inp createEventInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}
	err := h.services.Admins.CreateEvent(c.Request.Context(), service.CreateEventInput{
		Time:        inp.Time,
		Description: inp.Description,
		MapFile:     inp.MapFile,
		Zones:       inp.Zones,
	})
	if err != nil {
		newResponse(c, http.StatusBadRequest, "cold not create event")
		return
	}
	c.JSON(http.StatusCreated, idResponse{"success"})

}
