package handlers

import (
	"net/smtp"
	"fmt"
	"strings"
)

func (apiCfg *ApiConfig) sendEmail(from, to, subject, body string) error {
	var host string
	if apiCfg.Env == "dev" {
		host = "sandbox.smtp.mailtrap.io"
	}

	auth := smtp.PlainAuth(
		"",
		apiCfg.MailtrapUsername,
		apiCfg.MailtrapPassword,
		host,
	)

	msg := []byte(fmt.Sprintf(
    "From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
    from, to, subject, body,
	))

	return smtp.SendMail(
		"sandbox.smtp.mailtrap.io:587",
		auth,
		from,
		[]string{to},
		msg,
	)
}


func (apiCfg *ApiConfig) sendPasswordRecoveryEmail(magicLinkToken, to string) error {
	link := fmt.Sprintf("http://localhost:5173/resetPassword?token=%s", magicLinkToken)
	subject := "Password reset"
	body := fmt.Sprintf(`Hi,
	
	We received a request to reset your password. Click the link below to reset it:

	%s
	
	This link will expire in 1 hour.

	If you didn't request this, you can safely ignore this email.

	The Lighthouse team.`, link)
	
	body = strings.ReplaceAll(body, "\n", "\r\n")

	err := apiCfg.sendEmail(apiCfg.From, to, subject, body)
	if err != nil {
		return err
	}

	return nil
}