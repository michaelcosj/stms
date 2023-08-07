package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/michaelcosj/stms/framework"
	"github.com/michaelcosj/stms/models"
	"github.com/michaelcosj/stms/repository"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

// TODO: Better error handling

var ctx = context.Background()

func (h *handler) Register(c echo.Context) error {
	req := new(registerRequest)
	data := make(map[string]interface{})

	if err := c.Bind(req); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, newErrResp(ErrHandlingRequestMsg))
	}

	if data, ok := validateRegisterRequest(*req); !ok {
		data["detail"] = ErrValidatingRequestMsg
		return c.JSON(http.StatusBadRequest, newFailResp(data))
	}

	if _, err := h.userRepo.GetUserByEmail(req.Email); err != repository.ErrUserNotFound {
		data["detail"] = ErrUserAlreadyExists
		return c.JSON(http.StatusBadRequest, newFailResp(data))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, newErrResp(ErrHandlingRequestMsg))
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	id, err := h.userRepo.NewUser(user)
	if err != nil {
		c.Logger().Error(err)
		data["detail"] = ErrRegisteringUserMsg
		return c.JSON(http.StatusBadRequest, newFailResp(data))
	}

	user.ID = id
	data["user"] = user

	code, err := framework.CreateOTP(6)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, newErrResp(ErrHandlingRequestMsg))
	}

	expiry := time.Duration(1) * time.Hour
	if err := h.cache.Set(ctx, code, user.ID, expiry).Err(); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, newErrResp(ErrHandlingRequestMsg))
	}

	emailData := framework.EmailData{
		Code:    code,
		Subject: "Email Verification",
	}

	if err := framework.SendEmail(user.Email, emailData); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, newErrResp(ErrHandlingRequestMsg))
	}

	return c.JSON(http.StatusCreated, newSuccessResp(data))
}

func (h *handler) Login(c echo.Context) error {
	req := new(loginRequest)
	data := make(map[string]interface{})

	if err := c.Bind(req); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			newErrResp(("error handling request: " + err.Error())),
		)
	}

	user, err := h.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		data["detail"] = "error finding user: " + err.Error()
		return c.JSON(http.StatusNotFound, newFailResp(data))
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		data["detail"] = "incorrect email or password"
		return c.JSON(http.StatusNotFound, newFailResp(data))
	}

	secret := os.Getenv("ACCESS_TOKEN_SECRET")
	expiry, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRY_HOUR"))
	if err != nil {
		c.Logger().Warn("failed to parse ACCESS_TOKEN_EXPIRY_HOUR")
		expiry = 2 // default access token expiry if env isn't set
	}

	token, err := framework.CreateJwtToken(user.ID, secret, expiry)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			newErrResp(("error creating auth token: " + err.Error())),
		)
	}

	data["user"] = user
	data["token"] = token
	return c.JSON(http.StatusOK, newSuccessResp(data))
}

func (h *handler) VerifyUser(c echo.Context) error {
	data := make(map[string]interface{})
	req := new(struct {
		Code string `json:"code"`
	})

	if err := c.Bind(req); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			newErrResp(("error handling request: " + err.Error())),
		)
	}

	userIdStr, err := h.cache.Get(ctx, req.Code).Result()
	if err != nil {
		if err == redis.Nil {
			data["detail"] = "verification code expired: " + err.Error()
			return c.JSON(http.StatusNotFound, newFailResp(data))
		}

		return c.JSON(
			http.StatusInternalServerError,
			newErrResp(("error handling request: " + err.Error())),
		)
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newErrResp(fmt.Sprintf("error parsing userId %s request: %s", userIdStr, err.Error())))

	}

	user, err := h.userRepo.GetUserByID(int64(userId))
	if err != nil {
		data["detail"] = "user not found: " + err.Error()
		return c.JSON(http.StatusNotFound, newFailResp(data))
	}

	if user.IsVerified {
		data["detail"] = "user already verified"
		return c.JSON(http.StatusBadRequest, newFailResp(data))
	}

	user.IsVerified = true
	user.Username = "John"

	if err := h.userRepo.UpdateUser(user.ID, user); err != nil {
		return c.JSON(http.StatusBadRequest, newErrResp("user not updated: "+err.Error()))
	}

	data["user"] = user
	return c.JSON(http.StatusOK, newSuccessResp(data))
}
