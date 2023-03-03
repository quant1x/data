package internal

import (
	"gitee.com/quant1x/data/internal/dfcf"
	"reflect"
	"testing"
)

func Test_initTag(t *testing.T) {
	emkl := dfcf.KLine{}
	InitTag(reflect.TypeOf(emkl), "name")
}
