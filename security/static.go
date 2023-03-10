package security

import (
	"fmt"
	"gitee.com/quant1x/data/category"
	"gitee.com/quant1x/data/internal"
	"gitee.com/quant1x/data/security/futu"
	"gitee.com/quant1x/data/stock"
	"gitee.com/quant1x/pandas"
	"gitee.com/quant1x/pandas/stat"
	"github.com/mymmsc/gox/api"
	"github.com/mymmsc/gox/errors"
	"github.com/mymmsc/gox/logger"
	"os"
	"time"
)

var (
	// MarketName 市场名称
	MarketName = map[string]string{
		category.MARKET_SH: "上海",
		category.MARKET_SZ: "深圳",
		category.MARKET_HK: "香港",
	}
	MarketSecurity = map[string]int32{
		category.MARKET_SH: int32(futu.QotMarket_QotMarket_CNSH_Security),
		category.MARKET_SZ: int32(futu.QotMarket_QotMarket_CNSZ_Security),
		category.MARKET_HK: int32(futu.QotMarket_QotMarket_HK_Security),
	}
	// 证券市场缓存
	cacheSecurity = map[string]StaticBasic{}
)

var (
	// ErrCacheNotExist 没有缓存
	ErrCacheNotExist = errors.New("Cache not exist")
	// ErrCodeNotExist 证券代码不存在
	ErrCodeNotExist = errors.New("Securities code does not exist")
)

func GetStockCode(market string, code string) string {
	return fmt.Sprintf("%s%s", market, code)
}

// Security 两个字段确定一支股票
type Security struct {
	Market int32  `json:"market,omitempty"` //QotMarket,股票市场
	Code   string `json:"code,omitempty"`   //股票代码
}

type StaticBasic struct {
	Security      Security `json:"security,omitempty"`      //股票
	Id            int64    `json:"id,omitempty"`            //股票ID
	LotSize       int32    `json:"lotSize,omitempty"`       //每手数量,期权以及期货类型表示合约乘数
	SecType       int32    `json:"secType,omitempty"`       //Qot_Common.SecurityType,股票类型
	Name          string   `json:"name,omitempty"`          //股票名字
	ListTime      string   `json:"listTime,omitempty"`      //上市时间字符串
	Delisting     bool     `json:"delisting,omitempty"`     //是否退市
	ListTimestamp float64  `json:"listTimestamp,omitempty"` //上市时间戳, 秒数
	//CCass         string   `json:"cCass,omitempty"`         //CCASS股份编号
	//ASecurity     bool     `json:aSecurity,omitempty`       //是否A股通
}

// 生成指数静态信息
func genIndexStaticInfo(market, code, name, listTime string, id int64, lotSize int32) (*StaticBasic, error) {
	marketSecurity, ok := MarketSecurity[market]
	if !ok {
		return nil, ErrCodeNotExist
	}
	fullCode := GetStockCode(market, code)
	listTimestamp, err := internal.ParseTime(listTime)
	if err != nil {
		return nil, err
	}
	return &StaticBasic{
		Security: Security{
			Market: marketSecurity,
			Code:   fullCode,
		},
		Id:            id,
		LotSize:       lotSize,
		SecType:       int32(futu.SecurityType_SecurityType_Index),
		Name:          name,
		ListTime:      listTime,
		Delisting:     false,
		ListTimestamp: float64(listTimestamp.Unix()),
	}, nil
}

func init() {
	logger.Infof("开始初始化市场静态数据...")
	// 1.上海指数
	// 上证综合指数 000001.sh 1990-12-19
	index, err := genIndexStaticInfo(category.MARKET_SH, "000001", "上证指数", "1990-12-19", 1000001, 100)
	if err == nil && index != nil {
		cacheSecurity[index.Security.Code] = *index
	}

	// 中证500 sh000905 2007-01-15
	index, err = genIndexStaticInfo(category.MARKET_SH, "000905", "中证500", "2007-01-15", 1000905, 100)
	if err == nil && index != nil {
		cacheSecurity[index.Security.Code] = *index
	}

	// 2. 深圳指数
	// 深证成指, sz399001, 1993-01-03
	index, err = genIndexStaticInfo(category.MARKET_SZ, "399001", "深证成指", "1993-01-03", 2399001, 100)
	if err == nil && index != nil {
		cacheSecurity[index.Security.Code] = *index
	}
	// 创业板指, sz399006, 2010-06-02
	index, err = genIndexStaticInfo(category.MARKET_SZ, "399006", "创业板指", "2010-06-02", 2399006, 100)
	if err == nil && index != nil {
		cacheSecurity[index.Security.Code] = *index
	}
	// 板块信息
	df := stock.BlockList()
	if df.Nrow() > 0 {
		for i := 0; i < df.Nrow(); i++ {
			m := df.IndexOf(i)
			code := stat.AnyToString(m["code"])
			name := stat.AnyToString(m["name"])
			if len(code) == 0 || len(name) == 0 {
				continue
			}
			index, err = genIndexStaticInfo(category.MARKET_SH, code, name, "1990-12-19", int64(8000000+i), 100)
			if err == nil && index != nil {
				cacheSecurity[code] = *index
			}
		}
	}
	// 3. 香港指数 恒生指数时间不准确
	// 4. 加载 上海的个股信息
	for market, name := range MarketName {
		logger.Infof("开始加载 %s 个股静态信息...", name)
		list, err := GetStaticBasic(market)
		if err != nil {
			logger.Errorf(name + "个股信息加载失败")
		} else {
			for _, item := range list {
				if item.Delisting || item.ListTimestamp == 0.00 {
					// 跳过退市和时间戳为0的个股
					continue
				}
				code := GetStockCode(market, item.Security.Code)
				cacheSecurity[code] = item
			}
		}
		logger.Infof("开始加载 %s 个股静态信息...OK", name)
	}
	logger.Infof("开始初始化市场静态数据...OK")
}

func GetStaticBasic(market string) (list []StaticBasic, err error) {
	filename := fmt.Sprintf("%s/%s.csv", ResourcesPath, market)
	reader, err := resources.Open(filename)
	df := pandas.ReadCSV(reader)
	if df.Nrow() == 0 {
		return list, ErrCacheNotExist
	}
	for i := 0; i < df.Nrow(); i++ {
		row := df.IndexOf(i)
		marketId := MarketSecurity[market]
		stockId := row["stock_id"].(int64)
		stockType := row["stock_type"].(string)
		secType := futu.StockType[stockType]
		code := row["code"].(string)
		_, _, code = category.DetectMarket(code)
		name := row["name"].(string)
		listTime := row["listing_date"].(string)
		lotSize := row["lot_size"].(int64)
		delisting := row["delisting"].(string)
		listTimestamp, err := internal.ParseTime(listTime)
		if err != nil {
			listTimestamp = time.Unix(0, 0)
		}
		info := StaticBasic{
			Security:      Security{Market: marketId, Code: code},
			Id:            stockId,
			SecType:       secType,
			LotSize:       int32(lotSize),
			Name:          name,
			ListTime:      listTime,
			Delisting:     stat.AnyToBool(delisting),
			ListTimestamp: float64(listTimestamp.Unix()),
		}
		list = append(list, info)
	}

	return
}

func GetBasicInfo(code string) (*StaticBasic, error) {
	info, ok := cacheSecurity[code]
	if !ok || info.Delisting {
		return nil, ErrCodeNotExist
	}
	return &info, nil
}

func WriteBasicInfo(market string, array []byte) {
	filename := fmt.Sprintf("%s/%s.json", ResourcesPath, market)
	fw, err := os.OpenFile("data/security/"+filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, category.CACHE_FILE_MODE)
	if err != nil {
		logger.Debugf("filename[%s]: JSON文件打开失败", filename, err)
	}
	defer api.CloseQuietly(fw)
	_, err = fw.Write(array)
	if err != nil {
		logger.Debugf("filename[%s]: JSON文件写入失败", filename, err)
	}
}
