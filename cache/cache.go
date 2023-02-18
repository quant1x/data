package cache

import (
	"gitee.com/quant1x/data/category"
	"os"
	"path/filepath"
	"syscall"
)

const (
	DAY_PATH  = "day"  // 日线路径
	INFO_PATH = "info" // 信息路径
	TICK_PATH = "tick" // tick路径
)

type CacheType int

const (
	CACHE_TARS CacheType = iota
	CACHE_CSV
)

var (
	// CACHE_ROOT_PATH cache路径
	CACHE_ROOT_PATH = category.DATA_ROOT_PATH
	CACHE_TYPE      CacheType
)

func init() {
	CACHE_TYPE = CACHE_CSV
}

// CheckFilepath
//
//	检查filename 文件路径, 如果不存在就创建
func CheckFilepath(filename string) error {
	path := filepath.Dir(filename)
	dir, err := os.Stat(path)
	if verr, ok := err.(*os.PathError); ok {
		if verr.Err == syscall.ENOENT {
			err = os.MkdirAll(path, category.CACHE_DIR_MODE)
			if err != nil {
				return err
			}
		}
	}
	dir, err = os.Stat(path)
	if err != nil {
		return err
	}
	if dir.IsDir() {
		return nil
	}
	return os.MkdirAll(path, category.CACHE_DIR_MODE)
}

// GetInfoPath 证券信息路径
func GetInfoPath() string {
	return CACHE_ROOT_PATH + "/" + INFO_PATH
}

// GetDayPath 历史数据-日线缓存路径
func GetDayPath() string {
	return CACHE_ROOT_PATH + "/" + DAY_PATH
}
