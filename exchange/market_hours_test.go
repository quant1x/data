package exchange

import (
	"fmt"
	"testing"
	"time"

	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/pkg/yaml"
	"github.com/stretchr/testify/assert"
)

type tradingTimeRange struct {
	Session    TradingTimeRange `yaml:"session"`
	ChangeRate float64          `yaml:"change_rate" default:"0.01"`
}

func TestTradingTimeRangeWithParse(t *testing.T) {
	text := `time: "09:50:00~09:50:59,10:50:00~10:50:59"`
	text = `
session: "09:50:00~09:50:59"
name: "buyiwang"
`
	bytes := api.String2Bytes(text)
	v := tradingTimeRange{}
	err := yaml.Unmarshal(bytes, &v)
	fmt.Println(err, v)
}

func createTimeOfToday(timeString string) time.Time {
	layout := "2006-01-02 15:04:05" // 定义输入字符串的格式

	now := time.Now()
	today := now.Format("2006-01-02")
	dateTimeString := fmt.Sprintf("%s %s", today, timeString)
	parsedTime, err := time.ParseInLocation(layout, dateTimeString, time.Local)
	if err != nil {
		panic("解析时间字符串出错")
	}
	return parsedTime
}
func TestTimeKind(t *testing.T) {
	text := "T|09:15:00~09:26:59,09:15:00 ~ 09:19:59,09:25:00~11:29:59,13:00:00~14:59:59,09:00:00~09:14:59"
	var ts MarketHours
	ts.Parse(text)

	timestamp := createTimeOfToday("09:30:01").UnixMilli()
	kind, idx, err := ts.GetTimeKind(timestamp)
	assert.Equal(t, nil, err)
	assert.Equal(t, TradingAndCancel, kind)
	assert.Equal(t, 3, idx)

	timestamp = createTimeOfToday("09:20:00").UnixMilli()
	kind, idx, err = ts.GetTimeKind(timestamp)
	assert.Equal(t, nil, err)
	assert.Equal(t, Trading, kind)
	assert.Equal(t, 2, idx)

	timestamp = createTimeOfToday("08:20:00").UnixMilli()
	kind, idx, err = ts.GetTimeKind(timestamp)
	assert.NotNil(t, err)
	assert.Equal(t, TimeKind(0x0), kind)
	assert.Equal(t, -1, idx)
}
