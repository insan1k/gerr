package gerr

import (
	"errors"
	"reflect"
	"testing"
)

func Test_kind_Add(t *testing.T) {
	type fields struct {
		err       error
		separator string
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
				err:       errors.New("test"),
				separator: _separator(),
			},
			args: args{
				err: errors.New("added"),
			},
			want: wrapped{
				kind: kind{
					err:       errors.New("test"),
					separator: _separator(),
				},
				err: errors.New("added"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := kind{
				err:       tt.fields.err,
				separator: tt.fields.separator,
			}
			if got := k.Add(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_kind_Chain(t *testing.T) {
	type fields struct {
		err       error
		separator string
	}
	tests := []struct {
		name   string
		fields fields
		want   []error
	}{
		{
			name: "Chain",
			fields: fields{
				err:       errors.New("test"),
				separator: _separator(),
			},
			want: []error{errors.New("test")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := kind{
				err:       tt.fields.err,
				separator: tt.fields.separator,
			}
			if got := k.Chain(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Chain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_kind_Error(t *testing.T) {
	type fields struct {
		err       error
		separator string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Error",
			fields: fields{
				err:       errors.New("test"),
				separator: _separator(),
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := kind{
				err:       tt.fields.err,
				separator: tt.fields.separator,
			}
			if got := k.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_kind_Is(t *testing.T) {
	type fields struct {
		err       error
		separator string
	}
	type args struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Is",
			fields: fields{
				err:       errors.New("test"),
				separator: _separator(),
			},
			args: args{
				err: errors.New("test"),
			},
			want: true,
		},
		{
			name: "IsNot",
			fields: fields{
				err:       errors.New("test"),
				separator: _separator(),
			},
			args: args{
				err: errors.New("not"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := kind{
				err:       tt.fields.err,
				separator: tt.fields.separator,
			}
			if got := k.Is(tt.args.err); got != tt.want {
				t.Errorf("Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_kind_Sanitize(t *testing.T) {
	type fields struct {
		err       error
		separator string
	}
	tests := []struct {
		name   string
		fields fields
		want   Grr
	}{
		{
			name: "Sanitize",
			fields: fields{
				err:       errors.New("test"),
				separator: _separator(),
			},
			want: kind{
				err:       errors.New("test"),
				separator: _separator(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := kind{
				err:       tt.fields.err,
				separator: tt.fields.separator,
			}
			if got := k.Sanitize(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sanitize() = %v, want %v", got, tt.want)
			}
		})
	}
}
