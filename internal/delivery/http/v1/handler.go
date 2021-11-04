package v1

import (
	"github.com/cookienyancloud/back/internal/service"
	"github.com/cookienyancloud/back/pkg/auth"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services     *service.Services
	tokenManager auth.TokenManager
}

func NewHandler(services *service.Services, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		services:     services,
		tokenManager: tokenManager,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initUsersRoutes(v1)
		h.initZonesRoutes(v1)
		h.initEventsRoutes(v1)
		h.initCallbackRoutes(v1)
		h.initAdminRoutes(v1)
		h.initCallbackRoutes(v1)
		h.initAdminRoutes(v1)
	}
}


