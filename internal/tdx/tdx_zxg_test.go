package tdx

import (
	"fmt"
	"testing"
)

func TestGetZxgList(t *testing.T) {
	list := GetZxgList()
	fmt.Println(list)
}
