package requests

//RegisterByEmailRequest request by email
type RegisterByEmailRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required"`
}
