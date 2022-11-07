[![Go Reference](https://pkg.go.dev/badge/github.com/insan1k/gerr.svg)](https://pkg.go.dev/github.com/insan1k/gerr)
[![Go Report Card](https://goreportcard.com/badge/github.com/insan1k/gerr)](https://goreportcard.com/report/github.com/insan1k/gerr)
[![Coverage Status](https://coveralls.io/repos/github/insan1k/gerr/badge.svg?branch=main)](https://coveralls.io/github/insan1k/gerr?branch=main)
# gerr
Gerr is a simple, yet opinionated library to generate errors in go. 

The library is designed to be used in a microservice architecture, where
taking actions based on an error without having to know the exact error type is important, as it promotes decoupling
between layers. 

As a developer I would like to take an action based on the error, without having to know
the underlying implementation of the package and the error type, that created an error in the first place.
That being said I would still like to have the ability to check for the exact error type, if I need to. And to have the 
error type preserved, for logging and debugging purposes. 

Thus comes the motivation for writing this library, the name is a play on words, as it is a library to generate errors, 
so I thought, of the most common sound I make when doing anything fancier than the vanilla error handling in go. Grr!

## Installation
```bash
go get github.com/insan1k/gerr
```

## Usage
usage of the Gerr library is simple, you can use the `gerr.New` function to create a new error, or use the `gerr.Add` 
function add another error. The library supports wrapping errors, and adding enriching the error, with additional 
contextual information, not to be confused with `context.Context`.

### Important
Gerr does not add any overhead to the error, it is just a wrapper around the error, and does not add any significant 
size to the memory footprint of the error.  

#### How it works 

There are two types that implement the `Grr` interface, these types are `kind` and `wrapped`
- kind: is the vanilla error type. it has no additional information, other than the separator used for the error message.
- wrapped: is the error type that wraps another error, it has the same information as the `kind` type, but also has the 
`Unwrap` method, thus it fulfills the `errors.Unwrap` interface.

The `New` function returns a `Grr` interface type, behind the scenes it is either of type `wrapped` or type `kind`
depending on the arguments passed to the function.

The `Add` function is a convenience function to add an error to an existing error, if the underlying type is `kind`
then it will be converted to a `wrapped` type, and the new error will be added to the chain, if it's already a type
`wrapped` then the new error will be added to the chain.

The `Error` function exists to satisfy the `error` interface, and it returns the error message, which is a concatenation
of the error messages of the underlying errors, separated by the separator.

The `Chain` function returns the complete and sanitized error chain as `[]error` this is useful if more control is
needed, which is not exposed by the `Is` function.

The `Is` function not only is a convenience function to check if the error is of a specific type, but it also
if fulfills the `srrors.Is` interface therefore any calls to `errors.Is` will also work.

The `Sanitize` function simply removes all of the `wrapped` specific data and returns a `kind` type, one could add
more errors to the kind and thus it would become a `wrapped` type again.

The `Unwrap` function implements the `errors.Unwrap` interface, and returns the underlying error, if it exists. 
This method is only available on the `wrapped` type and it's not exposed on the `Grr` interface. 

#### Using `gerr.New`
```go
package main

import (
	"fmt"
	"errors"
	"github.com/insan1k/gerr"
)

func main() {
	top:=errors.New("this is an error")
	wrap:=errors.New("with a wrapped error")
	gerr.New(top,gerr.WithErr(wrap))
	gerr.New(top).Add(wrap) // same as above
	gerr.New(fmt.Errorf("%v %w",top,wrap)) // same as above
    fmt.Println(err)
}
```

#### Example `Grr` interface usage
```go
package main

import (
    "fmt"
    "errors"
    "github.com/insan1k/gerr"
)

func main() {
    err := gerr.New(errors.New("some error")).Add(errors.New("some other error"))
	
	// retrieves the chain of errors as []error
	err.Chain()
	
	// implements the error interface making Gerr compatible with the standard library
	err.Error()
	
	// overrides the default errors.Is implementation and checks if the error is present in the chain 
	err.Is(errors.New("some error")) // true
	err.Is(errors.New("some other error")) // true
	err.Is(errors.New("some error 2")) //false 
	
	// no mutability takes place in the library, all functions return a new instance of the error
	err = err.Add(errors.New("and another one"))
	err.Is(errors.New("and another one")) //  true

	// removes all errors from the chain and any additional contextual information
	err = err.Sanitize()
	err.Is(errors.New("some error")) // true
    

}
```
More examples can be found in the [examples](https://github.com/insan1k/gerr/blob/main/example_test.go) file.
There you can see an example of the Grr interface in action. How we can embed it to a struct and use it to enrich the 
error.

### Chain
The chain is constructed by using the `fmt.Errorf("some message: %w", err)` function, which wraps the error in a way 
that it can be retrieved by using the `errors.Unwrap` function.
however `errors.Unwrap` only keeps the bottom most error in the chain fully preserved. 

This is where a special function comes in handy, `gerr.Chain()` which returns the chain of errors as a slice of errors.
It does this by removing the next wrapped error from the current error, the subtraction of these strings, when 
effectively sanitized from any separators is the representation of the current error. Walk the chain and voila you have
the full chain of errors.
#### Example `gerr.Chain()` usage
```go
package main

import (
    "fmt"
    "errors"
    "github.com/insan1k/gerr"
)

func main() {
    err := gerr.New(errors.New("some error")).Add(errors.New("some other error"))
    
    // retrieves the chain of errors as []error
    err.Chain() // []error{errors.New("some error"), errors.New("some other error")}
	// it also works with your own wrapped errors
	err2 := fmt.Errorf("some error: %w", errors.New("some other error")) // some error: some other error
	gerr.Chain(err2, gerr.WithTrimColons(), gerr.WithTrimSpaces()) // []error{"some error", "some other error"}
	// there are several options to trim the error message, to make sure we accurately retrieve each individual error
	
	// in the gerr package you can specify the separator by calling the following function
	gerr.SetSpaceSeparator() // " " default
	gerr.SetCustomSeparator(" -- ") // custom separator
	gerr.SetColonSeparator() // ":" separator
	gerr.SetColonAndSpaceSeparator() // ": " separator
	
	//note that these functions only set the separator for the Grr interface, and do not affect the underlying 
	//gerr.Chain function
}
```
## Final Considerations
As it stands this library is a work in progress, I would like to keep a minimal API and as such I would not add a lot of
features on the `Grr` interface , if you have any suggestions please feel free to open an issue or a PR.

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License
[Apache 2.0](https://github.com/insan1k/gerr/blob/main/LICENSE)