package gerr

// newKind returns an error as the package type Kind
func newKind(err error) kind {
	return kind{
		err:       err,
		separator: separator(),
	}
}

// kind represents the type of error and does not contain any additional context to the error
type kind struct {
	err       error
	separator string
}

func (k kind) Sanitize() Grr {
	return k
}

// Grr implements the error interface
func (k kind) Error() string {
	return k.err.Error()
}

// Add adds the given error to the error chain and returns the error as the package type Wrapped
func (k kind) Add(err error) Grr {
	return newWrapped(k, WithErr(err))
}

// Is returns true if the error is of the given type, or exists in its chain
func (k kind) Is(err error) bool {
	return k.err.Error() == err.Error()
}

// Chain returns the error chain as a slice of errors
func (k kind) Chain() []error {
	return Chain(k, WithTrimCustom(k.separator))
}
