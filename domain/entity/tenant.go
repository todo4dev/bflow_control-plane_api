package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type TenantEntity struct {
	ID               uuid.UUID `json:"id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	DeletedAt        time.Time `json:"deleted_at"`
	Subdomain        string    `json:"subdomain"`
	Name             string    `json:"name"`
	Status           string    `json:"status"`
	Picture          *string   `json:"picture"`
	StripeCustomerID string    `json:"stripe_customer_id"`
}

func (t *TenantEntity) MarshalJSON() ([]byte, error) {
	type Alias TenantEntity
	return json.Marshal((*Alias)(t))
}

func (t *TenantEntity) UnmarshalJSON(data []byte) error {
	type Alias TenantEntity
	return json.Unmarshal(data, (*Alias)(t))
}
