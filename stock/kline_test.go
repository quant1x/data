package stock

import (
	"fmt"
	"testing"
)

func TestKLine(t *testing.T) {
	symbol := "sh600018"
	symbol = "sz002528"
	symbol = "600600.sh"
	df2 := KLine(symbol)
	fmt.Println(df2)
}
