// Code generated by "stringer -type=exampleErrors"; DO NOT EDIT.

package example

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ErrorExampleFailedToRead-0]
	_ = x[ErrorExampleFailedToWrite-1]
	_ = x[ErrorExampleFailedToParse-2]
}

const _exampleErrors_name = "ErrorExampleFailedToReadErrorExampleFailedToWriteErrorExampleFailedToParse"

var _exampleErrors_index = [...]uint8{0, 24, 49, 74}

func (i exampleErrors) String() string {
	if i < 0 || i >= exampleErrors(len(_exampleErrors_index)-1) {
		return "exampleErrors(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _exampleErrors_name[_exampleErrors_index[i]:_exampleErrors_index[i+1]]
}
