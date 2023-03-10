package tdx

import (
	"fmt"
	"testing"
)

func TestGetTickAll(t *testing.T) {
	df := GetTickData("sh881263", "2023-02-22")
	fmt.Println(df)
}
