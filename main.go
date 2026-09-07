package main

import (
	"fmt"

	gin "github.com/gin-gonic/gin"

	auth "go-api-learning/alice/auth"
	database "go-api-learning/alice/database"
	authHandler "go-api-learning/alice/handlers/auth"
	encryptHandlers "go-api-learning/alice/handlers/encrypt"
	hospitalHandler "go-api-learning/alice/handlers/hospital"
	keyHandler "go-api-learning/alice/handlers/key"
	consumer "go-api-learning/alice/repositories/consumer"
	hospitalRepositories "go-api-learning/alice/repositories/hospital"
	routes "go-api-learning/alice/routes"
)

func main() {

	// --------------------------------------------------------
	// Router
	// Membuat HTTP server menggunakan Gin.
	// --------------------------------------------------------

	router := gin.Default()

	// --------------------------------------------------------
	// Database
	// Test koneksi ke MySQL.
	// --------------------------------------------------------

	db, err := database.Connect()

	if err != nil {
		panic(err)
	}

	defer db.Close()

	fmt.Println("Database connected")

	// --------------------------------------------------------
	// Hospital Repository
	// Mengakses data rumah sakit dari MySQL.
	// --------------------------------------------------------

	hospitalRepository := &hospitalRepositories.HospitalRepository{
		DB: db,
	}

	// --------------------------------------------------------
	// Consumer Repository
	// Mengakses data consumer.
	// --------------------------------------------------------

	consumerRepository := &consumer.ConsumerRepository{
		DB: db,
	}

	// --------------------------------------------------------
	// Consumer Key Repository
	// Mengakses public key consumer.
	// --------------------------------------------------------

	consumerKeyRepository := &consumer.ConsumerKeyRepository{
		DB: db,
	}

	// --------------------------------------------------------
	// Hospital Handler
	// Menangani endpoint data rumah sakit.
	// --------------------------------------------------------

	hospitalHandler := &hospitalHandler.HospitalHandler{
		Repository:            hospitalRepository,
		ConsumerRepository:    consumerRepository,
		ConsumerKeyRepository: consumerKeyRepository,
	}

	// JWT Service
	jwtService, err := auth.NewJWTService()

	// Auth Handler
	authHandler := &authHandler.AuthHandler{
		ConsumerRepository: consumerRepository,
		JWTService:         jwtService,
	}

	// --------------------------------------------------------
	// Key Handler
	// Menangani proses generate key pair.
	// --------------------------------------------------------

	keyHandler := &keyHandler.KeyHandler{
		ConsumerRepository:    consumerRepository,
		ConsumerKeyRepository: consumerKeyRepository,
	}

	// --------------------------------------------------------
	// Message Store
	// Menyimpan sementara encrypted message di memory.
	// --------------------------------------------------------

	messageStore := encryptHandlers.NewMessageStore()

	// --------------------------------------------------------
	// Encrypt Handler
	// Menangani proses encryption data.
	// --------------------------------------------------------

	encryptHandler := &encryptHandlers.EncryptHandler{
		ConsumerRepository:    consumerRepository,
		ConsumerKeyRepository: consumerKeyRepository,
		MessageStore:          messageStore,
	}

	// --------------------------------------------------------
	// Message Handler
	// Menangani pengambilan encrypted message terbaru.
	// --------------------------------------------------------

	messageHandler := &encryptHandlers.MessageHandler{
		Store: messageStore,
	}

	// --------------------------------------------------------
	// Routes
	// Mendaftarkan endpoint ke Gin Router.
	// --------------------------------------------------------

	routes.Setup(
		router,
		authHandler,
		encryptHandler,
		keyHandler,
		messageHandler,
		hospitalHandler,
		jwtService,
	)

	// --------------------------------------------------------
	// Server
	// Menjalankan Alice pada port 8081.
	// --------------------------------------------------------

	fmt.Println(
		"Alice running on http://localhost:8081",
	)

	router.Run(":8081")
}
