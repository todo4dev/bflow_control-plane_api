package method_not_allowed

import "src/core/meta"

func Register() {
	e := New()
	meta.Describe(e,
		meta.Description("HTTP method is not allowed for the requested resource"),
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
