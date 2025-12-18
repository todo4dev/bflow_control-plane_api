package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type BillingInvoiceEntity struct {
	ID                    uuid.UUID             `json:"id"`
	CreatedAt             time.Time             `json:"created_at"`
	UpdatedAt             time.Time             `json:"updated_at"`
	Code                  string                `json:"code"`
	Name                  string                `json:"name"`
	Period                BillingPlanPeriodEnum `json:"period"`
	IsActiveFlag          bool                  `json:"is_active_flag"`
	Status                string                `json:"status"`
	CurrencyCode          string                `json:"currency_code"`
	TotalAmount           float64               `json:"total_amount"`
	TotalTaxAmount        float64               `json:"total_tax_amount"`
	TotalDiscountAmount   float64               `json:"total_discount_amount"`
	IssuedAt              time.Time             `json:"issued_at"`
	DueAt                 time.Time             `json:"due_at"`
	PaidAt                time.Time             `json:"paid_at"`
	StripeInvoiceID       string                `json:"stripe_invoice_id"`
	StripePaymentIntentID string                `json:"stripe_payment_intent_id"`
	TenantID              uuid.UUID             `json:"tenant_id"`
	TenantSubscriptionID  uuid.UUID             `json:"tenant_subscription_id"`
}

func (e *BillingInvoiceEntity) MarshalJSON() ([]byte, error) {
	type Alias BillingInvoiceEntity
	return json.Marshal((*Alias)(e))
}

func (e *BillingInvoiceEntity) UnmarshalJSON(data []byte) error {
	type Alias BillingInvoiceEntity
	return json.Unmarshal(data, (*Alias)(e))
}
