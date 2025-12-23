package precondition_failed

import "src/core/meta"

func Register() {
	e := New()
	meta.Describe(e,
		meta.Description("One or more preconditions required for this operation were not met"),
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
