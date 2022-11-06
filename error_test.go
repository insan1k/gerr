package gerr

import (
	"fmt"
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
				separator: separator(),
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
					separator: separator(),
				},
				err: fmt.Errorf("two"),
			},
		},
		{
			name: "WithWrappedErrors",
			args: args{
				kind: genWrappedErrWithSpace(10),
			},
			want: wrapped{
				kind: kind{
					err:       fmt.Errorf("top1"),
					separator: separator(),
				},
				err: genWrappedErrWithSpaceNoTop(10),
			},
		},
		{
			name: "WithWrappedErrorsAndWithErr",
			args: args{
				kind: genWrappedErrWithSpace(10),
				opts: []Option{WithErr(fmt.Errorf("thisThing"))},
			},
			want: wrapped{
				kind: kind{
					err:       fmt.Errorf("top1"),
					separator: separator(),
				},
				err: fmt.Errorf("%v%s%w", fmt.Errorf("thisThing"), separator(), genWrappedErrWithSpaceNoTop(10)),
			},
		},
		{
			name: "WithWrappedErrorsSimple",
			args: args{
				kind: fmt.Errorf("%v%s%w", fmt.Errorf("two"), separator(), fmt.Errorf("one")),
				opts: []Option{WithErr(fmt.Errorf("another"))},
			},
			want: wrapped{
				kind: kind{
					err:       fmt.Errorf("two"),
					separator: separator(),
				},
				err: fmt.Errorf("%v%s%w", fmt.Errorf("another"), separator(), fmt.Errorf("one")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.kind, tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func genWrappedErrWithSpaceNoTop(count int) error {
	return genWrappedErrNoTop(count, " ", "")
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
