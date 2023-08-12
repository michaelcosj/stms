package handlers

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/michaelcosj/stms/app"
	"github.com/michaelcosj/stms/framework"
)

type handler struct {
	app app.App
}

type Handler interface {
	Login(c echo.Context) error
	Register(c echo.Context) error
	AddTask(c echo.Context) error
	UpdateTask(c echo.Context) error
	GetTasks(c echo.Context) error
	RemoveTask(c echo.Context) error
	VerifyUser(c echo.Context) error
	StartVerification(c echo.Context) error
}

// TODO: use [https://echo.labstack.com/docs/error-handling]

func InitHandler(app app.App) Handler {
	return &handler{app}
}

type response struct {
	Status string                 `json:"status"`
	Data   map[string]interface{} `json:"data"`
}

type errorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func newSuccessResp(data map[string]interface{}) response {
	return response{Status: "success", Data: data}
}

func newFailResp(data map[string]interface{}) response {
	return response{Status: "fail", Data: data}
}

func newErrResp(message string, err error) errorResponse {
	return errorResponse{Status: "error", Message: fmt.Sprintf("%s: %s", message, err)}
}

func getAuthUserId(c echo.Context) int64 {
	token := c.Get("user").(*jwt.Token)
	return token.Claims.(*framework.CustomClaims).UserID
}
