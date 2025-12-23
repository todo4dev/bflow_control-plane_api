package forbidden

import "src/core/meta"

func Register() {
	e := New()
	meta.Describe(e,
		meta.Description("The authenticated user does not have permission to perform this operation"),
		meta.Example(e),
		meta.Field(&e.Code,
			meta.Description("Machine-readable error code"),
			meta.Example(DefaultCode)),
		meta.Field(&e.Message,
			meta.Description("Human-readable error message"),
			meta.Example(DefaultMessage)),
	)
}
