package cache

import "fmt"

// BlockFilename 板块缓存路径
func BlockFilename(ns ...string) string {
	// 默认取板块列表
	name := "block"
	if len(ns) > 0 {
		name = ns[0]
	}
	filename := fmt.Sprintf("%s/%s.csv", GetBkPath(), name)
	_ = CheckFilepath(filename)
	return filename
}
