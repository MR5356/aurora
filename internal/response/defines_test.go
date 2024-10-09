package response

import "testing"

func Test_message(t *testing.T) {
	type args struct {
		code string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test-success",
			args: args{
				code: CodeSuccess,
			},
			want: MessageSuccess,
		},
		{
			name: "test-not-found",
			args: args{
				code: CodeNotFound,
			},
			want: MessageNotFound,
		},
		{
			name: "test-params-error",
			args: args{
				code: CodeParamsError,
			},
			want: MessageParamsError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := message(tt.args.code); got != tt.want {
				t.Errorf("message() = %v, want %v", got, tt.want)
			}
		})
	}
}
