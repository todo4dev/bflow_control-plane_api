package jwt

import (
	"src/core/validator"
)

// JwtConfig representa a configuração necessária para geração e validação de JWT.
// É equivalente ao JwtConfig do serviço em TypeScript.
type JwtConfig struct {
	Algorithm  string // ex: "RS256"
	Audience   string
	Issuer     string
	AccessTTL  string // ex: "15m"
	RefreshTTL string // ex: "30d"
	PrivateKey string // chave privada (PEM)
	PublicKey  string // chave pública (PEM)
}

var _ validator.IValidable = (*JwtConfig)(nil)

func (c *JwtConfig) Validate() error {
	return validator.Object(c,
		validator.String(&c.Algorithm).Required().Allow("RS256", "HS256").Default("RS256"),
		validator.String(&c.Audience).Required(),
		validator.String(&c.Issuer).Required(),
		validator.String(&c.AccessTTL).Required().Default("15m"),
		validator.String(&c.RefreshTTL).Required().Default("30d"),
		validator.String(&c.PrivateKey).Required(),
		validator.String(&c.PublicKey).Required(),
	).Validate()
}
