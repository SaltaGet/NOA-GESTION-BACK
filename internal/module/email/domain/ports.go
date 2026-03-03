package domain


type EmailService interface {
	SendEmail(email, subject, body string) error

}
