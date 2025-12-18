package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type MembershipRoleEnum string

const (
	MembershipRole_Admin   MembershipRoleEnum = "ADMIN"
	MembershipRole_Manager MembershipRoleEnum = "MANAGER"
	MembershipRole_Member  MembershipRoleEnum = "MEMBER"
)

type MembershipStatusEnum string

const (
	MembershipStatus_Invited   MembershipStatusEnum = "INVITED"
	MembershipStatus_Active    MembershipStatusEnum = "ACTIVE"
	MembershipStatus_Suspended MembershipStatusEnum = "SUSPENDED"
	MembershipStatus_Removed   MembershipStatusEnum = "REMOVED"
)

type MembershipEntity struct {
	ID                    uuid.UUID `json:"id"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	RemovedAt             time.Time `json:"removed_at"`
	Role                  string    `json:"role"`
	Status                string    `json:"status"`
	InvitedByMembershipID uuid.UUID `json:"invited_by_membership_id"`
	TenantID              uuid.UUID `json:"tenant_id"`
	AccountID             uuid.UUID `json:"account_id"`
}

func (e *MembershipEntity) MarshalJSON() ([]byte, error) {
	type Alias MembershipEntity
	return json.Marshal((*Alias)(e))
}

func (e *MembershipEntity) UnmarshalJSON(data []byte) error {
	type Alias MembershipEntity
	return json.Unmarshal(data, (*Alias)(e))
}
