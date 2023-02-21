package tdx

import (
	"fmt"
	"gitee.com/quant1x/data/cache"
	"gitee.com/quant1x/data/category"
	"gitee.com/quant1x/data/category/date"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/pandas"
	"gitee.com/quant1x/pandas/stat"
	"github.com/mymmsc/gox/progressbar"
	"strconv"
	"strings"
)

var (
	stdApi *quotes.StdApi = nil
)

func prepare() *quotes.StdApi {
	if stdApi == nil {
		std_api, err := quotes.NewStdApi()
		if err != nil {
			return nil
		}
		stdApi = std_api
	}
	return stdApi
}

func startsWith(str string, prefixs []string) bool {
	if len(str) == 0 || len(prefixs) == 0 {
		return false
	}
	for _, prefix := range prefixs {
		if strings.HasPrefix(str, prefix) {
			return true
		}
	}
	return false
}

// 判断股票ID对应的证券市场匹配规则
//
// ['50', '51', '60', '90', '110'] 为 sh
// ['00', '12'，'13', '18', '15', '16', '18', '20', '30', '39', '115'] 为 sz
// ['5', '6', '9'] 开头的为 sh， 其余为 sz
func getStockMarket(symbol string) string {
	//:param string: False 返回市场ID，否则市场缩写名称
	//:param symbol: 股票ID, 若以 'sz', 'sh' 开头直接返回对应类型，否则使用内置规则判断
	//:return 'sh' or 'sz'

	market := "sh"
	if startsWith(symbol, []string{"sh", "sz", "SH", "SZ"}) {
		market = strings.ToLower(symbol[0:2])
	} else if startsWith(symbol, []string{"50", "51", "60", "68", "90", "110", "113", "132", "204"}) {
		market = "sh"
	} else if startsWith(symbol, []string{"00", "12", "13", "18", "15", "16", "18", "20", "30", "39", "115", "1318"}) {
		market = "sz"
	} else if startsWith(symbol, []string{"5", "6", "9", "7"}) {
		market = "sh"
	} else if startsWith(symbol, []string{"4", "8"}) {
		market = "bj"
	}
	return market
}

func getStockMarketId(symbol string) uint8 {
	market := getStockMarket(symbol)
	marketId := proto.MarketShangHai
	if market == "sh" {
		marketId = proto.MarketShangHai
	} else if market == "sz" {
		marketId = proto.MarketShenZhen
	} else if market == "bj" {
		marketId = proto.MarketBeiJing
	}
	//# logger.debug(f"market => {market}")

	return marketId
}

func GetTickAll(code string) {
	api := prepare()
	marketId, _, code := category.DetectMarket(code)
	info, err := api.GetFinanceInfo(marketId, code, 1)
	if err != nil {
		return
	}
	tStart := strconv.FormatInt(int64(info.IPODate), 10)
	tEnd := "20500101"
	dateRange := date.TradeRange(tStart, tEnd)
	bar := progressbar.NewBar(2, fmt.Sprintf("同步[%s]", code), len(dateRange))
	for _, tradeDate := range dateRange {
		bar.Add(1)
		//logger.Infof("同步[%s] %s tick...", code, tradeDate)
		fname := cache.TickFilename(code, tradeDate)
		if cache.FileExist(fname) {
			// 如果已经存在就跳过
			continue
		}
		df := GetTickData(code, tradeDate)
		_ = df
		//logger.Infof("同步[%s] %s tick...OK", code, tradeDate)
	}

	return
}

func GetTickData(code string, date string) pandas.DataFrame {
	api := prepare()
	marketId, marketName, code := category.DetectMarket(code)
	offset := uint16(1800)
	start := uint16(0)
	count := offset
	date = cache.CorrectDate(date)
	history := make([]quotes.HistoryTransaction, 0)
	hs := make([]quotes.HistoryTransactionReply, 0)
	for {
		iDate := stat.AnyToInt64(date)
		data, err := api.GetHistoryTransactionData(marketId, code, uint32(iDate), start, offset)
		if err != nil {
			panic("接口异常")
		}
		hs = append(hs, (*data))
		if data.Count < count {
			// 已经是最早的记录
			// 需要排序
			break
		}
		start += offset
	}
	hs = stat.Reverse(hs)
	for _, v := range hs {
		history = append(history, v.List...)
	}
	_ = marketName
	df := pandas.LoadStructs(history)
	df = df.Select([]string{"Time", "Price", "Vol", "Num", "BuyOrSell"})
	err := df.SetNames("time", "price", "vol", "num", "buyorsell")
	if err != nil {
		return pandas.DataFrame{}
	}
	tickFile := cache.TickFilename(code, date)
	if err != nil {
		return pandas.DataFrame{}
	}
	err = cache.CheckFilepath(tickFile)
	if err != nil {
		return pandas.DataFrame{}
	}
	err = df.WriteCSV(tickFile)
	if err != nil {
		return pandas.DataFrame{}
	}

	return df
}
