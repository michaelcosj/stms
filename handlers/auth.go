package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/michaelcosj/stms/framework"
	"github.com/michaelcosj/stms/models"
	"github.com/michaelcosj/stms/repository"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func validateRegisterRequest(r registerRequest) (map[string]interface{}, bool) {
	isValid := true
	data := make(map[string]interface{})

	switch {
	case !framework.IsValidEmail(r.Email):
		isValid = false
		data["email"] = "invalid email"
	case !framework.IsValidPassword(r.Password):
		isValid = false
		data["password"] = "invalid password"
	case !framework.IsValidUsername(r.Username):
		isValid = false
		data["username"] = "invalid username"
	}

	return data, isValid
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *handler) Register(c echo.Context) error {
	req := new(registerRequest)
	if err := c.Bind(req); err != nil {
		errMsg := fmt.Sprintf("error handling request: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, newErrorResponse(errMsg))
	}

	if data, ok := validateRegisterRequest(*req); !ok {
		data["detail"] = "invalid data"
		return c.JSON(http.StatusBadRequest, newFailResponse(data))
	}

	if _, err := h.userRepo.GetUserByEmail(req.Email); err != repository.ErrUserNotFound {
		data := map[string]interface{}{"detail": fmt.Sprintf("error registering user: %s", err.Error())}
		return c.JSON(http.StatusBadRequest, newFailResponse(data))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		errMsg := fmt.Sprintf("error handling request: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, newErrorResponse(errMsg))
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	id, err := h.userRepo.NewUser(user)
	if err != nil {
		data := map[string]interface{}{"detail": "error registering user: %s", "message": err.Error()}
		return c.JSON(http.StatusBadRequest, newFailResponse(data))
	}

	user.ID = id
	data := map[string]interface{}{"user": user}

	return c.JSON(http.StatusCreated, newSuccessResponse(data))
}

func (h *handler) Login(c echo.Context) error {
	req := new(loginRequest)
	if err := c.Bind(req); err != nil {
		errMsg := (fmt.Sprintf("error handling request: %s", err.Error()))
		return c.JSON(http.StatusInternalServerError, newErrorResponse(errMsg))
	}

	user, err := h.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		data := map[string]interface{}{"detail": err.Error()}
		return c.JSON(http.StatusNotFound, newFailResponse(data))
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		data := map[string]interface{}{"detail": "incorrect email or password"}
		return c.JSON(http.StatusNotFound, newFailResponse(data))
	}

	secret := os.Getenv("ACCESS_TOKEN_SECRET")
	expiry, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRY_HOUR"))
	if err != nil {
		// TODO: log this
		expiry = 2 // default access token expiry if env isn't set
	}

	token, err := framework.CreateJwtToken(user.ID, secret, expiry)
	if err != nil {
		errMsg := (fmt.Sprintf("error creating auth token: %s", err.Error()))
		c.JSON(http.StatusInternalServerError, newErrorResponse(errMsg))
	}

	data := map[string]interface{}{"user": user, "token": token}
	return c.JSON(http.StatusOK, newSuccessResponse(data))
}

func (h *handler) VerifyUser(c echo.Context) error {
	return nil
}
