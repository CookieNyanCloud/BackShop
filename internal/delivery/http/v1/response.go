package v1

import (
	"github.com/cookienyancloud/back/internal/domain"
	"github.com/cookienyancloud/back/pkg/logger"
	"github.com/gin-gonic/gin"
)

type allEventsResponse struct {
	Events []domain.Event `json:"events"`
}

type eventsResponse struct {
	Event domain.Event `json:"event"`
}

type takingResponse struct {
	State bool `json:"state"`
}

type dataResponse struct {
	Data interface{} `json:"data"`
}

type idResponse struct {
	ID interface{} `json:"id"`
}

type response struct {
	Message string `json:"message"`
}

type tokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type userInfoResponse struct {
	UserEmail string `json:"user_email"`
	Zones []domain.Zone `json:"zones"`
}

func newResponse(c *gin.Context, statusCode int, message string) {
	logger.Error(message)
	c.AbortWithStatusJSON(statusCode, response{message})
}
