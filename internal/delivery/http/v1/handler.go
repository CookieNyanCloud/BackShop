package v1

import (
	"errors"
	"github.com/cookienyancloud/back/internal/service"
	"github.com/cookienyancloud/back/pkg/auth"
	"github.com/gin-gonic/gin"
)

var (
	errUserAlreadyExists       = errors.New("user with such email already exists")
	errUserNotFound            = errors.New("user doesn't exists")
	errVerificationCodeInvalid = errors.New("verification code is invalid")
)

const (
	noId            = "no id"
	noCode          = "code is empty"
	needAccount     = "sign-in"
	errInvalidInput = "invalid input body"
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
		h.initEventsRoutes(v1)
		h.initZonesRoutes(v1)
		h.initCallbackRoutes(v1)
		h.initAdminRoutes(v1)
	}
}
