package usecase

import (
	"context"
	"crypto/subtle"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
)

type authUseCase struct {
	username string
	password string
	secret   []byte
	ttl      time.Duration
}

// NewAuthUseCase construye el caso de uso de autenticación. Las credenciales
// y el secreto de firma se inyectan desde configuración (variables de
// entorno), ya que esta API no mantiene una tabla de usuarios propia.
func NewAuthUseCase(username, password, secret string, ttl time.Duration) input.AuthService {
	return &authUseCase{
		username: username,
		password: password,
		secret:   []byte(secret),
		ttl:      ttl,
	}
}

func (uc *authUseCase) Login(ctx context.Context, username, password string) (string, error) {
	// Comparación en tiempo constante para no filtrar por timing si el
	// usuario o la contraseña son correctos parcialmente.
	usernameMatches := subtle.ConstantTimeCompare([]byte(username), []byte(uc.username)) == 1
	passwordMatches := subtle.ConstantTimeCompare([]byte(password), []byte(uc.password)) == 1
	if !usernameMatches || !passwordMatches {
		return "", domain.ErrInvalidCredentials
	}

	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   username,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(uc.ttl)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(uc.secret)
	if err != nil {
		return "", fmt.Errorf("signing jwt: %w", err)
	}

	return signed, nil
}

func (uc *authUseCase) ValidateToken(tokenString string) (domain.AuthClaims, error) {
	var claims jwt.RegisteredClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return uc.secret, nil
	})
	if err != nil || !token.Valid {
		return domain.AuthClaims{}, domain.ErrInvalidToken
	}

	return domain.AuthClaims{Subject: claims.Subject}, nil
}
