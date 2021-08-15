package v1

import (
	"errors"
	"github.com/cookienyancloud/back/internal/repository"
	"github.com/cookienyancloud/back/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) initUsersRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.POST("/sign-up", h.userSignUp)
		users.POST("/sign-in", h.userSignIn)
		users.POST("/auth/refresh", h.userRefresh)
		//users.get("/auth/refresh", h.userRefresh)

		authenticated := users.Group("/", h.userIdentity)
		{
			authenticated.POST("/verify/:code", h.userVerify)
			own := authenticated.Group("/own")
			{
				own.GET("/info", h.userGetOwnInfo)
				//schools.PUT("/:id", h.userUpdateZone)
			}
		}
	}
}

type userSignUpInput struct {
	Name     string `json:"name" binding:"required,min=2,max=64"`
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type signInInput struct {
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type refreshInput struct {
	Authorization string `header:"refreshToken" json:"refreshToken" binding:"required"`
	ContentType string `header:"ContentType" json:"ContentType"`
}

func (h *Handler) userSignUp(c *gin.Context) {
	var inp userSignUpInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}
	if err := h.services.Users.SignUp(c.Request.Context(), service.UserSignUpInput{
		Name:     inp.Name,
		Email:    inp.Email,
		Password: inp.Password,
	}); err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	res, err := h.services.Users.SignIn(c.Request.Context(), service.UserSignInInput{
		Email:    inp.Email,
		Password: inp.Password,
	})
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.SetCookie(
		jwt,
		res.AccessToken,
		3600, "/", "localhost", true, true)

	c.JSON(http.StatusCreated, tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
	println(res.AccessToken)
	println(res.RefreshToken)

}

func (h *Handler) userSignIn(c *gin.Context) {
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
		if errors.Is(err, service.ErrUserNotFound) {
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

func (h *Handler) userGetOwnInfo(c *gin.Context) {
	println("userinfo1")
	userId, _ := getUserId(c)
	println("userinfo2", userId)
	user, err := h.services.Users.GetUserInfo(c, userId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	println("userinfo1", user.ID)
	println("Sdsdsdsd")
	println(user.Email)
	c.JSON(http.StatusOK, userInfoResponse{
		user,
	})

}

func (h *Handler) userRefresh(c *gin.Context) {
	println("1")
	var inp refreshInput
	authorizationHeader:= c.GetHeader(refreshToken)
	println("asas",authorizationHeader)
	if err := c.BindHeader(&inp); err != nil {
	//if inp.refreshToken =c.GetHeader("refreshToken");inp.refreshToken ==""{
		println(err.Error())
		newResponse(c, http.StatusBadRequest, "invalid input body")
		println("2")
		println(inp.Authorization)
		return
	}
	println(inp.Authorization)
	res, err := h.services.Users.RefreshTokens(c.Request.Context(), inp.Authorization)
	//res, err := h.services.Users.RefreshTokens(c.Request.Context(), inpH)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		println("3")
		return
	}
	println("4")
	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}

func (h *Handler) userVerify(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		newResponse(c, http.StatusBadRequest, "code is empty")
		return
	}
	id, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	idInt := strconv.Itoa(id)
	if err := h.services.Users.Verify(c.Request.Context(), idInt, code); err != nil {
		if errors.Is(err, repository.ErrVerificationCodeInvalid) {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, response{"success"})
}
