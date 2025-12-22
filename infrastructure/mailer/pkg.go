package mailer

import (
	adapter "src/application/adapter/mailer"
	"src/core/di"
	"src/core/env"
	impl "src/infrastructure/mailer/smtp"
)

func init() {
	di.RegisterAs[adapter.IMailerAdapter](func() adapter.IMailerAdapter {
		config := &adapter.MailerConfig{
			DefaultFrom:   env.Get("MAILER_DEFAULT_FROM", "no-reply@localhost"),
			SMTPHostname:  env.Get("MAILER_SMTP_HOSTNAME", "localhost"),
			SMTPPort:      env.Get("MAILER_SMTP_PORT", 25),
			SMTPSecure:    env.Get("MAILER_SMTP_SECURE", false),
			SMTPIgnoreTLS: env.Get("MAILER_SMTP_IGNORE_TLS", false),
			SMTPUsername:  env.Get("MAILER_SMTP_USERNAME", "{{MAILER_SMTP_USERNAME}}"),
			SMTPPassword:  env.Get("MAILER_SMTP_PASSWORD", "{{MAILER_SMTP_PASSWORD}}"),
		}
		if err := config.Validate(); err != nil {
			panic(err)
		}
		return impl.NewSMTPMailerAdapter(config)
	})
}
