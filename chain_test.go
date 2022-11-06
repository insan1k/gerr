package gerr

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestChain(t *testing.T) {
	type args struct {
		curErr error
		opts   []SanitizeOpt
	}
	tests := []struct {
		name string
		args args
		want []error
	}{
		{
			name: "testWithColon10",
			args: args{
				curErr: genWrappedErrWithColon(10),
				opts:   []SanitizeOpt{WithTrimColons()},
			},
			want: []error{
				errors.New("bot10"),
				errors.New("mid9"),
				errors.New("mid8"),
				errors.New("mid7"),
				errors.New("mid6"),
				errors.New("mid5"),
				errors.New("mid4"),
				errors.New("mid3"),
				errors.New("mid2"),
				errors.New("top1"),
			},
		},
		{
			name: "testWithColon7",
			args: args{
				curErr: genWrappedErrWithColon(7),
				opts:   []SanitizeOpt{WithTrimColons()},
			},
			want: []error{
				errors.New("bot7"),
				errors.New("mid6"),
				errors.New("mid5"),
				errors.New("mid4"),
				errors.New("mid3"),
				errors.New("mid2"),
				errors.New("top1"),
			},
		},
		{
			name: "testWithColon3",
			args: args{
				curErr: genWrappedErrWithColon(3),
				opts:   []SanitizeOpt{WithTrimColons()},
			},
			want: []error{
				errors.New("bot3"),
				errors.New("mid2"),
				errors.New("top1"),
			},
		},
		{
			name: "testWithSpace10",
			args: args{
				curErr: genWrappedErrWithSpace(10),
				opts:   []SanitizeOpt{WithTrimSpaces()},
			},
			want: []error{
				errors.New("bot10"),
				errors.New("mid9"),
				errors.New("mid8"),
				errors.New("mid7"),
				errors.New("mid6"),
				errors.New("mid5"),
				errors.New("mid4"),
				errors.New("mid3"),
				errors.New("mid2"),
				errors.New("top1"),
			},
		},
		{
			name: "testWithSpace7",
			args: args{
				curErr: genWrappedErrWithSpace(7),
				opts:   []SanitizeOpt{WithTrimSpaces()},
			},
			want: []error{
				errors.New("bot7"),
				errors.New("mid6"),
				errors.New("mid5"),
				errors.New("mid4"),
				errors.New("mid3"),
				errors.New("mid2"),
				errors.New("top1"),
			},
		},
		{
			name: "testWithSpace3",
			args: args{
				curErr: genWrappedErrWithSpace(3),
				opts:   []SanitizeOpt{WithTrimSpaces()},
			},
			want: []error{
				errors.New("bot3"),
				errors.New("mid2"),
				errors.New("top1"),
			},
		},
		{
			name: "testWithSpaceAndColon10",
			args: args{
				curErr: genWrappedErrWithSpaceAndColon(10),
				opts:   []SanitizeOpt{WithTrimSpaces(), WithTrimColons()},
			},
			want: []error{
				errors.New("bot10"),
				errors.New("mid9"),
				errors.New("mid8"),
				errors.New("mid7"),
				errors.New("mid6"),
				errors.New("mid5"),
				errors.New("mid4"),
				errors.New("mid3"),
				errors.New("mid2"),
				errors.New("top1"),
			},
		},
		{
			name: "testWithSpaceAndColon7",
			args: args{
				curErr: genWrappedErrWithSpaceAndColon(7),
				opts:   []SanitizeOpt{WithTrimSpaces(), WithTrimColons()},
			},
			want: []error{
				errors.New("bot7"),
				errors.New("mid6"),
				errors.New("mid5"),
				errors.New("mid4"),
				errors.New("mid3"),
				errors.New("mid2"),
				errors.New("top1"),
			},
		},
		{
			name: "testWithSpaceAndColon3",
			args: args{
				curErr: genWrappedErrWithSpaceAndColon(3),
				opts:   []SanitizeOpt{WithTrimSpaces(), WithTrimColons()},
			},
			want: []error{
				errors.New("bot3"),
				errors.New("mid2"),
				errors.New("top1"),
			},
		},
		{
			name: "testWithSpaceAndColon2",
			args: args{
				curErr: genWrappedErrWithSpaceAndColon(2),
				opts:   []SanitizeOpt{WithTrimSpaces(), WithTrimColons()},
			},
			want: []error{
				errors.New("bot2"),
				errors.New("top1"),
			},
		},
		{
			name: "testWith1",
			args: args{
				curErr: errors.New("one"),
			},
			want: []error{
				errors.New("one"),
			},
		},
		{
			name: "testWithPrependedSpaceAndColon10",
			args: args{
				curErr: genWrappedErrWithSpaceAndColonAndPrepend(10),
				opts:   []SanitizeOpt{WithTrimSpaces(), WithTrimColons()},
			},
			want: []error{
				errors.New("bot10"),
				errors.New("mid9"),
				errors.New("mid8"),
				errors.New("mid7"),
				errors.New("mid6"),
				errors.New("mid5"),
				errors.New("mid4"),
				errors.New("mid3"),
				errors.New("mid2"),
				errors.New("top1"),
			},
		},
		{
			name: "testWithNormalOrder10",
			args: args{
				curErr: genWrappedErrWithSpaceAndColon(10),
				opts:   []SanitizeOpt{WithTrimSpaces(), WithTrimColons(), WithNormalOrder()},
			},
			want: []error{
				errors.New("top1"),
				errors.New("mid2"),
				errors.New("mid3"),
				errors.New("mid4"),
				errors.New("mid5"),
				errors.New("mid6"),
				errors.New("mid7"),
				errors.New("mid8"),
				errors.New("mid9"),
				errors.New("bot10"),
			},
		},
		{
			name: "testWithCustom10",
			args: args{
				curErr: genWrappedErrWith(10, "%custom%", "%custom%"),
				opts:   []SanitizeOpt{WithTrimCustom("%custom%")},
			},
			want: []error{
				errors.New("bot10"),
				errors.New("mid9"),
				errors.New("mid8"),
				errors.New("mid7"),
				errors.New("mid6"),
				errors.New("mid5"),
				errors.New("mid4"),
				errors.New("mid3"),
				errors.New("mid2"),
				errors.New("top1"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Chain(tt.args.curErr, tt.args.opts...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Chain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func genWrappedErrWith(count int, sep, pre string) error {
	if count < 2 {
		count = 2
	}
	var err error
	err = fmt.Errorf("%sbot%v", pre, count)
	for i := count - 1; i > 1; i-- {
		errNew := fmt.Errorf("mid%v", i)
		err = fmt.Errorf("%s%v%s%w", pre, errNew, sep, err)
	}
	err = fmt.Errorf("%stop1%s%w", pre, sep, err)
	return err
}

func genWrappedErrWithColon(count int) error {
	return genWrappedErrWith(count, ":", "")
}

func genWrappedErrWithSpace(count int) error {
	return genWrappedErrWith(count, " ", "")
}
func genWrappedErrWithSpaceAndColon(count int) error {
	return genWrappedErrWith(count, ": ", "")
}

func genWrappedErrWithSpaceAndColonAndPrepend(count int) error {
	return genWrappedErrWith(count, ": ", " :")
}
