package v1

import (
	"errors"
	"github.com/cookienyancloud/back/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)



func (h *Handler) initUsersRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.POST("/sign-up", h.signUp)
		users.POST("/sign-in", h.signIn)
		users.POST("/refresh", h.refresh)

		authenticated := users.Group("/", h.userIdentity)
		{
			authenticated.POST("/verify/:code", h.userVerify)
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
	user, err := h.services.Users.GetUserEmail(c, id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	//todo:users events and zones
	c.JSON(http.StatusOK, userInfoResponse{
		user,
	})

}
