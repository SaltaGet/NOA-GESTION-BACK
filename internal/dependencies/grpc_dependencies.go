package dependencies

import (
	"database/sql"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/repositories/grpc_repo"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/services/grpc_serv"
)

type GrpcMainContainer struct {
	TenantGrpcService *grpc_serv.GrpcTenantService
}

func NewGrpcApplication(mainDB *sql.DB) *GrpcMainContainer {
	mainRepo := &grpc_repo.GrpcMainRepository{DB: mainDB}

	tenantServ := &grpc_serv.GrpcTenantService{GrpcTenantRepository: mainRepo}

	return &GrpcMainContainer{
		TenantGrpcService: tenantServ,
	}


	// return &GrpcMainContainer{
	// 	TenantGrpcService: &server.GrpcTenantServer{GrpcTenantService: tenantServ},
	// }
}

