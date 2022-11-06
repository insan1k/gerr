package gerr

var separator = func() string {
	return defaultSeparator
}

const (
	defaultSeparator       = " "
	colonSeparator         = ":"
	colonAndSpaceSeparator = ": "
)

// SetCustomSeparator sets the separator for the error chain
func SetCustomSeparator(s string) {
	separator = func() string {
		return s
	}
}

// SetColonSeparator sets the separator for the error chain to a colon
func SetColonSeparator() {
	separator = func() string {
		return colonSeparator
	}
}

// SetColonAndSpaceSeparator sets the separator for the error chain to a colon and space
func SetColonAndSpaceSeparator() {
	separator = func() string {
		return colonAndSpaceSeparator
	}
}

// SetSpaceSeparator sets the separator for the error chain to a space
func SetSpaceSeparator() {
	separator = func() string {
		return defaultSeparator
	}
}
