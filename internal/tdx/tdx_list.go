package tdx

import (
	"gitee.com/quant1x/data/cache"
	"gitee.com/quant1x/data/category"
	"gitee.com/quant1x/data/category/date"
	"gitee.com/quant1x/data/internal"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/pandas"
	"os"
	"time"
)

const (
	fn_stock_list = "stocks.csv"
)

var (
	mapStock = map[string]string{}
)

func init() {
	fnList := cache.GetInfoPath() + "/" + fn_stock_list
	var df pandas.DataFrame
	bCreate := false
	if !cache.FileExist(fnList) {
		bCreate = true
	} else {
		// 获取文件创建时间
		finfo, _ := os.Stat(fnList)
		ftime := internal.DateZero(finfo.ModTime())
		fdate := ftime.Format(time.DateOnly)
		lastDay := date.LastTradeDate()
		if fdate < lastDay {
			bCreate = true
		}
	}
	if bCreate {
		df = get_stock_list()
		_ = df.WriteCSV(fnList)
	}
	df = pandas.ReadCSV(fnList)
	if df.Nrow() == 0 {
		return
	}
	for i := 0; i < df.Nrow(); i++ {
		m := df.IndexOf(i)
		code := m["Code"].(string)
		name := m["Name"].(string)
		mapStock[code] = name
	}
}

// GetStockName 获取证券名称
func GetStockName(code string) (string, bool) {
	_, mname, mcode := category.DetectMarket(code)
	code = mname + mcode
	name, ok := mapStock[code]
	return name, ok
}

// 证券列表
func get_stock_list() pandas.DataFrame {
	stdApi := prepare()
	offset := uint16(quotes.TDX_SECURITY_LIST_MAX)
	start := uint16(0)
	var df pandas.DataFrame
	for {
		reply, err := stdApi.GetSecurityList(proto.MarketShangHai, start)
		if err != nil {
			return df
		}
		for i := 0; i < int(reply.Count); i++ {
			reply.List[i].Code = "sh" + reply.List[i].Code
		}
		tmp := pandas.LoadStructs(reply.List)
		df = df.Concat(tmp)
		if reply.Count < offset {
			break
		}
		start += reply.Count
	}
	start = uint16(0)
	for {
		reply, err := stdApi.GetSecurityList(proto.MarketShenZhen, start)
		if err != nil {
			return df
		}
		for i := 0; i < int(reply.Count); i++ {
			reply.List[i].Name = "sz" + reply.List[i].Name
		}
		tmp := pandas.LoadStructs(reply.List)
		df = df.Concat(tmp)
		if reply.Count < offset {
			break
		}
		start += reply.Count
	}
	return df
}
