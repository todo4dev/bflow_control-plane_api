package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AccountStatusEnum string

const (
	AccountStatus_Pending  AccountStatusEnum = "PENDING"
	AccountStatus_Active   AccountStatusEnum = "ACTIVE"
	AccountStatus_Disabled AccountStatusEnum = "DISABLED"
	AccountStatus_Deleted  AccountStatusEnum = "DELETED"
)

type AccountEntity struct {
	ID        uuid.UUID         `json:"id"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	DeletedAt *time.Time        `json:"deleted_at"`
	Email     string            `json:"email"`
	Status    AccountStatusEnum `json:"status"`
}

func (e *AccountEntity) MarshalJSON() ([]byte, error) {
	type Alias AccountEntity
	return json.Marshal((*Alias)(e))
}

func (e *AccountEntity) UnmarshalJSON(data []byte) error {
	type Alias AccountEntity
	return json.Unmarshal(data, (*Alias)(e))
}
