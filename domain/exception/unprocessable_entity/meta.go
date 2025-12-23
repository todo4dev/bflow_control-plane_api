package unprocessable_entity

import "src/core/meta"

func Register() {
	e := New()
	meta.Describe(e,
		meta.Description("The request was well-formed but contains semantic errors and could not be processed"),
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
