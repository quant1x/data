package tdx

import (
	"fmt"
	"testing"
)

func TestGetKLineAll(t *testing.T) {
	df := GetKLineAll("000001.sh")
	fmt.Println(df)
	df = GetKLineAll("000001.sz")
	fmt.Println(df)
}
