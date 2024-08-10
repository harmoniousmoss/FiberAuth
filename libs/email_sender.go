package libs

import (
	"os"

	"gopkg.in/mail.v2"
)

// SendEmail sends an email with the given recipients, subject, and body
func SendEmail(to []string, subject, body string) error {
	m := mail.NewMessage()

	// Set the sender and recipient.
	m.SetHeader("From", os.Getenv("EMAIL_USER"))
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)

	// Set the body of the email to HTML
	m.SetBody("text/html", body)

	// Setup the mail server configuration. I'm using Office 365 SMTP server here.
	d := mail.NewDialer("smtp.office365.com", 587, os.Getenv("EMAIL_USER"), os.Getenv("EMAIL_PASS"))

	// Send the email.
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
