package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AccountCertificateEntity struct {
	ID                uuid.UUID  `json:"id"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	ExpiresAt         *time.Time `json:"expires_at"`
	Fingerprint       string     `json:"fingerprint"`
	Subject           string     `json:"subject"`
	Issuer            string     `json:"issuer"`
	SerialNumber      string     `json:"serial_number"`
	EncryptedPayload  string     `json:"encrypted_payload"`
	EncryptedPassword string     `json:"encrypted_password"`
	RevokedAt         *time.Time `json:"revoked_at"`
	AccountID         uuid.UUID  `json:"account_id"`
}

func (e *AccountCertificateEntity) MarshalJSON() ([]byte, error) {
	type Alias AccountCertificateEntity
	return json.Marshal((*Alias)(e))
}

func (e *AccountCertificateEntity) UnmarshalJSON(data []byte) error {
	type Alias AccountCertificateEntity
	return json.Unmarshal(data, (*Alias)(e))
}
