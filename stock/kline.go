package stock

import (
	"gitee.com/quant1x/data/internal/tdx"
	"gitee.com/quant1x/pandas"
)

// KLine 加载K线
func KLine(code string) pandas.DataFrame {
	return tdx.GetCacheKLine(code, true)
}
