package shopify

// Shopify Docs:
// https://shopify.dev/docs/apps/build/authentication-authorization/access-tokens/authorization-code-grant

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	Path_Authorize   = "/admin/oauth/authorize"
	Path_AccessToken = "/admin/oauth/access_token"
)

type AccessMode string

const (
	OnlineAccessMode  AccessMode = "ONLINE"
	OfflineAccessMode AccessMode = "OFFLINE"
)

func (a *App) GetAuthorizeUrl(shop Shop, accessMode AccessMode) (string, error) {
	state, err := a.GenerateNonce()
	if err != nil {
		return "", err
	}
	shopUrl, err := url.Parse(shop.BaseUrl())
	if err != nil {
		log.Printf("[goshopify] Unexpected error on AuthorizeUrl: %v", err)
		return "", err
	}
	shopUrl.Path = Path_Authorize

	q := shopUrl.Query()
	q.Set("client_id", a.cfg.ClientId)
	q.Set("redirect_uri", a.cfg.RedirectUrl)
	q.Set("scope", a.cfg.Scopes)
	q.Set("state", state)
	if accessMode == OnlineAccessMode {
		q.Set("grant_options[]", "per-user")
	}
	shopUrl.RawQuery = q.Encode()
	return shopUrl.String(), nil
}

func (a *App) GetAccessToken(ctx context.Context, shop Shop, code string) (*AccessToken, error) {
	reqBody := map[string]string{
		"client_id":     a.cfg.ClientId,
		"client_secret": a.cfg.ClientSecret,
		"code":          code,
	}
	b, err := json.Marshal(reqBody)
	if err != nil {
		// shouldn't be an error
		return nil, err
	}

	body := bytes.NewBuffer(b)
	requestUrl := shop.BaseUrl() + Path_AccessToken

	req, err := http.NewRequestWithContext(ctx, "POST", requestUrl, body)
	if err != nil {
		return nil, err
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 400 {
		log.Printf("%s returns with status code %d", requestUrl, resp.StatusCode)
		return nil, errors.New("shopify APi returns with error status code")
	}

	var accessToken AccessTokenBase
	if err := json.NewDecoder(resp.Body).Decode(&accessToken); err != nil {
		return nil, err
	}
	if accessToken.AccessToken == "" {
		return nil, errors.New("access token returned is empty")
	}
	ext := AccessToken{
		AccessTokenBase: accessToken,
	}
	if accessToken.ExpiresIn != nil {
		now := time.Now()
		eat := now.Add(time.Duration(*accessToken.ExpiresIn) * time.Second)
		ext.ExpiresAt = &eat
	}

	return &ext, nil
}

func (a *App) VerifyAuthorizationQuery(query url.Values) (bool, error) {
	hmacVal := query.Get("hmac")
	if hmacVal == "" {
		return false, nil
	}

	// Remove HMAC and Signature values from query
	query.Del("hmac")
	query.Del("signature")

	msg, err := url.QueryUnescape(query.Encode())
	if err != nil {
		return false, err
	}

	// Compare HMAC value with the message
	mac := hmac.New(sha256.New, []byte(a.cfg.ClientSecret))
	mac.Write([]byte(msg))
	expectedMAC := mac.Sum(nil)

	// shopify HMAC is in hex so it needs to be decoded
	actualMac, err := hex.DecodeString(hmacVal)
	if err != nil {
		return false, err
	}

	return hmac.Equal(actualMac, expectedMAC), nil
}
