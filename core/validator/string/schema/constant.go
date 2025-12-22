package schema

import "regexp"

var (
	RegexHostname    = regexp.MustCompile(`^(?i:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?)(?:\.(?i:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?))*$`)
	RegexAlphanum    = regexp.MustCompile(`^[A-Za-z0-9]+$`)
	RegexDataURI     = regexp.MustCompile(`^data:[^;]+(;base64)?,.+$`)
	RegexDomain      = regexp.MustCompile(`^([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`)
	RegexEmail       = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
	RegexGUID        = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)
	RegexHex         = regexp.MustCompile(`^[0-9a-fA-F]+$`)
	RegexIsoDuration = regexp.MustCompile(`^P(?:(?:\d+Y)?(?:\d+M)?(?:\d+W)?(?:\d+D)?(?:T(?:(?:\d+H)?(?:\d+M)?(?:\d+S)?)?)?)$`)
	RegexToken       = regexp.MustCompile(`^\w+$`)
)
