package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type BillingPlanPeriodEnum string

const (
	BillingPlanPeriod_Monthly BillingPlanPeriodEnum = "MONTHLY"
	BillingPlanPeriod_Yearly  BillingPlanPeriodEnum = "YEARLY"
)

type BillingPlanEntity struct {
	ID           uuid.UUID             `json:"id"`
	CreatedAt    time.Time             `json:"created_at"`
	UpdatedAt    time.Time             `json:"updated_at"`
	Code         string                `json:"code"`
	Name         string                `json:"name"`
	Period       BillingPlanPeriodEnum `json:"period"`
	IsActiveFlag bool                  `json:"is_active_flag"`
}

func (b *BillingPlanEntity) MarshalJSON() ([]byte, error) {
	type Alias BillingPlanEntity
	return json.Marshal((*Alias)(b))
}

func (b *BillingPlanEntity) UnmarshalJSON(data []byte) error {
	type Alias BillingPlanEntity
	return json.Unmarshal(data, (*Alias)(b))
}
