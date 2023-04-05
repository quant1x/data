package category

import (
	"github.com/mymmsc/gox/api"
	"strings"
)

type Market = uint8

const (
	MARKET_ID_SHENZHEN Market = iota // 深圳
	MARKET_ID_SHANGHAI Market = 1    // 上海
	MARKET_ID_BEIJING  Market = 2    // 北京
	MARKET_ID_HONGKONG Market = 21   // 香港
	MARKET_ID_USA      Market = 22   // 美国
)

const (
	MARKET_SH string = "sh" // 上海
	MARKET_SZ string = "sz" // 深圳
	MARKET_BJ string = "bj" // 北京
	MARKET_HK string = "hk" // 香港
	MARKET_US string = "us" // 美国
)

var (
	kMarketFlags = []string{"sh", "sz", "SH", "SZ", "bj", "BJ", "hk", "HK", "us", "US"}
)

// GetMarketName 通过市场ID取得市场名称缩写
func GetMarketName(marketId Market) string {
	switch marketId {
	case MARKET_ID_SHENZHEN:
		return MARKET_SZ
	case MARKET_ID_BEIJING:
		return MARKET_BJ
	case MARKET_ID_HONGKONG:
		return MARKET_HK
	case MARKET_ID_USA:
		return MARKET_US
	default:
		return MARKET_SH
	}
}

// DetectMarket 检测市场代码
func DetectMarket(symbol string) (marketId Market, market string, code string) {
	code = strings.TrimSpace(symbol)
	market = MARKET_SH
	if api.StartsWith(code, kMarketFlags) {
		// 前缀2位字母后面跟代码
		market = strings.ToLower(code[0:2])
		if code[2:3] == "." {
			code = code[3:]
		} else {
			code = code[2:]
		}
	} else if api.EndsWith(code, kMarketFlags) {
		length := len(code)
		// 后缀一个点号+2位字母, 代码在最前面
		market = strings.ToLower(code[length-2:])
		code = code[:length-3]
	} else if api.StartsWith(code, []string{"50", "51", "60", "68", "90", "110", "113", "132", "204"}) {
		market = MARKET_SH
	} else if api.StartsWith(code, []string{"00", "12", "13", "18", "15", "16", "18", "20", "30", "39", "115", "1318"}) {
		market = MARKET_SZ
	} else if api.StartsWith(code, []string{"5", "6", "9", "7"}) {
		market = MARKET_SH
	} else if api.StartsWith(code, []string{"88"}) {
		market = MARKET_SH
	} else if api.StartsWith(code, []string{"4", "8"}) {
		market = MARKET_BJ
	}
	marketId = MARKET_ID_SHANGHAI
	if market == MARKET_SH {
		marketId = MARKET_ID_SHANGHAI
	} else if market == MARKET_SZ {
		marketId = MARKET_ID_SHENZHEN
	} else if market == MARKET_BJ {
		marketId = MARKET_ID_BEIJING
	} else if market == MARKET_HK {
		marketId = MARKET_ID_HONGKONG
	}
	return marketId, market, code
}

// CodeIsIndex 证券代码是否指数
func CodeIsIndex(code string) bool {
	marketId, _, shortCode := DetectMarket(code)
	return IndexFromMarketAndCode(marketId, shortCode)
}

// IndexFromMarketAndCode 通过市场id和短码判断是否指数
func IndexFromMarketAndCode(marketId Market, code string) bool {
	if marketId == MARKET_ID_SHANGHAI && api.StartsWith(code, []string{"000", "88"}) {
		return true
	} else if marketId == MARKET_ID_SHENZHEN && api.StartsWith(code, []string{"399"}) {
		return true
	}
	return false
}

// MarketLimit 涨跌停板限制
func MarketLimit(code string) float64 {
	_, _, shortCode := DetectMarket(code)
	if api.StartsWith(shortCode, []string{"30", "68"}) {
		return 0.20
	}
	return 0.10
}
