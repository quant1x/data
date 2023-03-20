package date

import (
	"gitee.com/quant1x/data/cache"
	"gitee.com/quant1x/data/internal/js"
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
	ds := []string{ /*kCalendar*/ }
	for _, v := range ret.([]any) {
		ts := v.(time.Time)
		date := ts.Format(time.DateOnly)
		ds = append(ds, date)
	}
	td := pandas.NewSeries(stat.SERIES_TYPE_STRING, kCalendar, ds)
	df := pandas.NewDataFrame(td)
	err = df.WriteCSV(calendarFilename)
	if err != nil {
		return
	}
	err = os.Chtimes(calendarFilename, lastModified, lastModified)
	if err != nil {
		logger.Error(err)
	}
	gTradeDates = ds

}
