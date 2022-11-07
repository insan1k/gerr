package gerr

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		kind error
		opts []Option
	}
	tests := []struct {
		name string
		args args
		want Grr
	}{
		{
			name: "Error1",
			args: args{
				kind: fmt.Errorf("one"),
			},
			want: kind{
				err:       fmt.Errorf("one"),
				separator: _separator(),
			},
		},
		{
			name: "Error2",
			args: args{
				kind: fmt.Errorf("one"),
				opts: []Option{WithErr(fmt.Errorf("two"))},
			},
			want: wrapped{
				kind: kind{
					err:       fmt.Errorf("one"),
					separator: _separator(),
				},
				err: fmt.Errorf("two"),
			},
		},
		{
			name: "WithWrappedErrors",
			args: args{
				kind: genWrappedErrorsWithSeparator(10),
			},
			want: wrapped{
				kind: kind{
					err:       fmt.Errorf("top1"),
					separator: _separator(),
				},
				err: genWrappedWithSeparatorErrNoTop(10),
			},
		},
		{
			name: "WithWrappedErrorsAndWithErr",
			args: args{
				kind: genWrappedErrorsWithSeparator(10),
				opts: []Option{WithErr(fmt.Errorf("thisThing"))},
			},
			want: wrapped{
				kind: kind{
					err:       fmt.Errorf("top1"),
					separator: _separator(),
				},
				err: fmt.Errorf("%v%s%w", fmt.Errorf("thisThing"), _separator(), genWrappedWithSeparatorErrNoTop(10)),
			},
		},
		{
			name: "WithWrappedErrorsSimple",
			args: args{
				kind: fmt.Errorf("%v%s%w", fmt.Errorf("two"), _separator(), fmt.Errorf("one")),
				opts: []Option{WithErr(fmt.Errorf("another"))},
			},
			want: wrapped{
				kind: kind{
					err:       fmt.Errorf("two"),
					separator: _separator(),
				},
				err: fmt.Errorf("%v%s%w", fmt.Errorf("another"), _separator(), fmt.Errorf("one")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.kind, tt.args.opts...)
			originalGot := got
			originalWant := tt.want
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
				t.FailNow()
			}
			if !reflect.DeepEqual(got.Chain(), tt.want.Chain()) {
				t.Errorf("New().Chain() = %v, want %v", got.Chain(), tt.want.Chain())
			}
			if !reflect.DeepEqual(got.Error(), tt.want.Error()) {
				t.Errorf("New().Error() = %v, want %v", got.Error(), tt.want.Error())
			}
			addedErr := errors.New("addedErrNow")
			if !reflect.DeepEqual(got.Add(addedErr), tt.want.Add(addedErr)) {
				t.Errorf("New().Add() = %v, want %v", got.Add(addedErr), tt.want.Add(addedErr))
			}
			notAddedErr := errors.New("notAdded")
			if !reflect.DeepEqual(got.Is(notAddedErr), tt.want.Is(notAddedErr)) {
				t.Errorf("New().Is() = %v, want %v", got.Is(notAddedErr), tt.want.Is(notAddedErr))
			}
			if !reflect.DeepEqual(got.Sanitize(), tt.want.Sanitize()) {
				t.Errorf("New().Sanitize() = %v, want %v", got.Sanitize(), tt.want.Sanitize())
			}
			if !reflect.DeepEqual(originalGot, got) {
				t.Errorf("Mutability has occurred got.Grr = %v, want %v", got, originalGot)
			}
			if !reflect.DeepEqual(originalWant, tt.want) {
				t.Errorf("Mutability has occurred tt.want.Grr = %v, want %v", tt.want, originalWant)
			}
		})
	}
}

func genWrappedWithSeparatorErrNoTop(count int) error {
	return genWrappedErrNoTop(count, _separator(), "")
}

func genWrappedErrNoTop(count int, sep, pre string) error {
	if count < 2 {
		count = 2
	}
	var err error
	err = fmt.Errorf("%sbot%v", pre, count)
	for i := count - 1; i > 1; i-- {
		errNew := fmt.Errorf("mid%v", i)
		err = fmt.Errorf("%s%v%s%w", pre, errNew, sep, err)
	}
	return err
}

type someStuctWithGerr struct {
	Grr
	foo, bar string
}

func (s someStuctWithGerr) Error() string {
	return s.Grr.Error() + " " + s.foo + " " + s.bar
}

type someStructWithWrapped struct {
	wrapped
	foo, bar string
}

func (s someStructWithWrapped) Error() string {
	return s.wrapped.Error() + " " + s.foo + " " + s.bar
}

type someStructWithKind struct {
	kind
	foo, bar string
}

func (s someStructWithKind) Error() string {
	return s.kind.Error() + " " + s.foo + " " + s.bar
}

func TestAsGrr(t *testing.T) {
	type args struct {
		err error
	}
	type want struct {
		grr        Grr
		structType interface{}
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "someStructWithGerrKind",
			args: args{
				err: someStuctWithGerr{
					Grr: New(fmt.Errorf("one")),
					foo: "foo",
					bar: "bar",
				},
			},
			want: want{
				grr: kind{
					err:       fmt.Errorf("one"),
					separator: _separator(),
				},
				structType: someStuctWithGerr{},
			},
		},
		{
			name: "someStructWithGerrWrapped",
			args: args{
				err: someStuctWithGerr{
					Grr: wrapped{
						kind: kind{
							err:       fmt.Errorf("one"),
							separator: _separator(),
						},
						err: fmt.Errorf("two"),
					},
					foo: "foo",
					bar: "bar",
				},
			},
			want: want{
				grr: wrapped{
					kind: kind{
						err:       fmt.Errorf("one"),
						separator: _separator(),
					},
					err: fmt.Errorf("two"),
				},
				structType: someStuctWithGerr{},
			},
		},
		{
			name: "someStructWithKind",
			args: args{
				err: someStructWithKind{
					kind: kind{
						err:       fmt.Errorf("one"),
						separator: _separator(),
					},
					foo: "foo",
					bar: "bar",
				},
			},
			want: want{
				grr: kind{
					err:       fmt.Errorf("one"),
					separator: _separator(),
				},
				structType: someStructWithKind{},
			},
		},
		{
			name: "someStructWithWrapped",
			args: args{
				err: someStructWithWrapped{
					wrapped: wrapped{
						kind: kind{
							err:       fmt.Errorf("one"),
							separator: _separator(),
						},
						err: fmt.Errorf("two"),
					},
					foo: "foo",
					bar: "bar",
				},
			},
			want: want{
				grr: wrapped{
					kind: kind{
						err:       fmt.Errorf("one"),
						separator: _separator(),
					},
					err: fmt.Errorf("two"),
				},
				structType: someStructWithWrapped{},
			},
		},
		{
			name: "someOtherError",
			args: args{
				err: fmt.Errorf("one"),
			},
			want: want{
				grr: kind{
					err:       fmt.Errorf("one"),
					separator: _separator(),
				},
				structType: kind{},
			},
		},
		{
			name: "nil",
			args: args{
				err: nil,
			},
			want: want{
				grr:        nil,
				structType: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AsGrr(tt.args.err)
			if got != nil {
				if !assert.Implements(t, (*Grr)(nil), got) {
					t.Errorf("AsGrr() returned type that does not implement Grr")
					t.FailNow()
				}
				ok := errors.As(got, &tt.want.structType)
				if !ok && got != nil {
					t.Errorf("AsGrr() = %v, want %v", got, tt.want)
					t.FailNow()
				}
			}
			assert.EqualValues(t, tt.want.structType, got, "AsGrr() = %v, want %v", got, tt.want)
		})
	}
}
