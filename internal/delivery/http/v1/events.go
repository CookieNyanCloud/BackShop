package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) initEventsRoutes(api *gin.RouterGroup) {
	Events := api.Group("/events", )
	{
		Events.GET("/all", h.getAllEvents)
		Events.GET("/first", h.getFirstEvent)
		Events.GET("/:idevent", h.getEventById)
	}
}

func (h *Handler) getAllEvents(c *gin.Context) {
	events, err := h.services.Events.GetEvents(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	for i := range events {
		zones, err := h.services.Zones.GetZonesByEventId(c,events[i].Id)
		if err != nil {
			newResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		events[i].Zones = zones
	}

	c.JSON(http.StatusOK, allEventsResponse{events})
}

func (h *Handler) getEventById(c *gin.Context) {
	idEventStr := c.Param("idevent")
	idEvent, err := strconv.Atoi(idEventStr)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	event, err := h.services.Events.GetEventById(c,idEvent)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	zones, err := h.services.Zones.GetZonesByEventId(c,event.Id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	event.Zones = zones
	c.JSON(http.StatusOK, eventsResponse{event})

}

func (h *Handler) getFirstEvent(c *gin.Context) {
	event, err := h.services.Events.GetFirstEvent(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	zones, err := h.services.Zones.GetZonesByEventId(c,event.Id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	event.Zones = zones
	c.JSON(http.StatusOK, eventsResponse{event})

}