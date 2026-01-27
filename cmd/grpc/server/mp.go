package server

import (
	"context"

	"github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
	grpc_cache "github.com/SaltaGet/NOA-GESTION-BACK/internal/cache/grpc"
)

func (h *GrpcMPServer) SyncPurchasePayment(ctx context.Context, req *pb.DataInfoPay) (*pb.DataInfoPayResponse, error) {
	deps := grpc_cache.GetGrpcContainerFromContext(ctx)
	err := deps.Services.GrpcMPService.SyncPurchasePayment(ctx, req)
	if err != nil {
		return nil, err
	}

	return &pb.DataInfoPayResponse{}, nil
}