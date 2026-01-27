package grpc_repo

import (
	"context"
	"errors"
	"time"

	// "time"

	"github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// func (r *GrpcMainRepository) ListTenants() ([]models.Tenant, error) {
// 	var tenants []models.Tenant
// 	// if err := r.DB.
// 	// 	Preload("Setting").
// 	// 	Preload("Credentials", func(db *gorm.DB) *gorm.DB { 
// 	// 		return db.Select("tenant_id", "access_token_mp")
// 	// 	}).
// 	// 	Preload("Modules", func(db *gorm.DB) *gorm.DB { 
// 	// 		return db.Where("expiration > ?", time.Now())
// 	// 	}).
// 	// 	Where("expiration > ?", time.Now()).
// 	// 	Where("is_active = ?", true).
// 	// 	Find(&tenants).Error; err != nil {
// 	// 	return nil, status.Error(codes.Internal, "Error al obtener los tenants")
// 	// }
// 	now := time.Now().UTC()
// 	if err := r.DB.
//         // ---------------------------------------------------------
//         // 1. FILTRO ESTRUCTURAL (JOIN)
//         // Esto es lo que soluciona tu problema. 
//         // Une la tabla tenants con modules y exige que exista coincidencia.
//         // ---------------------------------------------------------
//         Joins("JOIN modules ON modules.tenant_id = tenants.id").
//         Where("modules.expiration > ?", now).
        
//         // Al hacer JOIN 1:N, si un tenant tiene 3 módulos válidos, 
//         // la base de datos devuelve 3 filas. Usamos Group para unificarlas.
//         Group("tenants.id"). 

//         // ---------------------------------------------------------
//         // 2. FILTROS DEL TENANT (Padre)
//         // Nota: especificamos 'tenants.expiration' para evitar ambigüedad
//         // porque 'modules' también tiene columna 'expiration'.
//         // ---------------------------------------------------------
//         Where("tenants.expiration > ?", now).
//         Where("tenants.is_active = ?", true).

//         // ---------------------------------------------------------
//         // 3. CARGA DE DATOS (Preloads)
//         // ---------------------------------------------------------
//         Preload("Setting").
//         Preload("Credentials", func(db *gorm.DB) *gorm.DB { 
//             return db.Select("tenant_id", "access_token_mp")
//         }).
//         // Seguimos filtrando el Preload para que, visualmente en el JSON,
//         // solo aparezcan los módulos vigentes dentro del array.
//         Preload("Modules", func(db *gorm.DB) *gorm.DB { 
//             return db.Where("expiration > ?", now)
//         }).

//         Find(&tenants).Error; err != nil {
//             return nil, status.Error(codes.Internal, "Error al obtener los tenants")
//     }

// 	return tenants, nil
// }
func (r *GrpcMainRepository) ListTenants() ([]models.Tenant, error) {
	var tenants []models.Tenant
	now := time.Now().UTC()

	// Asumo que el nombre en BBDD es 'ecommerce' (minúsculas). 
	// Ajusta si es 'Ecommerce' o usa ILIKE en Postgres / LOWER() en MySQL si dudas.
	targetModule := "ecommerce"

	if err := r.DB.
		// ---------------------------------------------------------
		// 1. LOS JOINS (El filtro principal)
		// ---------------------------------------------------------
		// Paso A: Unir Tenant con la tabla intermedia (tenant_modules)
		Joins("JOIN tenant_modules ON tenant_modules.tenant_id = tenants.id").
		// Paso B: Unir la intermedia con la definición del Módulo (modules)
		Joins("JOIN modules ON modules.id = tenant_modules.module_id").

		// ---------------------------------------------------------
		// 2. LOS FILTROS (Where)
		// ---------------------------------------------------------
		// Filtro 1: Que el nombre del módulo sea "ecommerce"
		Where("modules.name = ?", targetModule).
		// Filtro 2: Que la suscripción a ESE módulo no esté vencida
		Where("tenant_modules.expiration > ?", now).
		// Filtro 3: Que el Tenant en sí mismo no esté vencido ni inactivo
		Where("tenants.expiration > ?", now).
		Where("tenants.is_active = ?", true).

		// ---------------------------------------------------------
		// 3. GROUP (Evitar duplicados)
		// ---------------------------------------------------------
		Group("tenants.id").

		// ---------------------------------------------------------
		// 4. PRELOADS (Llenar los datos para el JSON)
		// ---------------------------------------------------------
		Preload("Setting").
		Preload("Credentials", func(db *gorm.DB) *gorm.DB {
			return db.Select("tenant_id", "access_token_mp", "token_email")
		}).
		// Aquí cargamos la relación para que se vea en el JSON.
		// Opcional: Si quieres que en el JSON SOLO aparezca el módulo ecommerce,
		// mantén el Where dentro del Preload también.
		Preload("Modules", func(db *gorm.DB) *gorm.DB {
			return db.Joins("JOIN modules ON modules.id = tenant_modules.module_id").
				Where("modules.name = ?", targetModule).
				Where("tenant_modules.expiration > ?", now)
		}).
		// IMPORTANTE: Preload anidado para ver el nombre del módulo dentro de la estructura TenantModule
		Preload("Modules.Module"). 

		Find(&tenants).Error; err != nil {
		return nil, status.Error(codes.Internal, "Error al obtener los tenants")
	}

	return tenants, nil
}

func (r *GrpcMainRepository) GetTenant(req *pb.TenantRequest) (*models.Tenant, error) {
	var tenant models.Tenant
	if err := r.DB.Preload("Setting").Where("identifier = ?", req.Identifier).First(&tenant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "Tenant no encontrado")
		}
		return nil, status.Error(codes.Internal, "Error al obtener el tenant")
	}

	if !tenant.IsActive {
		return nil, status.Error(codes.PermissionDenied, "El tenant no esta activo")
	}

	// if tenant.Expiration.Before(time.Now()) {
	// 	return nil, status.Error(codes.PermissionDenied, "El tenant ha expirado")
	// }

	return &tenant, nil
}

func (r *GrpcMainRepository) UpdateImageSetting(ctx context.Context, req *pb.TenantRequestImageSetting) (*pb.TenantUpdateImageResponse, error) {
	resp := &pb.TenantUpdateImageResponse{}

	var tenant models.Tenant
	if err := r.DB.Select("id", "is_active").Where("identifier = ?", req.TenantIdentifier).First(&tenant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "Tenant no encontrado")
		}
		return nil, status.Error(codes.Internal, "Error al obtener el tenant")
	}

	if !tenant.IsActive {
		return nil, status.Error(codes.PermissionDenied, "El tenant no esta activo")
	}

	var setting models.SettingTenant
	if err := r.DB.Where("tenant_id = ?", tenant.ID).First(&setting).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "Setting no encontrado")
		}
		return nil, status.Error(codes.Internal, "Error al obtener el setting")
	}
	// if tenant.Expiration.Before(time.Now()) {
	// 	return nil, status.Error(codes.PermissionDenied, "El tenant ha expirado")
	// }

	resp.LogoUuid = utils.Ternary(req.LogoUuid == nil, nil, setting.Logo)
	resp.FrontPageUuid = utils.Ternary(req.FrontPageUuid == nil, nil, setting.FrontPage)
	updateData := models.SettingTenant{
		Logo:      req.LogoUuid,
		FrontPage: req.FrontPageUuid,
	}

	if err := r.DB.Model(&models.SettingTenant{}).Where("tenant_id = ?", tenant.ID).Updates(updateData).Error; err != nil {
		return nil, status.Error(codes.Internal, "Error al actualizar el tenant")
	}

	return resp, nil
}
