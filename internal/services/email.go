package services

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/platform/event"
)

type ForgotPasswordPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

type WelcomeAdminPayload struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type WelcomeUserPayload struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (es *EmailService) SendForgotPasswordEmail(username, email, token string) error {
	payload := ForgotPasswordPayload{
		Username: username,
		Email:    email,
		Token:    token,
	}
	return event.Publish("email.forgot_password", payload)
}

func (es *EmailService) SendWelcomeAdminEmail(email, username, password string) error {
	payload := WelcomeAdminPayload{
		Email:    email,
		Username: username,
		Password: password,
	}
	return event.Publish("email.welcome_admin", payload)
}

func (es *EmailService) SendWelcomeUserEmail(email, username, password string) error {
	payload := WelcomeUserPayload{
		Email:    email,
		Username: username,
		Password: password,
	}
	return event.Publish("email.welcome_user", payload)
}
