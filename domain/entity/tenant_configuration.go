package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type TenantConfigurationEntity struct {
	ID                  uuid.UUID `json:"id"`
	UpdatedAt           time.Time `json:"updated_at"`
	DefaultTimezone     string    `json:"default_timezone"`
	DefaultCurrencyCode string    `json:"default_currency_code"`
	TenantID            uuid.UUID `json:"tenant_id"`
}

func (e *TenantConfigurationEntity) MarshalJSON() ([]byte, error) {
	type Alias TenantConfigurationEntity
	return json.Marshal((*Alias)(e))
}

func (e *TenantConfigurationEntity) UnmarshalJSON(data []byte) error {
	type Alias TenantConfigurationEntity
	return json.Unmarshal(data, (*Alias)(e))
}
