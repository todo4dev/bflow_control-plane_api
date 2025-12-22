package jwt

import (
	adapter "src/application/adapter/jwt"
	"src/core/di"
	"src/core/env"
	jwt_impl "src/infrastructure/jwt/jwt"
)

func init() {
	di.RegisterAs[adapter.IJwtAdapter](func() adapter.IJwtAdapter {
		config := &adapter.JwtConfig{
			Algorithm:  env.Get("JWT_ALGORITHM", "HS256"),
			Audience:   env.Get("JWT_AUDIENCE", "http://localhost:4000"),
			Issuer:     env.Get("JWT_ISSUER", "http://localhost:4000"),
			AccessTTL:  env.Get("JWT_ACCESS_TTL", "15m"),
			RefreshTTL: env.Get("JWT_REFRESH_TTL", "7d"),
			PrivateKey: env.Get("JWT_PRIVATE_KEY", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
			PublicKey:  env.Get("JWT_PUBLIC_KEY", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
		}
		if err := config.Validate(); err != nil {
			panic(err)
		}

		return jwt_impl.NewJwtAdapter(config)
	})
}
