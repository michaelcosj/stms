package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/michaelcosj/stms/models"
)

type userAuthRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *handler) Login(c echo.Context) error {
	req := new(userAuthRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusInternalServerError, newResponseError(err.Error()))
	}

	// TODO: validate fields
	user := models.User{Name: req.Username, Email: req.Email, Password: req.Password}
	h.userRepo.NewUser(user)

	data := map[string]interface{}{"user": user}
	return c.JSON(http.StatusCreated, newResponseSuccess(data))
}

func (h *handler) Register(c echo.Context) error {
	req := new(userAuthRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusInternalServerError, newResponseError(err.Error()))
	}

	user, err := h.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		data := map[string]interface{}{"detail": err.Error()}
		return c.JSON(http.StatusNotFound, newResponseFail(data))
	}

	if user.Password != req.Password {
		data := map[string]interface{}{"detail": "invalid email or password"}
		return c.JSON(http.StatusNotFound, newResponseFail(data))
	}

	data := map[string]interface{}{"user": user}
	return c.JSON(http.StatusOK, newResponseSuccess(data))
}
