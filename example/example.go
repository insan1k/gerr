package example

import (
	"fmt"
	"github.com/insan1k/gerr"
)

//go:generate stringer -type=exampleErrors
type exampleErrors int

func (e exampleErrors) Error() string {
	return e.String()
}

const (
	ErrorExampleFailedToRead exampleErrors = iota
	ErrorExampleFailedToWrite
	ErrorExampleFailedToParse
)

func NewLocalError(err error, user, group, reason string) *LocalError {
	return &LocalError{
		Grr:   gerr.New(err),
		user:  user,
		group: group,
	}
}

type LocalError struct {
	gerr.Grr
	user, group, reason string
}

func (l LocalError) Error() string {
	msg := l.Grr.Error() + fmt.Sprint(" user: ", l.user, " group: ", l.group)
	if l.reason != "" {
		msg += fmt.Sprint(" reason: ", l.reason)
	}
	return msg
}
