package date

import (
	"gitee.com/quant1x/data/cache"
	"gitee.com/quant1x/data/category"
	"gitee.com/quant1x/data/internal"
	"gitee.com/quant1x/data/internal/dfcf"
	"gitee.com/quant1x/data/util/js"
	"gitee.com/quant1x/data/util/unique"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/pandas"
	"gitee.com/quant1x/pandas/stat"
	"github.com/mymmsc/gox/api"
	"github.com/mymmsc/gox/http"
	"github.com/mymmsc/gox/logger"
	"golang.org/x/exp/slices"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	url_sina_klc_td_sh = "https://finance.sina.com.cn/realstock/company/klc_td_sh.txt"
	kCalendar          = "trade_date"
	kCalendarFormat    = "2006-01-02" // 交易日历日期格式
	kTradeDateFilename = "calendar"
	kIgnoreDate        = "1992-05-04" // TODO:已知缺失的交易日期, 现在已经能自动甄别缺失的交易日期
)

var (
	gTradeDates      []string // 交易日列表
	calendarFilename = cache.GetInfoPath() + "/" + kTradeDateFilename + ".csv"
)

func init() {
	updateCalendar()
	noDates, err := checkCalendar()
	if err == nil && len(noDates) > 0 {
		_ = os.Remove(calendarFilename)
		updateCalendar(noDates...)
	}
}

func getAllDates() []string {
	return slices.Clone(gTradeDates)
}

// IsHoliday 是否节假日
func IsHoliday(date string) bool {
	dates := getAllDates()
	iRet, found := sort.Find(len(dates), func(i int) int {
		return strings.Compare(date, dates[i])
	})
	_ = iRet
	return !found
}

func updateCalendar(noDates ...string) {
	if !cache.FileExist(calendarFilename) {
		err := cache.CheckFilepath(calendarFilename)
		if err != nil {
			panic("文件路径创建失败: " + calendarFilename)
		}
	}
	finfo, err := os.Stat(calendarFilename)
	var fileModTime time.Time
	if err == nil {
		fileModTime = finfo.ModTime()
	}

	header := map[string]any{
		http.IfModifiedSince: fileModTime,
	}
	data, lastModified, err := http.Request(url_sina_klc_td_sh, "get", header)
	if err != nil {
		panic("获取交易日历失败: " + url_sina_klc_td_sh)
	}
	if len(data) == 0 {
		df := pandas.ReadCSV(calendarFilename)
		ds := df.Col(kCalendar).Values().([]string)
		gTradeDates = slices.Clone(ds)
		return
	}
	ret, err := js.SinaJsDecode(api.Bytes2String(data))
	if err != nil {
		panic("js解码失败: " + url_sina_klc_td_sh)
	}
	dates := []string{}
	for _, v := range ret.([]any) {
		ts := v.(time.Time)
		date := ts.Format(kCalendarFormat)
		dates = append(dates, date)
	}
	for _, v := range noDates {
		dates = append(dates, v)
	}
	unique.Sort(unique.StringSlice{&dates})

	td := pandas.NewSeries(stat.SERIES_TYPE_STRING, kCalendar, dates)
	df := pandas.NewDataFrame(td)
	err = df.WriteCSV(calendarFilename)
	if err != nil {
		return
	}
	err = os.Chtimes(calendarFilename, lastModified, lastModified)
	if err != nil {
		logger.Error(err)
	}
	gTradeDates = dates
}

// // 校验缺失的日期, 返回没有的日期列表
func checkCalendar() (noDates []string, err error) {
	dateList := getShangHaiTradeDates()
	// 校验日期的缺失
	start := dateList[0]
	end := dateList[len(dateList)-1]
	dest := TradeRange(start, end)
	noDates = []string{}
	for _, v := range dateList {
		found := slices.Contains(dest, v)
		if !found {
			noDates = append(noDates, v)
		}
	}
	return noDates, nil
}

// 校验缺失的日期, 返回没有的日期列表
func dfcf_checkCalendar() (noDates []string, err error) {
	kls, err := dfcf.A("sh000001")
	if err != nil {
		return nil, err
	}
	df := pandas.LoadStructs(kls)
	if df.Nrow() == 0 {
		return nil, df.Err
	}
	dateList := df.Col("date").Strings()
	// 校验日期的缺失
	start := dateList[0]
	end := dateList[len(dateList)-1]
	dest := TradeRange(start, end)
	noDates = []string{}
	for _, v := range dateList {
		found := slices.Contains(dest, v)
		if !found {
			noDates = append(noDates, v)
		}
	}
	return noDates, nil
}

// 获取上证指数的交易日期, 目的是校验日期
func getShangHaiTradeDates() (dates []string) {
	fullCode := "sh000001"
	tdxApi, err := quotes.NewStdApi()
	if err != nil {
		return nil
	}
	defer tdxApi.Close()
	marketId, _, code := category.DetectMarket(fullCode)
	history := make([]quotes.SecurityBar, 0)
	step := uint16(quotes.TDX_SECURITY_BARS_MAX)
	start := uint16(0)
	hs := make([]quotes.SecurityBarsReply, 0)
	for {
		count := step
		var data *quotes.SecurityBarsReply
		var err error
		data, err = tdxApi.GetIndexBars(marketId, code, proto.KLINE_TYPE_RI_K, start, count)
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
	}
	hs = stat.Reverse(hs)
	for _, v := range hs {
		history = append(history, v.List...)
	}
	dates = []string{}
	for _, v := range history {
		date1 := v.DateTime
		dt, _ := internal.ParseTime(date1)
		date1 = dt.Format(kCalendarFormat)
		dates = append(dates, date1)
	}

	return dates
}
