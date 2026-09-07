package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	AccessTokenSecret     string
	RefreshTokenSecret    string
	AccessTokenExpiresIn  time.Duration
	RefreshTokenExpiresIn time.Duration
}

func NewJWTService() (*JWTService, error) {

	accessTokenSecret := os.Getenv("ACCESS_TOKEN_SECRET")
	refreshTokenSecret := os.Getenv("REFRESH_TOKEN_SECRET")

	accessTokenExpiresInString := os.Getenv(
		"ACCESS_TOKEN_EXPIRESIN",
	)

	refreshTokenExpiresInString := os.Getenv(
		"REFRESH_TOKEN_EXPIRESIN",
	)

	if accessTokenSecret == "" {
		return nil, fmt.Errorf(
			"ACCESS_TOKEN_SECRET is not configured",
		)
	}

	if refreshTokenSecret == "" {
		return nil, fmt.Errorf(
			"REFRESH_TOKEN_SECRET is not configured",
		)
	}

	accessTokenExpiresIn, err := time.ParseDuration(
		accessTokenExpiresInString,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"invalid ACCESS_TOKEN_EXPIRESIN: %w",
			err,
		)
	}

	refreshTokenExpiresIn, err := time.ParseDuration(
		refreshTokenExpiresInString,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"invalid REFRESH_TOKEN_EXPIRESIN: %w",
			err,
		)
	}

	return &JWTService{
		AccessTokenSecret:     accessTokenSecret,
		RefreshTokenSecret:    refreshTokenSecret,
		AccessTokenExpiresIn:  accessTokenExpiresIn,
		RefreshTokenExpiresIn: refreshTokenExpiresIn,
	}, nil
}

func (s *JWTService) GenerateAccessToken(
	consumerID int,
	userName string,
) (string, time.Time, time.Time, error) {

	now := time.Now()

	issuedAt := now
	expiredAt := now.Add(s.AccessTokenExpiresIn)

	claims := jwt.MapClaims{
		"id":       consumerID,
		"username": userName,
		"iat":      issuedAt.Unix(),
		"exp":      expiredAt.Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	accessToken, err := token.SignedString(
		[]byte(s.AccessTokenSecret),
	)

	if err != nil {
		return "", time.Time{}, time.Time{}, fmt.Errorf(
			"failed to sign access token: %w",
			err,
		)
	}

	return accessToken, issuedAt, expiredAt, nil
}

func (s *JWTService) GenerateRefreshToken(
	consumerID int,
	userName string,
) (string, time.Time, time.Time, error) {

	now := time.Now()

	issuedAt := now
	expiredAt := now.Add(s.RefreshTokenExpiresIn)

	claims := jwt.MapClaims{
		"id":       consumerID,
		"username": userName,
		"iat":      issuedAt.Unix(),
		"exp":      expiredAt.Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	refreshToken, err := token.SignedString(
		[]byte(s.RefreshTokenSecret),
	)

	if err != nil {
		return "", time.Time{}, time.Time{}, fmt.Errorf(
			"failed to sign refresh token: %w",
			err,
		)
	}

	return refreshToken, issuedAt, expiredAt, nil
}
