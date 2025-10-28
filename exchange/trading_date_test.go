package exchange

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLastNDate(t *testing.T) {
	dates := LastNDate(Today(), 5)
	fmt.Println(dates)
}

func TestNextTradeDate(t *testing.T) {
	date := NextTradeDate("20230403")
	fmt.Println(date)
}

func TestTradeRange(t *testing.T) {
	start := "2014-08-18"
	end := "2024-08-24"
	dates := TradeRange(start, end)
	fmt.Println(len(dates))
	fmt.Println(dates)
}

func TestGetLastDayForUpdate(t *testing.T) {
	fmt.Println(GetLastDayForUpdate())
}

func TestGetFrontTradeDay(t *testing.T) {
	fmt.Println(GetFrontTradeDay())
}

func TestGetCurrentDate(t *testing.T) {
	date := "20230721"
	v := GetCurrentDate(date)
	fmt.Println(v)
	type args struct {
		date []string
	}
	tests := []struct {
		name            string
		args            args
		wantCurrentDate string
	}{
		{
			name: "20230721 => 2023-07-21",
			args: args{
				date: []string{"20230721"},
			},
			wantCurrentDate: "2023-07-21",
		},
		{
			name: "2025-02-21 => 2025-02-21",
			args: args{
				date: []string{"2025-02-21"},
			},
			wantCurrentDate: "2025-02-21",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.wantCurrentDate, GetCurrentDate(tt.args.date...), "GetCurrentDate(%v)", tt.args.date)
		})
	}
}

func TestDateTimeRange(t *testing.T) {
	getTimeRanges()
	tr := DateTimeRange{Begin: trAMBegin, End: trAMEnd}
	fmt.Println(tr.Minutes())
}

func TestToU32Date(t *testing.T) {
	type args struct {
		datetime string
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		{
			name: "1970-01-01",
			args: args{datetime: "1970-01-01"},
			want: 19700101,
		},
		{
			name: "20231001",
			args: args{datetime: "20231001"},
			want: 20231001,
		},
		{
			name: "2023-10-01",
			args: args{datetime: "2023-10-01"},
			want: 20231001,
		},
		{
			name: "2023-10-01 09:10:11",
			args: args{datetime: "2023-10-01 09:10:11"},
			want: 20231001,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToUint32Date(tt.args.datetime); got != tt.want {
				t.Errorf("ToUint32Date() = %v, want %v", got, tt.want)
			}
		})
	}
}
