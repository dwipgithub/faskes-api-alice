package auth

import (
	"net/http"
	"time"

	gin "github.com/gin-gonic/gin"
	bcrypt "golang.org/x/crypto/bcrypt"

	auth "go-api-learning/alice/auth"
	authModel "go-api-learning/alice/models/auth"
	repositories "go-api-learning/alice/repositories/consumer"
)

type AuthHandler struct {
	ConsumerRepository *repositories.ConsumerRepository
	JWTService         *auth.JWTService
}

func (h *AuthHandler) Login(c *gin.Context) {

	// --------------------------------------------------------
	// Request
	// --------------------------------------------------------

	var request authModel.LoginRequest

	err := c.ShouldBindJSON(&request)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			authModel.LoginResponse{
				Status:  false,
				Message: "Invalid JSON",
				Data:    nil,
			},
		)

		return
	}

	// --------------------------------------------------------
	// Validasi request
	// --------------------------------------------------------

	if request.UserName == "" {
		c.JSON(
			http.StatusBadRequest,
			authModel.LoginResponse{
				Status:  false,
				Message: "username is required",
				Data:    nil,
			},
		)

		return
	}

	if request.Password == "" {
		c.JSON(
			http.StatusBadRequest,
			authModel.LoginResponse{
				Status:  false,
				Message: "password is required",
				Data:    nil,
			},
		)

		return
	}

	// --------------------------------------------------------
	// Cari consumer
	// --------------------------------------------------------

	consumerID, passwordHash, status, err :=
		h.ConsumerRepository.GetByUserName(
			request.UserName,
		)

	if err != nil {
		c.JSON(
			http.StatusNotFound,
			authModel.LoginResponse{
				Status:  false,
				Message: err.Error(),
				Data:    nil,
			},
		)

		return
	}

	// --------------------------------------------------------
	// Cek status consumer
	// --------------------------------------------------------

	if status != "1" {
		c.JSON(
			http.StatusForbidden,
			authModel.LoginResponse{
				Status:  false,
				Message: "Consumer is inactive",
				Data:    nil,
			},
		)

		return
	}

	// --------------------------------------------------------
	// Verifikasi password
	// --------------------------------------------------------

	err = bcrypt.CompareHashAndPassword(
		[]byte(passwordHash),
		[]byte(request.Password),
	)

	if err != nil {
		c.JSON(
			http.StatusUnauthorized,
			authModel.LoginResponse{
				Status:  false,
				Message: "Invalid password",
				Data:    nil,
			},
		)

		return
	}

	// --------------------------------------------------------
	// Generate access token
	// --------------------------------------------------------

	accessToken, issuedAt, expiredAt, err :=
		h.JWTService.GenerateAccessToken(
			consumerID,
			request.UserName,
		)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			authModel.LoginResponse{
				Status:  false,
				Message: err.Error(),
				Data:    nil,
			},
		)

		return
	}

	// --------------------------------------------------------
	// Generate refresh token
	// --------------------------------------------------------
	refreshToken, _, _, err :=
		h.JWTService.GenerateRefreshToken(
			consumerID,
			request.UserName,
		)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			authModel.LoginResponse{
				Status:  false,
				Message: err.Error(),
				Data:    nil,
			},
		)

		return
	}

	// --------------------------------------------------------
	// Simpan refresh token ke HTTP-only cookie
	// --------------------------------------------------------

	c.SetCookie(
		"refreshToken",
		refreshToken,
		int(h.JWTService.RefreshTokenExpiresIn.Seconds()),
		"/",
		"",
		true,  // HttpOnly
		false, // Secure - false karena masih localhost
	)

	// --------------------------------------------------------
	// Response
	// --------------------------------------------------------

	c.JSON(
		http.StatusCreated,
		authModel.LoginResponse{
			Status:  true,
			Message: "accesstoken created",
			Data: &authModel.LoginData{
				AccessToken: accessToken,
				IssuedAt:    issuedAt.Format(time.RFC3339),
				ExpiredAt:   expiredAt.Format(time.RFC3339),
				ExpiredIn:   h.JWTService.AccessTokenExpiresIn.String(),
			},
		},
	)
}
