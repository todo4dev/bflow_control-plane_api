package jwt

import "context"

// OpenIDInfo representa as claims públicas que vamos colocar dentro do JWT.
// É equivalente ao TOpenidInfo do código em TypeScript.
type OpenIDInfo struct {
	Subject    string `json:"sub"`
	Email      string `json:"email"`
	FamilyName string `json:"family_name"`
	GivenName  string `json:"given_name"`
	Language   string `json:"language"`
	Picture    string `json:"picture"`
	Theme      string `json:"theme"`
	Timezone   string `json:"timezone"`
}

// OpenIDToken representa o token emitido pela aplicação.
// É equivalente ao TOpenidToken do código em TypeScript.
type OpenIDToken struct {
	TokenType        string  `json:"token_type"`                   // ex: "Bearer"
	AccessToken      string  `json:"access_token"`                 // JWT de acesso
	AccessExpiresIn  int64   `json:"access_expires_in"`            // duração em milissegundos
	RefreshToken     *string `json:"refresh_token,omitempty"`      // JWT de refresh (opcional)
	RefreshExpiresIn *int64  `json:"refresh_expires_in,omitempty"` // duração em milissegundos
}

// DecodedToken é o resultado da operação de decode/verify de um JWT.
// É equivalente ao retorno do método decode do serviço em TypeScript.
type DecodedToken struct {
	Kind       string     `json:"kind"`        // "access" ou "refresh" (vem do header.typ)
	SessionKey string     `json:"session_key"` // jti
	OpenIDInfo OpenIDInfo `json:"openid_info"`
}

// IJwtAdapter define o contrato genérico para emissão e validação de JWTs.
// Implementações concretas ficam em infra/jwt/*.
type IJwtAdapter interface {
	// Create gera um par de tokens (access e, opcionalmente, refresh).
	// complete=true indica que o refresh_token também deve ser gerado.
	Create(ctx context.Context, sessionKey string, info *OpenIDInfo, complete bool) (OpenIDToken, error)

	// Decode valida o token (assinatura, issuer, audience, etc.) e devolve
	// as informações relevantes para a aplicação.
	Decode(ctx context.Context, token string) (DecodedToken, error)
}
