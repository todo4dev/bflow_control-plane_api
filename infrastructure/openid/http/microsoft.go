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
	microsoftAuthorizeURL = "https://login.microsoftonline.com/common/oauth2/v2.0/authorize"
	microsoftTokenURL     = "https://login.microsoftonline.com/common/oauth2/v2.0/token"
	microsoftUserinfoURL  = "https://graph.microsoft.com/oidc/userinfo"
	microsoftPhotoURL     = "https://graph.microsoft.com/v1.0/me/photo/$value"
	microsoftScope        = "openid profile email offline_access User.Read"
)

type MicrosoftProvider struct {
	config     *adapter.OpenIDConfig
	httpClient *http.Client
}

var _ adapter.IOpenIDProvider = (*MicrosoftProvider)(nil)

func NewMicrosoftProvider(config *adapter.OpenIDConfig) *MicrosoftProvider {
	return &MicrosoftProvider{config: config, httpClient: http.DefaultClient}
}

func (p *MicrosoftProvider) CreateRedirectURI(ctx context.Context, state string) (string, error) {
	if state == "" {
		panic("state cannot be empty")
	}

	redirectURI := p.config.BaseURI + p.config.MicrosoftCallbackURI

	query := url.Values{}
	query.Set("client_id", p.config.MicrosoftClientID)
	query.Set("redirect_uri", redirectURI)
	query.Set("response_type", "code")
	query.Set("scope", microsoftScope)
	query.Set("response_mode", "query")

	u, err := url.Parse(microsoftAuthorizeURL)
	if err != nil {
		return "", err
	}
	u.RawQuery = query.Encode()

	return u.String(), nil
}

func (p *MicrosoftProvider) GetRefreshToken(ctx context.Context, code string) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", p.config.BaseURI+p.config.MicrosoftCallbackURI)
	data.Set("client_id", p.config.MicrosoftClientID)
	data.Set("client_secret", p.config.MicrosoftClientSecret)
	data.Set("scope", microsoftScope)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, microsoftTokenURL, strings.NewReader(data.Encode()))
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
		return "", fmt.Errorf("microsoft GetRefreshToken: status %d: %s", resp.StatusCode, string(body))
	}

	var parsed struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", err
	}

	if parsed.RefreshToken == "" {
		return "", fmt.Errorf("microsoft GetRefreshToken: empty refresh_token")
	}

	return parsed.RefreshToken, nil
}

func (p *MicrosoftProvider) GetToken(ctx context.Context, refreshToken string) (*adapter.OpenIDToken, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", p.config.MicrosoftClientID)
	data.Set("client_secret", p.config.MicrosoftClientSecret)
	data.Set("scope", microsoftScope)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, microsoftTokenURL, strings.NewReader(data.Encode()))
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
		return nil, fmt.Errorf("microsoft GetToken: status %d: %s", resp.StatusCode, string(body))
	}

	var parsed adapter.OpenIDToken
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}

	return &parsed, nil
}

func (p *MicrosoftProvider) firstValue(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func (p *MicrosoftProvider) GetInfo(ctx context.Context, accessToken string) (*adapter.OpenIDInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, microsoftUserinfoURL, nil)
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
		return nil, fmt.Errorf("microsoft GetInfo: status %d: %s", resp.StatusCode, string(body))
	}

	var raw struct {
		Sub         string `json:"sub"`
		Email       string `json:"email"`
		GivenName   string `json:"given_name"`
		Givenname2  string `json:"givenname"`
		GivenName2  string `json:"givenName"`
		FamilyName  string `json:"family_name"`
		FamilyName2 string `json:"familyname"`
		FamilyName3 string `json:"familyName"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	return &adapter.OpenIDInfo{
		Sub:        raw.Sub,
		Email:      raw.Email,
		GivenName:  p.firstValue(raw.GivenName, raw.Givenname2, raw.GivenName2),
		FamilyName: p.firstValue(raw.FamilyName, raw.FamilyName2, raw.FamilyName3),
	}, nil
}

func (p *MicrosoftProvider) GetPicture(ctx context.Context, accessToken string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, microsoftPhotoURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		return nil, fmt.Errorf("microsoft GetPicture: status %d: %s", resp.StatusCode, string(body))
	}

	return resp.Body, nil
}
