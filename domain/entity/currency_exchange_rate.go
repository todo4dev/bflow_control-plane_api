package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type CurrencyExchangeRateSourceEnum string

const (
	CurrencyExchangeRateSource_Manual   CurrencyExchangeRateSourceEnum = "MANUAL"
	CurrencyExchangeRateSource_Provider CurrencyExchangeRateSourceEnum = "PROVIDER"
)

type CurrencyExchangeRateEntity struct {
	ID                uuid.UUID                      `json:"id"`
	CreatedAt         time.Time                      `json:"created_at"`
	EffectiveAt       time.Time                      `json:"effective_at"`
	Rate              float64                        `json:"rate"`
	Source            CurrencyExchangeRateSourceEnum `json:"source"`
	BaseCurrencyCode  string                         `json:"base_currency_code"`
	QuoteCurrencyCode string                         `json:"quote_currency_code"`
}

func (e *CurrencyExchangeRateEntity) MarshalJSON() ([]byte, error) {
	type Alias CurrencyExchangeRateEntity
	return json.Marshal((*Alias)(e))
}

func (e *CurrencyExchangeRateEntity) UnmarshalJSON(data []byte) error {
	type Alias CurrencyExchangeRateEntity
	return json.Unmarshal(data, (*Alias)(e))
}
