package shopify

import (
	"log"

	"github.com/golang-jwt/jwt/v5"
)

// Shopify Docs:
// https://shopify.dev/docs/apps/build/authentication-authorization/session-tokens#anatomy-of-a-session-token

type SessionTokenJWTClaims struct {
	jwt.RegisteredClaims
	// The shop's domain
	Destination string `json:"dest"`
	// A unique session ID per user and app
	SessionId string `json:"sid"`
	// Shopify signature
	Signature string `json:"sig"`
}

func ValidateSessionToken(sessionToken string) (*SessionTokenJWTClaims, error) {
	var claims SessionTokenJWTClaims
	_, err := jwt.ParseWithClaims(sessionToken, &claims, func(t *jwt.Token) (interface{}, error) {
		return nil, nil // TODO: validate with shopify API?
	})
	if err != nil {
		return nil, err
	}
	log.Printf("SessionTokenJWT Payload: %+v", claims)
	return &claims, nil
}
