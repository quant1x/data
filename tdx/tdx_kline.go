package tdx

import (
	"gitee.com/quant1x/data/cache"
	"gitee.com/quant1x/data/category"
	"gitee.com/quant1x/data/internal"
	"gitee.com/quant1x/data/security"
	"gitee.com/quant1x/data/security/date"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/pandas"
	"gitee.com/quant1x/pandas/stat"
	"strconv"
)

// GetKLine 获取日K线
func GetKLine(code string, start uint16, count uint16) pandas.DataFrame {
	api := prepare()

	marketId, _, code := security.DetectMarket(code)
	data, _ := api.GetKLine(marketId, code, proto.KLINE_TYPE_RI_K, start, count)
	df := pandas.LoadStructs(data.List)
	df = df.Select([]string{"Open", "Close", "High", "Low", "Vol", "Amount", "DateTime"})
	err := df.SetNames("open", "close", "high", "low", "volume", "amount", "date")
	if err != nil {
		return pandas.DataFrame{}
	}
	df = df.Select([]string{"date", "open", "close", "high", "low", "volume", "amount"})
	return df
}

// KLine 加载K线
func __kLine(code string) pandas.DataFrame {
	filename := cache.KLineFilename(code)
	var df pandas.DataFrame
	if !cache.FileExist(filename) {
		return df
	} else {
		df = pandas.ReadCSV(filename)
	}
	df = df.Select([]string{"date", "open", "close", "high", "low", "volume", "amount"})
	return df
}

// GetKLineAll GetKLine 获取日K线
func GetKLineAll(code string) pandas.DataFrame {
	api := prepare()
	startDate := "19901219"
	marketId, _, code := security.DetectMarket(code)
	df0 := __kLine(code)
	if df0.Nrow() > 0 {
		ds := df0.Col("date").Values().([]string)
		startDate = ds[len(ds)-1]
	} else {
		info, err := api.GetFinanceInfo(marketId, code, 1)
		if err != nil {
			return df0
		}
		startDate = strconv.FormatInt(int64(info.IPODate), 10)

	}
	endDate := cache.Today()
	ts := date.TradeRange(startDate, endDate)
	history := make([]quotes.SecurityBar, 0)
	step := uint16(800)
	total := uint16(len(ts))
	start := uint16(0)
	hs := make([]quotes.SecurityBarsReply, 0)
	for {
		count := step
		if total-start >= step {
			count = step
		} else {
			count = total - start
		}
		data, err := api.GetKLine(marketId, code, proto.KLINE_TYPE_RI_K, start, count)
		if err != nil {
			panic("接口异常")
		}
		hs = append(hs, *data)
		if data.Count < count {
			// 已经是最早的记录
			// 需要排序
			break
		}
		start += count
		if start >= total {
			break
		}
	}
	hs = stat.Reverse(hs)
	for _, v := range hs {
		history = append(history, v.List...)
	}

	//data, _ := api.GetKLine(marketId, code, proto.KLINE_TYPE_RI_K, start, count)
	df1 := pandas.LoadStructs(history)
	df1 = df1.Select([]string{"Open", "Close", "High", "Low", "Vol", "Amount", "DateTime"})
	err := df1.SetNames("open", "close", "high", "low", "volume", "amount", "date")
	if err != nil {
		return pandas.DataFrame{}
	}
	df1 = df1.Select([]string{"date", "open", "close", "high", "low", "volume", "amount"})
	ds1 := df1.Col("date", true)
	ds1.Apply2(func(idx int, v any) any {
		date1 := v.(string)
		dt, err := internal.ParseTime(date1)
		if err != nil {
			return date1
		}
		return dt.Format(category.INDEX_DATE)
	}, true)
	df := df0.Subset(0, df0.Nrow()-1)
	df = df.Concat(df1)

	return df
}
