package framework

import "net/mail"

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
