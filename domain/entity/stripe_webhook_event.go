package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type StripeWebhookEventEntity struct {
	ID               uuid.UUID `json:"id"`
	CreatedAt        time.Time `json:"created_at"`
	ProcessedAt      time.Time `json:"processed_at"`
	StripeEventID    string    `json:"stripe_event_id"`
	StripeEventType  string    `json:"stripe_event_type"`
	StripeObjectType string    `json:"stripe_object_type"`
	StripeObjectID   string    `json:"stripe_object_id"`
	TenantID         uuid.UUID `json:"tenant_id"`
}

func (e *StripeWebhookEventEntity) MarshalJSON() ([]byte, error) {
	type Alias StripeWebhookEventEntity
	return json.Marshal((*Alias)(e))
}

func (e *StripeWebhookEventEntity) UnmarshalJSON(data []byte) error {
	type Alias StripeWebhookEventEntity
	return json.Unmarshal(data, (*Alias)(e))
}
