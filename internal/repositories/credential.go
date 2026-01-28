package repositories

import (
	"errors"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *MainRepository) CredentialGetMPToken(tenantID int64) (*schemas.CredentialMPTokenResponse, error) {
	var credential models.Credential
	err := r.DB.Select("access_token_mp", "access_token_test_mp", "token_email").Where("tenant_id = ?", tenantID).First(&credential).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "Credenciales no encontradas", err)
		}
		return nil, schemas.ErrorResponse(500, "Error interno al obtener las credenciales", err)
	}

	response := &schemas.CredentialMPTokenResponse{
		AccessToken:     credential.AccessTokenMP,
		AccessTokenTest: credential.AccessTokenTestMP,
		TokenEmail:      credential.TokenEmail,
	}

	return response, nil
}

func (r *MainRepository) CredentialSetMPToken(tenantID int64, request *schemas.CredentialMPTokenRequest) error {
	credential := models.Credential{
		TenantID:          tenantID,
		AccessTokenMP:     &request.AccessToken,
		AccessTokenTestMP: &request.AccessTokenTest,
		// TokenEmail: &request.TokenEmail,
	}

	err := r.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "tenant_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"access_token_mp", "access_token_test_mp"}),
		// DoUpdates: clause.AssignmentColumns([]string{"access_token_mp", "access_token_test_mp", "token_email"}),
	}).Create(&credential).Error

	if err != nil {
		return schemas.ErrorResponse(500, "Error al actualizar credenciales", err)
	}

	return nil
}

func (r *MainRepository) CredentialGetArca(tenantID int64) (*schemas.CredentialArcaResponse, error) {
	var credential models.Credential
	if err := r.DB.Select("social_reason", "business_name", "address", "responsibility_front_iva", "cuit", "gross_income", "start_activities", "arca_certificate", "arca_key").
		Where("tenant_id = ?", tenantID).First(&credential).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "Credenciales no encontradas", err)
		}
		return nil, schemas.ErrorResponse(500, "Error interno al obtener las credenciales", err)
	}

	response := &schemas.CredentialArcaResponse{
		SocialReason:           credential.SocialReason,
		BusinessName:           credential.BusinessName,
		Address:                credential.Address,
		ResponsibilityFrontIVA: credential.ResponsibilityFrontIVA,
		GrossIncome:            credential.GrossIncome,
		StartActivities:        credential.StartActivities,
		Cuit:                   credential.Cuit,
		Concept:                credential.Concept,
		ArcaCertificate:        credential.ArcaCertificate,
		ArcaKey:                credential.ArcaKey,
	}

	return response, nil
}

func (r *MainRepository) CredentialSetArca(tenantID int64, request *schemas.CredentialArcaRequest) error {
	credential := models.Credential{
		TenantID:               tenantID,
		SocialReason:           &request.SocialReason,
		BusinessName:           &request.BusinessName,
		Address:                &request.Address,
		ResponsibilityFrontIVA: &request.ResponsibilityFrontIVA,
		GrossIncome:            &request.GrossIncome,
		StartActivities:        &request.StartActivities,
		Cuit:                   &request.Cuit,
		Concept:                &request.Concept,
		ArcaCertificate:        &request.ArcaCertificate,
		ArcaKey:                &request.ArcaKey,
	}

	err := r.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "tenant_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"social_reason", "business_name", "address", "responsibility_front_iva", "gross_income", "start_activities", "cuit", "concept", "arca_certificate", "arca_key"}),
	}).Create(&credential).Error

	if err != nil {
		return schemas.ErrorResponse(500, "Error al actualizar credenciales", err)
	}

	return nil
}
