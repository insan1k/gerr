package gerr

import (
	"fmt"
)

// newWrapped returns the error as a new package type Wrapped
func newWrapped(k error, opts ...Option) wrapped {
	k, opts = newWrappedFromWrappedError(k, opts...)
	w := wrapped{
		kind: newKind(k),
	}
	for _, opt := range opts {
		w = opt(w)
	}
	return w
}

func newWrappedFromWrappedError(k error, opts ...Option) (error, []Option) {
	if _, ok := k.(interface {
		Unwrap() error
	}); ok {
		chain := Chain(k, WithTrimCustom(separator()))
		if len(chain) > 1 {
			var newOpts []Option
			k = chain[len(chain)-1]
			chain = chain[:len(chain)-1]
			for i := 0; i < len(chain); i++ {
				newOpts = append(newOpts, WithErr(chain[i]))
			}
			opts = append(newOpts, opts...)
		}
	}
	return k, opts
}

type wrapped struct {
	kind kind
	err  error
}

// Sanitize removes all additional context from the Grr
func (w wrapped) Sanitize() Grr {
	return w.kind
}

// Add adds the given error to the error chain
func (w wrapped) Add(err error) Grr {
	return wrapped{
		kind: kind{
			err:       w.kind.err,
			separator: w.kind.separator,
		},
		err: fmt.Errorf("%v%s%w", err, w.kind.separator, w.err),
	}
}

// Grr implements the Wrapped interface
func (w wrapped) Error() string {
	return w.kind.err.Error() + w.kind.separator + w.err.Error()
}

// Is checks if the error is or contains the target error in its chain, therefore this function overrides the
// default implementation of errors.Is
func (w wrapped) Is(target error) bool {
	if w.kind.err.Error() == target.Error() {
		return true
	}
	chain := w.Chain()
	for _, err := range chain {
		if err.Error() == target.Error() {
			return true
		}
	}
	return false
}

// Chain builds the error chain as a slice of error for the package type Wrapped
func (w wrapped) Chain() []error {
	return Chain(w, WithTrimCustom(w.kind.separator))
}

// Unwrap implements the Unwrap interface
func (w wrapped) Unwrap() error {
	return w.err
}
