package tdx

import (
	"gitee.com/quant1x/data/cache"
	"gitee.com/quant1x/data/category"
	"gitee.com/quant1x/pandas"
)

const (
	BlockPath = "/T0002/blocknew"
	ZxgBlk    = "zxg.blk"
	BkltBlk   = "BKLT.blk"
)

func GetZxgList() []string {
	filename := cache.GetZxgFile()
	df := pandas.ReadCSV(filename)
	if df.Nrow() == 0 {
		return []string{}
	}
	rows := df.Col("code")
	if rows.Len() == 0 {
		return []string{}
	}
	// 校验证券代码, 统一格式前缀加代码
	cs := rows.Strings()
	for i, v := range cs {
		_, sid, code := category.DetectMarket(v)
		cs[i] = sid + code
	}
	return cs
}
