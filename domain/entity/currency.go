package entity

import (
	"encoding/json"
	"time"
)

type CurrencyEntity struct {
	Code      string     `json:"code"`
	Name      string     `json:"name"`
	MinorUnit int        `json:"minor_unit"`
	Symbol    string     `json:"symbol"`
	IsActive  bool       `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (e *CurrencyEntity) MarshalJSON() ([]byte, error) {
	type Alias CurrencyEntity
	return json.Marshal((*Alias)(e))
}

func (e *CurrencyEntity) UnmarshalJSON(data []byte) error {
	type Alias CurrencyEntity
	return json.Unmarshal(data, (*Alias)(e))
}
