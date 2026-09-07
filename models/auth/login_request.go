package auth

type LoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}
