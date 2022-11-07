package gerr_test

import (
	"errors"
	"fmt"
	"github.com/insan1k/gerr"
	"log"
	"testing"
)

type ExampleLocalError int

func (e ExampleLocalError) String() string {
	m := map[ExampleLocalError]string{
		ErrorExampleUnknown:       sErrorExampleUnknown,
		ErrorExampleFailedToWrite: sErrorExampleFailedToWrite,
		ErrorExampleFailedToRead:  sErrorExampleFailedToRead,
		ErrorExampleFailedToParse: sErrorExampleFailedToParse,
	}
	if v, ok := m[e]; ok {
		return v
	}
	return sErrorExampleUnknown
}

func (e ExampleLocalError) Error() string {
	return e.String()
}

const (
	sErrorExampleFailedToRead                    = "failed to read"
	sErrorExampleFailedToWrite                   = "failed to write"
	sErrorExampleFailedToParse                   = "failed to parse"
	sErrorExampleUnknown                         = "unknown error"
	ErrorExampleUnknown        ExampleLocalError = iota
	ErrorExampleFailedToRead
	ErrorExampleFailedToWrite
	ErrorExampleFailedToParse
)

func AsLocalError(err error) *LocalError {
	return gerr.AsGrr(err).(*LocalError)
}

func NewLocalError(err error, user, group string) *LocalError {
	return &LocalError{
		Grr:   gerr.New(err),
		user:  user,
		group: group,
	}
}

func (l *LocalError) AddReason(reason string) *LocalError {
	l.reason = reason
	return l
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

type AnotherLocalError gerr.Grr

func TestExampleErrors(t *testing.T) {
	if err := doReadGroup("someGroup", t); err != nil {
		log.Println("after sanitize")
		log.Println(err)
		log.Println("----------------------------------------")
	}
	if err := doWriteGroup("anotherGroup", t); err != nil {
		log.Println("after sanitize")
		log.Println(err)
		log.Println("----------------------------------------")
	}
	if err := doParseUser("someGroup", "someUser", t); err != nil {
		log.Println("after sanitize")
		log.Println(err)
		log.Println("----------------------------------------")
	}
	if err := doReadUser("someOtherGroup", "someOtherUser", t); err != nil {
		log.Println("after sanitize")
		log.Println(err)
		log.Println("----------------------------------------")
	}
	if err := doWriteUser("fooGroup", "fooUser", t); err != nil {
		log.Println("after sanitize")
		log.Println(err)
		log.Println("----------------------------------------")
	}
}

var doWriteGroup = func(group string, t *testing.T) error {
	err := writeGroup(group)
	writeGroupTestIs(err, t)
	log.Println("----------------------------------------")
	log.Println("before sanitize")
	log.Println(err)
	return gerr.AsGrr(err).Sanitize().Add(errors.New("internal error"))
}

var doReadGroup = func(group string, t *testing.T) error {
	err := readGroup(group)
	readGroupTestIs(err, t)
	log.Println("----------------------------------------")
	log.Println("before sanitize")
	log.Println(err)
	return gerr.AsGrr(err).Sanitize().Add(errors.New("internal error")).Add(fmt.Errorf("cannot read group %v", group))
}

var doWriteUser = func(user, group string, t *testing.T) error {
	err := writeUser(user, group)
	writeUserTestIs(err, t)
	lErr := AsLocalError(err)
	log.Println("----------------------------------------")
	log.Println("before sanitize")
	log.Println(err)
	return lErr.Sanitize().Add(fmt.Errorf("internal error"))
}

var doReadUser = func(user, group string, t *testing.T) error {
	err := readUser(user, group)
	readUserTestIs(err, t)
	lErr := AsLocalError(err)
	log.Println("----------------------------------------")
	log.Println("before sanitize")
	log.Println(err)
	return lErr.Sanitize().Add(errors.New("internal error")).Add(fmt.Errorf("cannot read %v/%v", group, user))
}

var doParseUser = func(user, group string, t *testing.T) error {
	err := parseUser(user, group)
	parseUserTestIs(err, t)
	lErr := AsLocalError(err)
	log.Println("----------------------------------------")
	log.Println("before sanitize")
	log.Println(err)
	return lErr.Sanitize().Add(errors.New(lErr.reason))
}

var fakeDbError = errors.New("fake DB error, cannot connect, sql://someUser@ConnectionString:10005")
var parseError = errors.New("failed to parse, invalid id format")
var fakeDb = func(_, _ string) error {
	return fakeDbError
}

var writeUser = func(user, group string) error {
	if err := fakeDb(user, group); err != nil {
		return NewLocalError(gerr.New(ErrorExampleFailedToWrite).Add(err), user, group)
	}
	return nil
}

var readUser = func(user, group string) error {
	if err := fakeDb(user, group); err != nil {
		return NewLocalError(gerr.New(ErrorExampleFailedToRead).Add(err), user, group)
	}
	return nil

}

var parseUser = func(user, group string) error {
	if err := fakeDb(user, group); err != nil {
		return NewLocalError(
			gerr.New(ErrorExampleFailedToParse).Add(parseError).Add(fakeDbError),
			user,
			group,
		).AddReason("invalid user id")
	}
	return nil
}

var writeGroup = func(group string) error {
	if err := fakeDb("", group); err != nil {
		return AnotherLocalError(gerr.New(ErrorExampleFailedToWrite).Add(err))
	}
	return nil
}

var readGroup = func(group string) error {
	if err := fakeDb("", group); err != nil {
		return AnotherLocalError(gerr.New(ErrorExampleFailedToRead).Add(err))
	}
	return nil

}

func writeGroupTestIs(err error, t *testing.T) {
	if !gerr.AsGrr(err).Is(ErrorExampleFailedToWrite) {
		t.Fatalf("expected %v, got %v", ErrorExampleFailedToWrite, err)
	}
	if !gerr.AsGrr(err).Is(fakeDbError) {
		t.Fatalf("expected %v, got %v", fakeDbError, err)
	}
}

func readGroupTestIs(err error, t *testing.T) {
	if !gerr.AsGrr(err).Is(ErrorExampleFailedToRead) {
		t.Fatalf("expected %v, got %v", ErrorExampleFailedToRead, err)
	}
	if !gerr.AsGrr(err).Is(fakeDbError) {
		t.Fatalf("expected %v, got %v", fakeDbError, err)
	}
}

func parseUserTestIs(err error, t *testing.T) {
	if !gerr.AsGrr(err).Is(ErrorExampleFailedToParse) {
		t.Fatalf("expected %v, got %v", ErrorExampleFailedToParse, err)
	}
	if !gerr.AsGrr(err).Is(parseError) {
		t.Fatalf("expected %v, got %v", parseError, err)
	}
	if !gerr.AsGrr(err).Is(fakeDbError) {
		t.Fatalf("expected %v, got %v", fakeDbError, err)
	}
}

func writeUserTestIs(err error, t *testing.T) {
	if !gerr.AsGrr(err).Is(ErrorExampleFailedToWrite) {
		t.Fatalf("expected %v, got %v", ErrorExampleFailedToWrite, err)
	}
	if !gerr.AsGrr(err).Is(fakeDbError) {
		t.Fatalf("expected %v, got %v", fakeDbError, err)
	}
}

func readUserTestIs(err error, t *testing.T) {
	if !gerr.AsGrr(err).Is(ErrorExampleFailedToRead) {
		t.Fatalf("expected %v, got %v", ErrorExampleFailedToRead, err)
	}
	if gerr.AsGrr(err).Is(ErrorExampleFailedToParse) {
		t.Fatalf("expected %v, got %v", ErrorExampleFailedToRead, err)
	}
	if gerr.AsGrr(err).Is(ErrorExampleFailedToWrite) {
		t.Fatalf("expected %v, got %v", ErrorExampleFailedToRead, err)
	}
	if gerr.AsGrr(err).Is(parseError) {
		t.Fatalf("expected %v, got %v", ErrorExampleFailedToRead, err)
	}
	if !gerr.AsGrr(err).Is(fakeDbError) {
		t.Fatalf("expected %v, got %v", fakeDbError, err)
	}
}
