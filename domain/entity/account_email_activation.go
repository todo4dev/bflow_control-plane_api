package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AccountEmailActivationEntity struct {
	ID        uuid.UUID  `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at"`
	UsedAt    *time.Time `json:"used_at"`
	TokenHash string     `json:"token_hash"`
	AccountID uuid.UUID  `json:"account_id"`
}

func (e *AccountEmailActivationEntity) MarshalJSON() ([]byte, error) {
	type Alias AccountEmailActivationEntity
	return json.Marshal((*Alias)(e))
}

func (e *AccountEmailActivationEntity) UnmarshalJSON(data []byte) error {
	type Alias AccountEmailActivationEntity
	return json.Unmarshal(data, (*Alias)(e))
}
