package utils

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
