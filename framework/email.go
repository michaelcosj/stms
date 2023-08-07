package framework

import "fmt"

type EmailData struct {
	Code    string
	Subject string
}

func SendEmail(userEmail string, emailData EmailData) error {
	fmt.Println(emailData)
	return nil
}
