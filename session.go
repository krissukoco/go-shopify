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

// Example SessionTokenJWTClaims
// {
// 	"iss"=>"https://exampleshop.myshopify.com/admin",
// 	"dest"=>"https://exampleshop.myshopify.com",
// 	"aud"=>"client-id-123",
// 	"sub"=>"42",
// 	"exp"=>1591765058,
// 	"nbf"=>1591764998,
// 	"iat"=>1591764998,
// 	"jti"=>"f8912129-1af6-4cad-9ca3-76b0f7621087",
// 	"sid"=>"aaea182f2732d44c23057c0fea584021a4485b2bd25d3eb7fd349313ad24c685",
// 	"sig"=>"f07cf3740270c17fb61c700b2f0f2e7f2f4fc8cc48426221738f7a39e4c475bf",
// }

func (a *App) ValidateSessionToken(sessionToken string) (*SessionTokenJWTClaims, error) {
	var claims SessionTokenJWTClaims
	_, err := jwt.ParseWithClaims(sessionToken, &claims, func(t *jwt.Token) (interface{}, error) {
		return nil, nil // TODO: validate the claims?
	})
	if err != nil {
		return nil, err
	}
	log.Printf("SessionTokenJWT Payload: %+v", claims)
	return &claims, nil
}
