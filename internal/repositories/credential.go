package repositories

import (
	"context"

	boilmodels "github.com/SaltaGet/NOA-GESTION-BACK/internal/models/boil"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (r *MainRepository) CredentialGetMPToken(tenantID int64) (*schemas.CredentialMPTokenResponse, error) {
	ctx := context.Background()

	c, err := boilmodels.Credentials(
		qm.Select(boilmodels.CredentialColumns.AccessTokenMP, boilmodels.CredentialColumns.AccessTokenTestMP, boilmodels.CredentialColumns.TokenEmail),
		boilmodels.CredentialWhere.TenantID.EQ(tenantID),
	).One(ctx, r.DB)

	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Credencial", schemas.Read)
	}

	response := &schemas.CredentialMPTokenResponse{}

	if c.AccessTokenMP.Valid {
		response.AccessToken = &c.AccessTokenMP.String
	}
	if c.AccessTokenTestMP.Valid {
		response.AccessTokenTest = &c.AccessTokenTestMP.String
	}
	if c.TokenEmail.Valid {
		response.TokenEmail = &c.TokenEmail.String
	}

	return response, nil
}

func (r *MainRepository) CredentialSetMPToken(tenantID int64, request *schemas.CredentialMPTokenRequest) error {
	ctx := context.Background()

	credential := boilmodels.Credential{
		TenantID:          tenantID,
		AccessTokenMP:     null.StringFrom(request.AccessToken),
		AccessTokenTestMP: null.StringFrom(request.AccessTokenTest),
	}

	err := credential.Upsert(ctx, r.DB, true,
		[]string{boilmodels.CredentialColumns.TenantID},
		boil.Whitelist(
			boilmodels.CredentialColumns.AccessTokenMP,
			boilmodels.CredentialColumns.AccessTokenTestMP,
			boilmodels.CredentialColumns.UpdatedAt,
		),
		boil.Infer(),
	)

	if err != nil {
		return schemas.HandlerErrorDB(err, "Credencial", schemas.Update)
	}

	return nil
}

func (r *MainRepository) CredentialGetArca(tenantID int64) (*schemas.CredentialArcaResponse, error) {
	ctx := context.Background()

	c, err := boilmodels.Credentials(
		qm.Select(
			boilmodels.CredentialColumns.SocialReason, boilmodels.CredentialColumns.BusinessName,
			boilmodels.CredentialColumns.Address, boilmodels.CredentialColumns.ResponsibilityFrontIva,
			boilmodels.CredentialColumns.Cuit, boilmodels.CredentialColumns.GrossIncome,
			boilmodels.CredentialColumns.StartActivities, boilmodels.CredentialColumns.ArcaCertificate,
			boilmodels.CredentialColumns.ArcaKey, boilmodels.CredentialColumns.Concept,
		),
		boilmodels.CredentialWhere.TenantID.EQ(tenantID),
	).One(ctx, r.DB)

	if err != nil {
		return nil, schemas.HandlerErrorDB(err, "Credencial", schemas.Read)
	}

	response := &schemas.CredentialArcaResponse{}

	if c.SocialReason.Valid {
		response.SocialReason = &c.SocialReason.String
	}
	if c.BusinessName.Valid {
		response.BusinessName = &c.BusinessName.String
	}
	if c.Address.Valid {
		response.Address = &c.Address.String
	}
	if c.ResponsibilityFrontIva.Valid {
		response.ResponsibilityFrontIVA = &c.ResponsibilityFrontIva.String
	}
	if c.GrossIncome.Valid {
		response.GrossIncome = &c.GrossIncome.String
	}
	if c.StartActivities.Valid {
		response.StartActivities = &c.StartActivities.String
	}
	if c.Cuit != "" {
		response.Cuit = &c.Cuit
	}
	if c.Concept.Valid {
		response.Concept = &c.Concept.String
	}
	if c.ArcaCertificate.Valid {
		response.ArcaCertificate = &c.ArcaCertificate.String
	}
	if c.ArcaKey.Valid {
		response.ArcaKey = &c.ArcaKey.String
	}

	return response, nil
}

func (r *MainRepository) CredentialSetArca(tenantID int64, request *schemas.CredentialArcaRequest) error {
	ctx := context.Background()

	credential := boilmodels.Credential{
		TenantID:               tenantID,
		SocialReason:           null.StringFrom(request.SocialReason),
		BusinessName:           null.StringFrom(request.BusinessName),
		Address:                null.StringFrom(request.Address),
		ResponsibilityFrontIva: null.StringFrom(request.ResponsibilityFrontIVA),
		GrossIncome:            null.StringFrom(request.GrossIncome),
		StartActivities:        null.StringFrom(request.StartActivities),
		Cuit:                   request.Cuit,
		Concept:                null.StringFrom(request.Concept),
		ArcaCertificate:        null.StringFrom(request.ArcaCertificate),
		ArcaKey:                null.StringFrom(request.ArcaKey),
	}

	err := credential.Upsert(ctx, r.DB, true,
		[]string{boilmodels.CredentialColumns.TenantID},
		boil.Whitelist(
			boilmodels.CredentialColumns.SocialReason,
			boilmodels.CredentialColumns.BusinessName,
			boilmodels.CredentialColumns.Address,
			boilmodels.CredentialColumns.ResponsibilityFrontIva,
			boilmodels.CredentialColumns.GrossIncome,
			boilmodels.CredentialColumns.StartActivities,
			boilmodels.CredentialColumns.Cuit,
			boilmodels.CredentialColumns.Concept,
			boilmodels.CredentialColumns.ArcaCertificate,
			boilmodels.CredentialColumns.ArcaKey,
			boilmodels.CredentialColumns.UpdatedAt,
		),
		boil.Infer(),
	)

	if err != nil {
		return schemas.HandlerErrorDB(err, "Credencial", schemas.Update)
	}

	return nil
}
