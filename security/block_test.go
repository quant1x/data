package security

import (
	"fmt"
	"testing"
)

func Test_get_block_file(t *testing.T) {
	get_block_file("block_fg.dat")
}

func TestGetBlockList(t *testing.T) {
	genBlockFile()
}

func TestBlockNameByType(t *testing.T) {
	code := 1
	name, ok := BlockTypeNameByCode(code)
	fmt.Println(name, ok)
}
