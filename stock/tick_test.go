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

func TestTick(t *testing.T) {
	symbol := "sz002528"
	symbol = "sz000506"
	//dates := date.TradeRange("2023-03-04", "2023-03-19")
	dates := date.TradeRange("2023-02-22", "2023-03-04")
	df := Tick(symbol, dates)
	bv := df.ColAsNDArray("bv")
	ba := df.ColAsNDArray("ba")
	va := ba.Div(bv).Div(100)
	df = df.Join(va)
	fmt.Println(df)
}
