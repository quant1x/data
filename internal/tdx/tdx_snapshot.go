package tdx

import (
	"gitee.com/quant1x/data/cache"
	"gitee.com/quant1x/data/category"
	"gitee.com/quant1x/data/category/date"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/pandas"
	"github.com/mymmsc/gox/api"
	"github.com/mymmsc/gox/logger"
	"time"
)

type QuoteSnapshot struct {
	Market         uint8   // 市场
	Code           string  // 代码
	Active1        uint16  // 活跃度
	Price          float64 // 现价
	LastClose      float64 // 昨收
	Open           float64 // 开盘
	High           float64 // 最高
	Low            float64 // 最低
	ServerTime     string  // 时间
	ReversedBytes0 int     // 保留(时间 ServerTime)
	ReversedBytes1 int     // 保留
	Vol            int     // 总量
	CurVol         int     // 现量
	Amount         float64 // 总金额
	SVol           int     // 内盘
	BVol           int     // 外盘
	ReversedBytes2 int     // 保留
	ReversedBytes3 int     // 保留
	//BidLevels      []quotes.Level
	//AskLevels      []quotes.Level
	Bid1           float64
	Ask1           float64
	BidVol1        int
	AskVol1        int
	Bid2           float64
	Ask2           float64
	BidVol2        int
	AskVol2        int
	Bid3           float64
	Ask3           float64
	BidVol3        int
	AskVol3        int
	Bid4           float64
	Ask4           float64
	BidVol4        int
	AskVol4        int
	Bid5           float64
	Ask5           float64
	BidVol5        int
	AskVol5        int
	ReversedBytes4 uint16  // 保留
	ReversedBytes5 int     // 保留
	ReversedBytes6 int     // 保留
	ReversedBytes7 int     // 保留
	ReversedBytes8 int     // 保留
	Rate           float64 // 涨速
	Active2        uint16  // 活跃度
}

// BatchSnapShot 批量获取即时行情数据快照
func BatchSnapShot(codes []string) []QuoteSnapshot {
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
	list := []QuoteSnapshot{}
	hq, err := tdxApi.GetSecurityQuotes(marketIds, symbols)
	if err != nil {
		logger.Errorf("获取即时行情数据失败", err)
		return list
	}
	//fmt.Printf("%+v\n", hq)
	lastTradeday := time.Now().Format(category.INDEX_DATE)
	td := date.TradeRange("2023-01-01", lastTradeday)
	lastTradeday = td[len(td)-1]
	for _, v := range hq.List {
		snapshot := QuoteSnapshot{}
		_ = api.Copy(&snapshot, &v)
		if snapshot.Code == proto.StockDelisting || snapshot.LastClose == float64(0) {
			continue
		}
		last := pandas.LoadStructs([]QuoteSnapshot{snapshot})
		if last.Nrow() == 0 {
			continue
		}
		fullCode := category.GetMarketName(v.Market) + v.Code
		fn := cache.SnapshotFilename(fullCode)
		err := last.WriteCSV(fn)
		if err != nil {
			logger.Errorf("更新快照数据文件失败:%s", v.Code)
		}
		list = append(list, snapshot)
	}
	return list
}
