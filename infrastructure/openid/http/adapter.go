package http

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	adapter "src/application/adapter/openid"
)

type OpenidAdapter struct {
	googleProvider    adapter.IOpenIDProvider
	microsoftProvider adapter.IOpenIDProvider
}

var _ adapter.IOpenIDAdapter = (*OpenidAdapter)(nil)

func NewOpenidAdapter(config *adapter.OpenIDConfig) *OpenidAdapter {
	return &OpenidAdapter{
		googleProvider:    NewGoogleOpenIDProvider(config),
		microsoftProvider: NewMicrosoftProvider(config),
	}
}

func (a *OpenidAdapter) GetProvider(name string) (adapter.IOpenIDProvider, error) {
	switch name {
	case "google":
		return a.googleProvider, nil
	case "microsoft":
		return a.microsoftProvider, nil
	default:
		return nil, fmt.Errorf("openid GetProvider: unsupported provider %q", name)
	}
}

func (a *OpenidAdapter) EncodeState(state map[string]string) (string, error) {
	if state == nil {
		state = make(map[string]string)
	}

	data, err := json.Marshal(state)
	if err != nil {
		return "", err
	}

	encoded := base64.RawURLEncoding.EncodeToString(data)
	return encoded, nil
}

func (a *OpenidAdapter) DecodeState(state string) (map[string]string, error) {
	if state == "" {
		return map[string]string{}, nil
	}

	data, err := base64.RawURLEncoding.DecodeString(state)
	if err != nil {
		return nil, err
	}

	var decoded map[string]string
	if err := json.Unmarshal(data, &decoded); err != nil {
		return nil, err
	}

	if decoded == nil {
		decoded = map[string]string{}
	}

	return decoded, nil
}
