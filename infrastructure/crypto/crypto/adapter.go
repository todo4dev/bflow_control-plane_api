package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	adapter "src/application/adapter/crypto"
)

type CryptoAdapter struct {
	key []byte
}

var _ adapter.ICryptoAdapter = (*CryptoAdapter)(nil)

// NewCryptoAdapter cria uma nova instância usando a chave do CryptoConfig.
func NewCryptoAdapter(config *adapter.CryptoConfig) *CryptoAdapter {
	return &CryptoAdapter{key: []byte(config.Key)}
}

// OTP gera um código curto pseudo-aleatório, 6 caracteres hex maiúsculos.
func (a *CryptoAdapter) OTP() string {
	b := make([]byte, 3) // 3 bytes → 6 hex chars
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	s := strings.ToUpper(hex.EncodeToString(b))
	return s[:6]
}

// Hash gera SHA-256 do texto em claro e retorna em base64url.
func (a *CryptoAdapter) Hash(plainText string) string {
	sum := sha256.Sum256([]byte(plainText))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

// Encrypt serializa o valor como JSON e cifra com AES-256-GCM.
// optionalIV, se não vazio, é usado como IV em texto (mesma semântica do TS).
// Retorna string no formato "iv.encrypted.tag".
func (a *CryptoAdapter) Encrypt(plainText any, optionalIV ...string) string {
	plainBytes, err := json.Marshal(plainText)
	if err != nil {
		panic(fmt.Errorf("crypto: failed to marshal plainText: %w", err))
	}

	block, err := aes.NewCipher(a.key)
	if err != nil {
		panic(fmt.Errorf("crypto: failed to create cipher: %w", err))
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(fmt.Errorf("crypto: failed to create GCM: %w", err))
	}

	var ivStr string
	if len(optionalIV) > 0 {
		ivStr = optionalIV[0]
	} else {
		nonce := make([]byte, 12)
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			panic(fmt.Errorf("crypto: failed to generate nonce: %w", err))
		}
		ivStr = base64.RawURLEncoding.EncodeToString(nonce)
	}

	nonce := []byte(ivStr)

	cipherWithTag := gcm.Seal(nil, nonce, plainBytes, nil)
	tagSize := gcm.Overhead()
	if len(cipherWithTag) < tagSize {
		panic("crypto: ciphertext too short")
	}

	cipherBytes := cipherWithTag[:len(cipherWithTag)-tagSize]
	tagBytes := cipherWithTag[len(cipherWithTag)-tagSize:]

	encrypted := base64.RawURLEncoding.EncodeToString(cipherBytes)
	tag := base64.RawURLEncoding.EncodeToString(tagBytes)

	return strings.Join([]string{ivStr, encrypted, tag}, ".")
}

// Decrypt faz o caminho inverso de Encrypt: separa iv/encrypted/tag, valida GCM e faz JSON unmarshal.
func (a *CryptoAdapter) Decrypt(cipherText string) (any, error) {
	parts := strings.Split(cipherText, ".")
	if len(parts) != 3 {
		return nil, errors.New("crypto: invalid cipherText format")
	}

	ivStr, encryptedStr, tagStr := parts[0], parts[1], parts[2]

	cipherBytes, err := base64.RawURLEncoding.DecodeString(encryptedStr)
	if err != nil {
		return nil, fmt.Errorf("crypto: invalid encrypted part: %w", err)
	}

	tagBytes, err := base64.RawURLEncoding.DecodeString(tagStr)
	if err != nil {
		return nil, fmt.Errorf("crypto: invalid tag part: %w", err)
	}

	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, fmt.Errorf("crypto: failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("crypto: failed to create GCM: %w", err)
	}

	nonce := []byte(ivStr)
	combined := append(cipherBytes, tagBytes...)

	plainBytes, err := gcm.Open(nil, nonce, combined, nil)
	if err != nil {
		return nil, fmt.Errorf("crypto: failed to decrypt: %w", err)
	}

	var out any
	if err := json.Unmarshal(plainBytes, &out); err != nil {
		return nil, fmt.Errorf("crypto: failed to unmarshal plain JSON: %w", err)
	}

	return out, nil
}
