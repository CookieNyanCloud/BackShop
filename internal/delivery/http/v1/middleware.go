package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

const (
	authorizationHeader = "authorization"
	userCtx = "userId"
	adminCtx   = "adminId"
	jwt = "jwt"
)

func getUserId(c*gin.Context) (int, error)  {
	return getIdByContext(c, userCtx)
}

func (h *Handler) userIdentity (c *gin.Context)  {
	id, err := h.parseAuthHeader(c)
	if err != nil {
		newResponse(c, http.StatusUnauthorized, err.Error())
	}
	c.Set(userCtx, id)
}

func (h *Handler) adminIdentity(c *gin.Context) {
	id, err := h.parseAuthHeader(c)
	if err != nil {
		newResponse(c, http.StatusUnauthorized, err.Error())
	}

	c.Set(adminCtx, id)
}

func (h *Handler) parseAuthHeader(c *gin.Context) (string, error) {
	header := c.GetHeader(authorizationHeader)
	println("A111111111111111111", header)

	println(header)
	if header == ""{
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

func getIdByContext(c *gin.Context, context string) ( int, error) {
	str:= c.GetString(context)
	strIntId, err:= strconv.Atoi(str)
	if err!=nil {
		return 0, errors.New("studentCtx not found")
	}
	return strIntId, nil
}
