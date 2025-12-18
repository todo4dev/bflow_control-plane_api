package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type MembershipInvitationRoleEnum string

const (
	MembershipInvitationRole_Admin   MembershipInvitationRoleEnum = "ADMIN"
	MembershipInvitationRole_Manager MembershipInvitationRoleEnum = "MANAGER"
	MembershipInvitationRole_Member  MembershipInvitationRoleEnum = "MEMBER"
)

type MembershipInvitationEntity struct {
	ID                    uuid.UUID                    `json:"id"`
	CreatedAt             time.Time                    `json:"created_at"`
	ExpiresAt             time.Time                    `json:"expires_at"`
	AcceptedAt            time.Time                    `json:"accepted_at"`
	RevokedAt             time.Time                    `json:"revoked_at"`
	InvitedEmailAddress   string                       `json:"invited_email_address"`
	Role                  MembershipInvitationRoleEnum `json:"role"`
	InvitedByMembershipID uuid.UUID                    `json:"invited_by_membership_id"`
	InvitationTokenHash   string                       `json:"invitation_token_hash"`
	TenantID              uuid.UUID                    `json:"tenant_id"`
}

func (e *MembershipInvitationEntity) MarshalJSON() ([]byte, error) {
	type Alias MembershipInvitationEntity
	return json.Marshal((*Alias)(e))
}

func (e *MembershipInvitationEntity) UnmarshalJSON(data []byte) error {
	type Alias MembershipInvitationEntity
	return json.Unmarshal(data, (*Alias)(e))
}
