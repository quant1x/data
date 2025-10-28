package quotes

import (
	"encoding/json"
	"fmt"
	"testing"

	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gox/api"
)

func TestSecurityListAPackage(t *testing.T) {
	stdApi, err := NewStdApi()
	if err != nil {
		panic(err)
	}
	defer stdApi.Close()
	reply, err := stdApi.GetSecurityListA(exchange.MarketIdBeiJing, 0, 1000)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	fmt.Printf("%+v\n", reply)
	fmt.Println("==========")
	data, _ := json.Marshal(reply)
	text := api.Bytes2String(data)
	fmt.Println(text)
}
