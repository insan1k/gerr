package gerr

// New returns the package type Wrapped as a standard error if the given arguments contain a wrapped error, or supply a
// error to wrap using WithErr, otherwise it returns Kind which implements the package Grr interface.
func New(kind error, opts ...Option) Grr {
	e := newWrapped(kind, opts...)
	if e.err != nil {
		return e
	}
	return e.kind
}

// WithErr embeds an error on the package type Wrapped, when calling WithErr multiple times note that the first call
// will be considered as the original error, subsequent calls will wrap that error. Additionally, if the given kind
// error contains a wrapped error, the given error will be added to the chain of the wrapped error and will not
// overwrite the original error.
func WithErr(err error) Option {
	return func(w wrapped) wrapped {
		if w.err != nil {
			w = w.Add(err).(wrapped)
		} else {
			w.err = err
		}
		return w
	}
}

// Option is the functional type for configuring a package type Grr
type Option func(w wrapped) wrapped

// Grr represents the interface that all error helper types implement
type Grr interface {
	// Add adds the given error to the error chain
	Add(err error) Grr
	// Error implements the error interface
	Error() string
	// Chain returns the error chain as a slice of errors
	Chain() []error
	// Is returns true if the error is of the given type, or exists in its chain
	Is(err error) bool
	// Sanitize removes all additional context from the Grr
	Sanitize() Grr
}

// AsGrr converts the given error to the package type Grr, if the given error is nil, or is not of the package type Grr
func AsGrr(err error) Grr {
	if err == nil {
		return nil
	}
	if g, ok := err.(Grr); ok {
		return g
	}
	return New(err)
}
