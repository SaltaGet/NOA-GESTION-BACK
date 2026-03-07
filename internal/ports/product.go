package ports

import (
	"mime/multipart"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models/tenant"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

type ProductRepository interface {
	ProductGetByID(id int64) (*tenant.Product, error)
	ProductGetByCode(code string) (*tenant.Product, error)
	ProductGetByCategoryID(categoryID int64) ([]*tenant.Product, error)
	ProductGetByName(name string) ([]*tenant.Product, error)
	ProductGetAll(page, limit int, isVisible *bool) ([]*tenant.Product, int64, error)
	ProductGetByCodeToQR(code string) (*tenant.Product, error)
	ProductCount() (int64, error)
	ProductCreate(memberID int64, productCreate *schemas.ProductCreate, plan *schemas.PlanResponseDTO) (int64, error)
	ProductInsertToExcel(memberID int64, products []schemas.ProductExcelCreate) ([]map[string]string, error)
	ProductUpdate(memberID int64, productUpdate *schemas.ProductUpdate) error
	ProductPriceUpdate(memberID int64, productUpdate *schemas.ListPriceUpdate) error
	ProductDelete(memberID int64, id int64) error
	ValidateProductImages(productValidateImage schemas.ProductValidateImage, plan *schemas.PlanResponseDTO) error
	ProductUpdateVisibility(productUpdate *schemas.ListVisibilityUpdate) error
}

type ProductService interface {
	ProductGetByID(id int64) (*schemas.ProductFullResponse, error)
	ProductGetByCode(code string) (*schemas.ProductFullResponse, error)
	ProductGetByName(name string) ([]*schemas.ProductFullResponse, error)
	ProductGetByCategoryID(categoryID int64) ([]*schemas.ProductFullResponse, error)
	ProductGetAll(page, limit int, isVisible *bool) ([]*schemas.ProductFullResponse, int64, error)
	ProductGenerateQR(code string, rows, cols int) ([]byte, error)
	ProductCreate(memberID int64, productCreate *schemas.ProductCreate, plan *schemas.PlanResponseDTO) (int64, error)
	ProductUpload(memberID int64, file *multipart.FileHeader, plan *schemas.PlanResponseDTO) ([]map[string]string, error)
	ProductUpdate(memberID int64, productUpdate *schemas.ProductUpdate) error
	ProductPriceUpdate(memberID int64, productUpdate *schemas.ListPriceUpdate) error
	ProductDelete(memberID int64, id int64) error
	ValidateProductImages(tenantIdentifier string, productValidateImage schemas.ProductValidateImage, plan *schemas.PlanResponseDTO) (string, error)
	ProductUpdateVisibility(productUpdate *schemas.ListVisibilityUpdate) error
}
