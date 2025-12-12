package data

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config 配置结构体
type Config struct {
	BaseDir string  `yaml:"basedir"` // 数据根目录
	Debug   bool    `yaml:"debug"`   // 是否调试模式运行
	Trader  Trader  `yaml:"trader"`  // 交易配置
	Runtime Runtime `yaml:"runtime"` // 运行时配置
	Data    Data    `yaml:"data"`    // 数据配置
}

// Trader 交易配置
type Trader struct {
	AccountId            int64      `yaml:"account_id"`               // 账户ID
	ProxyUrl             string     `yaml:"proxy_url"`                // 代理地址
	OrderPath            string     `yaml:"order_path"`               // 委托单保存路径
	StampDutyRateForBuy  float64    `yaml:"stamp_duty_rate_for_buy"`  // 买入印花税率
	StampDutyRateForSell float64    `yaml:"stamp_duty_rate_for_sell"` // 卖出印花税率
	TransferRate         float64    `yaml:"transfer_rate"`            // 过户费
	CommissionRate       float64    `yaml:"commission_rate"`          // 佣金率
	CommissionMin        float64    `yaml:"commission_min"`           // 佣金最低
	PositionRatio        float64    `yaml:"position_ratio"`           // 买入总占比
	KeepCash             float64    `yaml:"keep_cash"`                // 预留备用金
	BuyAmountMax         float64    `yaml:"buy_amount_max"`           // 最大可买金额
	BuyAmountMin         float64    `yaml:"buy_amount_min"`           // 最小可买金额
	Strategies           []Strategy `yaml:"strategies"`               // 策略列表
	AskTime              string     `yaml:"ask_time"`                 // 卖出时段
	TickTime             string     `yaml:"tick_time"`                // 盘中买入时段
}

// Strategy 策略配置
type Strategy struct {
	Id                          int64    `yaml:"id"`                             // 策略ID
	Name                        string   `yaml:"name"`                           // 策略名称
	Auto                        bool     `yaml:"auto"`                           // 是否自动交易
	Flag                        string   `yaml:"flag"`                           // 策略标识
	Time                        string   `yaml:"time"`                           // 策略运行时间
	Weight                      float64  `yaml:"weight"`                         // 策略权重
	Total                       int      `yaml:"total"`                          // 策略总仓位数
	PriceCageRatio              float64  `yaml:"price_cage_ratio"`               // 价格笼比例
	MinimumPriceFluctuationUnit float64  `yaml:"minimum_price_fluctuation_unit"` // 最小价格变动单位
	FeeMax                      float64  `yaml:"fee_max"`                        // 手续费封顶
	FeeMin                      float64  `yaml:"fee_min"`                        // 手续费最低
	Sectors                     []string `yaml:"sectors"`                        // 选股板块
	IgnoreMarginTrading         bool     `yaml:"ignore_margin_trading"`          // 忽略两融标的
	HoldingPeriod               int      `yaml:"holding_period"`                 // 持有周期(交易日)
	SellStrategy                int64    `yaml:"sell_strategy"`                  // 卖出策略
	Rules                       Rules    `yaml:"rules"`                          // 选股规则
	TakeProfitRatio             float64  `yaml:"take_profit_ratio,omitempty"`    // 止盈比例
	StopLossRatio               float64  `yaml:"stop_loss_ratio,omitempty"`      // 止损比例
}

// Rules 选股规则
type Rules struct {
	SectorsFilter     bool   `yaml:"sectors_filter"`                // 是否启用板块过滤
	SectorsTopN       int    `yaml:"sectors_top_n"`                 // 板块数量
	StockTopNInSector int    `yaml:"stock_top_n_in_sector"`         // 每个板块选股数量
	GapDown           bool   `yaml:"gap_down"`                      // 是否跳空低开
	IgnoreRuleGroup   []int  `yaml:"ignore_rule_group"`             // 忽略规则组
	Capital           string `yaml:"capital"`                       // 资金流入
	Price             string `yaml:"price"`                         // 价格
	OpenTurnZ         string `yaml:"open_turn_z,omitempty"`         // 开盘换手率Z值
	OpenChangeRate    string `yaml:"open_change_rate,omitempty"`    // 开盘涨跌幅
	OpenQuantityRatio string `yaml:"open_quantity_ratio,omitempty"` // 开盘量比
	VolumeRatio       string `yaml:"volume_ratio,omitempty"`        // 量比
	CheckSafetyScore  bool   `yaml:"check_safety_score"`            // 是否检查安全评分
	ChangeRate        string `yaml:"change_rate,omitempty"`         // 涨跌幅
	SafetyScore       string `yaml:"safety_score,omitempty"`        // 安全评分
}

// Runtime 运行时配置
type Runtime struct {
	Crontab map[string]CrontabItem `yaml:"crontab"` // 定时任务
}

// CrontabItem 定时任务
type CrontabItem struct {
	Enable  bool   `yaml:"enable"`            // 是否启用
	Trigger string `yaml:"trigger,omitempty"` // 触发器表达式
}

// Data 数据配置
type Data struct {
	Concurrency map[string]int `yaml:"concurrency"` // 并发控制
	Cache       CacheConfig    `yaml:"cache"`       // 缓存配置
	Trans       TransConfig    `yaml:"trans"`       // 成交数据配置
	Feature     FeatureConfig  `yaml:"feature"`     // 特征配置
	Snapshot    SnapshotConfig `yaml:"snapshot"`    // 快照配置
}

// CacheConfig 缓存配置
type CacheConfig struct {
	Kline map[string]bool `yaml:"kline"` // K线配置
	Chips ChipsConfig     `yaml:"chips"` // 筹码配置
}

// ChipsConfig 筹码配置
type ChipsConfig struct {
	IndexEnabled bool `yaml:"index_enabled"` // 是否启用指数筹码
}

// TransConfig 成交数据配置
type TransConfig struct {
	BeginDate string `yaml:"begin_date"` // 开始日期
}

// FeatureConfig 特征配置
type FeatureConfig struct {
	Tendency int        `yaml:"tendency"` // 策略是趋势主导还是股价主导
	Wave     WaveConfig `yaml:"wave"`     // 波浪配置
}

// WaveConfig 波浪配置
type WaveConfig struct {
	Fields  WaveFields `yaml:"fields"`  // 波浪-字段
	Periods int        `yaml:"periods"` // 波浪-K线周期
}

// WaveFields 波浪-字段
type WaveFields struct {
	Peak   string `yaml:"peak"`   // 波浪-字段-波峰
	Valley string `yaml:"valley"` // 波浪-字段-波谷
}

// SnapshotConfig 快照配置
type SnapshotConfig struct {
	Concurrency int `yaml:"concurrency"` // 快照-并发数
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		BaseDir: `d:\quant1x\data2`,
		Debug:   false, // 是否调试模式运行
		Trader: Trader{
			StampDutyRateForBuy:  0.0000,              // 买入印花税率
			StampDutyRateForSell: 0.0010,              // 卖出印花税率
			TransferRate:         0.0006,              // 过户费
			CommissionRate:       0.00025,             // 佣金率
			CommissionMin:        5.0000,              // 佣金最低
			PositionRatio:        1.0000,              // 买入总占比
			KeepCash:             1000.00,             // 预留备用金
			BuyAmountMax:         250000.00,           // 最大可买金额
			BuyAmountMin:         1000.00,             // 最小可买金额
			AskTime:              "09:50:00~14:59:30", // 卖出时段
			TickTime:             "09:39:00-14:56:30", // 盘中买入时段
		},
		Runtime: Runtime{
			Crontab: map[string]CrontabItem{
				"realtime_kline": {Enable: false},
				"sell_117":       {Enable: false, Trigger: "@every 1s"},
				"sell_861":       {Enable: true, Trigger: "@every 5s"},
			},
		},
		Data: Data{
			Concurrency: map[string]int{ // 并发控制, 默认5个线程
				"default": 5,
				"xdxr":    8,
				"day_raw": 2,
				"day":     2,
				"trans":   4,
				"f10":     2,
			},
			Cache: CacheConfig{
				Kline: map[string]bool{ // K线
					"5min": true, // 1分钟K线
				},
				Chips: ChipsConfig{
					IndexEnabled: false,
				},
			},
			Trans: TransConfig{
				BeginDate: "2023-10-01", // 成交数据开始日期
			},
			Feature: FeatureConfig{ // 特征
				Tendency: 0, // 策略是趋势主导还是股价主导
				Wave: WaveConfig{ // 波浪
					Fields: WaveFields{ // 波浪-字段
						Peak:   "close", // 波浪-字段-波峰, 默认收盘价close
						Valley: "close", // 波浪-字段-波谷, 默认收盘价close
					},
					Periods: 89, // 波浪-K线周期, 默认89天
				},
			},
			Snapshot: SnapshotConfig{ // 快照
				Concurrency: 2, // 快照-并发数为2
			},
		},
	}
}

func ParseConfig(filename string) (*Config, error) {
	config := DefaultConfig()
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
