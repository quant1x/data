package exchange

import (
	"strings"

	"gitee.com/quant1x/gox/api"
)

type MarketType = uint8

const (
	MarketIdShenZhen MarketType = iota // 深圳
	MarketIdShangHai MarketType = 1    // 上海
	MarketIdBeiJing  MarketType = 2    // 北京
	MarketIdHongKong MarketType = 21   // 香港
	MarketIdUSA      MarketType = 22   // 美国

	StockDelisting = "DELISTING" // 退市
)

const (
	MarketShangHai string = "sh" // 上海
	MarketShenZhen string = "sz" // 深圳
	MarketBeiJing  string = "bj" // 北京
	MarketHongKong string = "hk" // 香港
	MarketUSA      string = "us" // 美国
)

const (
	MARKET_CN_FIRST_DATE     = "19901219"   // 上证指数的第一个交易日
	MARKET_CH_FIRST_LISTTIME = "1990-12-19" // 个股上市日期
)

var (
	// 全部市场简写
	marketFlags = []string{"sh", "sz", "SH", "SZ", "bj", "BJ", "hk", "HK", "us", "US"}
	// A股市场代码
	marketAShareFlags = []string{"sh", "sz", "SH", "SZ", "bj", "BJ"}
)

var (
	// 上海证券交易所代码段
	//	主板: 60xxxx
	//	科创板: 688xxx
	//	B股: 900xxx
	//	优先股: 360xxx
	//	科创板存托凭证: 689xxx
	//	申购/配股/投票: 7xxxxx
	//	上海总规则: http://www.sse.com.cn/lawandrules/guide/stock/jyglywznylc/zn/a/20230209/4ae280c58535e0424b3a9c743c47e6b9.docx
	//	0: 国债/指数, 000 上证指数系列和中证指数系列, 00068x科创板指数
	//	1: 债券
	//	2: 回购
	//	3: 期货
	//	4: 备用
	//	5: 基金/权证
	//	6: A股
	//	7: 非交易业务(发行, 权益分配)
	//	8: 备用, 通达信编制板块指数占用880,881
	//	9: B股
	shanghaiPrefixes = []string{"50", "51", "60", "68", "689", "90", "110", "113", "132", "204", "000", "360", "880", "881", "7", "5", "6", "9"}

	// 上海证券交易所特殊代码段
	shanghaiSpecialPrefixes = []string{"5", "6", "9", "7"}

	// 深圳交易所代码段
	//	主板: 000,001
	//	中小板: 002,003,004
	//	创业板: 30xxxx
	//	优先股: 140xxx
	//	深圳总规则: https://zhuanlan.zhihu.com/p/63064991
	//	0: 股票
	//	1: 国债/基金
	//	2: B股
	//	30: 创业板
	//	36: 投票, 369999用于深交所认证业务的密码激活/密码挂失
	//	37: 增发/可转债申购
	//	38: 配股/可转债优先权
	//	395: 成家量统计指数
	//	399: 指数
	shenzhenPrefixes = []string{"00", "001", "002", "003", "004", "12", "13", "15", "16", "18", "20", "30", "36", "37", "38", "39", "115", "1318", "140", "395", "399", "159"}

	// 北京交易所证券代码段
	//	北交所指数: 899
	//	新三板: 40,43,83,87
	//	88开头: 通常表示公开发行的股票, 与新三板市场中的其他类型股票进行区分
	//	三板A: 400,430,830-839,870-873
	//	三板B: 420
	//	优先股: 820
	//	新代码段: 920
	beijingPrefixes = []string{"40", "43", "83", "87", "88", "420", "820", "899", "920"}
)

var (
	// 上海交易所指数前缀
	shanghaiIndexPrefixes = []string{"000"}
	// 深圳交易所指数前缀
	shenzhenIndexPrefixes = []string{"399"}
	// 北京交易所指数前缀
	beijingIndexPrefixes = []string{"899"}
)

var (
	// 通达信板块指数, 在上海交易所
	sectorPrefixes = []string{"880", "881"}
)

func GetSecurityCode(market MarketType, symbol string) (securityCode string) {
	switch market {
	case MarketIdUSA:
		return MarketUSA + symbol
	case MarketIdHongKong:
		return MarketHongKong + symbol[:5]
	case MarketIdBeiJing:
		return MarketBeiJing + symbol[:6]
	case MarketIdShenZhen:
		return MarketShenZhen + symbol[:6]
	default:
		return MarketShangHai + symbol[:6]
	}
}

// GetMarket 判断股票ID对应的证券市场匹配规则
//
//	['50', '51', '60', '90', '110'] 为 sh
//	['00', '12'，'13', '18', '15', '16', '18', '20', '30', '39', '115'] 为 sz
//	['5', '6', '9'] 开头的为 sh， 其余为 sz
func GetMarket(symbol string) string {
	symbol = strings.TrimSpace(symbol)
	market := "sh"
	if api.StartsWith(symbol, marketFlags) {
		market = strings.ToLower(symbol[0:2])
	} else if api.EndsWith(symbol, marketFlags) {
		length := len(symbol)
		// 后缀一个点号+2位字母, 代码在最前面
		market = strings.ToLower(symbol[length-2:])
	} else if api.StartsWith(symbol, shanghaiPrefixes) {
		market = "sh"
	} else if api.StartsWith(symbol, shenzhenPrefixes) {
		market = "sz"
	} else if api.StartsWith(symbol, sectorPrefixes) {
		market = "sh"
	} else if api.StartsWith(symbol, beijingPrefixes) {
		market = "bj"
	}
	return market
}

// GetMarketId 获得市场ID
func GetMarketId(symbol string) uint8 {
	market := GetMarket(symbol)
	marketId := MarketIdShangHai
	if market == "sh" {
		marketId = MarketIdShangHai
	} else if market == "sz" {
		marketId = MarketIdShenZhen
	} else if market == "bj" {
		marketId = MarketIdBeiJing
	}
	return marketId
}

func GetMarketFlag(marketId MarketType) string {
	switch marketId {
	case MarketIdShenZhen:
		return MarketShenZhen
	case MarketIdBeiJing:
		return MarketBeiJing
	case MarketIdHongKong:
		return MarketHongKong
	case MarketIdUSA:
		return MarketUSA
	default:
		return MarketShangHai
	}
}

// DetectMarket 检测市场代码
func DetectMarket(symbol string) (marketId MarketType, market string, code string) {
	code = strings.TrimSpace(symbol)
	market = MarketShangHai
	if api.StartsWith(code, marketFlags) {
		// 前缀2位字母后面跟代码
		market = strings.ToLower(code[0:2])
		if code[2:3] == "." {
			// SZ.000002
			code = code[3:]
		} else {
			// SZ000002
			code = code[2:]
		}
	} else if api.EndsWith(code, marketFlags) {
		length := len(code)
		// 后缀一个点号+2位字母, 代码在最前面
		// 600000.SH
		market = strings.ToLower(code[length-2:])
		code = code[:length-3]
	} else if api.StartsWith(code, shanghaiPrefixes) {
		market = MarketShangHai
	} else if api.StartsWith(code, shenzhenPrefixes) {
		market = MarketShenZhen
	} else if api.StartsWith(code, shanghaiSpecialPrefixes) {
		market = MarketShangHai
	} else if api.StartsWith(code, sectorPrefixes) {
		// 通达信板块指数, 在上海交易所
		market = MarketShangHai
	} else if api.StartsWith(code, beijingPrefixes) {
		market = MarketBeiJing
	}
	marketId = MarketIdShangHai
	if market == MarketShangHai {
		marketId = MarketIdShangHai
	} else if market == MarketShenZhen {
		marketId = MarketIdShenZhen
	} else if market == MarketBeiJing {
		marketId = MarketIdBeiJing
	} else if market == MarketHongKong {
		marketId = MarketIdHongKong
	}
	return marketId, market, code
}

// AssertIndexByMarketAndCode 通过市场id和短码判断是否指数
func AssertIndexByMarketAndCode(marketId MarketType, symbol string) (isIndex bool) {
	if marketId == MarketIdShangHai && api.StartsWith(symbol, shanghaiIndexPrefixes) {
		return true
	} else if marketId == MarketIdShangHai && api.StartsWith(symbol, sectorPrefixes) {
		return true
	} else if marketId == MarketIdShenZhen && api.StartsWith(symbol, shenzhenIndexPrefixes) {
		return true
	} else if marketId == MarketIdBeiJing && api.StartsWith(symbol, beijingIndexPrefixes) {
		return true
	}
	return false
}

// AssertIndexBySecurityCode 通过证券代码判断是否指数
func AssertIndexBySecurityCode(securityCode string) (isIndex bool) {
	marketId, _, code := DetectMarket(securityCode)
	return AssertIndexByMarketAndCode(marketId, code)
}

// AssertBlockBySecurityCode 断言证券代码是否板块
func AssertBlockBySecurityCode(securityCode *string) (isBlock bool) {
	marketId, flag, code := DetectMarket(*securityCode)
	if marketId != MarketIdShangHai {
		// 板块指数的数据在上海
		return false
	}
	if !api.StartsWith(code, sectorPrefixes) {
		return false
	}
	*securityCode = flag + code
	return true
}

// AssertETFByMarketAndCode 通过市场id和代码判断是否ETF
func AssertETFByMarketAndCode(marketId MarketType, symbol string) (isETF bool) {
	if marketId == MarketIdShangHai && api.StartsWith(symbol, []string{"51"}) {
		return true
	} else if marketId == MarketIdShenZhen && api.StartsWith(symbol, []string{"159"}) {
		return true
	}
	return false
}

// AssertStockByMarketAndCode 通过市场id和代码判断是否个股
func AssertStockByMarketAndCode(marketId MarketType, symbol string) (isStock bool) {
	if marketId == MarketIdShangHai && api.StartsWith(symbol, []string{"60", "68", "51"}) {
		return true
	} else if marketId == MarketIdShenZhen && api.StartsWith(symbol, []string{"00", "30", "159"}) {
		return true
	} else if marketId == MarketIdBeiJing && api.StartsWith(symbol, []string{"40", "43", "83", "87", "92"}) {
		return true
	}
	return false
}

// AssertStockBySecurityCode 通过证券代码判断是否个股
func AssertStockBySecurityCode(securityCode string) (isStock bool) {
	marketId, _, code := DetectMarket(securityCode)
	return AssertStockByMarketAndCode(marketId, code)
}

// CorrectSecurityCode 修正证券代码
func CorrectSecurityCode(securityCode string) string {
	if len(securityCode) == 0 {
		return ""
	}
	_, mFlag, mSymbol := DetectMarket(securityCode)
	return mFlag + mSymbol
}

// TargetKind 标的类型
type TargetKind int

const (
	STOCK = iota // 股票
	INDEX        // 指数
	BLOCK        // 板块
	ETF          // ETF
)

// AssertCode 判断一个代码类型
func AssertCode(securityCode string) TargetKind {
	marketId, _, code := DetectMarket(securityCode)
	// 板块, 板块指数的数据在上海
	if marketId == MarketIdShangHai && api.StartsWith(code, sectorPrefixes) {
		return BLOCK
	}
	// 上海代码, 000开头为指数
	if marketId == MarketIdShangHai && api.StartsWith(code, shanghaiIndexPrefixes) {
		return INDEX
	}
	// 深圳代码, 399开头为指数
	if marketId == MarketIdShenZhen && api.StartsWith(code, shenzhenIndexPrefixes) {
		return INDEX
	}
	if marketId == MarketIdBeiJing && api.StartsWith(code, beijingIndexPrefixes) {
		return INDEX
	}
	// ETF, 上海
	if marketId == MarketIdShangHai && api.StartsWith(code, []string{"51"}) {
		return ETF
	}
	// ETF, 深圳
	if marketId == MarketIdShenZhen && api.StartsWith(code, []string{"159"}) {
		return ETF
	}

	return STOCK
}
