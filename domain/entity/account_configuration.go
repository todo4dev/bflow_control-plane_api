package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AccountConfigurationEntity struct {
	ID        uuid.UUID `json:"id"`
	UpdatedAt time.Time `json:"updated_at"`
	Theme     string    `json:"theme"`
	AccountID uuid.UUID `json:"account_id"`
}

func (e *AccountConfigurationEntity) MarshalJSON() ([]byte, error) {
	type Alias AccountConfigurationEntity
	return json.Marshal((*Alias)(e))
}

func (e *AccountConfigurationEntity) UnmarshalJSON(data []byte) error {
	type Alias AccountConfigurationEntity
	return json.Unmarshal(data, (*Alias)(e))
}
