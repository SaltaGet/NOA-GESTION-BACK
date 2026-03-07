package grpc_repo

import (
	"database/sql"
)

type GrpcMainRepository struct {
	DB *sql.DB
}

type GrpcProductRepository struct {
	DB *sql.DB
}

type GrpcCategoryRepository struct {
	DB *sql.DB
}

type GrpcMPRepository struct {
	DB *sql.DB
}