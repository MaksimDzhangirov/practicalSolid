package badCode

import (
	"github.com/golang-jwt/jwt"
	"net/http"
)

func extractUsername(header http.Header) string {
	raw := header.Get("Authorization")
	parser := &jwt.Parser{}
	token, _, err := parser.ParseUnverified(raw, jwt.MapClaims{})
	if err != nil {
		return ""
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ""
	}

	return claims["username"].(string)
}
