package gerr

import (
	"fmt"
)

// newWrapped returns the error as a new package type Wrapped
func newWrapped(k error, opts ...Option) wrapped {
	sep := _separator()
	k, opts = newWrappedFromWrappedError(k, sep, opts...)
	w := wrapped{
		kind: kind{
			err:       k,
			separator: sep,
		},
	}
	for _, opt := range opts {
		w = opt(w)
	}
	return w
}

func newWrappedFromWrappedError(k error, sep string, opts ...Option) (error, []Option) {
	if _, ok := k.(interface {
		Unwrap() error
	}); ok {
		chain := Chain(k, WithTrimCustom(sep), WithBottomFirst())
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
	if w.err == nil {
		return w.kind.err.Error()
	}
	return w.kind.err.Error() + w.kind.separator + w.err.Error()
}

// Is checks if the error is or contains the target error in its chain, therefore this function overrides the
// default implementation of errors.Is
func (w wrapped) Is(target error) bool {
	// check if the target is the same as the error
	if g, ok := target.(Grr); ok {
		if g.Error() == w.Error() {
			return true
		}
	}
	// check if the error is the same as the kind
	if w.kind.err.Error() == target.Error() {
		return true
	}
	// finally check if the error is in the chain
	chain := w.Chain()
	for _, err := range chain {
		if err.Error() == target.Error() {
			return true
		}
	}
	return false
}

// Chain builds the error chain as a slice of error for the package type Wrapped, ordered with the last element being
// the bottom most error in the chain
func (w wrapped) Chain() []error {
	return Chain(w, WithTrimCustom(w.kind.separator))
}

// Unwrap implements the Unwrap interface
func (w wrapped) Unwrap() error {
	return w.err
}
