package stock

import (
	"gitee.com/quant1x/data/category/date"
	"gitee.com/quant1x/data/internal"
	"gitee.com/quant1x/data/internal/tdx"
	"gitee.com/quant1x/pandas"
	"gitee.com/quant1x/pandas/stat"
)

// KLine 加载K线
func KLine(code string, argv ...bool) pandas.DataFrame {
	weekday := false
	if len(argv) == 1 {
		weekday = stat.AnyToBool(argv[0])
	}
	df9 := tdx.GetCacheKLine(code, true)
	if !weekday {
		return df9
	}
	return KLineToWeekly(df9)
}

// KLineToWeekly 日线转周线
func KLineToWeekly(kline pandas.DataFrame) pandas.DataFrame {
	// 周线
	var df pandas.DataFrame
	//date,open,close,high,low,volume,amount,up,down
	var wdate string
	var o, c, h, l, v, a stat.DType
	var bv, sv, ba, sa stat.DType
	var prevClose stat.DType
	for i := 0; i < kline.Nrow(); i++ {
		m := kline.IndexOf(i)
		//date,open,close,high,low,volume,amount,up,down
		// 周线日期以最后一天的日期为准
		_date, ok := m["date"].(string)
		if ok {
			wdate = _date
		}
		// 周线开盘价以第一天OPEN为准
		_open, ok := m["open"].(stat.DType)
		if ok && o == stat.DType(0) {
			o = _open
		}
		// 周线的收盘价以本周最后一个交易日的CLOSE为准
		_close, ok := m["close"].(stat.DType)
		if ok {
			c = _close
		}
		// 涨幅
		zf := (c/prevClose - 1.00) * 100.00
		_high, ok := m["high"].(stat.DType)
		if ok && h == stat.DType(0) {
			h = _high
		}
		if h < _high {
			h = _high
		}
		_low, ok := m["low"].(stat.DType)
		if ok && l == stat.DType(0) {
			l = _low
		}
		if l > _low {
			l = _low
		}
		_vol, ok := m["volume"]
		if ok {
			v += stat.Any2DType(_vol)
		}
		_amount, ok := m["amount"].(stat.DType)
		if ok {
			a += _amount
		}
		_bv, ok := m["bv"]
		if ok {
			bv += stat.Any2DType(_bv)
		}
		_sv, ok := m["sv"]
		if ok {
			sv += stat.Any2DType(_sv)
		}
		_ba, ok := m["ba"]
		if ok {
			ba += stat.Any2DType(_ba)
		}
		_sa, ok := m["sa"]
		if ok {
			sa += stat.Any2DType(_sa)
		}
		dt, _ := internal.ParseTime(wdate)
		w := int(dt.Weekday())
		last := false
		today := date.IndexToday()
		if wdate == today {
			last = true
		}
		// 如果是周五
		if !last && w == 5 {
			last = true
		}
		if !last {
			nextDate := date.NextTradeDate(wdate)
			ndt, _ := internal.ParseTime(nextDate)
			nw := int(ndt.Weekday())
			if nw < w || internal.DifferDays(ndt, dt) >= 7 {
				last = true
			}
		}
		if last {
			df0 := pandas.NewDataFrame(
				pandas.NewSeries(stat.SERIES_TYPE_STRING, "date", wdate),
				pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "open", o),
				pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "close", c),
				pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "high", h),
				pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "low", l),
				pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "volume", v),
				pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "amount", a),
				pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "bv", bv),
				pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "sv", sv),
				pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "ba", ba),
				pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "sa", sa),
				pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "zf", zf),
			)
			df = df.Concat(df0)
			wdate = ""
			prevClose = c
			o = stat.DType(0)
			c = stat.DType(0)
			h = stat.DType(0)
			l = stat.DType(0)
			v = stat.DType(0)
			a = stat.DType(0)
		}
	}
	return df
}
