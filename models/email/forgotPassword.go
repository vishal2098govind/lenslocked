package models

import (
	"fmt"
	"strings"
)

type SendForgotPasswordEmailRequest struct {
	To       string
	ResetUrl string
}

func (es *EmailService) SendForgotPasswordEmail(r SendForgotPasswordEmailRequest) error {

	buf := &strings.Builder{}
	err := es.Templates.ForgotPasswordTpl.Execute(
		buf,
		struct{ ResetUrl string }{ResetUrl: r.ResetUrl},
	)
	if err != nil {
		return fmt.Errorf("send forget password email: %w", err)
	}

	email := Email{
		Subject:   "Reset your password",
		To:        r.To,
		PlainText: "To reset your password, please visit the following link: " + r.ResetUrl,
		HTML:      buf.String(),
	}

	err = es.Send(email)
	if err != nil {
		return fmt.Errorf("send forgot password email: %w", err)
	}

	return nil
}
