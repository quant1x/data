package tdx

import (
	"fmt"
	"gitee.com/quant1x/gotdx/proto"
	"testing"
)

func TestGetKLineAll(t *testing.T) {
	df := GetKLineAll("000001.sh")
	fmt.Println(df)
	df = GetKLineAll("000001.sz", proto.KLINE_TYPE_WEEKLY)
	fmt.Println(df)
}
