package email

import (
	"fmt"
	"net/smtp"
	"os"
)

// SendResetEmail sends a password reset email to the specified address with the given token.
func SendResetEmail(to, token string) error{

	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")

	from := "no-reply@digitalcontract.com"
	subject := "Subject: Password Reset Request\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf("<h1>Reset Your Password</h1><p>Use this token to reset your password: <b>%s</b></p>", token)

	msg := []byte(subject + mime + "\n" + body)
	addr := fmt.Sprintf("%s:%s", host, port)

	//In Mailhog(dev), we don't need authentication so we pass nil
	return smtp.SendMail(addr, nil, from, []string{to}, msg)
}