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

func TestTradeRange(t *testing.T) {
	start := "2023-04-12"
	end := "2023-04-13"
	dates := TradeRange(start, end)
	fmt.Println(len(dates))
	fmt.Println(dates)
}
