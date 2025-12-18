package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AccountPasswordResetEntity struct {
	ID        uuid.UUID  `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at"`
	UsedAt    *time.Time `json:"used_at"`
	TokenHash string     `json:"token_hash"`
	AccountID uuid.UUID  `json:"account_id"`
}

func (e *AccountPasswordResetEntity) MarshalJSON() ([]byte, error) {
	type Alias AccountPasswordResetEntity
	return json.Marshal((*Alias)(e))
}

func (e *AccountPasswordResetEntity) UnmarshalJSON(data []byte) error {
	type Alias AccountPasswordResetEntity
	return json.Unmarshal(data, (*Alias)(e))
}
