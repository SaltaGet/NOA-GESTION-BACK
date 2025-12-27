package dependencies

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/repositories/grpc_repo"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/services/grpc_serv"
	"gorm.io/gorm"
)

type GrpcMainContainer struct {
	TenantGrpcService *grpc_serv.GrpcTenantService
}

func NewGrpcApplication(mainDB *gorm.DB) *GrpcMainContainer {
	mainRepo := &grpc_repo.GrpcMainRepository{DB: mainDB}

	tenantServ := &grpc_serv.GrpcTenantService{GrpcTenantRepository: mainRepo}

	return &GrpcMainContainer{
		TenantGrpcService: tenantServ,
	}


	// return &GrpcMainContainer{
	// 	TenantGrpcService: &server.GrpcTenantServer{GrpcTenantService: tenantServ},
	// }
}

