package tdx

import (
	"fmt"
	"testing"
)

func TestGetMinuteHistory(t *testing.T) {
	hs := GetMinuteHistory("sh881266", "2023-03-16")
	fmt.Println(hs)
}

func TestQuantityRelativeRatio(t *testing.T) {
	code := "sh603789"
	code = "sz002483"
	ratio := QuantityRelativeRatio(code)
	fmt.Println(code, ratio)
}
