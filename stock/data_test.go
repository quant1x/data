package stock

import (
	"fmt"
	"testing"
)

func TestTick(t *testing.T) {
	symbol := "600600"
	df := Tick(symbol)
	fmt.Println(df)
}

func TestKLine(t *testing.T) {
	symbol2 := "sh600018"
	symbol2 = "sz002528"
	df2 := KLine(symbol2)
	fmt.Println(df2)
}
