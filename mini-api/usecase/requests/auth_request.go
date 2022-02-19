package requests

//RegisterByEmailRequest request by email
type RegisterByEmailRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// UserLoginRequest ....
type UserLoginRequest struct {
	RegisterType string `json:"register_type" validate:"required"`
	User         string `json:"user"`
	Password     string `json:"password" validate:"required"`
}
