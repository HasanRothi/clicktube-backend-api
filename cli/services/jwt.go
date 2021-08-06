package services

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var JWT_SECRET = os.Getenv("JWT_SECRET")

type jwtCustomClaims struct {
	Gmail string `json:"gmail"`
	jwt.StandardClaims
}

func GenarateJwtToken(gmail string) string {
	claims := &jwtCustomClaims{
		gmail,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token using the secret signing key
	t, err := token.SignedString([]byte(JWT_SECRET))
	if err != nil {
		panic(err)
	}
	return t

}
func VerifyJwtToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Signing method validation
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret signing key
		return []byte(JWT_SECRET), nil
	})

}
