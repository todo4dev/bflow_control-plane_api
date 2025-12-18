package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type BillingPaymentStatusEnum string

const (
	BillingPaymentStatus_Pending   BillingPaymentStatusEnum = "PENDING"
	BillingPaymentStatus_Succeeded BillingPaymentStatusEnum = "SUCCEEDED"
	BillingPaymentStatus_Failed    BillingPaymentStatusEnum = "FAILED"
	BillingPaymentStatus_Refunded  BillingPaymentStatusEnum = "REFUNDED"
)

type BillingPaymentEntity struct {
	ID                    uuid.UUID             `json:"id"`
	CreatedAt             time.Time             `json:"created_at"`
	UpdatedAt             time.Time             `json:"updated_at"`
	Code                  string                `json:"code"`
	Name                  string                `json:"name"`
	Period                BillingPlanPeriodEnum `json:"period"`
	IsActiveFlag          bool                  `json:"is_active_flag"`
	Status                string                `json:"status"`
	Amount                float64               `json:"amount"`
	PaidAt                time.Time             `json:"paid_at"`
	StripePaymentIntentID string                `json:"stripe_payment_intent_id"`
	BillingInvoiceID      uuid.UUID             `json:"billing_invoice_id"`
}

func (e *BillingPaymentEntity) MarshalJSON() ([]byte, error) {
	type Alias BillingPaymentEntity
	return json.Marshal((*Alias)(e))
}

func (e *BillingPaymentEntity) UnmarshalJSON(data []byte) error {
	type Alias BillingPaymentEntity
	return json.Unmarshal(data, (*Alias)(e))
}
