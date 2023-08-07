package framework

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"net/mail"
)

func CreateOTP(maxDigits int) (string, error) {
	code, err := rand.Int(
		rand.Reader,
		big.NewInt(int64(math.Pow(10, float64(maxDigits)))),
	)

	if err != nil {
		return "", fmt.Errorf("error generating code:%v", err)
	}
	return fmt.Sprintf("%0*d", maxDigits, code), nil
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsValidUsername(name string) bool {
	return (len(name) > 3 && name[0] != ' ')
}

// TODO: change to something more secure
func IsValidPassword(pwd string) bool {
	return len(pwd) >= 8
}
