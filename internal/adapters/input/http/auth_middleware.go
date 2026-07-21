package httpadapter

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/galenos-pro/appointments-api/internal/ports/input"
)

const bearerPrefix = "Bearer "

// RequireBearerToken construye un middleware Gin que exige un JWT válido en
// el header "Authorization: Bearer <token>" para todas las rutas del grupo
// donde se registre. El subject del token queda disponible en el contexto
// bajo la clave "authSubject" para los handlers que lo necesiten.
func RequireBearerToken(authService input.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, bearerPrefix) {
			respondError(c, http.StatusUnauthorized, "MISSING_TOKEN", "se requiere un token Bearer en el header Authorization")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(header, bearerPrefix)
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			respondError(c, http.StatusUnauthorized, "INVALID_TOKEN", "el token es invalido o expiro")
			c.Abort()
			return
		}

		c.Set("authSubject", claims.Subject)
		c.Next()
	}
}
