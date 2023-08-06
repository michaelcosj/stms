package handlers

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/michaelcosj/stms/framework"
	"github.com/michaelcosj/stms/repository"
)

type handler struct {
	userRepo repository.UserRepo
}

type Handler interface {
	Login(c echo.Context) error
	Register(c echo.Context) error
	AddTask(c echo.Context) error
	UpdateTask(c echo.Context) error
	GetTasks(c echo.Context) error
	RemoveTask(c echo.Context) error
	VerifyUser(c echo.Context) error
}

func InitHandler(userRepo repository.UserRepo) Handler {
	return &handler{userRepo}
}

type response struct {
	Status string                 `json:"status"`
	Data   map[string]interface{} `json:"data"`
}

func newSuccessResponse(data map[string]interface{}) response {
	return response{Status: "success", Data: data}
}

func newFailResponse(data map[string]interface{}) response {
	return response{Status: "fail", Data: data}
}

type errorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func newErrorResponse(message string) errorResponse {
	return errorResponse{Status: "error", Message: message}
}

func getUserIdFromContext(c echo.Context) int64 {
	token := c.Get("user").(*jwt.Token)
	return token.Claims.(*framework.CustomClaims).UserID
}
