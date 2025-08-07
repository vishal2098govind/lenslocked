package main

import (
	"fmt"
	"html/template"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	emailM "github.com/vishal2098govind/lenslocked/models/email"
	"github.com/vishal2098govind/lenslocked/templates"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic(err)
	}
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	emailS := emailM.NewEmailService(emailM.SMTPConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	})
	emailS.Templates.ForgotPasswordTpl = template.Must(
		template.ParseFS(templates.FS, "forgot_password.gohtml"),
	)

	err = emailS.SendForgotPasswordEmail(
		emailM.SendForgotPasswordEmailRequest{
			To:       "vishal.govind2098@gmail.com",
			ResetUrl: "http://localhost:3000/",
		},
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Message sent")
}
