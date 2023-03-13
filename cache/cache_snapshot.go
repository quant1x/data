package cache

import "fmt"

// SnapshotFilename snapshot缓存路径
func SnapshotFilename(code string) string {
	cacheId := CacheId(code)
	length := len(cacheId)
	filepath := fmt.Sprintf("%s/%s/%s.csv", GetSnapshotPath(), cacheId[:length-3], cacheId)
	_ = CheckFilepath(filepath)
	return filepath
}
