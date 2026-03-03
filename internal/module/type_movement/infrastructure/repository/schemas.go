package repository


type TypeMovementCreate struct {
	Name         string `json:"name" validate:"required"`
	TypeMovement string `json:"type_movement" validate:"required,oneof=income expense"`
}

type TypeMovementUpdate struct {
	ID           int64  `json:"id" validate:"required"`
	Name         string `json:"name" validate:"required"`
	TypeMovement string `json:"type_movement" validate:"required,oneof=income expense"`
}

type TypeMovementResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
