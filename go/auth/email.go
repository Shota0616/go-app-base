package auth

import (
	"fmt"
	"net/smtp"
	"os"
	"log"
)

func SendEmail(to string, subject string, body string) error {
	from := os.Getenv("EMAIL_FROM") // Use a more generic name like EMAIL_FROM
	password := os.Getenv("EMAIL_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	// メールのヘッダーと本文を構築
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, to, subject, body)

	var auth smtp.Auth
	// If username and password are provided, use PlainAuth
	if from != "" && password != "" {
		auth = smtp.PlainAuth("", from, password, smtpHost)
	}

	log.Printf("Sending email to %s via %s:%s", to, smtpHost, smtpPort)

	// メール送信
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(message))
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	log.Println("Email sent successfully")
	return nil
}
