package not_found

import "src/core/meta"

func Register() {
	notFoundException := New()
	meta.Describe(notFoundException,
		meta.Description("The requested resource could not be found"),
		meta.Example(notFoundException),
		meta.Field(
			&notFoundException.Code,
			meta.Description("Machine-readable error code"),
			meta.Example(DefaultCode)),
		meta.Field(
			&notFoundException.Message,
			meta.Description("Human-readable error message"),
			meta.Example(DefaultMessage)),
	)
}
