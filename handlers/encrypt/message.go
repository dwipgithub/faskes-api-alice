package encrypt

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	"go-api-learning/alice/crypto"
)

type StoredMessage struct {
	ID                 int    `json:"id"`
	Ciphertext         string `json:"ciphertext"`
	Nonce              string `json:"nonce"`
	EphemeralPublicKey string `json:"ephemeral_public_key"`
}

type MessageStore struct {
	mu sync.RWMutex

	LatestKeyID string

	Messages map[string][]StoredMessage
}

type MessageHandler struct {
	Store *MessageStore
}

func NewMessageStore() *MessageStore {

	return &MessageStore{
		Messages: make(
			map[string][]StoredMessage,
		),
	}
}

func (s *MessageStore) AddMessage(
	message *crypto.EncryptedMessage,
) {

	s.mu.Lock()
	defer s.mu.Unlock()

	keyID := message.KeyID

	s.LatestKeyID = keyID

	messages := s.Messages[keyID]

	id := len(messages) + 1

	storedMessage := StoredMessage{
		ID:                 id,
		Ciphertext:         message.Ciphertext,
		Nonce:              message.Nonce,
		EphemeralPublicKey: message.EphemeralPublicKey,
	}

	s.Messages[keyID] = append(
		messages,
		storedMessage,
	)
}

func (s *MessageStore) GetLatestMessages() (
	string,
	[]StoredMessage,
) {

	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.LatestKeyID == "" {
		return "", []StoredMessage{}
	}

	messages := s.Messages[s.LatestKeyID]

	return s.LatestKeyID, messages
}

func (h *MessageHandler) GetLatestMessages(
	c *gin.Context,
) {

	keyID, messages := h.Store.GetLatestMessages()

	if keyID == "" {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error": "No encrypted messages available",
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"key_id":   keyID,
			"messages": messages,
		},
	)
}
