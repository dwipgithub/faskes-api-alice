package routes

import (
	gin "github.com/gin-gonic/gin"

	auth "go-api-learning/alice/auth"
	authHandler "go-api-learning/alice/handlers/auth"
	encryptHandlers "go-api-learning/alice/handlers/encrypt"
	hospitalHandlers "go-api-learning/alice/handlers/hospital"
	keyHandlers "go-api-learning/alice/handlers/key"
)

func Setup(
	router *gin.Engine,
	authHandler *authHandler.AuthHandler,
	encryptHandler *encryptHandlers.EncryptHandler,
	keyHandler *keyHandlers.KeyHandler,
	messageHandler *encryptHandlers.MessageHandler,
	hospitalHandler *hospitalHandlers.HospitalHandler,
	jwtService *auth.JWTService,
) {

	// ========================================================
	// Faskes API
	// ========================================================
	faskes := router.Group("/faskes")

	// ========================================================
	// Public Routes
	// ========================================================
	faskes.POST(
		"/login",
		authHandler.Login,
	)

	// ========================================================
	// Protected Routes
	// ========================================================
	protected := faskes.Group("/")

	protected.Use(
		auth.JWTMiddleware(jwtService),
	)

	protected.GET(
		"/hospitals",
		hospitalHandler.GetHospital,
	)

	protected.POST(
		"/keys/generate",
		keyHandler.GenerateKey,
	)

	router.POST(
		"/encrypt",
		encryptHandler.Encrypt,
	)

	router.GET(
		"/messages/latest",
		messageHandler.GetLatestMessages,
	)
}
