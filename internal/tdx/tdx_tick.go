package tdx

import (
	"fmt"
	"gitee.com/quant1x/data/cache"
	"gitee.com/quant1x/data/category"
	"gitee.com/quant1x/data/category/date"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/pandas"
	"gitee.com/quant1x/pandas/stat"
	"github.com/mymmsc/gox/logger"
	"github.com/mymmsc/gox/progressbar"
	"strconv"
)

// GetTickAll 下载全部tick数据
func GetTickAll(code string) {
	api := prepare()
	marketId, _, code := category.DetectMarket(code)
	info, err := api.GetFinanceInfo(marketId, code, 1)
	if err != nil {
		return
	}
	tStart := strconv.FormatInt(int64(info.IPODate), 10)
	fixStart := cache.TickStartDate
	if tStart < fixStart {
		tStart = fixStart
	}
	tEnd := cache.Today()
	logger.Errorf("[%s]tick数据范围: %s<->%s", code, tStart, tEnd)
	dateRange := date.TradeRange(tStart, tEnd)
	// 反转切片
	dateRange = stat.Reverse(dateRange)
	bar := progressbar.NewBar(2, fmt.Sprintf("同步[%s]", code), len(dateRange))
	today := date.IndexToday()
	ignore := false
	for _, tradeDate := range dateRange {
		bar.Add(1)
		if ignore {
			continue
		}
		fname := cache.TickFilename(code, tradeDate)
		if tradeDate != today && cache.FileExist(fname) {
			// 如果已经存在, 假定之前的数据已经下载过了, 不需要继续
			ignore = true
			continue
		}
		df := GetTickData(code, tradeDate)
		if df.Nrow() == 0 && tradeDate != today {
			// 如果数据为空, 且不是当前日期, 认定为从这天起往前是没有分笔成交数据的
			if logger.IsDebug() {
				logger.Errorf("[%s]tick数据[%s<->%s]空, 后面忽略", code, tradeDate, today)
			}
			ignore = true
		}
	}

	return
}

// GetTickData 获取指定日期的分笔成交记录
func GetTickData(code string, date string) pandas.DataFrame {
	tdxApi := prepare()
	marketId, marketName, code := category.DetectMarket(code)
	offset := uint16(quotes.TDX_TRANSACTION_MAX)
	start := uint16(0)
	date = cache.CorrectDate(date)
	history := make([]quotes.TickTransaction, 0)
	hs := make([]quotes.TransactionReply, 0)
	for {
		var data *quotes.TransactionReply
		var err error
		if date == cache.Today() {
			data, err = tdxApi.GetTransactionData(marketId, code, start, offset)
		} else {
			iDate := stat.AnyToInt64(date)
			data, err = tdxApi.GetHistoryTransactionData(marketId, code, uint32(iDate), start, offset)
		}

		if err != nil {
			panic("接口异常")
		}
		hs = append(hs, (*data))
		if data.Count < offset {
			// 已经是最早的记录
			// 需要排序
			break
		}
		start += offset
	}
	hs = stat.Reverse(hs)
	for _, v := range hs {
		history = append(history, v.List...)
	}
	_ = marketName
	df := pandas.LoadStructs(history)
	df = df.Select([]string{"Time", "Price", "Vol", "Num", "BuyOrSell"})
	err := df.SetNames("time", "price", "vol", "num", "buyorsell")
	if err != nil {
		return pandas.DataFrame{}
	}
	tickFile := cache.TickFilename(code, date)
	if err != nil {
		return pandas.DataFrame{}
	}
	err = cache.CheckFilepath(tickFile)
	if err != nil {
		return pandas.DataFrame{}
	}
	err = df.WriteCSV(tickFile)
	if err != nil {
		return pandas.DataFrame{}
	}

	return df
}
