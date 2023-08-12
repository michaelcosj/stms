package app

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/michaelcosj/stms/framework"
	"github.com/michaelcosj/stms/models"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

func (a *app) NewUser(username, email, password string) (models.User, error) {
	switch {
	case !framework.IsValidEmail(email):
		return models.User{}, fmt.Errorf("invalid email")
	case !framework.IsValidPassword(password):
		return models.User{}, fmt.Errorf("invalid password")
	case !framework.IsValidUsername(username):
		return models.User{}, fmt.Errorf("invalid username")
	case a.repo.UserEmailExists(email):
		return models.User{}, fmt.Errorf("user does not exist")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, fmt.Errorf("error hashing password: %v", err)
	}

	user := models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	id, err := a.repo.NewUser(user)
	if err != nil {
		return models.User{}, fmt.Errorf("error inserting user to database: %v", err)
	}

	user.ID = id
	return user, nil
}

func (a *app) SendVerificationCode(email string) error {
	if !framework.IsValidEmail(email) {
		return fmt.Errorf("invalid email")
	}

	code, err := framework.CreateOTP(4)
	if err != nil {
		return fmt.Errorf("error creating code: %v", err)
	}

	exp_hrs, err := strconv.Atoi(os.Getenv("OTP_EXPIRY_HOURS"))
	if err != nil {
		exp_hrs = 1 // default expiry if env isn't set
	}
	expiry := time.Duration(exp_hrs) * time.Hour

	if err := a.cache.Set(ctx, code, email, expiry).Err(); err != nil {
		return fmt.Errorf("error caching code: %v", err)
	}

	emailData := framework.EmailData{
		Code:    code,
		Subject: "Email Verification",
	}

	if err := framework.SendEmail(email, emailData); err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}

	return nil
}

func (a *app) VerifyUser(code string) (models.User, error) {
	email, err := a.cache.Get(ctx, code).Result()
	if err != nil {
		if err == redis.Nil {
			return models.User{}, fmt.Errorf("code expired or invalid")
		}
		return models.User{}, fmt.Errorf("error getting user from cache: %v", err)
	}

	user, err := a.repo.GetUserByEmail(email)
	if err != nil {
		return models.User{}, fmt.Errorf("error getting user from database: %v", err)
	}

	if user.IsVerified {
		return models.User{}, fmt.Errorf("user already verified")
	}

	user.IsVerified = true
	if err := a.repo.UpdateUser(user.ID, user); err != nil {
		return models.User{}, fmt.Errorf("error updating user from database: %v", err)
	}

	return user, nil
}

func (a *app) GetUser(email, password string) (models.User, string, error) {
	user, err := a.repo.GetUserByEmail(email)
	if err != nil {
		return models.User{}, "", fmt.Errorf("error getting user from database: %v", err)
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return models.User{}, "", fmt.Errorf("invalid email or password")
	}

	secret := os.Getenv("ACCESS_TOKEN_SECRET")
	expiry, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRY_HOUR"))
	if err != nil {
		expiry = 2 // default access token expiry if env isn't set
	}

	token, err := framework.CreateJwtToken(user.ID, secret, expiry)
	if err != nil {
		return models.User{}, "", fmt.Errorf("error creating jwt token: %v", err)
	}

	return user, token, nil
}
