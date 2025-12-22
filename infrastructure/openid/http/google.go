package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	adapter "src/application/adapter/openid"
)

const (
	googleAuthorizeURL = "https://accounts.google.com/o/oauth2/v2/auth"
	googleTokenURL     = "https://oauth2.googleapis.com/token"
	googleUserinfoURL  = "https://openidconnect.googleapis.com/v1/userinfo"
	googleScope        = "openid profile email"
)

type GoogleOpenIDProvider struct {
	config     *adapter.OpenIDConfig
	httpClient *http.Client
}

var _ adapter.IOpenIDProvider = (*GoogleOpenIDProvider)(nil)

func NewGoogleOpenIDProvider(config *adapter.OpenIDConfig) *GoogleOpenIDProvider {
	return &GoogleOpenIDProvider{config: config, httpClient: http.DefaultClient}
}
func (p *GoogleOpenIDProvider) CreateRedirectURI(ctx context.Context, state string) (string, error) {
	if state == "" {
		panic("state cannot be empty")
	}

	redirectURI := p.config.BaseURI + p.config.GoogleCallbackURI

	query := url.Values{}
	query.Set("client_id", p.config.GoogleClientID)
	query.Set("redirect_uri", redirectURI)
	query.Set("response_type", "code")
	query.Set("scope", googleScope)
	query.Set("access_type", "offline")
	query.Set("prompt", "consent")

	u, _ := url.Parse(googleAuthorizeURL)
	u.RawQuery = query.Encode()

	return u.String(), nil
}

func (p *GoogleOpenIDProvider) GetRefreshToken(ctx context.Context, code string) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", p.config.BaseURI+p.config.GoogleCallbackURI)
	data.Set("client_id", p.config.GoogleClientID)
	data.Set("client_secret", p.config.GoogleClientSecret)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, googleTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "app/x-www-form-urlencoded")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("google GetRefreshToken: status %d: %s", resp.StatusCode, string(body))
	}

	var parsed struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", err
	}

	if parsed.RefreshToken == "" {
		return "", fmt.Errorf("google GetRefreshToken: empty refresh_token")
	}

	return parsed.RefreshToken, nil
}

func (p *GoogleOpenIDProvider) GetToken(ctx context.Context, refreshToken string) (*adapter.OpenIDToken, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", p.config.GoogleClientID)
	data.Set("client_secret", p.config.GoogleClientSecret)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, googleTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "app/x-www-form-urlencoded")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("google GetToken: status %d: %s", resp.StatusCode, string(body))
	}

	var parsed adapter.OpenIDToken
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}

	return &parsed, nil
}

func (p *GoogleOpenIDProvider) GetInfo(ctx context.Context, accessToken string) (*adapter.OpenIDInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, googleUserinfoURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("google GetInfo: status %d: %s", resp.StatusCode, string(body))
	}

	var raw adapter.OpenIDInfo
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}
	return &raw, nil
}

func (p *GoogleOpenIDProvider) GetPicture(ctx context.Context, accessToken string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, googleUserinfoURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("google GetPicture(userinfo): status %d: %s", resp.StatusCode, string(body))
	}

	var raw struct {
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}
	if raw.Picture == "" {
		return nil, fmt.Errorf("google GetPicture: empty picture url")
	}

	picReq, err := http.NewRequestWithContext(ctx, http.MethodGet, raw.Picture, nil)
	if err != nil {
		return nil, err
	}

	picResp, err := p.httpClient.Do(picReq)
	if err != nil {
		return nil, err
	}

	if picResp.StatusCode < 200 || picResp.StatusCode >= 300 {
		body, _ := io.ReadAll(picResp.Body)
		_ = picResp.Body.Close()
		return nil, fmt.Errorf("google GetPicture(fetch): status %d: %s", picResp.StatusCode, string(body))
	}

	return picResp.Body, nil
}
