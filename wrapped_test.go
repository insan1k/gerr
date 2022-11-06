package gerr

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func Test_wrapped_Add(t *testing.T) {
	type fields struct {
		kind kind
		err  error
	}
	type args struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Grr
	}{
		{
			name: "Add",
			fields: fields{
				kind: kind{
					err:       errors.New("test"),
					separator: separator(),
				},
				err: errors.New("wasThere"),
			},
			args: args{
				err: errors.New("add"),
			},
			want: wrapped{
				kind: kind{
					err:       errors.New("test"),
					separator: separator(),
				},
				err: fmt.Errorf("%v%s%w", errors.New("add"), separator(), errors.New("wasThere")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := wrapped{
				kind: tt.fields.kind,
				err:  tt.fields.err,
			}
			if got := w.Add(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_wrapped_Chain(t *testing.T) {
	type fields struct {
		kind kind
		err  error
	}
	tests := []struct {
		name   string
		fields fields
		want   []error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := wrapped{
				kind: tt.fields.kind,
				err:  tt.fields.err,
			}
			if got := w.Chain(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Chain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_wrapped_Error(t *testing.T) {
	type fields struct {
		kind kind
		err  error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := wrapped{
				kind: tt.fields.kind,
				err:  tt.fields.err,
			}
			if got := w.Error(); got != tt.want {
				t.Errorf("Grr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_wrapped_Is(t *testing.T) {
	type fields struct {
		kind kind
		err  error
	}
	type args struct {
		target error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := wrapped{
				kind: tt.fields.kind,
				err:  tt.fields.err,
			}
			if got := w.Is(tt.args.target); got != tt.want {
				t.Errorf("Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_wrapped_Sanitize(t *testing.T) {
	type fields struct {
		kind kind
		err  error
	}
	tests := []struct {
		name   string
		fields fields
		want   Grr
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := wrapped{
				kind: tt.fields.kind,
				err:  tt.fields.err,
			}
			if got := w.Sanitize(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sanitize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_wrapped_Unwrap(t *testing.T) {
	type fields struct {
		kind kind
		err  error
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := wrapped{
				kind: tt.fields.kind,
				err:  tt.fields.err,
			}
			if err := w.Unwrap(); (err != nil) != tt.wantErr {
				t.Errorf("Unwrap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
