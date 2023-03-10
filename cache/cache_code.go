package cache

import (
	"fmt"
	"gitee.com/quant1x/data/category"
)

// CacheId 通过代码构建目录结构
func CacheId(code string) string {
	marketId, marketName, code := category.DetectMarket(code)
	cacheId := fmt.Sprintf("%s%s", marketName, code)

	_ = marketId
	return cacheId
}

// CacheIdPath code从后保留3位, 市场缩写+从头到倒数第3的代码, 确保每个目录只有000~999个代码
func CacheIdPath(code string) string {
	N := 3
	cacheId := CacheId(code)
	length := len(cacheId)

	prefix := cacheId[:length-N]
	return fmt.Sprintf("%s/%s", prefix, cacheId)
}
