package gerr

import (
	"errors"
	"strings"
)

// WithTrimColons connfigures the sanitizer to remove the colons from the beginning and end of the error string
func WithTrimColons() SanitizeOpt {
	return func(c *sanitizeConfig) {
		c.trimColonsFunc = trimColons
	}
}

// WithTrimSpaces configures the sanitizer to remove the spaces from the beginning and end of the error string
func WithTrimSpaces() SanitizeOpt {
	return func(c *sanitizeConfig) {
		c.trimSpacesFunc = trimSpaces
	}
}

// WithTrimCustom configures the sanitizer to remove the given string from the beginning and end of the error string
func WithTrimCustom(toTrim string) SanitizeOpt {
	return func(c *sanitizeConfig) {
		c.trimCustomFunc = func(s string) string {
			return trimCustom(s, toTrim)
		}
	}
}

// WithNormalOrder configures the sanitizer to using the normal order of the error chain instead of the reverse order
func WithNormalOrder() SanitizeOpt {
	return func(c *sanitizeConfig) {
		c.noInverseOrder = true
	}
}

// SanitizeOpt is the functional type for configuring the sanitization taking place in the Chain function
type SanitizeOpt func(c *sanitizeConfig)

// Chain builds and returns the error chain as a slice of errors, note that for this function to work as intended the
// the error must be wrapped according to the go 1.13 error wrapping guidelines
func Chain(err error, opts ...SanitizeOpt) []error {
	c := sanitizeConfig{}
	for _, opt := range opts {
		opt(&c)
	}
	if u, ok := err.(interface {
		Unwrap() error
	}); ok {
		var chain []error
		currentErr := err
		nextErr := u.Unwrap()
		hasNext := true
		for hasNext {
			currentErr, nextErr, hasNext = unwrapBuildChain(currentErr, nextErr, &chain, c)
		}
		if c.noInverseOrder == false {
			reverse[[]error, error](chain)
		}
		return chain
	}
	return []error{errors.New(sanitizeString(err.Error(), c))}
}

// unwrapBuildChain unwraps the error chain and builds it as a slice of errors returns the current error, the next error
// and a boolean indicating if there are more errors to continue
func unwrapBuildChain(currentErr error, nextErr error, chain *[]error, conf sanitizeConfig) (error, error, bool) {
	current, currOk := currentErr.(interface {
		Unwrap() error
	})
	next, nextOk := nextErr.(interface {
		Unwrap() error
	})
	if currOk && nextOk {
		sanitizeAppend(currentErr, nextErr, chain, conf)
		return current.Unwrap(), next.Unwrap(), true
	}
	if currOk && !nextOk {
		sanitizeAppend(currentErr, nextErr, chain, conf)
		*chain = append(*chain, errors.New(sanitizeString(nextErr.Error(), conf)))
	}
	return nil, nil, false
}

// sanitizeAppend removes the string of the next error from the current error and appends it if the result is not an
// empty string, this is done to support poorly implemented error wrapping
func sanitizeAppend(currentErr, nextErr error, chain *[]error, conf sanitizeConfig) {
	s := removeEqualPartFromError(currentErr, nextErr)
	if s != "" {
		s = sanitizeString(s, conf)
		*chain = append(*chain, errors.New(s))
	}
}

// sanitizeString applies the sanitization options to the given string
func sanitizeString(s string, conf sanitizeConfig) string {
	if conf.trimCustomFunc != nil {
		s = conf.trimCustomFunc(s)
	}
	if conf.trimSpacesFunc != nil {
		s = conf.trimSpacesFunc(s)
	}
	if conf.trimColonsFunc != nil {
		s = conf.trimColonsFunc(s)
	}
	return s
}

// removeEqualPartFromError removes the string of the next error from the current error
func removeEqualPartFromError(currentErr, nextErr error) string {
	s := strings.Replace(currentErr.Error(), nextErr.Error(), "", 1)
	return s
}

// reverse reverses the given slice of type E
func reverse[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// trimCustom trims the given string from the beginning and end of the error string
func trimCustom(s string, toTrim string) string {
	return strings.TrimPrefix(strings.TrimSuffix(s, toTrim), toTrim)
}

// trimSpaces removes the spaces from the beginning and end of the error string
func trimSpaces(s string) string {
	return strings.TrimSpace(s)
}

// trimColons removes the colons from the beginning and end of the error string
func trimColons(s string) string {
	return trimCustom(s, ":")
}

// sanitizeConfig is the configuration for the sanitization taking place in the Chain function
type sanitizeConfig struct {
	trimSpacesFunc func(s string) string
	trimColonsFunc func(s string) string
	trimCustomFunc func(s string) string
	noInverseOrder bool
}
