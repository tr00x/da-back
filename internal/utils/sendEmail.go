package utils

import (
	"net/smtp"
)

func SendEmail(title, subject, receiver string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	from := "contact.dubai.auto@gmail.com"
	password := "sloh hauz bouf smlq"
	message := []byte("Subject: " + title + "\r\n" +
		"\r\n" +
		subject + "\r\n")

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{receiver}, message)
	return err
}
