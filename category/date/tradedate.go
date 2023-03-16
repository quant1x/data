package date

import (
	"gitee.com/quant1x/data/internal"
	"gitee.com/quant1x/pandas/stat"
	"golang.org/x/exp/slices"
	"sort"
	"time"
)

const (
	kIndexDate = "2006-01-02" // 索引日期格式
)

func fixTradeDate(date string) string {
	dt, err := internal.ParseTime(date)
	if err != nil {
		panic(err)
	}
	return dt.Format(kIndexDate)
}

// IndexToday 当天
func IndexToday() string {
	now := time.Now()
	return now.Format(kIndexDate)
}

// TradeRange 输出交易日范围
func TradeRange(start, end string) []string {
	start = fixTradeDate(start)
	end = fixTradeDate(end)

	is := sort.SearchStrings(gTradeDates, start)
	ie := sort.SearchStrings(gTradeDates, end)
	today := IndexToday()
	lastDay := gTradeDates[ie]
	if lastDay > today {
		ie = ie - 1
	}
	return slices.Clone(gTradeDates[is : ie+1])
}

// LastTradeDate 获得最后一个交易日
func LastTradeDate() string {
	today := IndexToday()
	end := sort.SearchStrings(gTradeDates, today)
	lastDay := gTradeDates[end]
	if lastDay > today {
		end = end - 1
	}
	return gTradeDates[end]
}

// LastNDate 获得指定日期前的N个交易日期数组
func LastNDate(date string, n ...int) []string {
	__opt_end := 0
	if len(n) > 0 {
		__opt_end = n[0]
	}
	r := stat.RangeFinite(-__opt_end)
	date = fixTradeDate(date)
	end := sort.SearchStrings(gTradeDates, date)
	lastDay := gTradeDates[end]
	if lastDay > date {
		end = end - 1
	}
	date_length := len(gTradeDates[0:end])
	s, e, err := r.Limits(date_length)
	if err != nil {
		return nil
	}
	return slices.Clone(gTradeDates[s : e+1])
}
