package framework

import "fmt"

type EmailData struct {
	Code    string
	Subject string
}

func SendEmail(userEmail string, emailData EmailData) error {
	// TODO: integrate a mailing service
	fmt.Println(emailData)
	return nil
}
