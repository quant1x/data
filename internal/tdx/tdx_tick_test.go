package tdx

import (
	"fmt"
	"testing"
)

func TestGetTickAll(t *testing.T) {
	df := GetTickData("sh600600", "2023-04-04")
	fmt.Println(df)
}
