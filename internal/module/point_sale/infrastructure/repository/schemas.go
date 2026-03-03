package repository


type PointSaleCreate struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description"`
	Number			int64   `json:"number" validate:"required"`
	IsDeposit   *bool   `json:"is_deposit" validate:"required"`
}

type PointSaleUpdate struct {
	ID          int64   `json:"id" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description"`
	Number			int64   `json:"number" validate:"required"`
	IsDeposit   *bool   `json:"is_deposit" validate:"required"`
}

type PointSaleResponse struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Number			int64   `json:"number"`
	IsDeposit   bool    `json:"is_deposit"`
	IsMain      bool    `json:"is_main"`
}

type PointSaleUpdateMain struct {
	ID      int64 `json:"id" validate:"required" example:"1"`
	NewMain int64 `json:"new_main" validate:"required,nefield=ID" example:"2"`
}

