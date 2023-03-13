package cache

import (
	"gitee.com/quant1x/data/internal"
	"time"
)

// LastDate 最后一个有效交易日期
//
//	Deprecated: 不推荐使用, 这个用法不准确
func LastDate() time.Time {
	return internal.CanUpdateTime()
}
