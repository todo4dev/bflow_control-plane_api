package session

import (
	"src/application/usecase/session/command/logout"
	"src/application/usecase/session/command/refresh_authentication_token"
	"src/application/usecase/session/query/get_current_session"
)

func Register() {
	logout.Register()
	refresh_authentication_token.Register()

	get_current_session.Register()
}
