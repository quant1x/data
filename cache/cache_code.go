package cache

import (
	"fmt"
	"gitee.com/quant1x/data/security"
)

// CacheId 通过代码构建目录结构
func CacheId(code string) string {
	marketId, marketName, code := security.DetectMarket(code)
	cacheId := fmt.Sprintf("%s%s", marketName, code)

	_ = marketId
	return cacheId
}

// IdPath code从后保留3位, 市场缩写+从头到倒数第3的代码, 确保每个木有只有000~999个代码
func IdPath(code string) string {
	N := 3
	cacheId := CacheId(code)
	length := len(cacheId)

	prefix := cacheId[:length-N]
	return fmt.Sprintf("%s/%s", prefix, cacheId)
}
