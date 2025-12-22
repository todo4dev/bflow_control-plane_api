package crypto

import (
	adapter "src/application/adapter/crypto"
	"src/core/di"
	"src/core/env"
	impl "src/infrastructure/crypto/crypto"
)

func init() {
	di.RegisterAs[adapter.ICryptoAdapter](func() adapter.ICryptoAdapter {
		config := &adapter.CryptoConfig{
			Key: env.Get("CRYPTO_KEY", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
		}
		if err := config.Validate(); err != nil {
			panic(err)
		}
		return impl.NewCryptoAdapter(config)
	})
}
