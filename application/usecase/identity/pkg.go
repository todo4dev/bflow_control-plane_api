package identity

import (
	"src/application/usecase/identity/command/activate_email"
	"src/application/usecase/identity/command/complete_sso_callback"
	"src/application/usecase/identity/command/delete_account"
	"src/application/usecase/identity/command/login_with_email_and_password"
	"src/application/usecase/identity/command/login_with_email_otp"
	"src/application/usecase/identity/command/login_with_sso_token"
	"src/application/usecase/identity/command/register_account_with_email"
	"src/application/usecase/identity/command/resend_activation_email"
	"src/application/usecase/identity/command/reset_password"
	"src/application/usecase/identity/command/start_password_recovery"
	"src/application/usecase/identity/command/start_sso_login"
	"src/application/usecase/identity/query/check_email_availability"
	"src/application/usecase/identity/query/get_account_by_id"
)

func Register() {
	activate_email.Register()
	complete_sso_callback.Register()
	delete_account.Register()
	login_with_email_and_password.Register()
	login_with_email_otp.Register()
	login_with_sso_token.Register()
	register_account_with_email.Register()
	resend_activation_email.Register()
	reset_password.Register()
	start_password_recovery.Register()
	start_sso_login.Register()

	check_email_availability.Register()
	get_account_by_id.Register()
}
