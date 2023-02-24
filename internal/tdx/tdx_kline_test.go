package tdx

import (
	"fmt"
	"testing"
)

func TestGetKLineAll(t *testing.T) {
	df := GetKLineAll("sh000001")
	fmt.Println(df)
}
