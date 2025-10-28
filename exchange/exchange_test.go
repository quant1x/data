package exchange

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetExchange(t *testing.T) {
	sessions := GetExchange(CN)
	fmt.Println(sessions)
	fmt.Println(sessions.Minutes())
}

func TestGetSecurityCode(t *testing.T) {
	fmt.Println(GetSecurityCode(1, "600600"))
	fmt.Println(GetSecurityCode(0, "399001"))
	fmt.Println(GetSecurityCode(2, "399001"))
}

func TestGetMarket(t *testing.T) {
	code := "sh600600"
	fmt.Println(GetMarket(code))
	code = "sh.600600"
	fmt.Println(GetMarket(code))
	code = "600600.sh"
	fmt.Println(GetMarket(code))

	code = "880818"
	v := AssertBlockBySecurityCode(&code)
	fmt.Println(v)
}

func TestCorrectSecurityCode(t *testing.T) {
	correctedCode := CorrectSecurityCode("")
	assert.Equal(t, 0, len(correctedCode))
}

func TestGetExchangeMarketHours(t *testing.T) {
	sessions := GetExchangeMarketHours(CN)
	fmt.Println(sessions)
	fmt.Println(sessions.Minutes())

	opt := GetOperator()
	kind, _ := opt.Kind()
	fmt.Println(kind)

	timestamp := createTimeOfToday("09:18:01").UnixMilli()
	kind, idx := opt.Kind(timestamp)
	assert.Equal(t, CallAuctionAndCancel, kind)
	assert.Equal(t, 0, idx)

	timestamp = createTimeOfToday("09:21:00").UnixMilli()
	kind, idx = opt.Kind(timestamp)
	assert.Equal(t, CallAuction, kind)
	assert.Equal(t, 1, idx)

	timestamp = createTimeOfToday("10:20:00").UnixMilli()
	kind, idx = opt.Kind(timestamp)
	assert.Equal(t, TradingAndCancel, kind)
	assert.Equal(t, 2, idx)

	timestamp = createTimeOfToday("13:20:00").UnixMilli()
	kind, idx = opt.Kind(timestamp)
	assert.Equal(t, TradingAndCancel, kind)
	assert.Equal(t, 3, idx)

	timestamp = createTimeOfToday("14:58:00").UnixMilli()
	kind, idx = opt.Kind(timestamp)
	assert.Equal(t, CallAuctionAndTrading, kind)
	assert.Equal(t, 4, idx)

	timestamp = createTimeOfToday("08:20:00").UnixMilli()
	kind, idx = opt.Kind(timestamp)
	assert.Equal(t, TimeKind(0x0), kind)
	assert.Equal(t, -1, idx)

}
