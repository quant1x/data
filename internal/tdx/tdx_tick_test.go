package tdx

import (
	"fmt"
	"testing"
)

func TestGetTickData(t *testing.T) {
	df := GetTickData("sh600600", "2023-04-04")
	fmt.Println(df)
}

func TestGetTickAll(t *testing.T) {
	code := "sz002351"
	GetTickAll(code)
}
