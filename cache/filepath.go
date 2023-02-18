package cache

import (
	"gitee.com/quant1x/data/category"
	"strings"
)

func KLinePath(fc string) (string, string, int) {
	fc = strings.TrimSpace(fc)
	fcLen := len(fc)
	if fcLen != 7 && fcLen != 8 {
		return fc, "", category.D_ECODE | category.D_ECODE
	}
	pos := len(fc) - 3
	fc = strings.ToLower(fc)
	// 组织存储路径
	filename := GetDayPath() + "/" + fc[0:pos] + "/" + fc + ".csv"
	return fc, filename, category.D_OK
}
