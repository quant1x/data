package dfcf

import (
	"fmt"
	"gitee.com/quant1x/pandas"
	"testing"
)

func Test_stock_hist(t *testing.T) {
	ks, err := A("600600")
	if err != nil {
		_ = fmt.Errorf("error :%v", err.Error())
	}
	fmt.Println(ks)
	df := pandas.LoadStructs(ks)
	fmt.Println(df)
}
