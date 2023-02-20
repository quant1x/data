package tdx

import (
	"gitee.com/quant1x/data/category"
	"time"
)

const (
	TDX_DATE = "20060102" // 通达信日期格式
)

func tdxToday() string {
	now := time.Now()
	return now.Format(category.TDX_DATE)
}
