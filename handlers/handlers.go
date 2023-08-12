package handlers

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/michaelcosj/stms/framework"
	"github.com/michaelcosj/stms/repository"
	"github.com/redis/go-redis/v9"
)

type handler struct {
	userRepo repository.UserRepo
	cache    *redis.Client
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

// TODO: use [https://echo.labstack.com/docs/error-handling]

var (
	ErrHandlingRequestMsg   = errors.New("error handling request")
	ErrValidatingRequestMsg = errors.New("error validating request")
	ErrSigningInUserMsg     = errors.New("error signing in user")
	ErrRegisteringUserMsg   = errors.New("error registering user")
	ErrUserAlreadyExistsMsg = errors.New("user already exists")
	ErrUserNotFoundMsg      = errors.New("user not found")
)

func InitHandler(userRepo repository.UserRepo, cache *redis.Client) Handler {
	return &handler{userRepo, cache}
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

func validateRegisterRequest(r registerRequest) (map[string]interface{}, bool) {
	isValid := true
	data := make(map[string]interface{})

	if !framework.IsValidEmail(r.Email) {
		isValid = false
		data["email"] = "invalid email"
	}

	if !framework.IsValidPassword(r.Password) {
		isValid = false
		data["email"] = "invalid email"
	}

	if !framework.IsValidUsername(r.Username) {
		isValid = false
		data["username"] = "invalid username"
	}

	return data, isValid
}
