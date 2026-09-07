package encrypt

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	crypto "go-api-learning/alice/crypto"
	consumer "go-api-learning/alice/repositories/consumer"
)

type EncryptHandler struct {
	ConsumerRepository    *consumer.ConsumerRepository
	ConsumerKeyRepository *consumer.ConsumerKeyRepository
	MessageStore          *MessageStore
}

func (h *EncryptHandler) Encrypt(c *gin.Context) {

	// --------------------------------------------------------
	// Baca request
	// --------------------------------------------------------

	var request struct {
		ConsumerCode string                 `json:"consumer_code"`
		Data         map[string]interface{} `json:"data"`
	}

	err := c.ShouldBindJSON(&request)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "Invalid JSON",
			},
		)

		return
	}

	// --------------------------------------------------------
	// Validasi consumer_code
	// --------------------------------------------------------

	if request.ConsumerCode == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "consumer_code is required",
			},
		)

		return
	}

	// --------------------------------------------------------
	// Cari consumer
	// --------------------------------------------------------

	consumerID, err := h.ConsumerRepository.GetByCode(
		request.ConsumerCode,
	)

	if err != nil {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	// --------------------------------------------------------
	// Ambil active public key
	// --------------------------------------------------------

	keyID, publicKeyString, err :=
		h.ConsumerKeyRepository.GetActiveKey(
			consumerID,
		)

	if err != nil {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	// --------------------------------------------------------
	// Convert public key
	// --------------------------------------------------------

	publicKey, err := crypto.ParsePublicKey(
		publicKeyString,
	)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "Invalid public key",
			},
		)

		return
	}

	// --------------------------------------------------------
	// JSON → []byte
	// --------------------------------------------------------

	plaintext, err := json.Marshal(
		request.Data,
	)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "Failed to encode data",
			},
		)

		return
	}

	// --------------------------------------------------------
	// Encrypt
	// --------------------------------------------------------

	encryptedMessage, err := crypto.Encrypt(
		publicKey,
		keyID,
		plaintext,
	)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "Encryption failed",
			},
		)

		return
	}

	// --------------------------------------------------------
	// Simpan encrypted message
	// --------------------------------------------------------

	h.MessageStore.AddMessage(
		encryptedMessage,
	)

	// --------------------------------------------------------
	// Response
	// --------------------------------------------------------

	c.JSON(
		http.StatusOK,
		gin.H{
			"consumer_code": request.ConsumerCode,
			"key_id":        keyID,
			"message":       encryptedMessage,
		},
	)
}
