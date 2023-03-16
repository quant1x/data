package tdx

import (
	"gitee.com/quant1x/data/cache"
	"gitee.com/quant1x/data/category"
	"gitee.com/quant1x/data/category/date"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/pandas/stat"
)

// GetMinuteHistory 获得指定日期的分时数据
func GetMinuteHistory(code, date string) *quotes.HistoryMinuteTimeReply {
	tdxApi := prepare()
	marketId, _, code := category.DetectMarket(code)
	date = cache.CorrectDate(date)
	iDate := stat.AnyToInt64(date)
	hs, err := tdxApi.GetHistoryMinuteTimeData(marketId, code, uint32(iDate))
	if err != nil {
		return nil
	}
	return hs
}

// QuantityRelativeRatio 获得指定日期的分时数据
//
//	量比是衡量相对成交量的指标。
//	它是指股市开市后平均每分钟的成交量与过去5个交易日平均每分钟成交量之比。
//	其计算公式为：量比=（现成交总手数 / 现累计开市时间(分) ）/ 过去5日平均每分钟成交量
func QuantityRelativeRatio(code string) float64 {
	today := cache.Today()
	dates := date.LastNDate(today, 5)
	hs_num := 0
	hs_vols := 0
	for _, dt := range dates {
		hs := GetMinuteHistory(code, dt)
		if hs != nil && hs.Count > 0 {
			for _, v := range hs.List {
				hs_vols += v.Vol
				hs_num += 1
			}
		}
	}
	ratio5 := float64(hs_vols) / float64(hs_num)
	hs := GetMinuteHistory(code, date.LastTradeDate())
	ts_num := 0
	ts_vols := 0
	if hs != nil && hs.Count > 0 {
		for _, v := range hs.List {
			ts_vols += v.Vol
			ts_num += 1
		}
	}
	todayRatio := float64(ts_vols) / float64(ts_num)
	return todayRatio / ratio5
}
