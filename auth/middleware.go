package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(jwtService *JWTService) gin.HandlerFunc {

	return func(c *gin.Context) {

		// --------------------------------------------------------
		// Ambil Authorization Header
		// --------------------------------------------------------

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"status":  false,
					"message": "Authorization header is required",
					"data":    nil,
				},
			)

			c.Abort()
			return
		}

		// --------------------------------------------------------
		// Pastikan format:
		//
		// Authorization: Bearer <token>
		// --------------------------------------------------------

		parts := strings.SplitN(
			authHeader,
			" ",
			2,
		)

		if len(parts) != 2 ||
			strings.ToLower(parts[0]) != "bearer" {

			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"status":  false,
					"message": "Invalid authorization format",
					"data":    nil,
				},
			)

			c.Abort()
			return
		}

		tokenString := parts[1]

		if tokenString == "" {
			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"status":  false,
					"message": "Access token is required",
					"data":    nil,
				},
			)

			c.Abort()
			return
		}

		// --------------------------------------------------------
		// Parse & Verify JWT
		// --------------------------------------------------------

		token, err := jwt.Parse(
			tokenString,
			func(token *jwt.Token) (interface{}, error) {

				// Pastikan menggunakan HMAC
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}

				return []byte(
					jwtService.AccessTokenSecret,
				), nil
			},
		)

		if err != nil || !token.Valid {

			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"status":  false,
					"message": "Invalid or expired access token",
					"data":    nil,
				},
			)

			c.Abort()
			return
		}

		// --------------------------------------------------------
		// Ambil Claims
		// --------------------------------------------------------

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"status":  false,
					"message": "Invalid token claims",
					"data":    nil,
				},
			)

			c.Abort()
			return
		}

		// --------------------------------------------------------
		// Ambil username
		// --------------------------------------------------------

		username, ok :=
			claims["username"].(string)

		if !ok || username == "" {
			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"status":  false,
					"message": "username not found in token",
					"data":    nil,
				},
			)

			c.Abort()
			return
		}

		// --------------------------------------------------------
		// Simpan data JWT ke Gin Context
		// --------------------------------------------------------

		c.Set(
			"username",
			username,
		)

		if id, ok := claims["id"].(float64); ok {
			c.Set(
				"consumer_id",
				int(id),
			)
		}

		// --------------------------------------------------------
		// Lanjut ke handler
		// --------------------------------------------------------

		c.Next()
	}
}
