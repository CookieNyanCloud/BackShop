package v1

import (
	"errors"
	"github.com/cookienyancloud/back/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

//todo:oidc

func (h *Handler) initUsersRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.POST("/sign-up", h.signUp)
		users.POST("/sign-in", h.signIn)
		users.POST("/refresh", h.refresh)

		authenticated := users.Group("/", h.userIdentity)
		{
			authenticated.POST("/verify/:code", h.userVerify)
			authenticated.POST("/order", )
			account := authenticated.Group("/account")
			{
				account.GET("/page", h.getOwnInfo)
			}
		}
	}
}

type userSignUpInput struct {
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type signInInput struct {
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type refreshInput struct {
	Authorization string `header:"refreshToken" json:"refreshToken" binding:"required"`
}

type orderInput struct {
	orderID uuid.UUID `json:"order_id"`
	eventId int       `json:"event_id"`
	zonesId []int     `json:"zones_id"`
}

type generatePaymentLinkResponse struct {
	URL string `json:"url"`
}

func (h *Handler) signUp(c *gin.Context) {
	var inp userSignUpInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, errInvalidInput)
		return
	}
	if err := h.services.Users.SignUp(c.Request.Context(), service.UserSignUpInput{
		Email:    inp.Email,
		Password: inp.Password,
	}); err != nil {
		if errors.Is(err, errUserAlreadyExists) {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, response{"verify"})

}

func (h *Handler) signIn(c *gin.Context) {
	var inp signInInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}
	res, err := h.services.Users.SignIn(c.Request.Context(), service.UserSignInInput{
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

func (h *Handler) refresh(c *gin.Context) {
	var inp refreshInput
	if err := c.BindHeader(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, errInvalidInput)
		return
	}
	res, err := h.services.Users.RefreshTokens(c.Request.Context(), inp.Authorization)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}

func (h *Handler) userVerify(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		newResponse(c, http.StatusBadRequest, noCode)
		return
	}
	id := getUserId(c)
	if id == "" {
		newResponse(c, http.StatusInternalServerError, noId)
		return
	}
	if err := h.services.Users.Verify(c.Request.Context(), id, code); err != nil {
		if errors.Is(err, errVerificationCodeInvalid) {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, response{"success"})
}

func (h *Handler) getOwnInfo(c *gin.Context) {
	id := getUserId(c)
	userEmail, err := h.services.Users.GetUserEmail(c.Request.Context(), id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	zones, err := h.services.Zones.GetZonesByUserId(c.Request.Context(), id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	//todo:users events and zones
	c.JSON(http.StatusOK, userInfoResponse{
		userEmail,
		zones,
	})

}

func (h *Handler) studentCreateOrder(c *gin.Context) {
	id := getUserId(c)
	if id == "" {
		newResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	var inp orderInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}
	err := h.services.Orders.Create(c.Request.Context(), inp.orderID, id, inp.eventId, inp.zonesId)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	url, err := h.services.Payments.GeneratePaymentLink(c.Request.Context(), inp.orderID)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, generatePaymentLinkResponse{url})
}
