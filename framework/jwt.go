package framework

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func CreateJwtToken(userID int64, secret string, expiry int) (string, error) {
	expTime := time.Now().Add(time.Hour * time.Duration(expiry))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &CustomClaims{
		userID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	})

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return t, nil
}
