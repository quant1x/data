package date

import (
	"fmt"
	"gitee.com/quant1x/data/cache"
	"testing"
)

func TestLastNDate(t *testing.T) {
	dates := LastNDate(cache.Today(), 5)
	fmt.Println(dates)
}

func TestNextTradeDate(t *testing.T) {
	date := NextTradeDate("20230403")
	fmt.Println(date)
}
