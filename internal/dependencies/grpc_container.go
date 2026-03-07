// pkg/dependencies/container.go
package dependencies

import (
	"database/sql"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/repositories/grpc_repo"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/services/grpc_serv"
)

type GrpcContainer struct {
	DB      *sql.DB
	Services struct {
		GrpcProductService *grpc_serv.GrpcProductService
		GrpcCategoryService *grpc_serv.GrpcCategoryService
		GrpcMPService *grpc_serv.GrpcMPService
	}
	Repositories struct {
		GrpcProductRepository *grpc_repo.GrpcProductRepository
		GrpcCategoryRepository *grpc_repo.GrpcCategoryRepository
		GrpcMPRepository *grpc_repo.GrpcMPRepository
	}
}

func NewGrpcContainer(db *sql.DB) *GrpcContainer {
	c := &GrpcContainer{DB: db}

	// Inicializar repositorios
	c.Repositories.GrpcProductRepository = &grpc_repo.GrpcProductRepository{DB: c.DB}
	c.Repositories.GrpcCategoryRepository = &grpc_repo.GrpcCategoryRepository{DB: c.DB}
	c.Repositories.GrpcMPRepository = &grpc_repo.GrpcMPRepository{DB: c.DB}
	// Inicializar servicios

	c.Services.GrpcProductService = &grpc_serv.GrpcProductService{
		GrpcProductRepository: c.Repositories.GrpcProductRepository,
	}
	c.Services.GrpcCategoryService = &grpc_serv.GrpcCategoryService{
		GrpcCategoryRepository: c.Repositories.GrpcCategoryRepository,
	}
	c.Services.GrpcMPService = &grpc_serv.GrpcMPService{
		GrpcMPRepository: c.Repositories.GrpcMPRepository,
	}

	return c
}
