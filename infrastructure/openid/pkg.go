package openid

import (
	adapter "src/application/adapter/openid"
	"src/core/di"
	"src/core/env"
	http_impl "src/infrastructure/openid/http"
)

func init() {
	di.RegisterAs[adapter.IOpenIDAdapter](func() adapter.IOpenIDAdapter {
		config := &adapter.OpenIDConfig{
			BaseURI:               env.Get("OPENID_BASE_URI", "http://localhost:4000"),
			MicrosoftClientID:     env.Get("OPENID_MICROSOFT_CLIENT_ID", "{{OPENID_MICROSOFT_CLIENT_ID}}"),
			MicrosoftClientSecret: env.Get("OPENID_MICROSOFT_CLIENT_SECRET", "{{OPENID_MICROSOFT_CLIENT_ID}}"),
			MicrosoftCallbackURI:  env.Get("OPENID_MICROSOFT_CALLBACK_URI", "{{OPENID_MICROSOFT_CLIENT_ID}}"),
			GoogleClientID:        env.Get("OPENID_GOOGLE_CLIENT_ID", "{{OPENID_MICROSOFT_CLIENT_ID}}"),
			GoogleClientSecret:    env.Get("OPENID_GOOGLE_CLIENT_SECRET", "{{OPENID_MICROSOFT_CLIENT_ID}}"),
			GoogleCallbackURI:     env.Get("OPENID_GOOGLE_CALLBACK_URI", "{{OPENID_MICROSOFT_CLIENT_ID}}"),
		}
		if err := config.Validate(); err != nil {
			panic(err)
		}

		return http_impl.NewOpenidAdapter(config)
	})
}
