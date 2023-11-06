package email

import (
	"fmt"
	"net/smtp"
	"os"
)

// Sends email using SMTP
func SendEmail(recipient, subject, body string) error {
	// Set up authentication information.
	smtpHost := os.Getenv("SMTP_HOST") // SMTP server host
	smtpPort := os.Getenv("SMTP_PORT") // SMTP server port, often 587 for STARTTLS or 465 for SSL
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	from := "support@telukbirukarya.com"

	// Usually, the auth will be PlainAuth for SMTP servers requiring authentication.
	auth := smtp.PlainAuth("", username, password, smtpHost)

	// Set MIME and other headers
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = recipient
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"utf-8\""

	// Build the headers
	header := ""
	for k, v := range headers {
		header += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	// The msg parameter should be an RFC 822-style email with headers first,
	// a blank line, and then the message body.
	msg := []byte(
		header + "\r\n" +
			body + "\r\n")
	// []byte(
	// 	"To: " + recipient + "\r\n" +
	// 		"Subject: " + subject + "\r\n" +
	// 		"\r\n" +
	// 		body + "\r\n")

	// Combine host and port for the smtp.SendMail() call
	addr := smtpHost + ":" + smtpPort

	// This sends the email with a plain auth setup
	err := smtp.SendMail(addr, auth, from, []string{recipient}, msg)
	if err != nil {
		return fmt.Errorf("smtp.SendMail() failed with: %s", err)
	}
	fmt.Println("Email sent successfully!")
	return nil
}
