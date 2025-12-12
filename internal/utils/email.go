package utils

import (
	"io"
	"os"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/assets"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"gopkg.in/gomail.v2"
)

func ForgotPassword(username, email, token string) string {
	htmlBody := `
	<!DOCTYPE html>
	<html>
	<head>
		<style>
			.container {
				max-width: 600px;
				margin: auto;
				padding: 20px;
				background: #ffffff;
				border-radius: 10px;
				font-family: Arial, sans-serif;
			}
			h1 {
				color: #333;
				text-align: center;
			}
			p {
				font-size: 16px;
				color: #555;
			}
			.footer {
				margin-top: 20px;
				text-align: center;
				font-size: 12px;
				color: #777;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h1>Gil como vas a olvidar la contraseña !!! salame!  ` + email + ` : ` + username + `</h1>
			<p>Este es un correo de prueba con HTML, estilos y una imagen incrustada.</p>

			<div style="text-align:center;">
				<img src="cid:logo" style="width:200px; margin-top:20px;" />
			</div>

			<p>¡Gracias por utilizar nuestro servicio!</p>
			<p>Para restablecer tu contraseña, haz clic en el siguiente enlace y que sea la ultima vez no gil?:</p>
			<p><a href="https://example.com/reset-password?token=` + token + `">Restablecer Contraseña</a></p>
			<p>Si no solicitaste este cambio, puedes ignorar este correo.</p>

			<div class="footer">
				© 2025 Mi Empresa. Todos los derechos reservados.
			</div>
		</div>
	</body>
	</html>
	`

	return htmlBody
}

func WelcomeAdmin(email, username, pass string) string {
	htmlBody := `
	<!DOCTYPE html>
	<html>
	<head>
		<style>
			.container {
				max-width: 600px;
				margin: auto;
				padding: 20px;
				background: #ffffff;
				border-radius: 10px;
				font-family: Arial, sans-serif;
			}
			h1 {
				color: #333;
				text-align: center;
			}
			p {
				font-size: 16px;
				color: #555;
				line-height: 1.5;
			}
			.credentials {
				background: #f7f7f7;
				padding: 15px;
				border-radius: 8px;
				margin-top: 20px;
				font-size: 15px;
				color: #333;
			}
			.footer {
				margin-top: 30px;
				text-align: center;
				font-size: 12px;
				color: #777;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h1>Bienvenido a NOA-GESTION</h1>

			<p>Hola <strong>` + username + `</strong>,</p>
			<p>Tu cuenta de administrador ha sido creada exitosamente.</p>

			<div class="credentials">
				<p><strong>Usuario:</strong> ` + username + `</p>
				<p><strong>Email:</strong> ` + email + `</p>
				<p><strong>Contraseña temporal:</strong> ` + pass + `</p>
			</div>

			<p>Te recomendamos iniciar sesión y cambiar tu contraseña por una más segura.</p>

			<div style="text-align:center;">
				<img src="cid:logo" style="width:200px; margin-top:20px;" />
			</div>

			<div class="footer">
				© 2025 NOA-GESTION. Todos los derechos reservados.
			</div>
		</div>
	</body>
	</html>
	`

	return htmlBody
}



func SendEmail(email, subject, body string, cfg *schemas.EmailConfig) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", os.Getenv("EMAIL"))
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

	dialer := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)

	if err := dialer.DialAndSend(msg); err != nil {
		return schemas.ErrorResponse(500, "error al enviar email", err)
	}

	return nil
}