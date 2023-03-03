package tdx

import (
	"gitee.com/quant1x/data/cache"
	"gitee.com/quant1x/data/category"
	"gitee.com/quant1x/pandas"
	"github.com/mymmsc/gox/logger"
)

// GetXdxrInfo 除权除息数据
func GetXdxrInfo(code string) pandas.DataFrame {
	marketId, _, code := category.DetectMarket(code)
	tdxApi := prepare()
	xdxrInfos, err := tdxApi.GetXdxrInfo(marketId, code)
	if err != nil {
		logger.Errorf("获取除权除息数据失败", err)
	}
	//fmt.Printf("%+v\n", hq)
	df := pandas.LoadStructs(xdxrInfos)
	filename := cache.XdxrFilename(code)
	_ = df.WriteCSV(filename)
	return df
}

func GetCacheXdxr(code string) pandas.DataFrame {
	filename := cache.XdxrFilename(code)
	var df pandas.DataFrame
	if !cache.FileExist(filename) {
		// 缓存文件不存在, 下载
		df = GetXdxrInfo(code)
	}
	if !cache.FileExist(filename) {
		return df
	} else {
		df = pandas.ReadCSV(filename)
	}
	return df
}
