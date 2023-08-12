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
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

var ctx = context.Background()

func (h *handler) Register(c echo.Context) error {
	req := new(registerRequest)
	data := make(map[string]interface{})

	if err := c.Bind(req); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, newErrResp(ErrHandlingRequestMsg))
	}

	if data, ok := validateRegisterRequest(*req); !ok {
		data["detail"] = ErrValidatingRequestMsg
		return c.JSON(http.StatusBadRequest, newFailResp(data))
	}

	if h.userRepo.CheckUserEmailExists(req.Email) {
		data["detail"] = ErrUserAlreadyExistsMsg
		return c.JSON(http.StatusBadRequest, newFailResp(data))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, newErrResp(ErrHandlingRequestMsg))
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	id, err := h.userRepo.NewUser(user)
	if err != nil {
		c.Logger().Error(err.Error())
		data["detail"] = ErrRegisteringUserMsg
		return c.JSON(http.StatusBadRequest, newFailResp(data))
	}

	user.ID = id
	data["user"] = user

	return c.JSON(http.StatusCreated, newSuccessResp(data))
}

func (h *handler) StartVerification(c echo.Context) error {
	data := make(map[string]interface{})

	user_email := c.QueryParam("email")
	if !framework.IsValidEmail(user_email) {
		data["email"] = "Invalid email"
		data["detail"] = ErrValidatingRequestMsg
		return c.JSON(http.StatusBadRequest, newFailResp(data))

	}

	code, err := framework.CreateOTP(4)
	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, newErrResp(ErrHandlingRequestMsg))
	}

	expiry := time.Duration(1) * time.Hour
	if err := h.cache.Set(ctx, code, user_email, expiry).Err(); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, newErrResp(ErrHandlingRequestMsg))
	}

	emailData := framework.EmailData{
		Code:    code,
		Subject: "Email Verification",
	}

	if err := framework.SendEmail(user_email, emailData); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, newErrResp(ErrHandlingRequestMsg))
	}

	data["detail"] = fmt.Sprintf("Verification email sent to %s. Expires in %d", user_email, expiry)
	return c.JSON(http.StatusOK, newSuccessResp(data))
}

func (h *handler) VerifyUser(c echo.Context) error {
	data := make(map[string]interface{})
	req := new(struct {
		Code string `json:"code"`
	})

	if err := c.Bind(req); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, newErrResp(ErrHandlingRequestMsg))
	}

	userIdStr, err := h.cache.Get(ctx, req.Code).Result()
	if err != nil {
		c.Logger().Error(err.Error())
		if err == redis.Nil {
			data["detail"] = "verification code expired"
			return c.JSON(http.StatusNotFound, newFailResp(data))
		}
		return c.JSON(http.StatusInternalServerError, newErrResp(ErrHandlingRequestMsg))
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.Logger().Error(err.Error())
		errMsg := fmt.Errorf("error parsing userId %s request: %v", userIdStr, err)
		return c.JSON(http.StatusBadRequest, newErrResp(errMsg))
	}

	user, err := h.userRepo.GetUserByID(int64(userId))
	if err != nil {
		c.Logger().Error(err.Error())
		data["detail"] = ErrUserNotFoundMsg
		return c.JSON(http.StatusNotFound, newFailResp(data))
	}

	if user.IsVerified {
		data["detail"] = "user already verified"
		return c.JSON(http.StatusBadRequest, newFailResp(data))
	}

	user.IsVerified = true
	if err := h.userRepo.UpdateUser(user.ID, user); err != nil {
		c.Logger().Error(err.Error())
		errMsg := fmt.Errorf("error updating user data: %v", err)
		return c.JSON(http.StatusInternalServerError, newErrResp(errMsg))
	}

	data["user"] = user
	return c.JSON(http.StatusOK, newSuccessResp(data))
}

func (h *handler) Login(c echo.Context) error {
	req := new(loginRequest)
	data := make(map[string]interface{})

	if err := c.Bind(req); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, newErrResp(ErrHandlingRequestMsg))
	}

	user, err := h.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		c.Logger().Error(err.Error())
		data["detail"] = ErrUserNotFoundMsg
		return c.JSON(http.StatusNotFound, newFailResp(data))
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		data["detail"] = "invalid email or password"
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
		c.Logger().Error(err.Error())
		errMsg := fmt.Errorf("error creating auth token: %v", err)
		c.JSON(http.StatusInternalServerError, newErrResp(errMsg))
	}

	data["user"] = user
	data["token"] = token

	return c.JSON(http.StatusOK, newSuccessResp(data))
}
