package services

import (
	"io"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/assets"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"gopkg.in/gomail.v2"
)

func (es *EmailService) SendEmail(email, subject, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", es.Dialer.Username)
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	data, err := assets.LogoFS.ReadFile("logo.png")
	if err != nil {
		return schemas.ErrorResponse(500, "error al leer logo", err)
	}

	msg.Embed("logo.png", gomail.SetCopyFunc(func(w io.Writer) error {
		_, err := w.Write(data)
		return err
	}), gomail.SetHeader(map[string][]string{"Content-ID": {"<logo>"}}))

	if err := es.Dialer.DialAndSend(msg); err != nil {
		return schemas.ErrorResponse(500, "error al enviar email", err)
	}

	return nil
}
