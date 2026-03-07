package grpc_serv

import (
	"context"
	"strings"

	pb "github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models/tenant"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *GrpcProductService) ProductGetByCode(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	if req.Code == "" {
		return nil, status.Error(codes.InvalidArgument, "Code is required")
	}

	prod, err := s.GrpcProductRepository.ProductGetByCode(req.Code)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Product not found")
	}

	return mapModelToProto(prod), nil
}

func (s *GrpcProductService) ProductList(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	products, total, err := s.GrpcProductRepository.ProductList(req)
	if err != nil {
		return nil, status.Error(codes.Internal, "Error listing products")
	}

	var dtos []*pb.ProductDTO
	for _, p := range products {
		dtos = append(dtos, mapModelToDTO(p))
	}

	return &pb.ListProductsResponse{
		Products: dtos,
		Total:    int32(total),
		TenantId: "",
	}, nil
}

func (s *GrpcProductService) SaveUrlImage(ctx context.Context, req *pb.SaveImageRequest) (*pb.SaveImageResponse, error) {
	err := s.GrpcProductRepository.SaveUrlImage(req)
	if err != nil {
		return nil, status.Error(codes.Internal, "Error saving image")
	}
	return &pb.SaveImageResponse{
		Success: true,
	}, nil
}

func (s *GrpcProductService) ProductGetByID(ctx context.Context, req *pb.ProductRequest) (*pb.Product, error) {
	prod, err := s.GrpcProductRepository.ProductGetByID(req.ProductId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Product not found")
	}

	return mapModelToProto(prod), nil
}

// Helpers de mapeo (puedes moverlos a otro archivo)
func mapModelToProto(m *tenant.Product) *pb.Product {
	var second []string = []string{}
	if m.SecondaryImages.Valid && m.SecondaryImages.String != "" {
		second = strings.Split(m.SecondaryImages.String, ",")
	}

	stock, _ := 0.0, false
	if m.R != nil && len(m.R.Deposits) > 0 {
		stock, _ = m.R.Deposits[0].Stock.Float64()
	}

	catID := int64(m.CategoryID)
	catName := ""
	if m.R != nil && m.R.Category != nil {
		catName = m.R.Category.Name
	}

	price, _ := m.Price.Float64()

	return &pb.Product{
		Id:              m.ID,
		Code:            m.Code,
		Name:            m.Name,
		Description:     m.Description.Ptr(),
		Price:           price,
		Stock:           float32(stock),
		PrimaryImage:    m.PrimaryImage.Ptr(),
		SecondaryImages: second,
		Category: &pb.Category{
			Id:   catID,
			Name: catName,
		},
	}
}

func mapModelToDTO(m *tenant.Product) *pb.ProductDTO {
	stock, _ := 0.0, false
	if m.R != nil && len(m.R.Deposits) > 0 {
		stock, _ = m.R.Deposits[0].Stock.Float64()
	}

	catID := int64(m.CategoryID)
	catName := ""
	if m.R != nil && m.R.Category != nil {
		catName = m.R.Category.Name
	}

	price, _ := m.Price.Float64()

	return &pb.ProductDTO{
		Id:           m.ID,
		Code:         m.Code,
		Name:         m.Name,
		Price:        price,
		Stock:        float32(stock),
		PrimaryImage: m.PrimaryImage.Ptr(),
		Category: &pb.Category{
			Id:   catID,
			Name: catName,
		},
	}
}

// func (s *GrpcProductService) ValidateProducts(ctx context.Context, req *pb.ProductValidateRequest) (*pb.ProductValidateResponse, error) {
// 	products, err := s.GrpcProductRepository.ValidateProducts(ctx, req)
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, "Error validating products")
// 	}

// 	var dtos *pb.ProductValidateResponse
// 	for _, p := range products {
// 		var prod *pb.ProductValidate = &pb.ProductValidate{
// 			Id:    p.ID,
// 			Price: p.Price,
// 			Stock: p.StockDeposit.Stock,
// 		}
// 		dtos.Products = append(dtos.Products, prod)
// 	}

//		return dtos, nil
//	}
func (s *GrpcProductService) ValidateProducts(ctx context.Context, req *pb.ProductValidateRequest) (*pb.ProductValidateResponse, error) {
	products, err := s.GrpcProductRepository.ValidateProducts(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, "Error validating products")
	}

	// FIX 1: Initialize the struct pointer
	res := &pb.ProductValidateResponse{
		Products: make([]*pb.ProductValidate, 0, len(products)),
	}

	for _, p := range products {
		stock, _ := 0.0, false
		if p.R != nil && len(p.R.Deposits) > 0 {
			stock, _ = p.R.Deposits[0].Stock.Float64()
		}

		price, _ := p.Price.Float64()

		prod := &pb.ProductValidate{
			Id:    p.ID,
			Price: price,
			Stock: stock,
		}

		res.Products = append(res.Products, prod)
	}

	return res, nil
}
