package cache

import (
	"fmt"
	"testing"
)

func TestKLine(t *testing.T) {
	df := KLine("002528")
	fmt.Println(df)
}
