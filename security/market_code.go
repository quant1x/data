package security

import (
	"gitee.com/quant1x/gotdx/util"
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
	MARKET_HK string = "hk" // 香港
	MARKET_US string = "us" // 美国
)

var (
	kMarketFlags = []string{"sh", "sz", "SH", "SZ", "bj", "BJ", "hk", "HK", "us", "US"}
)

// GetMarket 判断股票ID对应的证券市场匹配规则
//
//	['50', '51', '60', '90', '110'] 为 sh
//	['00', '12'，'13', '18', '15', '16', '18', '20', '30', '39', '115'] 为 sz
//	['5', '6', '9'] 开头的为 sh， 其余为 sz
//
// Deprecated: 不推荐使用
func GetMarket(symbol string) string {
	market := "sh"
	if util.StartsWith(symbol, []string{"sh", "sz", "SH", "SZ"}) {
		market = strings.ToLower(symbol[0:2])
	} else if util.StartsWith(symbol, []string{"50", "51", "60", "68", "90", "110", "113", "132", "204"}) {
		market = "sh"
	} else if util.StartsWith(symbol, []string{"00", "12", "13", "18", "15", "16", "18", "20", "30", "39", "115", "1318"}) {
		market = "sz"
	} else if util.StartsWith(symbol, []string{"5", "6", "9", "7"}) {
		market = "sh"
	} else if util.StartsWith(symbol, []string{"4", "8"}) {
		market = "bj"
	}
	return market
}

// GetMarketId 获得市场ID
// Deprecated: 不推荐使用
func GetMarketId(symbol string) Market {
	market := GetMarket(symbol)
	marketId := MARKET_ID_SHANGHAI
	if market == "sh" {
		marketId = MARKET_ID_SHANGHAI
	} else if market == "sz" {
		marketId = MARKET_ID_SHENZHEN
	} else if market == "bj" {
		marketId = MARKET_ID_BEIJING
	}
	return marketId
}

// DetectMarket 检测市场代码
func DetectMarket(symbol string) (Market, string, string) {
	code := strings.TrimSpace(symbol)
	market := "sh"
	if api.StartsWith(code, kMarketFlags) {
		// 前缀2位字母后面跟代码
		market = strings.ToLower(code[0:2])
		code = code[2:]
	} else if api.EndsWith(code, kMarketFlags) {
		length := len(code)
		// 后缀一个点号+2位字母, 代码在最前面
		market = strings.ToLower(code[length-2:])
		code = code[:length-3]
	} else if api.StartsWith(code, []string{"50", "51", "60", "68", "90", "110", "113", "132", "204"}) {
		market = "sh"
	} else if api.StartsWith(code, []string{"00", "12", "13", "18", "15", "16", "18", "20", "30", "39", "115", "1318"}) {
		market = "sz"
	} else if api.StartsWith(code, []string{"5", "6", "9", "7"}) {
		market = "sh"
	} else if api.StartsWith(code, []string{"4", "8"}) {
		market = "bj"
	}
	marketId := MARKET_ID_SHANGHAI
	if market == "sh" {
		marketId = MARKET_ID_SHANGHAI
	} else if market == "sz" {
		marketId = MARKET_ID_SHENZHEN
	} else if market == "bj" {
		marketId = MARKET_ID_BEIJING
	}
	return marketId, market, code
}
