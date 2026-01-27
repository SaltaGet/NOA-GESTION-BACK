package grpc_ports

import (
	"context"

	"github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
)

type GrpcMPRepository interface {
	SyncPurchasePayment(ctx context.Context, req *pb.DataInfoPay) (error)
}

type GrpcMPService interface {
	SyncPurchasePayment(ctx context.Context, req *pb.DataInfoPay) (error)
}