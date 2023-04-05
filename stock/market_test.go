package stock

import (
	"fmt"
	"testing"
)

func Test_needIgnore(t *testing.T) {
	code := "sz002022"
	fmt.Println(needIgnore(code))
}
