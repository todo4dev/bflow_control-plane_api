package unauthorized

import "src/core/meta"

func Register() {
	e := New()
	meta.Describe(e,
		meta.Description("Authentication is required or the provided credentials are invalid"),
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
