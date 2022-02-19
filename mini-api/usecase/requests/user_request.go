package requests

// UserRequest ...
type UserRequest struct {
	Name           string      `json:"name"`
	Email          string      `json:"email"`
	Password       string      `json:"password"`
	Status         string      `json:"status"`
	RegisterType   string      `json:"register_type"`
	RegisterDetail interface{} `json:"register_detail"`
}
