package conflict

import "src/core/meta"

func Register() {
	e := New()
	meta.Describe(e,
		meta.Description("Request could not be completed due to a conflict with the current state of the resource"),
		meta.Example(e),
		meta.Field(&e.Code,
			meta.Description("Machine-readable error code"),
			meta.Example(DefaultCode)),
		meta.Field(&e.Message,
			meta.Description("Human-readable error message"),
			meta.Example(DefaultMessage)),
	)
}
