package stock

import (
	"fmt"
	"gitee.com/quant1x/pandas/stat"
	"testing"
)

func TestKLine1(t *testing.T) {
	symbol := "sh600018"
	symbol = "sz002528"
	symbol = "600600.sh"
	df := KLine(symbol, false)
	fmt.Println(df)
}

func TestKLine2(t *testing.T) {
	symbol := "sh600018"
	symbol = "sz002528"
	symbol = "600600.sh"
	df := KLine(symbol, true)
	fmt.Println(df.SelectRows(stat.RangeFinite(-10)))
}

func TestKLineToWeekly(t *testing.T) {
	symbol := "600600.sh"
	symbol = "000001.sh"
	symbol = "880081.sh"
	df := KLine(symbol)
	fmt.Println(df)
	df = KLineToWeekly(df)
	fmt.Println(df.SelectRows(stat.RangeFinite(-10)))
}
