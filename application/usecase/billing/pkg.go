package billing

import (
	"src/application/usecase/billing/command/cancel_subscription"
	"src/application/usecase/billing/command/change_subscription_plan"
	"src/application/usecase/billing/command/handle_stripe_webhook_event"
	"src/application/usecase/billing/command/start_subscription"
	"src/application/usecase/billing/query/search_invoice"
	"src/application/usecase/billing/query/search_payment"
)

func Register() {
	cancel_subscription.Register()
	change_subscription_plan.Register()
	handle_stripe_webhook_event.Register()
	start_subscription.Register()

	search_invoice.Register()
	search_payment.Register()
}
