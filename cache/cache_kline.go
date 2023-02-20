package cache

import (
	"fmt"
	"strings"
)

// KLineFilename KLine缓存路径
func KLineFilename(code string) string {
	cacheId := CacheId(code)
	length := len(cacheId)
	tickPath := fmt.Sprintf("%s/%s/%s.csv", GetDayPath(), cacheId[:length-3], cacheId)
	return tickPath
}

// GetKLineFilename 获取缓存的文件名
//
//	Deprecated: 不推荐使用, 建议使用 KLineFilename
func GetKLineFilename(fullCode string) string {
	fullCode = strings.TrimSpace(fullCode)
	if len(fullCode) != 7 && len(fullCode) != 8 {
		return fullCode
	}
	pos := len(fullCode) - 3
	fullCode = strings.ToLower(fullCode)
	// 组织存储路径
	filename := GetDayPath() + "/" + fullCode[0:pos] + "/" + fullCode
	if CACHE_TYPE == CACHE_CSV {
		filename += ".csv"
	}
	err := CheckFilepath(filename)
	if err != nil {
		panic(fmt.Errorf("create file %s, failed", fullCode))
	}
	return filename
}
