package cache

import "testing"

func TestIdPath(t *testing.T) {
	type args struct {
		code string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "sh",
			args: args{code: "600600"},
			want: "sh600/sh600600",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IdPath(tt.args.code); got != tt.want {
				t.Errorf("IdPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
