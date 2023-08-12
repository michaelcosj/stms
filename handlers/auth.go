package handlers

// TODO: move error response to model

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func (h *handler) Register(c echo.Context) error {
	req := new(registerRequest)
	data := make(map[string]interface{})

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusInternalServerError, newErrResp("error handling request", err))
	}

	user, err := h.app.NewUser(req.Username, req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newErrResp("error registering user", err))
	}

	data["user"] = user
	return c.JSON(http.StatusCreated, newSuccessResp(data))
}

func (h *handler) StartVerification(c echo.Context) error {
	data := make(map[string]interface{})

	user_email := c.QueryParam("email")
	if err := h.app.SendVerificationCode(user_email); err != nil {
		return c.JSON(http.StatusInternalServerError, newErrResp("error sending verification code", err))
	}

	exp_hrs, err := strconv.Atoi(os.Getenv("OTP_EXPIRY_HOURS"))
	if err != nil {
		exp_hrs = 1 // default expiry if env isn't set
	}
	expiry := time.Now().Add(time.Duration(exp_hrs) * time.Hour)

	data["detail"] = fmt.Sprintf("code sent to email %s. Expires in %s", user_email, expiry.Format(time.ANSIC))
	return c.JSON(http.StatusOK, newSuccessResp(data))
}

func (h *handler) VerifyUser(c echo.Context) error {
	data := make(map[string]interface{})
	req := new(struct {
		Code string `json:"code"`
	})

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusInternalServerError, newErrResp("error handling request", err))
	}

	user, err := h.app.VerifyUser(req.Code)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newErrResp("error verifying user", err))
	}

	data["user"] = user
	return c.JSON(http.StatusOK, newSuccessResp(data))
}

func (h *handler) Login(c echo.Context) error {
	req := new(loginRequest)
	data := make(map[string]interface{})

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusInternalServerError, newErrResp("error handling request", err))
	}

	user, token, err := h.app.GetUser(req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newErrResp("error signing in user: %v", err))
	}

	data["user"] = user
	data["token"] = token

	return c.JSON(http.StatusOK, newSuccessResp(data))
}
