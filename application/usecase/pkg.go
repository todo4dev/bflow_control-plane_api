package usecase

import (
	"src/application/usecase/activity"
	"src/application/usecase/billing"
	"src/application/usecase/currency"
	"src/application/usecase/document"
	"src/application/usecase/identity"
	"src/application/usecase/notification"
	"src/application/usecase/profile"
	"src/application/usecase/session"
	"src/application/usecase/system"
	"src/application/usecase/tenant"
)

func Register() {
	activity.Register()
	billing.Register()
	currency.Register()
	document.Register()
	identity.Register()
	notification.Register()
	profile.Register()
	session.Register()
	system.Register()
	tenant.Register()
}
