package repository


type CategoryResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type CategoryResponseStock struct {
	ID   int64  `json:"id" gorm:"column:category_id"`
	Name string `json:"name" gorm:"column:category_name"`
}

type CategoryCreate struct {
	Name string `json:"name" validate:"required" example:"Categoria1"`
}

type CategoryUpdate struct {
	ID   int64  `json:"id" validate:"required" example:"1"`
	Name string `json:"name" validate:"required" example:"Categoria1"`
}
