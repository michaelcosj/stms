package handlers

import (
	"github.com/labstack/echo/v4"
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
}

func InitHandler(userRepo repository.UserRepo) Handler {
	return &handler{userRepo}
}

type response struct {
	Status string                 `json:"status"`
	Data   map[string]interface{} `json:"data"`
}

func newResponseSuccess(data map[string]interface{}) response {
	return response{Status: "success", Data: data}
}

func newResponseFail(data map[string]interface{}) response {
	return response{Status: "fail", Data: data}
}

type responseError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func newResponseError(message string) responseError {
	return responseError{Status: "error", Message: message}
}
