package mailer

import "src/core/validator"

type MailerConfig struct {
	DefaultFrom   string
	SMTPHostname  string
	SMTPPort      int
	SMTPSecure    bool
	SMTPIgnoreTLS bool
	SMTPUsername  string
	SMTPPassword  string
}

var _ validator.IValidable = (*MailerConfig)(nil)

func (c *MailerConfig) Validate() error {
	return validator.Object(c,
		validator.String(&c.DefaultFrom).Required(),
		validator.String(&c.SMTPHostname).Required(),
		validator.Number(&c.SMTPPort).Integer().Positive().Required(),
		validator.Boolean(&c.SMTPSecure).Default(false),
		validator.Boolean(&c.SMTPIgnoreTLS).Default(false),
		validator.String(&c.SMTPUsername).Required(),
		validator.String(&c.SMTPPassword).Required(),
	).Validate()
}
