package openid

import (
	"src/core/validator"
)

type OpenIDConfig struct {
	BaseURI               string
	MicrosoftClientID     string
	MicrosoftClientSecret string
	MicrosoftCallbackURI  string
	GoogleClientID        string
	GoogleClientSecret    string
	GoogleCallbackURI     string
}

var _ validator.IValidable = (*OpenIDConfig)(nil)

func (config *OpenIDConfig) Validate() error {
	return validator.Object(config,
		validator.String(&config.BaseURI).Required().Default("http://localhost:4000"),
		validator.String(&config.MicrosoftClientID).Required().Default("{{MicrosoftClientID}}"),
		validator.String(&config.MicrosoftClientSecret).Required().Default("{{MicrosoftClientSecret}}"),
		validator.String(&config.MicrosoftCallbackURI).Required().Default("{{MicrosoftCallbackURI}}"),
		validator.String(&config.GoogleClientID).Required().Default("{{GoogleClientID}}"),
		validator.String(&config.GoogleClientSecret).Required().Default("{{GoogleClientSecret}}"),
		validator.String(&config.GoogleCallbackURI).Required().Default("{{GoogleCallbackURI}}"),
	).Validate()
}
