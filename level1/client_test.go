package level1

import (
	"fmt"
	"testing"
	"time"

	"gitee.com/quant1x/data/level1/proto"
)

func TestReOpen(t *testing.T) {
	api := GetApi()
	v, _ := api.GetXdxrInfo("sh600072")
	fmt.Println(v)
	time.Sleep(20 * time.Second)
	ReOpen()
	v, _ = api.GetXdxrInfo("sh600072")
	fmt.Println(v)
	fmt.Println(api.NumOfServers())
	klines, _ := api.GetKLine("sh600600", proto.KLINE_TYPE_RI_K, 0, 1)
	fmt.Println(klines)
}
