package cache

import (
	"fmt"
)

// TickFilename tick文件比较多, 目录结构${tick}/${YYYY}/${YYYYMMDD}/${CacheIdPath}
func TickFilename(code string, date string) string {
	date = CorrectDate(date)
	cacheId := CacheId(code)
	tickPath := fmt.Sprintf("%s/%s/%s/%s.csv", GetTickPath(), date[0:4], date, cacheId)
	return tickPath
}
