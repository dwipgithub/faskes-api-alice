package auth

type LoginResponse struct {
	Status  bool       `json:"status"`
	Message string     `json:"message"`
	Data    *LoginData `json:"data"`
}

type LoginData struct {
	AccessToken string `json:"accessToken"`
	IssuedAt    string `json:"issuedAt"`
	ExpiredAt   string `json:"expiredAt"`
	ExpiredIn   string `json:"expiredIn"`
}
