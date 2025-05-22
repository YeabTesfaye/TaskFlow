package utils

import (
	"os"

	"gopkg.in/gomail.v2"
)

func SendEmail(to string, subject string, htmlBody string) error {
	from := os.Getenv("APP_EMAIL")
	password := os.Getenv("APP_EMAIL_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := 587 // default fallback

	// Construct the email
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	// Dial and send
	d := gomail.NewDialer(smtpHost, smtpPort, from, password)
	return d.DialAndSend(m)
}
