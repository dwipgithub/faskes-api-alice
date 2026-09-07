package key

type KeyResponse struct {
	Status  bool    `json:"status"`
	Message string  `json:"message"`
	Data    KeyData `json:"data"`
}

type KeyData struct {
	UserName    string `json:"userName"`
	KeyID       string `json:"keyId"`
	PrivateKey  string `json:"privateKey"`
	PublicKey   string `json:"publicKey"`
	GeneratedAt string `json:"generatedAt"`
}
