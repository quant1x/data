package tdx

import (
	"gitee.com/quant1x/data/cache"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/pandas/utils"
	"os"
)

func getBlockFile(filename string) {
	api := prepare()
	fn := cache.GetBkPath() + "/" + filename
	stat, err := os.Stat(fn)
	if err == nil || os.IsExist(err) {
		today, _ := utils.ParseTime(cache.Today())
		if stat.ModTime().After(today) {
			return
		}
	}
	resp, err := api.GetBlockInfo(filename)
	if err == nil {
		fn := cache.GetBkPath() + "/" + filename
		fp, err := os.Create(fn)
		if err == nil {
			_, _ = fp.Write(resp.Data)
			_ = fp.Close()
		}
	}
}

func init() {
	getBlockFile(quotes.BLOCK_DEFAULT)
	getBlockFile(quotes.BLOCK_GAINIAN)
	getBlockFile(quotes.BLOCK_FENGGE)
	getBlockFile(quotes.BLOCK_ZHISHU)
}
