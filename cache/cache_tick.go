package cache

import (
	"fmt"
	"gitee.com/quant1x/data/utils"
)

const (
	__date_format = "20060102"
)

// TickFilename tick文件比较多, 目录结构${tick}/${YYYY}/${YYYYMMDD}/${IdPath}
func TickFilename(code string, date string) (string, error) {
	dt, err := utils.ParseTime(date)
	if err != nil {
		return "", err
	}
	date = dt.Format(__date_format)
	cacheId := CacheId(code)
	tickPath := fmt.Sprintf("%s/%s/%s/%s.csv", GetTickPath(), date[0:4], date, cacheId)
	return tickPath, nil
}
