package exchange

// 交易日, 固定格式YYYY-MM-DD
type TradingDay string

func NewTradingDay(date string) TradingDay {
	date = FixTradeDate(date)
	return TradingDay(date)
}

func (t TradingDay) String() string {
	return string(t)
}
