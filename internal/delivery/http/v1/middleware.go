package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "authorization"
	userCtx             = "userId"
	adminCtx            = "adminId"
	jwt                 = "jwt"
	refreshToken        = "refreshToken"
)

func getUserId(c *gin.Context) string {
	return c.GetString(userCtx)
}

func (h *Handler) userIdentity(c *gin.Context) {
	id, err := h.parseAuthHeader(c)
	if err != nil {
		println(err.Error())
		newResponse(c, http.StatusUnauthorized, err.Error())
	}
	c.Set(userCtx, id)
}

func (h *Handler) adminIdentity(c *gin.Context) {
	id, err := h.parseAuthHeader(c)
	if err != nil {
		println(err.Error())
		newResponse(c, http.StatusUnauthorized, err.Error())
	}

	c.Set(adminCtx, id)
}

func (h *Handler) parseAuthHeader(c *gin.Context) (string, error) {
	header := c.GetHeader(authorizationHeader)
	println(header)
	if header == "" {
		println("no header")
		return "", nil
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}
	if len(headerParts[1]) == 0 {
		return "", nil
	}
	return h.tokenManager.Parse(headerParts[1])
}
