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
					separator: _separator(),
				},
				err: errors.New("wasThere"),
			},
			args: args{
				err: errors.New("add"),
			},
			want: wrapped{
				kind: kind{
					err:       errors.New("test"),
					separator: _separator(),
				},
				err: fmt.Errorf("%v%s%w", errors.New("add"), _separator(), errors.New("wasThere")),
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
		{
			name: "ChainTwoErrors",
			fields: fields{
				kind: kind{
					err:       errors.New("test"),
					separator: _separator(),
				},
				err: fmt.Errorf("%v%s%w", errors.New("added"), _separator(), errors.New("original")),
			},
			want: []error{errors.New("test"), errors.New("added"), errors.New("original")},
		},
		{
			name: "ChainOneError",
			fields: fields{
				kind: kind{
					err:       errors.New("test"),
					separator: _separator(),
				},
				err: errors.New("original"),
			},
			want: []error{errors.New("test"), errors.New("original")},
		},
		{
			name: "ChainSeveralErrors",
			fields: fields{
				kind: kind{
					err:       errors.New("test"),
					separator: _separator(),
				},
				err: fmt.Errorf("%v%s%w", errors.New("added"), _separator(), fmt.Errorf("%v%s%w", errors.New("added2"), _separator(), errors.New("original"))),
			},
			want: []error{errors.New("test"), errors.New("added"), errors.New("added2"), errors.New("original")},
		},
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
		{
			name: "ChainTwoErrors",
			fields: fields{
				kind: kind{
					err:       errors.New("test"),
					separator: _separator(),
				},
				err: fmt.Errorf("%v%s%w", errors.New("added"), _separator(), errors.New("original")),
			},
			want: "test" + _separator() + "added" + _separator() + "original",
		},
		{
			name: "ChainOneError",
			fields: fields{
				kind: kind{
					err:       errors.New("test"),
					separator: _separator(),
				},
				err: errors.New("original"),
			},
			want: "test" + _separator() + "original",
		},
		{
			name: "ChainSeveralErrors",
			fields: fields{
				kind: kind{
					err:       errors.New("test"),
					separator: _separator(),
				},
				err: fmt.Errorf("%v%s%w", errors.New("added"), _separator(), fmt.Errorf("%v%s%w", errors.New("added2"), _separator(), errors.New("original"))),
			},
			want: "test" + _separator() + "added" + _separator() + "added2" + _separator() + "original",
		},
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
	nn := New(genWrappedErrorsWithSeparator(100)).(wrapped)
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "IsTarget",
			fields: fields{
				kind: kind{
					err:       errors.New("test"),
					separator: _separator(),
				},
				err: errors.New("original"),
			},
			args: args{
				target: errors.New("test"),
			},
			want: true,
		},
		{
			name: "IsTarget2",
			fields: fields{
				kind: kind{
					err:       errors.New("test"),
					separator: _separator(),
				},
				err: errors.New("original"),
			},
			args: args{
				target: errors.New("original"),
			},
			want: true,
		},
		{
			name: "IsNotTarget",
			fields: fields{
				kind: kind{
					err:       errors.New("test"),
					separator: _separator(),
				},
				err: errors.New("original"),
			},
			args: args{
				target: errors.New("not"),
			},
			want: false,
		},
		{
			name: "IsTargetBot100",
			fields: fields{
				kind: kind{
					err:       nn.kind.err,
					separator: nn.kind.separator,
				},
				err: nn.err,
			},
			args: args{
				target: errors.New("bot100"),
			},
			want: true,
		},
		{
			name: "IsTargetMid50",
			fields: fields{
				kind: kind{
					err:       nn.kind.err,
					separator: nn.kind.separator,
				},
				err: nn.err,
			},
			args: args{
				target: errors.New("mid50"),
			},
			want: true,
		},
		{
			name: "IsTargetTop1",
			fields: fields{
				kind: kind{
					err:       nn.kind.err,
					separator: nn.kind.separator,
				},
				err: nn.err,
			},
			args: args{
				target: errors.New("top1"),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := wrapped{
				kind: tt.fields.kind,
				err:  tt.fields.err,
			}
			w2 := New(tt.fields.kind).Add(tt.fields.err)
			if got := w.Is(tt.args.target); got != tt.want {
				t.Errorf("Is() = %v, want %v", got, tt.want)
			}
			if errors.Is(w, tt.args.target) != tt.want {
				t.Errorf("errors.Is() = %v, want %v", w, tt.want)
			}
			if !errors.Is(w, w2) {
				t.Errorf("errors.Is() = %v, want %v", w, tt.want)
			}
		})
	}
}

func Test_wrapped_Sanitize(t *testing.T) {
	type fields struct {
		kind kind
		err  error
	}
	nn := New(genWrappedErrorsWithSeparator(100)).(wrapped)
	tests := []struct {
		name   string
		fields fields
		want   Grr
	}{
		{
			name: "Sanitize",
			fields: fields{
				kind: kind{
					err:       errors.New("test"),
					separator: _separator(),
				},
				err: errors.New("original"),
			},
			want: kind{
				err:       errors.New("test"),
				separator: _separator(),
			},
		},
		{
			name: "Sanitize100",
			fields: fields{
				kind: kind{
					err:       nn.kind.err,
					separator: nn.kind.separator,
				},
				err: nn.err,
			},
			want: kind{
				err:       errors.New("top1"),
				separator: nn.kind.separator,
			},
		},
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
		wantErr error
	}{
		{
			name: "Unwrap",
			fields: fields{
				kind: kind{
					err:       errors.New("test"),
					separator: _separator(),
				},
				err: errors.New("original"),
			},
			wantErr: errors.New("original"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := wrapped{
				kind: tt.fields.kind,
				err:  tt.fields.err,
			}
			if err := w.Unwrap(); err.Error() != tt.wantErr.Error() {
				t.Errorf("Unwrap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func genWrappedErrorsWithSeparator(count int) error {
	return genWrappedErrWith(count, _separator(), "")
}
