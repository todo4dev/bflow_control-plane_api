package system

import "src/application/usecase/system/query/healthcheck"

func Register() {
	healthcheck.Register()
}
