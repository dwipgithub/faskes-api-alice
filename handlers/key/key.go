package handlers

import (
	"net/http"
	"time"

	gin "github.com/gin-gonic/gin"

	crypto "go-api-learning/alice/crypto"
	keyModel "go-api-learning/alice/models/key"
	consumerRepository "go-api-learning/alice/repositories/consumer"
)

type KeyHandler struct {
	ConsumerRepository    *consumerRepository.ConsumerRepository
	ConsumerKeyRepository *consumerRepository.ConsumerKeyRepository
}

func (h *KeyHandler) GenerateKey(c *gin.Context) {

	// ========================================================
	// 1. Ambil username dari JWT Middleware
	// ========================================================

	usernameValue, exists := c.Get("username")

	if !exists {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"status":  false,
				"message": "username not found in authentication context",
				"data":    nil,
			},
		)

		return
	}

	consumerCode, ok := usernameValue.(string)

	if !ok || consumerCode == "" {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"status":  false,
				"message": "Invalid consumer_code",
				"data":    nil,
			},
		)

		return
	}

	// ========================================================
	// 2. Cari consumer berdasarkan consumer_code
	// ========================================================

	consumerID, err := h.ConsumerRepository.GetByCode(
		consumerCode,
	)

	if err != nil {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  false,
				"message": "Consumer not found",
				"data":    nil,
			},
		)

		return
	}

	// ========================================================
	// 3. Generate Key ID
	// ========================================================

	keyID, err := h.ConsumerKeyRepository.GenerateKeyID()

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  false,
				"message": "Failed to generate key ID",
				"data":    nil,
			},
		)

		return
	}

	// ========================================================
	// 4. Generate ECC Key Pair
	// ========================================================

	generatedKey, err := crypto.GenerateKeyPair(
		keyID,
	)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  false,
				"message": "Failed to generate ECC key pair",
				"data":    nil,
			},
		)

		return
	}

	// ========================================================
	// 5. Nonaktifkan key sebelumnya
	// ========================================================

	err = h.ConsumerKeyRepository.DeactivateActiveKey(
		consumerID,
	)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  false,
				"message": "Failed to deactivate previous key",
				"data":    nil,
			},
		)

		return
	}

	// ========================================================
	// 6. Simpan public key ke database
	// ========================================================

	err = h.ConsumerKeyRepository.Create(
		consumerID,
		generatedKey.KeyID,
		generatedKey.PublicKey,
	)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  false,
				"message": "Failed to save public key",
				"data":    nil,
			},
		)

		return
	}

	// ========================================================
	// 7. Response
	// ========================================================

	c.JSON(
		http.StatusCreated,
		keyModel.KeyResponse{
			Status:  true,
			Message: "Encryption key generated successfully",
			Data: keyModel.KeyData{
				UserName:    consumerCode,
				KeyID:       generatedKey.KeyID,
				PrivateKey:  generatedKey.PrivateKey,
				PublicKey:   generatedKey.PublicKey,
				GeneratedAt: time.Now().Format(time.RFC3339),
			},
		},
	)
}
