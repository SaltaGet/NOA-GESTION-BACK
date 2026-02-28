package schemas

type ClientCreate struct {
	FirstName   string  `json:"first_name" validate:"required" example:"Jorge"`
	LastName    string  `json:"last_name" validate:"required" example:"Lopez"`
	CompanyName *string `json:"company_name" example:"John Company | null"`
	Identifier  *string `json:"identifier" example:"30000000 | null"`
	Email       *string `json:"email" validate:"omitempty,email" example:"john@example.com | null"`
	Phone       *string `json:"phone" example:"1111111111 | null"`
	Address     *string `json:"address" example:" Calle 123 | null"`
	ResponsabilityFrontIVA *string `json:"responsability_front_iva" example:"responsable_inscripto | monotributo | consumidor_final | null"`
}

type ClientUpdate struct {
	ID          int64   `json:"id" validate:"required"`
	FirstName   string  `json:"first_name" validate:"required" example:"Jorge"`
	LastName    string  `json:"last_name" validate:"required" example:"Lopez"`
	CompanyName *string `json:"company_name" example:"John Company | null"`
	Identifier  *string `json:"identifier" example:"30000000 | null"`
	Email       *string `json:"email" validate:"omitempty,email" example:"john@example.com | null"`
	Phone       *string `json:"phone" example:"1111111111 | null"`
	Address     *string `json:"address" example:" Calle 123 | null"`
	ResponsabilityFrontIVA *string `json:"responsability_front_iva" example:"responsable_inscripto | monotributo | consumidor_final | null"`
}

type ClientUpdateCredit struct {
	ID        int64        `json:"id" validate:"required" example:"2"`
	PayCredit []*PayCredit `json:"pay_credit" validate:"required,min=1,dive"`
	Total     float64      `json:"total" validate:"required"`
}

type PayCredit struct {
	CreditID  int64   `json:"credit_id" validate:"required"`
	MethodPay string  `json:"method_pay" validate:"oneof=cash card transfer" example:"cash card transfer"`
	Total     float64 `json:"total" validate:"required"`
}

type ClientResponseDTO struct {
	ID          int64    `json:"id"`
	FirstName   string   `json:"first_name"`
	LastName    string   `json:"last_name"`
	CompanyName *string  `json:"company_name,omitempty"`
	Identifier  *string  `json:"identifier,omitempty"`
	Email       *string  `json:"email,omitempty"`
	Phone       *string  `json:"phone,omitempty"`
	Debt      *float64 `json:"debt,omitempty"`
	ResponsabilityFrontIVA *string `json:"responsability_front_iva,omitempty"`
}

type ClientResponse struct {
	ID           int64             `json:"id"`
	FirstName    string            `json:"first_name"`
	LastName     string            `json:"last_name"`
	CompanyName  *string           `json:"company_name,omitempty"`
	Identifier   *string           `json:"identifier,omitempty"`
	Email        *string           `json:"email,omitempty"`
	Phone        *string           `json:"phone,omitempty"`
	Address      *string           `json:"address,omitempty"`
	MemberCreate *MemberSimpleDTO  `json:"member,omitempty"`
	ResponsabilityFrontIVA *string `json:"responsability_front_iva,omitempty"`
	Pay          []PayDebtResponse `json:"pay"`
}

type ClientSimpleDTO struct {
	ID          int64   `json:"id"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	CompanyName *string `json:"company_name,omitempty"`
}
