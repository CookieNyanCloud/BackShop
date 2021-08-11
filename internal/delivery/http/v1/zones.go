package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) initZonesRoutes(api *gin.RouterGroup) {

	Zones := api.Group("/zones", h.userIdentity)
	{
		//Zones.GET("",h.getAllZones)
		Zones.GET("/",h.getAllZones)
		Zones.GET("/:idevent/:id",h.takeZone)
	}
}

func (h *Handler) getAllZones(c *gin.Context) {
	events, err := h.services.Events.GetEvent()
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	zones, err := h.services.Zones.GetZonesByEventId(events[0].Id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	userId, _:= getUserId(c)
	println("tyt",userId)
	c.JSON(http.StatusOK, dataResponse{zones})
	c.Set(userCtx, userId)
}


func (h *Handler) takeZone(c *gin.Context) {

	idEventStr:= c.Param("idevent")
	idZoneStr:= c.Param("id")
	userId, _:= getUserId(c)


	if idEventStr == "" || idZoneStr== ""|| userId == 0{
		newResponse(c, http.StatusBadRequest, "empty id param")
		return
	}
	idEventInt, err:= strconv.Atoi(idEventStr)
	if err != nil {
		newResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	idZoneInt, err:= strconv.Atoi(idZoneStr)
	if err != nil {
		newResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	zones, err := h.services.Zones.TakeZoneById(idEventInt,idZoneInt,userId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, dataResponse{zones})
	c.Set(userCtx, userId)
	return
}