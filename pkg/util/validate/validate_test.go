package validate

import (
	"errors"
	"testing"
)

type User struct {
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"required,gt=18"`
}

func TestValidate(t *testing.T) {
	type args struct {
		obj any
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "test-nil",
			args:    args{obj: &User{}},
			wantErr: errors.New("Name为必填字段, Age为必填字段"),
		},
		{
			name:    "test-with-name",
			args:    args{obj: &User{Name: "李华"}},
			wantErr: errors.New("Age为必填字段"),
		},
		{
			name:    "test-with-name-and-age",
			args:    args{obj: &User{Name: "李华", Age: 18}},
			wantErr: errors.New("Age必须大于18"),
		},
		{
			name:    "test-pass",
			args:    args{obj: &User{Name: "李华", Age: 19}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Validate(tt.args.obj); (err != nil && tt.wantErr != nil) && err.Error() != tt.wantErr.Error() {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
