package structutil

import (
	"reflect"
	"testing"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func TestStruct2String(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test",
			args: args{
				v: User{
					ID:   "1",
					Name: "test",
				},
			},
			want: `{
    "id": "1",
    "name": "test"
}`,
		},
		{
			name: "test-with-nil",
			args: args{
				v: User{
					ID: "1",
				},
			},
			want: `{
    "id": "1",
    "name": ""
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Struct2String(tt.args.v); got != tt.want {
				t.Errorf("Struct2String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStruct2Map(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name string
		args args
		want map[string]any
	}{
		{
			name: "test",
			args: args{
				v: &User{
					Name: "test",
				},
			},
			want: map[string]any{
				"Name": "test",
			},
		},
		{
			name: "test-with-nil",
			args: args{
				v: &User{
					ID: "1",
				},
			},
			want: map[string]any{
				"Name": "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Struct2Map(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Struct2Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyIsNil(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test-string-nil",
			args: args{
				v: "",
			},
			want: true,
		},
		{
			name: "test-string-not-nil",
			args: args{
				v: "test",
			},
			want: false,
		},
		{
			name: "test-int-nil",
			args: args{
				v: 0,
			},
			want: true,
		},
		{
			name: "test-int-not-nil",
			args: args{
				v: 1,
			},
			want: false,
		},
		{
			name: "test-float-nil",
			args: args{
				v: 0.0,
			},
			want: true,
		},
		{
			name: "test-float-not-nil",
			args: args{
				v: 1.0,
			},
			want: false,
		},
		{
			name: "test-bool-nil",
			args: args{
				v: false,
			},
			want: true,
		},
		{
			name: "test-bool-not-nil",
			args: args{
				v: true,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AnyIsNil(tt.args.v); got != tt.want {
				t.Errorf("AnyIsNil() = %v, want %v", got, tt.want)
			}
		})
	}
}
