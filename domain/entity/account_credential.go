package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AccountCredentialEntity struct {
	ID              uuid.UUID  `json:"id"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
	CredentialType  string     `json:"credential_type"`
	ProviderSubject string     `json:"provider_subject"`
	PasswordHash    string     `json:"password_hash"`
	AccountID       uuid.UUID  `json:"account_id"`
}

func (e *AccountCredentialEntity) MarshalJSON() ([]byte, error) {
	type Alias AccountCredentialEntity
	return json.Marshal((*Alias)(e))
}

func (e *AccountCredentialEntity) UnmarshalJSON(data []byte) error {
	type Alias AccountCredentialEntity
	return json.Unmarshal(data, (*Alias)(e))
}
