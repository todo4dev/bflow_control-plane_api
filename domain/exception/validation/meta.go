package validation

import "src/core/meta"

func Register() {
	e := New()
	meta.Describe(e,
		meta.Description("The request contains invalid or malformed data"),
		meta.Example(e),
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
