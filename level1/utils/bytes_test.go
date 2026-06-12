package utils

import (
	"testing"

	"github.com/quant1x/data/cache"
)

func TestUnzipPreserveTimes(t *testing.T) {
	filename := "zhb.zip"
	srcZip := cache.GetBlockPath() + "/" + filename
	dest := "testdata2"
	err := UnzipPreserveTimes(srcZip, dest, "tdxhy.cfg", "tdxzs.cfg", "tdxzs3.cfg")
	if err != nil {
		t.Fatal(err)
	}
}
