package notification

import "src/application/usecase/notification/query/search_notification"

func Register() {
	search_notification.Register()
}
