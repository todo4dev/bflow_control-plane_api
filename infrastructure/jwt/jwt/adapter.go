package jwt

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"

	adapter "src/application/adapter/jwt"
)

//#region Jwt

type JwtAdapter struct {
	config        *adapter.JwtConfig
	privateKey    *rsa.PrivateKey
	publicKey     *rsa.PublicKey
	signingMethod jwtlib.SigningMethod
}

var _ adapter.IJwtAdapter = (*JwtAdapter)(nil)

func NewJwtAdapter(config *adapter.JwtConfig) *JwtAdapter {
	privateKey, err := parseRSAPrivateKeyFromPEM(config.PrivateKey)
	if err != nil {
		panic(fmt.Errorf("invalid private key: %w", err))
	}

	publicKey, err := parseRSAPublicKeyFromPEM(config.PublicKey)
	if err != nil {
		panic(fmt.Errorf("invalid public key: %w", err))
	}

	return &JwtAdapter{config: config, privateKey: privateKey, publicKey: publicKey}
}
func (a *JwtAdapter) Create(ctx context.Context, sessionKey string, info *adapter.OpenIDInfo, complete bool) (adapter.OpenIDToken, error) {
	if err := ctx.Err(); err != nil {
		return adapter.OpenIDToken{}, err
	}

	accessTTLms, err := a.parseTTLToMilliseconds(a.config.AccessTTL)
	if err != nil {
		return adapter.OpenIDToken{}, fmt.Errorf("invalid AccessTTL: %w", err)
	}

	now := time.Now().UTC()

	// claims base para o access token
	accessClaims := a.buildBaseClaims(info, sessionKey)
	accessClaims["iat"] = now.Unix()
	accessClaims["nbf"] = now.Unix()
	accessClaims["exp"] = now.Add(time.Duration(accessTTLms) * time.Millisecond).Unix()

	accessToken, err := a.signWithType(accessClaims, "access")
	if err != nil {
		return adapter.OpenIDToken{}, err
	}

	result := adapter.OpenIDToken{
		TokenType:       "Bearer",
		AccessToken:     accessToken,
		AccessExpiresIn: accessTTLms,
	}

	if complete {
		refreshTTLms, err := a.parseTTLToMilliseconds(a.config.RefreshTTL)
		if err != nil {
			return adapter.OpenIDToken{}, fmt.Errorf("invalid RefreshTTL: %w", err)
		}

		refreshClaims := a.buildBaseClaims(info, sessionKey)
		refreshClaims["iat"] = now.Unix()
		refreshClaims["nbf"] = now.Unix()
		refreshClaims["exp"] = now.Add(time.Duration(refreshTTLms) * time.Millisecond).Unix()

		refreshToken, err := a.signWithType(refreshClaims, "refresh")
		if err != nil {
			return adapter.OpenIDToken{}, err
		}

		result.RefreshToken = &refreshToken
		result.RefreshExpiresIn = &refreshTTLms
	}

	return result, nil
}
func (a *JwtAdapter) Decode(ctx context.Context, tokenString string) (adapter.DecodedToken, error) {
	if err := ctx.Err(); err != nil {
		return adapter.DecodedToken{}, err
	}

	parsedToken, err := jwtlib.Parse(
		tokenString,
		func(token *jwtlib.Token) (any, error) {
			// garante que não mudaram o algoritmo
			if token.Method.Alg() != a.signingMethod.Alg() {
				return nil, fmt.Errorf("unexpected signing method: %s", token.Method.Alg())
			}
			return a.publicKey, nil
		},
		jwtlib.WithAudience(a.config.Audience),
		jwtlib.WithIssuer(a.config.Issuer),
	)
	if err != nil {
		return adapter.DecodedToken{}, err
	}

	if !parsedToken.Valid {
		return adapter.DecodedToken{}, errors.New("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwtlib.MapClaims)
	if !ok {
		return adapter.DecodedToken{}, errors.New("invalid claims type")
	}

	kind, _ := parsedToken.Header["typ"].(string)

	sessionKey, _ := claims["jti"].(string)
	if sessionKey == "" {
		return adapter.DecodedToken{}, errors.New("missing jti (sessionKey) in token")
	}

	openIDInfo := adapter.OpenIDInfo{
		Subject:    a.toString(claims["sub"]),
		Email:      a.toString(claims["email"]),
		FamilyName: a.toString(claims["family_name"]),
		GivenName:  a.toString(claims["given_name"]),
		Language:   a.toString(claims["language"]),
		Picture:    a.toString(claims["picture"]),
		Theme:      a.toString(claims["theme"]),
		Timezone:   a.toString(claims["timezone"]),
	}

	return adapter.DecodedToken{
		Kind:       kind,
		SessionKey: sessionKey,
		OpenIDInfo: openIDInfo,
	}, nil
}
func (a *JwtAdapter) buildBaseClaims(info *adapter.OpenIDInfo, sessionKey string) jwtlib.MapClaims {
	return jwtlib.MapClaims{
		// OpenIDInfo
		"sub":         info.Subject,
		"email":       info.Email,
		"family_name": info.FamilyName,
		"given_name":  info.GivenName,
		"language":    info.Language,
		"picture":     info.Picture,
		"theme":       info.Theme,
		"timezone":    info.Timezone,
		// padrão
		"aud": a.config.Audience,
		"iss": a.config.Issuer,
		"jti": sessionKey,
	}
}
func (a *JwtAdapter) signWithType(claims jwtlib.MapClaims, typ string) (string, error) {
	token := jwtlib.NewWithClaims(a.signingMethod, claims)
	// força o header.typ = "access" | "refresh"
	token.Header["typ"] = typ

	return token.SignedString(a.privateKey)
}
func (a *JwtAdapter) toString(value any) string {
	if value == nil {
		return ""
	}
	if s, ok := value.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", value)
}
func (a *JwtAdapter) parseTTLToMilliseconds(ttl string) (int64, error) {
	ttl = strings.TrimSpace(ttl)
	if ttl == "" {
		return 0, errors.New("ttl is empty")
	}

	multiplier := int64(1) // default: milissegundos
	valuePart := ttl

	last := ttl[len(ttl)-1]
	switch last {
	case 's', 'S':
		multiplier = 1000
		valuePart = ttl[:len(ttl)-1]
	case 'm', 'M':
		multiplier = 60 * 1000
		valuePart = ttl[:len(ttl)-1]
	case 'h', 'H':
		multiplier = 60 * 60 * 1000
		valuePart = ttl[:len(ttl)-1]
	case 'd', 'D':
		multiplier = 24 * 60 * 60 * 1000
		valuePart = ttl[:len(ttl)-1]
	default:
		// sem sufixo -> assume que já veio em milissegundos
	}

	value, err := strconv.ParseInt(valuePart, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid ttl %q: %w", ttl, err)
	}
	if value <= 0 {
		return 0, fmt.Errorf("ttl must be positive, got %d", value)
	}

	return value * multiplier, nil
}
func parseRSAPrivateKeyFromPEM(pemString string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemString))
	if block == nil {
		return nil, errors.New("failed to parse PEM block from private key")
	}

	switch block.Type {
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	case "PRIVATE KEY":
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		privateKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("not an RSA private key")
		}
		return privateKey, nil
	default:
		return nil, fmt.Errorf("unsupported private key type %q", block.Type)
	}
}
func parseRSAPublicKeyFromPEM(pemString string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemString))
	if block == nil {
		return nil, errors.New("failed to parse PEM block from public key")
	}

	switch block.Type {
	case "PUBLIC KEY":
		pubAny, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		publicKey, ok := pubAny.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("not an RSA public key")
		}
		return publicKey, nil
	case "RSA PUBLIC KEY":
		return x509.ParsePKCS1PublicKey(block.Bytes)
	default:
		return nil, fmt.Errorf("unsupported public key type %q", block.Type)
	}
}

//#endregion
