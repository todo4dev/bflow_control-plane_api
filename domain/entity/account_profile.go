package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AccountProfileEntity struct {
	ID        uuid.UUID `json:"id"`
	UpdatedAt time.Time `json:"updated_at"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Picture   *string   `json:"picture"`
	AccountID uuid.UUID `json:"account_id"`
}

func (e *AccountProfileEntity) MarshalJSON() ([]byte, error) {
	type Alias AccountProfileEntity
	return json.Marshal((*Alias)(e))
}

func (e *AccountProfileEntity) UnmarshalJSON(data []byte) error {
	type Alias AccountProfileEntity
	return json.Unmarshal(data, (*Alias)(e))
}
