package security

import "testing"

func Test_get_block_file(t *testing.T) {
	get_block_file("block_fg.dat")
}

func TestGetBlockList(t *testing.T) {
	genBlockFile()
}
