package cache

import "testing"

func TestTickPath(t *testing.T) {
	type args struct {
		code string
		date string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "yyyy-mm-dd",
			args: args{
				code: "600600.sh",
				date: "2020-01-02",
			},
			want:    GetTickPath() + "/2020/20200102/sh600600.csv",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TickFilename(tt.args.code, tt.args.date)
			if got != tt.want {
				t.Errorf("TickPath() got = %v, want %v", got, tt.want)
			}
		})
	}
}
