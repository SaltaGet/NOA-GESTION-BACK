package schemas

type Response struct {
	Status bool `json:"status"`
	Body any `json:"body"`
	Message string `json:"message"`
}