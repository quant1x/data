package stock

import (
	"fmt"
	"gitee.com/quant1x/data/category/date"
	"testing"
)

func TestTickByDate(t *testing.T) {
	symbol := "600600"
	symbol = "sz002528"
	df := TickByDate(symbol, "2023-03-03")
	fmt.Println(df)
}

func TestKLine(t *testing.T) {
	symbol := "sh600018"
	symbol = "sz002528"
	symbol = "600600.sh"
	df2 := KLine(symbol)
	fmt.Println(df2)
}

func TestTick(t *testing.T) {
	symbol := "sz002528"
	symbol = "sz002209"
	dates := date.TradeRange("2023-01-17", "2023-01-19")
	df := Tick(symbol, dates)
	bv := df.ColAsNDArray("bv")
	ba := df.ColAsNDArray("ba")
	va := ba.Div(bv).Div(100)
	df = df.Join(va)
	fmt.Println(df)
}
