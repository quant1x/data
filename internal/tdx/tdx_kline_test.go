package tdx

import (
	"fmt"
	"testing"
)

func TestGetKLineAll(t *testing.T) {
	df := GetKLineAll("sh881396")
	fmt.Println(df)
}
