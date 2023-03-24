package tdx

import (
	"fmt"
	"testing"
)

func TestStockList(t *testing.T) {
	df := get_stock_list()
	fmt.Println(df)
}

func TestGetStockName(t *testing.T) {
	fmt.Println(GetStockName("sh.600600"))
}
