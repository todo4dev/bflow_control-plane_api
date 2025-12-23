package profile

import (
	"src/application/usecase/profile/command/delete_picture"
	"src/application/usecase/profile/command/update_profile"
	"src/application/usecase/profile/command/upload_picture"
	"src/application/usecase/profile/query/get_picture"
)

func Register() {
	delete_picture.Register()
	update_profile.Register()
	upload_picture.Register()

	get_picture.Register()
}
