package gerr

import (
	"reflect"
	"testing"
)

func Test_separator(t *testing.T) {
	if got := _separator(); got != " " {
		t.Errorf("_separator() = %v, want %v", got, " ")
		t.FailNow()
	}
	SetColonSeparator()
	if got := _separator(); got != ":" {
		t.Errorf("_separator() = %v, want %v", got, ":")
		t.FailNow()
	}
	SetSpaceSeparator()
	if got := _separator(); got != " " {
		t.Errorf("_separator() = %v, want %v", got, " ")
		t.FailNow()
	}
	SetColonAndSpaceSeparator()
	if got := _separator(); got != ": " {
		t.Errorf("_separator() = %v, want %v", got, ": ")
		t.FailNow()
	}
	SetCustomSeparator("--(custom)--")
	if got := _separator(); got != "--(custom)--" {
		t.Errorf("_separator() = %v, want %v", got, "--(custom)--")
		t.FailNow()
	}
	customSeparatorTestfunc(t)
	colonSeparatorTestFunc(t)
	colonAndSpaceSeparatorTestFunc(t)
	spaceSeparatorTestFunc(t)
	defaultSeparatorTestFunc(t)
}

func customSeparatorTestfunc(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "SetDumbSeparator",
			args: args{
				s: "%Very%Dumb%Separator%_X401030_:LLL__",
			},
			want: "%Very%Dumb%Separator%_X401030_:LLL__",
		},
		{
			name: "SetCustomSeparator",
			args: args{
				s: "-custom-",
			},
			want: "-custom-",
		},
		{
			name: "SetEmptySeparator",
			args: args{
				s: "",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_separator = func() string {
				return ""
			}
			SetCustomSeparator(tt.args.s)
			if got := _separator(); got != tt.want {
				t.Errorf("_separator() = %v, want %v", got, tt.want)
				t.FailNow()
			}
			TestNew(t)
			Test_wrapped_Add(t)
			Test_wrapped_Error(t)
			Test_wrapped_Chain(t)
			Test_wrapped_Is(t)
			Test_wrapped_Sanitize(t)
			Test_wrapped_Unwrap(t)
			Test_kind_Add(t)
			Test_kind_Error(t)
			Test_kind_Chain(t)
			Test_kind_Is(t)
			Test_kind_Sanitize(t)
		})
	}
}

func colonSeparatorTestFunc(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "SetColonSeparator",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_separator = func() string {
				return ""
			}
			SetColonSeparator()
			if got := _separator(); got != ":" {
				t.Errorf("_separator() = %v, want %v", got, ":")
				t.FailNow()
			}
			TestNew(t)
			Test_wrapped_Add(t)
			Test_wrapped_Error(t)
			Test_wrapped_Chain(t)
			Test_wrapped_Is(t)
			Test_wrapped_Sanitize(t)
			Test_wrapped_Unwrap(t)
			Test_kind_Add(t)
			Test_kind_Error(t)
			Test_kind_Chain(t)
			Test_kind_Is(t)
			Test_kind_Sanitize(t)
		})
	}
}

func colonAndSpaceSeparatorTestFunc(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "SetColonAndSpaceSeparator",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_separator = func() string {
				return ""
			}
			SetColonAndSpaceSeparator()
			if got := _separator(); got != ": " {
				t.Errorf("_separator() = %v, want %v", got, ": ")
				t.FailNow()
			}
			TestNew(t)
			Test_wrapped_Add(t)
			Test_wrapped_Error(t)
			Test_wrapped_Chain(t)
			Test_wrapped_Is(t)
			Test_wrapped_Sanitize(t)
			Test_wrapped_Unwrap(t)
			Test_kind_Add(t)
			Test_kind_Error(t)
			Test_kind_Chain(t)
			Test_kind_Is(t)
			Test_kind_Sanitize(t)
		})
	}
}

func spaceSeparatorTestFunc(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "SetSpaceSeparator",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_separator = func() string {
				return ""
			}
			SetSpaceSeparator()
			if got := _separator(); got != " " {
				t.Errorf("_separator() = %v, want %v", got, " ")
				t.FailNow()
			}
			TestNew(t)
			Test_wrapped_Add(t)
			Test_wrapped_Error(t)
			Test_wrapped_Chain(t)
			Test_wrapped_Is(t)
			Test_wrapped_Sanitize(t)
			Test_wrapped_Unwrap(t)
			Test_kind_Add(t)
			Test_kind_Error(t)
			Test_kind_Chain(t)
			Test_kind_Is(t)
			Test_kind_Sanitize(t)
		})
	}
}

func defaultSeparatorTestFunc(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "DefaultSeparator",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defaultLike := func() string {
				return " "
			}
			if !reflect.DeepEqual(_separator(), defaultLike()) {
				t.Errorf("_separator() = %v, want %v", _separator(), defaultLike())
				t.FailNow()
			}
			if got := _separator(); got != " " {
				t.Errorf("_separator() = %v, want %v", got, " ")
				t.FailNow()
			}
			TestNew(t)
			Test_wrapped_Add(t)
			Test_wrapped_Error(t)
			Test_wrapped_Chain(t)
			Test_wrapped_Is(t)
			Test_wrapped_Sanitize(t)
			Test_wrapped_Unwrap(t)
			Test_kind_Add(t)
			Test_kind_Error(t)
			Test_kind_Chain(t)
			Test_kind_Is(t)
			Test_kind_Sanitize(t)
		})
	}
}
