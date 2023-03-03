package tdx

import (
	"fmt"
	"gitee.com/quant1x/data/cache"
	"gitee.com/quant1x/data/category"
	"gitee.com/quant1x/data/category/date"
	"gitee.com/quant1x/data/internal/dfcf"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/pandas"
	"github.com/mymmsc/gox/logger"
	"time"
)

// RealTime 即时行情数据
func RealTime(code string) {
	marketId, _, code := category.DetectMarket(code)
	tdxApi := prepare()
	hq, err := tdxApi.GetSecurityQuotes([]proto.Market{marketId}, []string{code})
	if err != nil {
		logger.Errorf("获取即时行情数据失败", err)
	}
	fmt.Printf("%+v\n", hq)
}

func BatchRealtime(codes []string) {
	marketIds := []proto.Market{}
	symbols := []string{}

	for _, code := range codes {
		id, _, symbol := category.DetectMarket(code)
		if len(symbol) == 6 {
			marketIds = append(marketIds, id)
			symbols = append(symbols, symbol)
		}
	}
	tdxApi := prepare()
	hq, err := tdxApi.GetSecurityQuotes(marketIds, symbols)
	if err != nil {
		logger.Errorf("获取即时行情数据失败", err)
		return
	}
	//fmt.Printf("%+v\n", hq)
	today := time.Now().Format(category.INDEX_DATE)
	for _, v := range hq.List {
		if v.Code == proto.StockDelisting || v.LastClose == float64(0) {
			continue
		}
		kl := dfcf.KLine{
			Date:   today,
			Open:   v.Open,
			Close:  v.Price,
			High:   v.High,
			Low:    v.Low,
			Volume: int64(v.Vol),
			Amount: v.Amount,
		}
		last := pandas.LoadStructs([]dfcf.KLine{kl})
		fullCode := category.GetMarketName(v.Market) + v.Code
		df := GetCacheKLine(fullCode)
		if df.Nrow() == 0 || last.Nrow() == 0 {
			continue
		}
		lastDay := df.Col("date").IndexOf(-1).(string)
		today := date.IndexToday()
		ts := date.TradeRange(lastDay, today)
		if len(ts) > 2 {
			// 超过2天的差距, 不能用realtime更新K线数据
			continue
		}
		if lastDay == today {
			// 如果最后一条数据和当前日期相同, 那么去掉缓存中的最后一条, 用实时数据填补
			df = df.Subset(0, df.Nrow()-1)
		}
		// 连接缓存和实时数据
		df = df.Concat(last)
		fn := cache.KLineFilename(fullCode)
		err := df.WriteCSV(fn)
		if err != nil {
			logger.Errorf("更新K线数据文件失败:%s", v.Code)
		}
	}
}
