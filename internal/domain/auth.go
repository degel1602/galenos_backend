package domain

// AuthClaims son los datos verificados extraídos de un token Bearer JWT
// válido.
type AuthClaims struct {
	Subject string
}
