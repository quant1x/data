package tdx

import (
	"fmt"
	"gitee.com/quant1x/data/category/date"
	"gitee.com/quant1x/gotdx/proto"
	"golang.org/x/exp/slices"
	"testing"
)

func TestGetKLineAll(t *testing.T) {
	df := GetKLineAll("000001.sh")
	fmt.Println(df)
	dates := df.Col("date").Strings()
	fmt.Println(len(dates))
	// 校验日期的缺失
	start := "1990-12-19"
	end := "2023-04-05"
	dest := date.TradeRange(start, end)
	fmt.Println(len(dest))
	for i, v := range dates {
		found := slices.Contains(dest, v)
		if !found {
			fmt.Println(v)
			tmp := df.IndexOf(i)
			fmt.Println(tmp)
		}
	}
}

func TestGetKLineAll2(t *testing.T) {
	df := GetKLineAll("000001.sh")
	fmt.Println(df)
	df = GetKLineAll("000001.sz", proto.KLINE_TYPE_WEEKLY)
	fmt.Println(df)
}

func TestGetCacheKLine(t *testing.T) {
	code := "sh688981"
	df := GetCacheKLine(code, true)
	fmt.Println(df)
}
