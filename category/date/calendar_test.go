package date

import (
	"fmt"
	"testing"
)

func Test_updateHoliday(t *testing.T) {
	updateHoliday()
}

func TestIsHoliday(t *testing.T) {
	type args struct {
		date string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "周末",
			args: args{date: "2023-02-18"},
			want: true,
		},
		{
			name: "周末",
			args: args{date: "2023-02-19"},
			want: true,
		},
		{
			name: "春节",
			args: args{date: "2023-01-23"},
			want: true,
		},
		{
			name: "工作日",
			args: args{date: "2023-02-20"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsHoliday(tt.args.date); got != tt.want {
				t.Errorf("IsHoliday() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTradeRange(t *testing.T) {
	ds := TradeRange("2023-02-17", "2023-03-01")
	fmt.Println(len(ds))
	for _, v := range ds {
		fmt.Println(v)
	}
}
