package ports

type EmailService interface {
	SendForgotPasswordEmail(username, email, token string) error
	SendWelcomeAdminEmail(email, username, password string) error
	SendWelcomeUserEmail(email, username, password string) error
}