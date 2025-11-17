package ports

import (
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/xuri/excelize/v2"
)

type ReportRepository interface {
	ReportMovementByDate(fromDate, toDate time.Time, form string) (any, error)
	ReportProfitableProducts(start, end time.Time) ([]schemas.ReportProfitableProducts, error)
	ReportStockProducts() ([]*models.Product, error)
}

type ReportService interface {
	ReportExcelGet(start, end time.Time) (*excelize.File, error)
	ReportMovementByDate(fromDate, toDate time.Time, form string) (any, error)
	ReportProfitableProducts(start, end time.Time) ([]schemas.ReportProfitableProducts, error)
}