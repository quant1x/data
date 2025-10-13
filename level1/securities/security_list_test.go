package securities

import (
	"fmt"
	"testing"

	"gitee.com/quant1x/data/level1/internal"
	"gitee.com/quant1x/exchange"
)

func TestGetStockName(t *testing.T) {
	code := "sh880635"
	v := GetStockName(code)
	fmt.Println(v)
}

func TestAllCodeList(t *testing.T) {
	v := AllCodeList()
	fmt.Println(v)
}

func TestBaseUnit(t *testing.T) {
	marketId := exchange.MarketIdShangHai
	code := "000001"
	v := internal.BaseUnit(marketId, code)
	fmt.Println(v)
}
