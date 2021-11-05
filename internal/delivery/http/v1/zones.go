package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) initZonesRoutes(api *gin.RouterGroup) {

	Zones := api.Group("/zones", h.userIdentity)
	{
		Zones.POST("/take",h.takeZone)
	}
}

type takeZonesInput struct {
	EventId int `json:"event_id"`
	ZonesId    []int `json:"zones_id"`
}

func (h *Handler) takeZone(c *gin.Context) {
	var inp takeZonesInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, errInvalidInput)
		return
	}
	userId:= getUserId(c)

	if userId == "" {
		newResponse(c, http.StatusBadRequest, needAccount)
		return
	}

	_, err := h.services.Zones.TakeZonesById(c,inp.EventId,inp.ZonesId,userId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, takingResponse{true})
	return
}