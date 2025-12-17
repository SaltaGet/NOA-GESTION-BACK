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
				line-height: 1.5;
			}
			.info-box {
				background: #f7f7f7;
				padding: 15px;
				border-radius: 8px;
				margin-top: 20px;
				font-size: 15px;
				color: #333;
			}
			.warning {
				background: #fff3cd;
				border-left: 4px solid #ffc107;
				padding: 12px 15px;
				margin: 20px 0;
				border-radius: 4px;
				color: #856404;
			}
			.btn-reset {
				display: inline-block;
				padding: 16px 40px;
				background-color: #048bfa;
				color: #ffffff !important;
				text-decoration: none;
				font-size: 18px;
				font-weight: 600;
				border-radius: 8px;
				transition: all 0.3s ease;
				box-shadow: 0 4px 6px rgba(4, 139, 250, 0.3);
				margin: 20px 0;
			}
			.btn-reset:hover {
				background-color: #0370d1;
				color: #ffffff !important;
				transform: translateY(-2px);
				box-shadow: 0 6px 12px rgba(4, 139, 250, 0.4);
			}
			.btn-container {
				text-align: center;
				margin: 25px 0;
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
			<h1>Recuperación de Contraseña</h1>

			<p>Hola <strong>` + username + `</strong>,</p>
			<p>Hemos recibido una solicitud para restablecer la contraseña de tu cuenta.</p>

			<div class="info-box">
				<p><strong>Usuario:</strong> ` + username + `</p>
				<p><strong>Email:</strong> ` + email + `</p>
			</div>

			<div class="warning">
				<strong>⏱️ Importante:</strong> Este enlace tiene una validez de <strong>30 minutos</strong>. Después de este tiempo deberás solicitar un nuevo enlace.
			</div>

			<p>Para restablecer tu contraseña, haz clic en el siguiente botón:</p>

			<div class="btn-container">
				<a href="https://noagestion.com.ar/reset-password/` + token + `" target="_blank" class="btn-reset">Restablecer Contraseña</a>
			</div>

			<p>Si no solicitaste este cambio, puedes ignorar este correo. Tu contraseña permanecerá sin cambios.</p>

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
			.btn-login {
				display: inline-block;
				padding: 16px 40px;
				background-color: #048bfa;
				color: #ffffff !important;
				text-decoration: none;
				font-size: 18px;
				font-weight: 600;
				border-radius: 8px;
				transition: all 0.3s ease;
				box-shadow: 0 4px 6px rgba(4, 139, 250, 0.3);
				margin: 20px 0;
			}
			.btn-login:hover {
				background-color: #0370d1;
				color: #ffffff !important;
				transform: translateY(-2px);
				box-shadow: 0 6px 12px rgba(4, 139, 250, 0.4);
			}
			.btn-container {
				text-align: center;
				margin: 25px 0;
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

			<div class="btn-container">
				<a href="https://noagestion.com.ar/login-admin" target="_blank" class="btn-login">Iniciar sesión</a>
			</div>

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

func WelcomeUser(username, pass string) string {
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
			.btn-login {
				display: inline-block;
				padding: 16px 40px;
				background-color: #048bfa;
				color: #ffffff !important;
				text-decoration: none;
				font-size: 18px;
				font-weight: 600;
				border-radius: 8px;
				transition: all 0.3s ease;
				box-shadow: 0 4px 6px rgba(4, 139, 250, 0.3);
				margin: 20px 0;
			}
			.btn-login:hover {
				background-color: #0370d1;
				color: #ffffff !important;
				transform: translateY(-2px);
				box-shadow: 0 6px 12px rgba(4, 139, 250, 0.4);
			}
			.btn-container {
				text-align: center;
				margin: 25px 0;
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
			<p>Tu cuenta de usuario ha sido creada exitosamente.</p>

			<div class="credentials">
				<p><strong>Usuario:</strong> ` + username + `</p>
				<p><strong>Contraseña temporal:</strong> ` + pass + `</p>
			</div>

			<p>Recuerda que para ingresar al sistema el usuario debe tener el formato <strong>username@negocio</strong></p>
			
			<p>Te recomendamos iniciar sesión y cambiar tu contraseña por una más segura.</p>

			<div class="btn-container">
				<a href="https://noagestion.com.ar/login" target="_blank" class="btn-login">Iniciar sesión</a>
			</div>

			<p>¡Gracias por confiar en NOA-GESTION!</p>

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