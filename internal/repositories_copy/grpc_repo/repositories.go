package grpc_repo_copy

import "gorm.io/gorm"

type GrpcMainRepository struct {
	DB *gorm.DB
}

type GrpcProductRepository struct {
	DB *gorm.DB
}

type GrpcCategoryRepository struct {
	DB *gorm.DB
}

type GrpcMPRepository struct {
	DB *gorm.DB
}