package goodCode

import (
	"github.com/golang-jwt/jwt"
	"net/http"
)

func extractRawToken(header http.Header) string {
	return header.Get("Authorization")
}

func extractClaims(raw string) jwt.MapClaims {
	parser := &jwt.Parser{}
	token, _, err := parser.ParseUnverified(raw, jwt.MapClaims{})
	if err != nil {
		return nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil
	}

	return claims
}

func extractUsername(header http.Header) string {
	raw := extractRawToken(header)
	claims := extractClaims(raw)
	if claims == nil {
		return ""
	}

	return claims["username"].(string)
}
