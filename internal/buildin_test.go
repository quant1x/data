package internal

import (
	"fmt"
	"gitee.com/quant1x/data/dfcf"
	"reflect"
	"testing"
)

func Test_initTag(t *testing.T) {
	emkl := dfcf.KLine{}
	initTag(reflect.TypeOf(emkl), "name")
}

func TestToCSV(t *testing.T) {
	_, err := dfcf.A("600600")
	if err != nil {
		_ = fmt.Errorf("error :%v", err.Error())
	}
}
