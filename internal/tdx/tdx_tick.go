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
	marketId, marketName, code := category.DetectMarket(code)
	code = marketName + code
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

// 获取指定日期的分笔成交记录
func tickData(code string, date string, ignore bool) pandas.DataFrame {
	var df pandas.DataFrame
	if ignore {
		// 在默认日期之前的数据直接返回空
		startDate := cache.CorrectDate(cache.TickStartDate)
		if date < startDate {
			return df
		}
	}
	filename := cache.GetTickFilename(code, date, false)
	if cache.FileExist(filename) {
		df = pandas.ReadCSV(filename)
		return df
	}

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
		hs = append(hs, *data)
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
	df = pandas.LoadStructs(history)
	df = df.Select([]string{"Time", "Price", "Vol", "Num", "BuyOrSell"})
	err := df.SetNames("time", "price", "vol", "num", "buyorsell")
	if err != nil {
		return pandas.DataFrame{}
	}

	return df
}

// GetTickData 获取指定日期的分笔成交记录
func GetTickData(code string, date string) pandas.DataFrame {
	df := tickData(code, date, false)
	tickFile := cache.TickFilename(code, date)
	err := cache.CheckFilepath(tickFile)
	if err != nil {
		return pandas.DataFrame{}
	}
	err = df.WriteCSV(tickFile)
	if err != nil {
		return pandas.DataFrame{}
	}

	return df
}

// 附加成交量
func attachVolume(df pandas.DataFrame, code string) pandas.DataFrame {
	dates := df.Col("date").Strings()
	if len(dates) == 0 {
		return df
	}
	buyVolumes := []stat.DType{}
	sellVolumes := []stat.DType{}
	buyAmounts := []stat.DType{}
	sellAmounts := []stat.DType{}
	for _, tradedate := range dates {
		buyVolume := stat.DType(0)
		sellVolume := stat.DType(0)
		buyAmount := stat.DType(0)
		sellAmount := stat.DType(0)
		tmp := tickData(code, tradedate, true)
		if tmp.Nrow() > 0 {
			for i := 0; i < tmp.Nrow(); i++ {
				m := tmp.IndexOf(i)
				t := stat.AnyToInt32(m["buyorsell"])
				p := stat.AnyToFloat64(m["price"])
				v := stat.AnyToFloat64(m["vol"])
				if t == 1 {
					// 卖出
					sellVolume += v
					sellAmount += v * p * 100
				} else {
					buyVolume += v
					buyAmount += v * p * 100
				}
			}
		}
		buyVolumes = append(buyVolumes, buyVolume)
		sellVolumes = append(sellVolumes, sellVolume)
		buyAmounts = append(buyAmounts, buyAmount)
		sellAmounts = append(sellAmounts, sellAmount)
	}
	bv := pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "bv", buyVolumes)
	sv := pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "sv", sellVolumes)
	ba := pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "ba", buyAmounts)
	sa := pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "sa", sellAmounts)
	df = df.Join(bv, sv, ba, sa)
	return df
}
