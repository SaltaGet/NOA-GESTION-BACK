package grpc_serv

import (
	"context"

	"github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
)

func (s *GrpcMPService) SyncPurchasePayment(ctx context.Context, req *pb.DataInfoPay) (error) {
	return s.GrpcMPRepository.SyncPurchasePayment(ctx, req)
}	