package tdx

import (
	"fmt"
	"gitee.com/quant1x/data/cache"
	"gitee.com/quant1x/data/category"
	"gitee.com/quant1x/data/category/date"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/pandas"
	"github.com/mymmsc/gox/logger"
	"time"
)

// RTSecurityBar K线数据
type RTSecurityBar struct {
	Date      string
	Open      float64
	Close     float64
	High      float64
	Low       float64
	Volume    float64
	Amount    float64
	UpCount   int     // 指数有效, 上涨家数
	DownCount int     // 指数有效, 下跌家数
	BuyVol    float64 // 外盘
	SellVol   float64 // 内盘
	BuyAmt    float64 // 外盘成交金额
	SellAmt   float64 // 内盘成交金额
	//Lb        float64 // 量比
	//Zf        float64 // 涨幅
}

var (
	//RTBarsRaw         = []string{"Date", "Open", "Close", "High", "Low", "Volume", "Amount", "UpCount", "DownCount"}
	RTBarsRename      = []string{"date", "open", "close", "high", "low", "volume", "amount", "up", "down", "bv", "sv", "ba", "sa"}
	RTBarsStockFields = []string{"date", "open", "close", "high", "low", "volume", "amount", "bv", "sv", "ba", "sa"}
	RTBarsIndexFields = []string{"date", "open", "close", "high", "low", "volume", "amount", "up", "down", "bv", "sv", "ba", "sa"}
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

// BatchRealtime 批量获取实时行情数据
func BatchRealtime(codes []string) {
	if len(codes) > int(quotes.TDX_SECURITY_QUOTES_MAX) {
		panic(fmt.Sprintf("BatchRealtime: codes count > %d", quotes.TDX_SECURITY_QUOTES_MAX))
	}
	now := time.Now()
	lastTradeday := now.Format(category.INDEX_DATE)
	nowServerTime := now.Format(date.CN_SERVERTIME_FORMAT)
	td := date.TradeRange("2023-01-01", lastTradeday)
	lastTradeday = td[len(td)-1]
	today := date.Today()
	if lastTradeday != today {
		// 当天非交易日, 不更新, 直接返回
		return
	}
	if nowServerTime < date.CN_StartTime || nowServerTime > date.CN_StopTime {
		// 非交易时间, 不更新, 直接返回
		return
	}

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
	for _, v := range hq.List {
		if v.Code == proto.StockDelisting || v.LastClose == float64(0) {
			continue
		}
		kl := RTSecurityBar{
			Date:      lastTradeday,
			Open:      v.Open,
			Close:     v.Price,
			High:      v.High,
			Low:       v.Low,
			Volume:    float64(v.Vol),
			Amount:    v.Amount,
			UpCount:   v.BidVol1,
			DownCount: v.AskVol1,
			BuyVol:    float64(v.BVol),
			SellVol:   float64(v.SVol),
		}
		last := pandas.LoadStructs([]RTSecurityBar{kl})
		fullCode := category.GetMarketName(v.Market) + v.Code
		isIndex := category.IndexFromMarketAndCode(v.Market, v.Code)
		newFields := RTBarsRename
		_ = last.SetNames(newFields...)
		fields := RTBarsStockFields
		if isIndex {
			fields = RTBarsIndexFields
		}
		df := GetCacheKLine(fullCode)
		if df.Nrow() == 0 || last.Nrow() == 0 {
			continue
		}
		df = df.Select(fields)
		lastDay := df.Col("date").IndexOf(-1).(string)
		ts := date.TradeRange(lastDay, lastTradeday)
		if len(ts) > 2 {
			// 超过2天的差距, 不能用realtime更新K线数据
			// 只能是当天更新 或者是新增, 跨越2个以上的交易日不更新
			continue
		}
		// 数据差异数
		diffDays := 0
		// 当日的K线数据已经存在
		if lastDay == lastTradeday {
			// 如果最后一条数据和最后一个交易日相同, 那么去掉缓存中的最后一条, 用实时数据填补
			// 这种情况的出现是K线被更新过了, 现在做的是用快照更新K线
			diffDays = 1
		} else if nowServerTime > v.ServerTime {
			diffDays = 0
		}
		if diffDays > 0 {
			df = df.Subset(0, df.Nrow()-diffDays)
		}
		// 连接缓存和实时数据
		tmp := last.Select(fields)
		df = df.Concat(tmp)
		fn := cache.KLineFilename(fullCode)
		err := df.WriteCSV(fn)
		if err != nil {
			logger.Errorf("更新K线数据文件失败:%s", v.Code)
		}
	}
}
