package cache

import (
	"fmt"
	"gitee.com/quant1x/data/category"
	"gitee.com/quant1x/data/internal"
)

var (
	// TickStartDate 最早的时间
	TickStartDate = "20220101"
)

// TickFilename tick文件比较多, 目录结构${tick}/${YYYY}/${YYYYMMDD}/${CacheIdPath}
func TickFilename(code string, date string) string {
	return GetTickFilename(code, date, true)
}

func GetTickFilename(code string, date string, createPath bool) string {
	date = CorrectDate(date)
	cacheId := CacheId(code)
	tickPath := fmt.Sprintf("%s/%s/%s/%s.csv", GetTickPath(), date[0:4], date, cacheId)
	if createPath {
		_ = CheckFilepath(tickPath)
	}
	return tickPath
}

// UpdateTickStartDate 修改tick数据开始下载的日期
func UpdateTickStartDate(date string) {
	dt, err := internal.ParseTime(date)
	if err != nil {
		return
	}
	date = dt.Format(category.CACHE_DATE)
	TickStartDate = date
}
