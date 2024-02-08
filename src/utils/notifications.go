package utils

import (
	"bytes"
	"crypto/tls"
	"log"
	"os"
	"text/template"

	"github.com/Numostanley/d8er_app/env"
	"github.com/Numostanley/d8er_app/models"
	"gopkg.in/gomail.v2"
)

func VerificationEmail(user *models.User) {
	htmlContent, err := os.ReadFile("templates/verification.html")
	if err != nil {
		log.Fatal(err)
	}
	htmlTemplate := string(htmlContent)
	tmpl, err := template.New("verification").Parse(htmlTemplate)
	if err != nil {
		log.Fatal(err)
	}

	data := struct {
		Code string
		Name string
	}{
		Code: GenerateOTP(),
		Name: user.GetFullName(),
	}
	var resultBytes bytes.Buffer
	err = tmpl.Execute(&resultBytes, data)
	if err != nil {
		log.Fatal(err)
	}
	resultString := resultBytes.String()
	subject := "Your Email Verification Code"
	sendEmail(resultString, user.Email, subject)
}

func sendEmail(htmlBody string, toEmail string, subject string) {
	enV := env.GetEnv{}
	enV.LoadEnv()

	email := gomail.NewMessage()
	email.SetHeader("From", enV.EmailUser)
	email.SetHeader("To", toEmail)
	email.SetHeader("Subject", subject)
	email.SetBody("text/html", htmlBody)

	dialer := gomail.NewDialer("smtp.hostinger.com", 465, enV.EmailUser, enV.EmailPassword)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: false, ServerName: "smtp.hostinger.com"}

	if err := dialer.DialAndSend(email); err != nil {
		log.Panicln(err)
	}
}
