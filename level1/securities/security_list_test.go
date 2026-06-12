package securities

import (
	"fmt"
	"testing"

	"github.com/quant1x/data/exchange"
	"github.com/quant1x/data/level1/utils"
)

func TestGetStockName(t *testing.T) {
	code := "sh880635"
	v := GetStockName(code)
	fmt.Println(v)
}

func TestAllCodeList(t *testing.T) {
	v := AllCodeList()
	fmt.Println(v)
	_ = v
}

func TestBaseUnit(t *testing.T) {
	marketId := exchange.MarketIdShangHai
	code := "000001"
	v := utils.BaseUnit(marketId, code)
	fmt.Println(v)
}
