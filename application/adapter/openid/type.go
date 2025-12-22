package openid

import (
	"context"
	"io"
)

type OpenIDToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

type OpenIDInfo struct {
	Sub        string `json:"sub"`
	Email      string `json:"email"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
}

type IOpenIDProvider interface {
	CreateRedirectURI(ctx context.Context, state string) (string, error)
	GetRefreshToken(ctx context.Context, code string) (string, error)
	GetToken(ctx context.Context, refreshToken string) (*OpenIDToken, error)
	GetInfo(ctx context.Context, accessToken string) (*OpenIDInfo, error)
	GetPicture(ctx context.Context, accessToken string) (io.ReadCloser, error)
}

type IOpenIDAdapter interface {
	GetProvider(name string) (IOpenIDProvider, error)
	EncodeState(state map[string]string) (string, error)
	DecodeState(state string) (map[string]string, error)
}
