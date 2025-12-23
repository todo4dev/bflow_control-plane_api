package currency

import (
	"src/application/usecase/currency/command/update_exchange_rate"
	"src/application/usecase/currency/query/get_exchange_rate"
	"src/application/usecase/currency/query/list_currency"
)

func Register() {
	update_exchange_rate.Register()

	get_exchange_rate.Register()
	list_currency.Register()
}
