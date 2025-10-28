package exchange

import (
	"fmt"
	"testing"
)

func TestTradingTimeRange_Minutes(t *testing.T) {
	tr := ExchangeTradingTimeRange(Trading, "09:30:00", "11:30:00")
	minutes := tr.Minutes()
	fmt.Println(minutes)
}

func TestTradingTimeRange(t *testing.T) {
	text := " 09:30:00 ~ 14:56:30 "
	text = " 14:56:30 - 09:30:00 "
	var tr TradingTimeRange
	tr.Parse(text)
	fmt.Println(tr)
	text = "09:15:00~09:26:59,09:15:00~09:19:59,09:25:00~11:29:59,13:00:00~14:59:59,09:00:00~09:14:59"
	var ts MarketHours
	ts.Parse(text)
	fmt.Println(ts)
	fmt.Println(ts.IsTrading())
}
