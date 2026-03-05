package ports

import (
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/xuri/excelize/v2"
)

type ReportRepository interface {
	ReportMovementByDate(start, end time.Time, form string) (any, error)
	ReportMovementByDatePointSale(start, end time.Time, form string) (any, error)
	ReportProfitableProducts(start, end time.Time) ([]schemas.ReportProfitableProducts, error)
	ReportStockProducts() ([]schemas.ReportStockProduct, error)
}

type ReportService interface {
	ReportExcelGet(start, end time.Time) (*excelize.File, error)

	ReportMovementByDate(start, end time.Time, form string) (any, error)
	ReportMovementByDatePointSale(start, end time.Time, form string) (any, error)
	ReportProfitableProducts(start, end time.Time) ([]schemas.ReportProfitableProducts, error)
}
