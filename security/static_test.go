package security

import (
	"fmt"
	"testing"
)

func TestGetStaticBasic(t *testing.T) {
	list, err := GetStaticBasic("sh")

	fmt.Println(list)
	_ = err
}
