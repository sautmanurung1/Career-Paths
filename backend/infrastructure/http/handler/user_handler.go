package handler

import (
	"career-paths/dto"
	"career-paths/entities/seeds"
	"career-paths/interfaces"
	"net/http"

	"github.com/labstack/echo/v4"
)

type userHandler struct {
	s interfaces.UserService
}

func NewUserController(userService interfaces.UserService) interfaces.UserHandler {
	return &userHandler{
		s: userService,
	}
}

// LoginUser implements interfaces.UserHandler
func (h *userHandler) LoginUser(c echo.Context) error {
	user := &dto.Login{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
			"status":  false,
			"code":    http.StatusBadRequest,
		})
	}

	res, fail := h.s.LoginUserService(user)
	if fail != nil {
		return c.JSON(fail.Code, map[string]interface{}{
			"message": fail.Message,
			"status":  false,
			"code":    fail.Code,
		})
	}

	cookiesAccessToken := http.Cookie{
		Name:     "access_token",
		Value:    res.AccessToken,
		MaxAge:   int(res.ExpiredRefreshToken),
		HttpOnly: true,
		Secure:   true,
		Domain:   "localhost",
	}

	cookiesRefreshToken := http.Cookie{
		Name:     "refresh_token",
		Value:    res.RefreshToken,
		MaxAge:   int(res.ExpiredRefreshToken),
		HttpOnly: true,
		Secure:   true,
		Domain:   "localhost",
	}

	c.SetCookie(&cookiesAccessToken)
	c.SetCookie(&cookiesRefreshToken)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":    http.StatusOK,
		"message": "login successfully",
		"result":  res,
	})
}

// RegisterUser implements interfaces.UserHandler
func (h *userHandler) RegisterUser(c echo.Context) error {
	user := &dto.Register{
		RoleID: seeds.UserRoleID,
	}

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
			"status":  false,
			"code":    http.StatusBadRequest,
		})
	}

	fail := h.s.RegisterUserService(user)
	if fail != nil {
		return c.JSON(fail.Code, map[string]interface{}{
			"message": fail.Message,
			"status":  false,
			"code":    fail.Code,
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"code":    http.StatusCreated,
		"message": "register user successfully! Please check your email to activate your account",
	})
}

// VerifyEmail implements interfaces.UserHandler
func (h *userHandler) VerifyEmail(c echo.Context) error {
	email := c.Param("email")
	token := c.QueryParam("token")
	expires := c.QueryParam("expires")

	fail := h.s.VerifyEmail(email, token, expires)
	if fail != nil {
		return c.JSON(fail.Code, map[string]interface{}{
			"code":    fail.Code,
			"message": fail.Message,
			"status":  false,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":    http.StatusOK,
		"message": "verify email successfully",
		"status":  true,
	})
}
