package shopify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	GrantType_TokenExchange = "urn:ietf:params:oauth:grant-type:token-exchange"

	SubjectTokenType_IdToken = "urn:ietf:params:oauth:token-type:id_token"

	TokenType_OfflineAccessToken = "urn:shopify:params:oauth:token-type:offline-access-token"
	TokenType_OnlineAccessToken  = "urn:shopify:params:oauth:token-type:online-access-token"
)

func (a *App) ExchangeToken(shop Shop, sessionToken string, accessMode AccessMode) (*AccessToken, error) {
	reqTokenType := TokenType_OfflineAccessToken
	if accessMode == OnlineAccessMode {
		reqTokenType = TokenType_OnlineAccessToken
	}
	payload := map[string]string{
		"client_id":            a.cfg.ClientId,
		"client_secret":        a.cfg.ClientSecret,
		"grant_type":           GrantType_TokenExchange,
		"subject_token":        sessionToken,
		"subject_token_type":   SubjectTokenType_IdToken,
		"requested_token_type": reqTokenType,
	}

	// Perform HTTP request to Shopify API
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	body := bytes.NewBuffer(b)
	reqUrl := shop.BaseUrl() + Path_Authorize

	req, err := http.NewRequest("POST", reqUrl, body)
	if err != nil {
		return nil, err
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("%w: %d", ErrShopifyAPIStatusCode, resp.StatusCode)
	}

	var base AccessTokenBase
	if err := json.NewDecoder(resp.Body).Decode(&base); err != nil {
		return nil, err
	}

	if base.AccessToken == "" {
		return nil, errors.New("emplty access token")
	}

	accessToken := AccessToken{AccessTokenBase: base}
	if base.ExpiresIn != nil {
		now := time.Now()
		eat := now.Add(time.Duration(*base.ExpiresIn) * time.Second)
		accessToken.ExpiresAt = &eat
	}

	return &accessToken, nil
}
