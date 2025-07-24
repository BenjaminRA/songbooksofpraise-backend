package email

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail(to string, subject string, msg string) error {
	from := os.Getenv("MAIL_USERNAME")

	msg = fmt.Sprintf("From: Songbooks of Praise\nTo: %s\nSubject: %s\n\n%s", to, subject, msg)

	err := smtp.SendMail(
		fmt.Sprintf("%s:587", os.Getenv("MAIL_HOST")),
		smtp.PlainAuth(
			"",
			os.Getenv("MAIL_USERNAME"),
			os.Getenv("MAIL_PASSWORD"),
			os.Getenv("MAIL_HOST"),
		),
		from,
		[]string{to},
		[]byte(msg),
	)

	if err != nil {
		return err
	}

	return nil
}
