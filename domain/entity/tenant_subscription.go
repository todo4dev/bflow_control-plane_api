package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type TenantSubscriptionStatusEnum string

const (
	TenantSubscriptionStatus_Trialing TenantSubscriptionStatusEnum = "TRIALING"
	TenantSubscriptionStatus_Active   TenantSubscriptionStatusEnum = "ACTIVE"
	TenantSubscriptionStatus_PastDue  TenantSubscriptionStatusEnum = "PAST_DUE"
	TenantSubscriptionStatus_Canceled TenantSubscriptionStatusEnum = "CANCELED"
)

type TenantSubscriptionEntity struct {
	ID                   uuid.UUID                    `json:"id"`
	CreatedAt            time.Time                    `json:"created_at"`
	UpdatedAt            time.Time                    `json:"updated_at"`
	Status               TenantSubscriptionStatusEnum `json:"status"`
	StartAt              time.Time                    `json:"start_at"`
	EndAt                time.Time                    `json:"end_at"`
	CancelAt             time.Time                    `json:"cancel_at"`
	CanceledAt           time.Time                    `json:"canceled_at"`
	IsCancelAtPeriodEnd  bool                         `json:"is_cancel_at_period_end"`
	StripeSubscriptionID string                       `json:"stripe_subscription_id"`
	TenantID             uuid.UUID                    `json:"tenant_id"`
	BillingPlanID        uuid.UUID                    `json:"billing_plan_id"`
}

func (t *TenantSubscriptionEntity) MarshalJSON() ([]byte, error) {
	type Alias TenantSubscriptionEntity
	return json.Marshal((*Alias)(t))
}

func (t *TenantSubscriptionEntity) UnmarshalJSON(data []byte) error {
	type Alias TenantSubscriptionEntity
	return json.Unmarshal(data, (*Alias)(t))
}
