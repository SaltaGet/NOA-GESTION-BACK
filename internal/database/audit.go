package database

import (
	"encoding/json"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"gorm.io/gorm"
	"github.com/rs/zerolog/log"
)

func SaveAuditAsync(db *gorm.DB, audit models.AuditLog, oldValue, newValue any) {
	func(audit models.AuditLog, oldVal, newVal any) {
		if oldVal != nil {
			data, err := json.Marshal(oldVal)
			if err != nil {
				log.Err(err).Msg("error al serializar OldValue en auditorí")
			} else {
				s := string(data)
				audit.OldValue = &s
			}
		}

		if newVal != nil {
			data, err := json.Marshal(newVal)
			if err != nil {
				log.Err(err).Msg("error al serializar NewValue en auditoría")
			} else {
				s := string(data)
				audit.NewValue = &s
			}
		}

		if err := db.Create(&audit).Error; err != nil {
			log.Err(err).Msg("error al guardar auditoría")
		}
	}(audit, oldValue, newValue)
}

func SaveAuditAdminAsync(db *gorm.DB, audit models.AuditLogAdmin, oldValue, newValue any) {
		if oldValue != nil {
			data, err := json.Marshal(oldValue)
			if err != nil {
				log.Err(err).Msg("error al serializar OldValue en auditoría")
			} else {
				s := string(data)
				audit.OldValue = &s
			}
		}

		if newValue != nil {
			data, err := json.Marshal(newValue)
			if err != nil {
				log.Err(err).Msg("error al serializar NewValue en auditoría")
			} else {
				s := string(data)
				audit.NewValue = &s
			}
		}

		if err := db.Create(&audit).Error; err != nil {
			log.Err(err).Msg("error al guardar auditoría")
		}
}