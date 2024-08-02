package shopify

import (
	"context"
	"errors"
	"time"
)

var (
	ErrAccessTokenNotFound = errors.New("access token not found")
)

type AccessTokenBase struct {
	AccessToken         string          `json:"access_token"`
	Scope               string          `json:"scope"`                 // e.g. "write_orders,read_customers"
	ExpiresIn           *int64          `json:"expires_in"`            // Online Access Mode only
	AssociatedUserScope *string         `json:"associated_user_scope"` // Online Access Mode only
	AssociatedUser      *AssociatedUser `json:"associated_user"`       // Online Access Mode only
}

type AccessToken struct {
	AccessTokenBase
	ExpiresAt *time.Time `json:"expires_at"`
}

type AssociatedUser struct {
	ID            int64  `json:"id"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	AccountOwner  bool   `json:"account_owner"`
	Locale        string `json:"locale"`
	Collaborator  bool   `json:"collaborator"`
}

type AccessTokenRepository interface {
	FindByShop(ctx context.Context, shop string) ([]AccessToken, error)
	Upsert(ctx context.Context, accessToken AccessToken) error
}
