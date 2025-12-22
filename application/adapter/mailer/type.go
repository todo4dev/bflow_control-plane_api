package mailer

import "context"

// MailAttachment represents an email attachment.
type MailAttachment struct {
	Filename    string
	Content     []byte
	Encoding    string
	ContentType string
}

// MailPayload is the email sending payload.
type MailPayload struct {
	To          []string
	Subject     string
	HTML        string
	Text        string
	CC          []string
	BCC         []string
	ReplyTo     []string
	From        string
	HTTP        string
	Attachments []MailAttachment
}

// IMailerAdapter defines the generic contract for email sending.
type IMailerAdapter interface {
	// Send sends an email based on the provided input.
	Send(ctx context.Context, input MailPayload) error
}
