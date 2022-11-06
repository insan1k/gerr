package example

import (
	"errors"
	"fmt"
	"github.com/insan1k/gerr"
	"testing"
)

type AnotherLocalError gerr.Grr

func Test_Example(t *testing.T) {
	fakeDbError := errors.New("fake DB error, cannot connect")
	parseError := errors.New("failed to parse, invalid id format")
	fakeDb := func(_, _ string) error {
		return fakeDbError
	}

	writeUser := func(user, group string) *LocalError {
		if err := fakeDb(user, group); err != nil {
			return NewLocalError(gerr.New(ErrorExampleFailedToWrite).Add(err), user, group, "internal")
		}
		return nil
	}
	readUser := func(user, group string) *LocalError {
		if err := fakeDb(user, group); err != nil {
			return NewLocalError(gerr.New(ErrorExampleFailedToRead).Add(err), user, group, "internal")
		}
		return nil

	}
	parseUser := func(user, group string) *LocalError {
		if err := fakeDb(user, group); err != nil {
			return NewLocalError(
				gerr.New(ErrorExampleFailedToParse).Add(parseError).Add(fakeDbError),
				user,
				group,
				"invalid user id",
			)
		}
		return nil
	}
	if err := writeUser("user", "group"); err != nil {
		fmt.Printf("Is(%v)%v\n", ErrorExampleFailedToWrite, err.Is(ErrorExampleFailedToWrite))
		fmt.Printf("Is(%v)%v\n", ErrorExampleFailedToRead, err.Is(ErrorExampleFailedToRead))
		fmt.Printf("Is(%v)%v\n", ErrorExampleFailedToParse, err.Is(ErrorExampleFailedToParse))
		fmt.Printf("Is(%v)%v\n", parseError, err.Is(parseError))
		fmt.Printf("Is(%v)%v\n", fakeDbError, err.Is(fakeDbError))
	}
	if err := readUser("user", "group"); err != nil {
		fmt.Printf("Is(%v)%v\n", ErrorExampleFailedToWrite, err.Is(ErrorExampleFailedToWrite))
		fmt.Printf("Is(%v)%v\n", ErrorExampleFailedToRead, err.Is(ErrorExampleFailedToRead))
		fmt.Printf("Is(%v)%v\n", ErrorExampleFailedToParse, err.Is(ErrorExampleFailedToParse))
		fmt.Printf("Is(%v)%v\n", parseError, err.Is(parseError))
		fmt.Printf("Is(%v)%v\n", fakeDbError, err.Is(fakeDbError))
	}
	if err := parseUser("user", "group"); err != nil {
		fmt.Printf("Is(%v)%v\n", ErrorExampleFailedToParse, err.Is(ErrorExampleFailedToParse))
		fmt.Printf("Is(%v)%v\n", parseError, err.Is(parseError))
		fmt.Printf("Is(%v)%v\n", fakeDbError, err.Is(fakeDbError))
	}

}
