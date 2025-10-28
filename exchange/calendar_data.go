package exchange

import (
	"strings"
	"time"

	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/http"
	"gitee.com/quant1x/gox/logger"
)

func downloadCalendar(fileModTime time.Time) ([]calendar, time.Time) {
	header := map[string]any{
		http.IfModifiedSince: fileModTime,
	}
	data, lastModified, err := http.Request(urlSinaRealstockCompanyKlcTdSh, http.MethodGet, "", header)
	if err != nil {
		logger.Fatal("获取交易日历失败: " + urlSinaRealstockCompanyKlcTdSh)
	}
	if len(data) == 0 {
		return nil, lastModified
	}
	text := api.Bytes2String(data)
	tmp := strings.Split(text, "=")
	if len(tmp) > 1 {
		text = tmp[1]
	}
	text = strings.Split(text, ";")[0]
	text = strings.ReplaceAll(text, "\"", "")

	decoder := NewCalendarDecoder(text)
	out := decoder.Decode()
	// if out == nil {
	// 	logger.Fatal("js解码失败: " + urlSinaRealstockCompanyKlcTdSh)
	// }
	var dates []calendar
	for _, date := range out.([]string) {
		e := calendar{
			Date:   date,
			Source: "sina",
		}
		dates = append(dates, e)
	}
	return dates, lastModified
}
