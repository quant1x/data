package tdx

import (
	"fmt"
	"testing"
)

func TestGetKLineAll(t *testing.T) {
	df := GetKLineAll("sz002528")
	fmt.Println(df)
}
