package date

import (
	"gitee.com/quant1x/data/cache"
	"gitee.com/quant1x/data/internal/js"
	"gitee.com/quant1x/pandas"
	"gitee.com/quant1x/pandas/stat"
	"github.com/mymmsc/gox/api"
	"github.com/mymmsc/gox/http"
	"sort"
	"strings"
	"time"
)

const (
	url_sina_klc_td_sh = "https://finance.sina.com.cn/realstock/company/klc_td_sh.txt"
	kHoliday           = "trade_date"
)

var (
	tradeDates      []string // 交易日列表
	holidayFilename = cache.GetInfoPath() + "/holiday.csv"
)

func init() {
	updateHoliday()
}

// IsHoliday 是否节假日
func IsHoliday(date string) bool {
	iRet, found := sort.Find(len(tradeDates), func(i int) int {
		return strings.Compare(date, tradeDates[i])
	})
	_ = iRet
	return !found
}

func updateHoliday() {
	if !cache.FileExist(holidayFilename) {
		err := cache.CheckFilepath(holidayFilename)
		if err != nil {
			panic("文件路径创建失败: " + holidayFilename)
		}
		data, err := http.HttpGet(url_sina_klc_td_sh)
		if err != nil {
			panic("获取交易日历失败: " + url_sina_klc_td_sh)
		}
		ret, err := js.SinaJsDecode(api.Bytes2String(data))
		if err != nil {
			panic("js解码失败: " + url_sina_klc_td_sh)
		}
		ds := []string{ /*kHoliday*/ }
		for _, v := range ret.([]any) {
			ts := v.(time.Time)
			date := ts.Format(time.DateOnly)
			ds = append(ds, date)
		}
		td := pandas.NewSeries(stat.SERIES_TYPE_STRING, kHoliday, ds)
		//df := pandas.LoadRecords([][]string{ds, ds})
		df := pandas.NewDataFrame(td)
		err = df.WriteCSV(holidayFilename)
		if err != nil {
			return
		}
		tradeDates = ds
	} else {
		df := pandas.ReadCSV(holidayFilename)
		ds := df.Col(kHoliday).Values().([]string)
		tradeDates = ds
	}
}
