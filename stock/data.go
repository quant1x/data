package stock

import (
	"gitee.com/quant1x/data/cache"
	"gitee.com/quant1x/data/category"
	"gitee.com/quant1x/data/internal/tdx"
	"gitee.com/quant1x/pandas"
)

// KLine 加载K线
func KLine(code string) pandas.DataFrame {
	return tdx.GetCacheKLine(code, true)
}

// Tick 加载tick缓存数据
func Tick(code string, args ...string) pandas.DataFrame {
	var date = ""
	if len(args) > 0 {
		date = cache.CorrectDate(args[0])
	} else {
		date = cache.Today()
	}
	var df pandas.DataFrame
	filename := cache.TickFilename(code, date)
	if !cache.FileExist(filename) {
		df = tdx.GetTickData(code, date)
	} else {
		df = pandas.ReadCSV(filename)
	}
	if df.Nrow() == 0 {
		tm := cache.LastDate()
		stm := tm.Format(category.CACHE_DATE)
		return Tick(code, stm)
	}
	return df
}
