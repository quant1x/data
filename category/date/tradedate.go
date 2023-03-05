package date

import (
	"gitee.com/quant1x/data/internal"
	"golang.org/x/exp/slices"
	"sort"
	"time"
)

const (
	kIndexDate = "2006-01-02" // 索引日期格式
)

// IndexToday 当天
func IndexToday() string {
	now := time.Now()
	return now.Format(kIndexDate)
}

// TradeRange 输出交易日范围
func TradeRange(start, end string) []string {
	dt, err := internal.ParseTime(start)
	if err != nil {
		return []string{}
	}
	start = dt.Format(time.DateOnly)

	dt, err = internal.ParseTime(end)
	if err != nil {
		return []string{}
	}
	end = dt.Format(time.DateOnly)

	is := sort.SearchStrings(tradeDates, start)
	ie := sort.SearchStrings(tradeDates, end)
	today := IndexToday()
	lastDay := tradeDates[ie]
	if lastDay > today {
		ie = ie - 1
	}
	return slices.Clone(tradeDates[is : ie+1])
}
