package not_acceptable

import "src/core/meta"

func Register() {
	e := New()
	meta.Describe(e,
		meta.Description("Requested representation cannot be served (content negotiation failed)"),
		meta.Field(
			&e.Code,
			meta.Description("Machine-readable error code"),
			meta.Example(DefaultCode)),
		meta.Field(
			&e.Message,
			meta.Description("Human-readable error message"),
			meta.Example(DefaultMessage)),
	)
}
