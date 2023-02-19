package security

import "testing"

func TestDetectMarket(t *testing.T) {
	type args struct {
		symbol string
	}
	tests := []struct {
		name  string
		args  args
		want  Market
		want1 string
		want2 string
	}{
		{
			name:  "测试裸代码",
			args:  args{symbol: "600600"},
			want:  MARKET_ID_SHANGHAI,
			want1: MARKET_SH,
			want2: "600600",
		},
		{
			name:  "测试裸代码",
			args:  args{symbol: "002528"},
			want:  MARKET_ID_SHENZHEN,
			want1: MARKET_SZ,
			want2: "002528",
		},
		{
			name:  "测试前缀",
			args:  args{symbol: "sz002528"},
			want:  MARKET_ID_SHENZHEN,
			want1: MARKET_SZ,
			want2: "002528",
		},
		{
			name:  "测试前缀",
			args:  args{symbol: "sh600600"},
			want:  MARKET_ID_SHANGHAI,
			want1: MARKET_SH,
			want2: "600600",
		},
		{
			name:  "测试后缀",
			args:  args{symbol: "002528.sz"},
			want:  MARKET_ID_SHENZHEN,
			want1: MARKET_SZ,
			want2: "002528",
		},
		{
			name:  "测试后缀",
			args:  args{symbol: "600600.sh"},
			want:  MARKET_ID_SHANGHAI,
			want1: MARKET_SH,
			want2: "600600",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := DetectMarket(tt.args.symbol)
			if got != tt.want {
				t.Errorf("DetectMarket() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("DetectMarket() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("DetectMarket() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
