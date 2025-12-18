package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AccountSessionEntity struct {
	ID               uuid.UUID  `json:"id"`
	CreatedAt        time.Time  `json:"created_at"`
	ExpiresAt        *time.Time `json:"expires_at"`
	RevokedAt        *time.Time `json:"revoked_at"`
	RefreshTokenHash string     `json:"refresh_token_hash"`
	UserAgent        string     `json:"user_agent"`
	IP               string     `json:"ip"`
	AccountID        uuid.UUID  `json:"account_id"`
}

func (e *AccountSessionEntity) MarshalJSON() ([]byte, error) {
	type Alias AccountSessionEntity
	return json.Marshal((*Alias)(e))
}

func (e *AccountSessionEntity) UnmarshalJSON(data []byte) error {
	type Alias AccountSessionEntity
	return json.Unmarshal(data, (*Alias)(e))
}
