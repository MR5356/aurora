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
			name: "test-int8-nil",
			args: args{
				v: int8(0),
			},
			want: true,
		},
		{
			name: "test-int8-not-nil",
			args: args{
				v: int8(1),
			},
			want: false,
		},
		{
			name: "test-int16-nil",
			args: args{
				v: int16(0),
			},
			want: true,
		},
		{
			name: "test-int16-not-nil",
			args: args{
				v: int16(1),
			},
			want: false,
		},
		{
			name: "test-int32-nil",
			args: args{
				v: int32(0),
			},
			want: true,
		},
		{
			name: "test-int32-not-nil",
			args: args{
				v: int32(1),
			},
			want: false,
		},
		{
			name: "test-int64-nil",
			args: args{
				v: int64(0),
			},
			want: true,
		},
		{
			name: "test-int64-not-nil",
			args: args{
				v: int64(1),
			},
			want: false,
		},
		{
			name: "test-uint-nil",
			args: args{
				v: uint(0),
			},
			want: true,
		},
		{
			name: "test-uint-not-nil",
			args: args{
				v: uint(1),
			},
			want: false,
		},
		{
			name: "test-uint8-nil",
			args: args{
				v: uint8(0),
			},
			want: true,
		},
		{
			name: "test-uint8-not-nil",
			args: args{
				v: uint8(1),
			},
			want: false,
		},
		{
			name: "test-uint16-nil",
			args: args{
				v: uint16(0),
			},
			want: true,
		},
		{
			name: "test-uint16-not-nil",
			args: args{
				v: uint16(1),
			},
			want: false,
		},
		{
			name: "test-uint32-nil",
			args: args{
				v: uint32(0),
			},
			want: true,
		},
		{
			name: "test-uint32-not-nil",
			args: args{
				v: uint32(1),
			},
			want: false,
		},
		{
			name: "test-uint64-nil",
			args: args{
				v: uint64(0),
			},
			want: true,
		},
		{
			name: "test-uint64-not-nil",
			args: args{
				v: uint64(1),
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
			name: "test-float32-nil",
			args: args{
				v: float32(0.0),
			},
			want: true,
		},
		{
			name: "test-float32-not-nil",
			args: args{
				v: float32(1.0),
			},
			want: false,
		},
		{
			name: "test-float64-nil",
			args: args{
				v: float64(0.0),
			},
			want: true,
		},
		{
			name: "test-float64-not-nil",
			args: args{
				v: float64(1.0),
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
		{
			name: "test-default-nil",
			args: args{
				v: User{},
			},
			want: true,
		},
		{
			name: "test-default-not-nil",
			args: args{
				v: User{
					ID: "1",
				},
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

func TestStruct2Bytes(t *testing.T) {
	type args[T any] struct {
		v T
	}
	type testCase[T any] struct {
		name    string
		args    args[T]
		want    []byte
		wantErr bool
	}
	tests := []testCase[User]{
		{
			name: "test",
			args: args[User]{
				v: User{
					ID:   "1",
					Name: "test",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Struct2Bytes(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("Struct2Bytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestBytes2Struct(t *testing.T) {
	user := User{
		ID:   "1",
		Name: "test",
	}

	b, _ := Struct2Bytes(user)

	type args struct {
		data []byte
	}
	type testCase[T any] struct {
		name    string
		args    args
		want    T
		wantErr bool
	}
	tests := []testCase[User]{
		{
			name: "test",
			args: args{
				data: b,
			},
			want: User{
				ID:   "1",
				Name: "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Bytes2Struct[User](tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Bytes2Struct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bytes2Struct() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytesArray2Struct(t *testing.T) {
	users := make([][]byte, 0)
	user := User{
		ID:   "1",
		Name: "test",
	}

	b, _ := Struct2Bytes(user)
	users = append(users, b)

	type args struct {
		data [][]byte
	}
	type testCase[T any] struct {
		name    string
		args    args
		want    []T
		wantErr bool
	}
	tests := []testCase[User]{
		{
			name: "test",
			args: args{
				data: users,
			},
			want: []User{
				{
					ID:   "1",
					Name: "test",
				},
			},
			wantErr: false,
		},
		{
			name: "test",
			args: args{
				data: [][]byte{
					[]byte("test"),
				},
			},
			want:    make([]User, 0),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BytesArray2Struct[User](tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("BytesArray2Struct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BytesArray2Struct() got = %v, want %v", got, tt.want)
			}
		})
	}
}
