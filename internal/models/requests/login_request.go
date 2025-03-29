package requests

type LoginRequest struct {
	ID       int64  `json:"id"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required,min=5"`
}
