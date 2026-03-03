package repository

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/module/email/domain"
)

type emailRepository struct {
	// TODO: Add DB connections or clients
}

func NewEmailRepository() domain.EmailRepository {
	return &emailRepository{}
}
