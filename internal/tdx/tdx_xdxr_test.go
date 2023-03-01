package tdx

import (
	"fmt"
	"testing"
)

func TestGetXdxrInfo(t *testing.T) {
	df := GetCacheXdxr("sh600018")
	fmt.Println(df)
}

func TestGetXdxrInfo1(t *testing.T) {
	df := GetXdxrInfo("sh600018")
	fmt.Println(df)
}
