package stock

import (
	"fmt"
	"gitee.com/quant1x/pandas/stat"
	"testing"
)

func TestKLine(t *testing.T) {
	symbol := "sh600018"
	symbol = "sz002528"
	symbol = "600600.sh"
	df2 := KLine(symbol, true)
	fmt.Println(df2.SelectRows(stat.RangeFinite(-10)))
}

func TestKLineToWeekly(t *testing.T) {
	symbol := "600600.sh"
	df := KLine(symbol)
	df = KLineToWeekly(df)
	fmt.Println(df.SelectRows(stat.RangeFinite(-10)))
}
