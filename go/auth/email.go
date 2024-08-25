package auth

import (
	"fmt"
	"net/smtp"
	"os"
	"log"
)

func SendEmail(to string, subject string, body string) error {
	from := os.Getenv("EMAIL_ADDRESS")
	password := os.Getenv("EMAIL_PASSWORD")

	// GmailのSMTPサーバー設定
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// メールのヘッダーと本文を構築
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, to, subject, body)

	// SMTP認証情報
	auth := smtp.PlainAuth("", from, password, smtpHost)

	log.Println(from)
	log.Println(password)
	log.Println(message)
	log.Println(auth)


	// メール送信
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(message))
	if err != nil {
		return err
	}

	return nil
}
