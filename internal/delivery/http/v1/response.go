package v1

import (
	"github.com/cookienyancloud/back/pkg/logger"
	"github.com/gin-gonic/gin"
)

type dataResponse struct {
	Data interface{} `json:"data"`
	//Id   int         `json:"id"`
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
	UserInfo interface{} `json:"userInfo"`

}


func newResponse(c *gin.Context, statusCode int, message string) {
	logger.Error(message)
	c.AbortWithStatusJSON(statusCode, response{message})
}
