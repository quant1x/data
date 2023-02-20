package tdx

import (
	"fmt"
	"gitee.com/quant1x/data/cache"
	"gitee.com/quant1x/data/category"
	"gitee.com/quant1x/data/dfcf"
	"gitee.com/quant1x/data/security"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/pandas"
	"github.com/mymmsc/gox/logger"
	"time"
)

// RealTime 即时行情数据
func RealTime(code string) {
	marketId, _, code := security.DetectMarket(code)
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
		id, _, symbol := security.DetectMarket(code)
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
		//fmt.Printf("%+v\n", kl)
		last := pandas.LoadStructs([]dfcf.KLine{kl})
		df := __kLine(v.Code)
		df = df.Subset(0, df.Nrow()-1)
		df = df.Concat(last)
		fn := cache.KLineFilename(v.Code)
		err := df.WriteCSV(fn)
		if err != nil {
			logger.Errorf("更新K线数据文件失败:%s", v.Code)
		}
	}
}
