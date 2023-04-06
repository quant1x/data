package date

import (
	"fmt"
	"gitee.com/quant1x/data/cache"
	"gitee.com/quant1x/data/internal/dfcf"
	"gitee.com/quant1x/data/util/js"
	"gitee.com/quant1x/data/util/unique"
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
	kTradeDateFilename = "calendar"
	kIgnoreDate        = "1992-05-04"
)

var (
	gTradeDates      []string // 交易日列表
	calendarFilename = cache.GetInfoPath() + "/" + kTradeDateFilename + ".csv"
)

func init() {
	updateCalendar()
}

// IsHoliday 是否节假日
func IsHoliday(date string) bool {
	iRet, found := sort.Find(len(gTradeDates), func(i int) int {
		return strings.Compare(date, gTradeDates[i])
	})
	_ = iRet
	return !found
}

func updateCalendar() {
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
		date := ts.Format(time.DateOnly)
		dates = append(dates, date)
	}
	dates = append(dates, kIgnoreDate)
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

func checkCalendar() (dates []string, err error) {
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
	start := "1990-12-19"
	end := "2023-04-05"
	dest := TradeRange(start, end)
	fmt.Println(len(dest))
	for i, v := range dates {
		found := slices.Contains(dest, v)
		if !found {
			fmt.Println(v)
			tmp := df.IndexOf(i)
			fmt.Println(tmp)
		}
	}
	return dateList, nil
}
