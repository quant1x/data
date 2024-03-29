package tdx

import (
	"gitee.com/quant1x/data/cache"
	"gitee.com/quant1x/data/category"
	"gitee.com/quant1x/data/category/date"
	"gitee.com/quant1x/data/internal"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/pandas"
	"gitee.com/quant1x/pandas/formula"
	"gitee.com/quant1x/pandas/stat"
	"reflect"
	"strconv"
)

var (
	FBarsStockRaw    = []string{"Open", "Close", "High", "Low", "Vol", "Amount", "DateTime"}
	FBarsStockRename = []string{"open", "close", "high", "low", "volume", "amount", "date"}
	FBarsStockFields = []string{"date", "open", "close", "high", "low", "volume", "amount", "bv", "sv", "ba", "sa"}

	FBarsIndexRaw    = []string{"Open", "Close", "High", "Low", "Vol", "Amount", "DateTime", "UpCount", "DownCount"}
	FBarsIndexRename = []string{"open", "close", "high", "low", "volume", "amount", "date", "up", "down"}
	FBarsIndexFields = []string{"date", "open", "close", "high", "low", "volume", "amount", "up", "down", "bv", "sv", "ba", "sa"}
)

// SecurityBar K线数据
//type SecurityBar struct {
//	Open      float64
//	Close     float64
//	High      float64
//	Low       float64
//	Vol       float64
//	Amount    float64
//	Year      int
//	Month     int
//	Day       int
//	Hour      int
//	Minute    int
//	DateTime  string
//	UpCount   uint16  // 指数有效, 上涨家数
//	DownCount uint16  // 指数有效, 下跌家数
//	BuyVol    float64 // 外盘
//	SellVol   float64 // 内盘
//	BuyAmt    float64 // 外盘成交金额
//	SellAmt   float64 // 内盘成交金额
//}

// getKLine 获取日K线
//func getKLine(code string, start uint16, count uint16) pandas.DataFrame {
//	api := prepare()
//	marketId, _, code := category.DetectMarket(code)
//	data, _ := api.GetKLine(marketId, code, proto.KLINE_TYPE_RI_K, start, count)
//	df := pandas.LoadStructs(data.List)
//	df = df.Select([]string{"Open", "Close", "High", "Low", "Vol", "Amount", "DateTime"})
//	err := df.SetNames("open", "close", "high", "low", "volume", "amount", "date")
//	if err != nil {
//		return pandas.DataFrame{}
//	}
//	df = df.Select([]string{"date", "open", "close", "high", "low", "volume", "amount"})
//	return df
//}

// GetCacheKLine 加载K线
//
//	第2个参数, 是否前复权
func GetCacheKLine(code string, argv ...bool) pandas.DataFrame {
	// 默认不复权
	qfq := false
	if len(argv) > 0 {
		qfq = argv[0]
	}
	isIndex := category.CodeIsIndex(code)
	fields := FBarsStockFields
	if isIndex {
		fields = FBarsIndexFields
	}
	filename := cache.KLineFilename(code)
	var df pandas.DataFrame
	if !cache.FileExist(filename) {
		return df
	} else {
		df = pandas.ReadCSV(filename)
	}
	df = df.Select(fields)
	if df.Nrow() == 0 {
		return df
	}
	if qfq {
		drdf := GetCacheXdxr(code)
		for i := 0; i < drdf.Nrow(); i++ {
			m0 := drdf.IndexOf(i)
			if m0["Category"].(int64) != 1 {
				continue
			}
			end := m0["Date"].(string)
			songZhuangu := stat.AnyToFloat64(m0["SongZhuanGu"])
			peiGu := stat.AnyToFloat64(m0["PeiGu"])
			suoGu := stat.AnyToFloat64(m0["SuoGu"])
			xdxrGuShu := (songZhuangu + peiGu - suoGu) / 10
			fenHong := stat.AnyToFloat64(m0["FenHong"])
			peiGuJia := stat.AnyToFloat64(m0["PeiGuJia"])
			xdxrFenHong := (peiGuJia*peiGu - fenHong) / 10
			for i := 0; i < df.Nrow(); i++ {
				m1 := df.IndexOf(i, true)
				dt := m1["date"].(reflect.Value).String()
				if dt > end {
					break
				}
				po := m1["open"].(reflect.Value)
				po.SetFloat((po.Float() + xdxrFenHong) / (1 + xdxrGuShu))
				pc := m1["close"].(reflect.Value)
				pc.SetFloat((pc.Float() + xdxrFenHong) / (1 + xdxrGuShu))
				ph := m1["high"].(reflect.Value)
				ph.SetFloat((ph.Float() + xdxrFenHong) / (1 + xdxrGuShu))
				pl := m1["low"].(reflect.Value)
				pl.SetFloat((pl.Float() + xdxrFenHong) / (1 + xdxrGuShu))
				if dt == end {
					break
				}
			}
		}
	}
	// 取出成交量序列
	VOL := df.Col("volume")
	DATES := df.Col("date")
	CLOSE := df.Col("close")
	lastDay := DATES.IndexOf(-1).(string)
	total := df.Nrow()
	// 计算5日均量
	mv5 := formula.MA(VOL, 5)
	mav := formula.REF(mv5, 1)
	lb := VOL.Div(mav)
	lb = lb.Apply2(func(idx int, v any) any {
		if idx+1 < total {
			return v
		} else {
			tmp := stat.Any2DType(v)
			ms := stat.DType(date.Minutes(lastDay)) / float64(date.CN_TOTALFZNUM)
			tmp /= ms
			return tmp
		}
	}, true)
	prevClose := stat.NewSeries(stat.Shift(CLOSE.DTypes(), 1)...)
	zf := CLOSE.Div(prevClose).Sub(1).Mul(100)

	// 链接量比序列
	oLB := pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "lb", lb)
	oZF := pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "zf", zf)
	df = df.Join(oLB, oZF)
	return df
}

// GetKLineAll getKLine 获取日K线
func GetKLineAll(fullCode string, argv ...int) pandas.DataFrame {
	kType := uint16(proto.KLINE_TYPE_RI_K)
	if len(argv) == 1 {
		kType = uint16(argv[0])
	}
	tdxApi := prepare()
	startDate := "19901219"
	marketId, market, code := category.DetectMarket(fullCode)
	dfCache := GetCacheKLine(market + code)
	isIndex := category.CodeIsIndex(fullCode)
	fields := FBarsStockFields
	rawFields := FBarsStockRaw
	newFields := FBarsStockRename
	if isIndex {
		fields = FBarsIndexFields
		rawFields = FBarsIndexRaw
		newFields = FBarsIndexRename
	}
	// 尝试选择一次字段, 如果出现异常, 则清空dataframe, 重新下载
	dfCache = dfCache.Select(fields)
	if dfCache.Nrow() == 0 {
		dfCache = pandas.DataFrame{}
	}
	var info *quotes.FinanceInfo
	var err error
	if dfCache.Nrow() > 0 {
		ds := dfCache.Col("date").Strings()
		startDate = ds[len(ds)-1]
	} else {
		info, err = tdxApi.GetFinanceInfo(marketId, code, 1)
		if err != nil {
			return dfCache
		}
		if info.IPODate > 0 {
			startDate = strconv.FormatInt(int64(info.IPODate), 10)
		}
		if info.IPODate == 0 && info.LiuTongGuBen > 0 && info.ZongGuBen > 0 && info.BaoLiu2 > 0 {
			isIndex = true
		}
	}
	if !isIndex {
		if info == nil {
			info, err = tdxApi.GetFinanceInfo(marketId, code, 1)
			if err != nil {
				return dfCache
			}
		}
		if info.IPODate == 0 && info.LiuTongGuBen > 0 && info.ZongGuBen > 0 && info.BaoLiu2 > 0 {
			isIndex = true
		}
	}
	endDate := cache.Today()
	ts := date.TradeRange(startDate, endDate)
	history := make([]quotes.SecurityBar, 0)
	step := uint16(quotes.TDX_SECURITY_BARS_MAX)
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
		var data *quotes.SecurityBarsReply
		var err error
		if isIndex {
			data, err = tdxApi.GetIndexBars(marketId, code, kType, start, count)
		} else {
			data, err = tdxApi.GetKLine(marketId, code, kType, start, count)
		}
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
		for _, row := range v.List {
			dateTime := row.DateTime
			dt, err := internal.ParseTime(dateTime)
			if err != nil {
				dateTime = row.DateTime[0:len(category.INDEX_DATE)]
			} else {
				dateTime = dt.Format(category.INDEX_DATE)
			}
			if dateTime < startDate {
				continue
			}
			history = append(history, row)
		}
	}

	df1 := pandas.LoadStructs(history)
	df1 = df1.Select(rawFields)
	err = df1.SetNames(newFields...)
	if err != nil {
		return pandas.DataFrame{}
	}
	ds1 := df1.Col("date", true)
	ds1.Apply2(func(idx int, v any) any {
		date1 := v.(string)
		dt, err := internal.ParseTime(date1)
		if err != nil {
			return date1
		}
		return dt.Format(category.INDEX_DATE)
	}, true)
	df1 = attachVolume(df1, fullCode)
	df1 = df1.Select(fields)
	// 计算新增的天数
	//tmpDates := df1.Col("date").Strings()
	//tmpDays := df1.Nrow()
	//fixDays := int(total) - tmpDays
	df := dfCache.Subset(0, dfCache.Nrow()-1)
	if df.Nrow() > 0 {
		df = df.Concat(df1)
	} else {
		df = df1
	}

	return df
}
