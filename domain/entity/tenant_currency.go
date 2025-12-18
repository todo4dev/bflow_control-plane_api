package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type TenantCurrencyEntity struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CurrencyCode string    `json:"currency_code"`
	IsEnabled    bool      `json:"is_enabled"`
	TenantID     uuid.UUID `json:"tenant_id"`
}

func (t *TenantCurrencyEntity) MarshalJSON() ([]byte, error) {
	type Alias TenantCurrencyEntity
	return json.Marshal((*Alias)(t))
}

func (t *TenantCurrencyEntity) UnmarshalJSON(data []byte) error {
	type Alias TenantCurrencyEntity
	return json.Unmarshal(data, (*Alias)(t))
}
