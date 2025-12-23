package tenant

import (
	"src/application/usecase/tenant/command/add_membership_to_tenant"
	"src/application/usecase/tenant/command/delete_tenant"
	"src/application/usecase/tenant/command/delete_tenant_picture"
	"src/application/usecase/tenant/command/remove_membership_from_tenant"
	"src/application/usecase/tenant/command/update_membership_role_in_tenant"
	"src/application/usecase/tenant/command/update_tenant"
	"src/application/usecase/tenant/command/update_tenant_configuration"
	"src/application/usecase/tenant/command/upload_tenant_picture"
	"src/application/usecase/tenant/query/check_subdomain_availability"
	"src/application/usecase/tenant/query/get_tenant_configuration"
	"src/application/usecase/tenant/query/get_tenant_picture"
	"src/application/usecase/tenant/query/list_tenant_membership"
	"src/application/usecase/tenant/query/search_tenant"
)

func Register() {
	add_membership_to_tenant.Register()
	delete_tenant.Register()
	delete_tenant_picture.Register()
	remove_membership_from_tenant.Register()
	update_membership_role_in_tenant.Register()
	update_tenant.Register()
	update_tenant_configuration.Register()
	upload_tenant_picture.Register()

	check_subdomain_availability.Register()
	get_tenant_configuration.Register()
	get_tenant_picture.Register()
	list_tenant_membership.Register()
	search_tenant.Register()
}
