package tdx

import (
	"fmt"
	"testing"
)

func TestGetXdxrInfo(t *testing.T) {
	df := GetCacheXdxr("sz002900")
	fmt.Println(df)
}
