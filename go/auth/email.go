package auth

import (
	"context"
	"fmt"
	"log"
	"math/rand" // For GenerateVerificationCode
	"net/smtp"
	"os"
	"time" // For GenerateVerificationCode and Redis expiration

	"go-app-base/config" // For Redis client
)

// EmailSenderInterface defines the methods for sending emails.
type EmailSenderInterface interface {
	SendEmail(to, subject, body string) error
}

// defaultEmailSender is the default implementation of EmailSenderInterface.
type defaultEmailSender struct{}

// SendEmail sends an email using SMTP.
func (d *defaultEmailSender) SendEmail(to string, subject string, body string) error {
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

// EmailSender is a global variable that holds the current email sender implementation.
// It can be replaced with a mock for testing.
var EmailSender EmailSenderInterface = &defaultEmailSender{}

// SendEmail is the public function to send emails, which delegates to the global EmailSender.
func SendEmail(to string, subject string, body string) error {
	return EmailSender.SendEmail(to, subject, body)
}

// GenerateVerificationCode generates a random 6-digit verification code.
func GenerateVerificationCode() (string, error) {
	rand.Seed(time.Now().UnixNano())
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	return code, nil
}

// StoreVerificationCode stores the verification code in Redis with an expiration.
func StoreVerificationCode(email, code string, expiration time.Duration) error {
	ctx := context.Background()
	err := config.RDB.Set(ctx, "verification:"+email, code, expiration).Err()
	if err != nil {
		log.Printf("Failed to store verification code for %s in Redis: %v", email, err)
		return err
	}
	return nil
}

// SendVerificationEmail sends a verification email to the user.
func SendVerificationEmail(to, code string) error {
	subject := "Verify your email address"
	appURL := os.Getenv("APP_URL")
	verificationLink := fmt.Sprintf("%s/auth/verify?email=%s&code=%s", appURL, to, code)
	body := fmt.Sprintf("Your verification code is: %s\n\nThis code will expire in 5 minutes.\n\nOr click the link to verify: %s", code, verificationLink)
	return SendEmail(to, subject, body)
}
