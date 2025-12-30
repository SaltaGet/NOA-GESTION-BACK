package grpc_serv

import (
	"context"
	"strings"

	pb "github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
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
	// Defaults
	if req.Page <= 0 { req.Page = 1 }
	if req.PageSize <= 0 { req.PageSize = 10 }

    // Obtenemos categoryId de forma segura si es opcional
    var catID int32 = 0
    if req.CategoryId != nil {
        catID = *req.CategoryId
    }
    
    // Obtenemos search
    search := ""
    if req.Search != nil {
        search = *req.Search
    }
    
    sort := pb.ListProductsRequest_PRICE_LOW_TO_HIGH
    if req.Sort != nil {
        sort = *req.Sort
    }

	products, total, err := s.GrpcProductRepository.ProductList(int32(req.Page), int32(req.PageSize), &catID, &search, int32(sort))
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
func mapModelToProto(m *models.Product) *pb.Product {
	var second []string = []string{}
	if m.SecondaryImages != nil && *m.SecondaryImages != "" {
        second = strings.Split(*m.SecondaryImages, ",")
    }

	return &pb.Product{
		Id:          int64(m.ID),
		Code:        m.Code,
		Name:        m.Name,
		Description: m.Description,
		Price:       m.Price,
		Stock:       float32(m.StockDeposit.Stock),
		PrimaryImage: m.PrimaryImage, 
		SecondaryImages: second,
		Category: &pb.Category{
			Id:   int64(m.Category.ID),
			Name: m.Category.Name,
		},
	}
}

func mapModelToDTO(m *models.Product) *pb.ProductDTO {
	var stock float32 = 0.0
    if m.StockDeposit != nil {
        stock = float32(m.StockDeposit.Stock)
    }

	return &pb.ProductDTO{
		Id:    int64(m.ID),
		Code:  m.Code,
		Name:  m.Name,
		Price: m.Price,
		Stock: stock,
		PrimaryImage: m.PrimaryImage,
		Category: &pb.Category{
			Id:   int64(m.Category.ID),
			Name: m.Category.Name,
		},
	}
}