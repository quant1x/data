package cache

import (
	"gitee.com/quant1x/data/internal"
	"time"
)

// LastDate 最后一个有效交易日期
func LastDate() time.Time {
	return internal.CanUpdateTime()
}
