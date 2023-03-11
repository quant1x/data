package stock

import (
	"fmt"
	"gitee.com/quant1x/pandas/stat"
)

var (
// mapIndexes = map[string]bool{}
)

// GetCodeList 加载全部股票代码
func GetCodeList() []string {
	fullCodes := make([]string, 0)
	// 指数
	indexes := []string{"sh000001",
		"sh000905", "sz399001", "sz399006"}
	fullCodes = append(fullCodes, indexes...)
	//for _, v := range indexes {
	//	mapIndexes[v] = true
	//}

	// 板块信息
	df := BlockList()
	if df.Nrow() > 0 {
		for i := 0; i < df.Nrow(); i++ {
			m := df.IndexOf(i)
			code := stat.AnyToString(m["code"])
			name := stat.AnyToString(m["name"])
			if len(code) == 0 || len(name) == 0 {
				continue
			}
			code = "sh" + code
			fullCodes = append(fullCodes, code)
			//mapIndexes[code] = true
		}
	}

	// 更新代码
	// 上海
	// sh600000-sh600999
	{
		var (
			codeBegin = 600000
			codeEnd   = 600999
		)
		for i := codeBegin; i <= codeEnd; i++ {
			fc := fmt.Sprintf("sh%d", i)
			fullCodes = append(fullCodes, fc)
		}
	}
	// sh601000-sh601999
	{
		var (
			codeBegin = 601000
			codeEnd   = 601999
		)
		for i := codeBegin; i <= codeEnd; i++ {
			fc := fmt.Sprintf("sh%d", i)
			fullCodes = append(fullCodes, fc)
		}
	}
	// sh603000-sh603999
	{
		var (
			codeBegin = 603000
			codeEnd   = 603999
		)
		for i := codeBegin; i <= codeEnd; i++ {
			fc := fmt.Sprintf("sh%d", i)
			fullCodes = append(fullCodes, fc)
		}
	}
	// sh688000-sh688999
	{
		var (
			codeBegin = 688000
			codeEnd   = 688999
		)
		for i := codeBegin; i <= codeEnd; i++ {
			fc := fmt.Sprintf("sh%d", i)
			fullCodes = append(fullCodes, fc)
		}
	}
	// 深圳证券交易所
	// 深圳主板: sz000000-sz000999
	{
		var (
			codeBegin = 0
			codeEnd   = 999
		)
		for i := codeBegin; i <= codeEnd; i++ {
			fc := fmt.Sprintf("sz000%03d", i)
			fullCodes = append(fullCodes, fc)
		}
	}
	// 中小板: sz002000-sz002999
	{
		var (
			codeBegin = 2000
			codeEnd   = 2999
		)
		for i := codeBegin; i <= codeEnd; i++ {
			fc := fmt.Sprintf("sz00%04d", i)
			fullCodes = append(fullCodes, fc)
		}
	}
	// 创业板: sz300000-sz300999
	{
		var (
			codeBegin = 300000
			codeEnd   = 300999
		)
		for i := codeBegin; i <= codeEnd; i++ {
			fc := fmt.Sprintf("sz%06d", i)
			fullCodes = append(fullCodes, fc)
		}
	}
	//fullCodes = fullCodes[0:0]
	// 港股: hk00001-hk09999
	//{
	//	var (
	//		codeBegin = 1
	//		codeEnd   = 9999
	//	)
	//	for i := codeBegin; i <= codeEnd; i++ {
	//		fc := fmt.Sprintf("hk%05d", i)
	//		fullCodes = append(fullCodes, fc)
	//	}
	//}

	return fullCodes
}
