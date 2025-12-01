package ports

import "github.com/SaltaGet/NOA-GESTION-BACK/internal/models"

type EmailService interface {
	SendEmail(member *models.Member, subject, body string) error
}