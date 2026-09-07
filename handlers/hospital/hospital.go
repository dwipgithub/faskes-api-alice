package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	gin "github.com/gin-gonic/gin"

	crypto "go-api-learning/alice/crypto"
	hospitalModel "go-api-learning/alice/models/hospital"
	consumer "go-api-learning/alice/repositories/consumer"
	hospital "go-api-learning/alice/repositories/hospital"
)

type HospitalHandler struct {
	Repository            *hospital.HospitalRepository
	ConsumerRepository    *consumer.ConsumerRepository
	ConsumerKeyRepository *consumer.ConsumerKeyRepository
}

// GetHospital menangani GET /hospitals
func (h *HospitalHandler) GetHospital(c *gin.Context) {

	// ========================================================
	// 1. Ambil consumer_code dari JWT Middleware
	// ========================================================

	consumerCodeValue, exists := c.Get("username")

	if !exists {
		c.JSON(
			http.StatusUnauthorized,
			hospitalModel.HospitalResponse{
				Status:  false,
				Message: "username not found in authentication context",
				Data:    nil,
			},
		)
		return
	}

	consumerCode, ok := consumerCodeValue.(string)

	if !ok || consumerCode == "" {
		c.JSON(
			http.StatusUnauthorized,
			hospitalModel.HospitalResponse{
				Status:  false,
				Message: "Invalid username",
				Data:    nil,
			},
		)
		return
	}

	// ========================================================
	// 2. Cari consumer berdasarkan username
	// ========================================================

	consumerID, err := h.ConsumerRepository.GetByCode(
		consumerCode,
	)

	if err != nil {
		c.JSON(
			http.StatusNotFound,
			hospitalModel.HospitalResponse{
				Status:  false,
				Message: "Consumer not found",
				Data:    nil,
			},
		)
		return
	}

	// ========================================================
	// 3. Cek apakah consumer menggunakan encryption
	// ========================================================

	encryptionEnabled, err :=
		h.ConsumerRepository.GetEncryptionStatus(
			consumerID,
		)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			hospitalModel.HospitalResponse{
				Status:  false,
				Message: "Failed to get encryption configuration",
				Data:    nil,
			},
		)
		return
	}

	// ========================================================
	// 4. Ambil parameter pagination
	// ========================================================

	page := 1
	limit := 100

	// Jika page dikirim, gunakan nilai tersebut
	if pageValue := c.Query("page"); pageValue != "" {

		page, err = strconv.Atoi(pageValue)

		if err != nil || page < 1 {
			c.JSON(
				http.StatusBadRequest,
				hospitalModel.HospitalResponse{
					Status:  false,
					Message: "Invalid page",
					Data:    nil,
				},
			)
			return
		}
	}

	// Jika limit dikirim, gunakan nilai tersebut
	if limitValue := c.Query("limit"); limitValue != "" {

		limit, err = strconv.Atoi(limitValue)

		if err != nil || limit < 1 {
			c.JSON(
				http.StatusBadRequest,
				hospitalModel.HospitalResponse{
					Status:  false,
					Message: "Invalid limit",
					Data:    nil,
				},
			)
			return
		}

		// Maksimal 100 data per halaman
		if limit > 100 {
			limit = 100
		}
	}

	// ========================================================
	// 5. Ambil filter
	// ========================================================

	filter := hospitalModel.HospitalFilter{
		Page:            page,
		Limit:           limit,
		PelayananNama:   c.Query("pelayanan"),
		ProvinsiID:      c.Query("provinsiId"),
		KabKotaID:       c.Query("kabKotaId"),
		Nama:            c.Query("nama"),
		Aktive:          c.Query("aktive"),
		StartModifiedAt: c.Query("startModifiedAt"),
		EndModifiedAt:   c.Query("endModifiedAt"),
	}

	// ========================================================
	// 6. Ambil data rumah sakit dari repository
	// ========================================================

	hospitals, totalRowCount, err :=
		h.Repository.GetAll(filter)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			hospitalModel.HospitalResponse{
				Status:  false,
				Message: "Invalid limit",
				Data:    nil,
			},
		)
		return
	}

	// ========================================================
	// 7. Hitung jumlah halaman
	// ========================================================

	totalNumberOfPages := 0

	if totalRowCount > 0 {
		totalNumberOfPages =
			(totalRowCount + limit - 1) / limit
	}

	// ========================================================
	// 8. Tentukan halaman berikutnya
	// ========================================================

	var next *hospitalModel.NextPage

	if page < totalNumberOfPages {

		next = &hospitalModel.NextPage{
			Page:  page + 1,
			Limit: limit,
		}
	}

	// ========================================================
	// 9. Buat object pagination
	// ========================================================

	pagination := hospitalModel.Pagination{
		TotalNumberOfPages: totalNumberOfPages,
		CurrentPage:        page,
		Next:               next,
	}

	// ========================================================
	// 10. Jika data tidak ditemukan
	// ========================================================

	if totalRowCount == 0 {

		c.JSON(
			http.StatusOK,
			hospitalModel.HospitalResponse{
				Status:     true,
				Message:    "data not found",
				Pagination: pagination,
				Data:       hospitals,
			},
		)

		return
	}

	// ========================================================
	// 11. Encryption OFF
	// ========================================================

	if !encryptionEnabled {

		c.JSON(
			http.StatusOK,
			hospitalModel.HospitalResponse{
				Status:     true,
				Message:    "data found",
				Pagination: pagination,
				Data:       hospitals,
			},
		)

		return
	}

	// ========================================================
	// 12. Encryption ON
	// ========================================================

	// --------------------------------------------------------
	// Ambil public key yang aktif untuk consumer
	// --------------------------------------------------------

	keyID, publicKeyString, err :=
		h.ConsumerKeyRepository.GetActiveKey(
			consumerID,
		)

	if err != nil {
		c.JSON(
			http.StatusNotFound,
			hospitalModel.HospitalResponse{
				Status:  false,
				Message: "Active encryption key not found",
				Data:    nil,
			},
		)
		return
	}

	// --------------------------------------------------------
	// Convert public key dari string ke EC PublicKey
	// --------------------------------------------------------

	publicKey, err :=
		crypto.ParsePublicKey(publicKeyString)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			hospitalModel.HospitalResponse{
				Status:  false,
				Message: "Invalid public key",
				Data:    nil,
			},
		)
		return
	}

	// --------------------------------------------------------
	// Data yang akan dienkripsi
	// --------------------------------------------------------

	plaintextData := gin.H{
		"pagination": pagination,
		"data":       hospitals,
	}

	// --------------------------------------------------------
	// Convert data menjadi JSON
	// --------------------------------------------------------

	plaintext, err :=
		json.Marshal(plaintextData)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			hospitalModel.HospitalResponse{
				Status:  false,
				Message: "Failed to encode hospital data",
				Data:    nil,
			},
		)
		return
	}

	// --------------------------------------------------------
	// Compress
	// --------------------------------------------------------

	



	// --------------------------------------------------------
	// Encrypt data menggunakan public key
	// --------------------------------------------------------

	encryptedMessage, err :=
		crypto.Encrypt(
			publicKey,
			keyID,
			plaintext,
		)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			hospitalModel.HospitalResponse{
				Status:  false,
				Message: "Encryption failed",
				Data:    nil,
			},
		)
		return
	}

	// ========================================================
	// 13. Response jika encryption ON
	// ========================================================

	c.JSON(
		http.StatusOK,
		hospitalModel.HospitalResponse{
			Status:     true,
			Message:    "data found",
			Pagination: pagination,
			Data:       encryptedMessage,
		},
	)
}
