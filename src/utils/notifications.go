package utils

import (
	"bytes"
	"crypto/tls"
	"log"
	"os"
	"text/template"

	"github.com/Numostanley/auth_app/db"
	"github.com/Numostanley/auth_app/env"
	"github.com/Numostanley/auth_app/models"
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

	database := db.Database.DB
	vCode := models.CreateVerificationCode(user, database)

	data := struct {
		Code string
		Name string
	}{
		Code: vCode.Code,
		Name: user.GetFullName(),
	}

	var resultBytes bytes.Buffer
	err = tmpl.Execute(&resultBytes, data)
	if err != nil {
		log.Fatal(err)
	}
	resultString := resultBytes.String()
	subject := "Verification Code"
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

	dialer := gomail.NewDialer(enV.EmailHost, 465, enV.EmailUser, enV.EmailPassword)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: false, ServerName: enV.EmailHost}

	if err := dialer.DialAndSend(email); err != nil {
		log.Println(err)
	}
}
