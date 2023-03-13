package cache

import (
	"gitee.com/quant1x/data/category"
	"gitee.com/quant1x/data/internal"
	"os"
	"path/filepath"
	"time"
)

const (
	DAY_PATH      = "day"      // 日线路径
	INFO_PATH     = "info"     // 信息路径
	TICK_PATH     = "tick"     // tick路径
	XDXR_PATH     = "xdxr"     // 除权除息路径
	BLOCK_PATH    = "bk"       // 板块数据
	SNAPSHOT_PATH = "snapshot" // 快照数据路径
)

type CacheType int

const (
	CACHE_TARS  CacheType = iota
	CACHE_CSV             = 1 // CSV
	CACHE_EXCEL           = 2 // EXCEL
)

var (
	// CACHE_ROOT_PATH cache路径
	CACHE_ROOT_PATH           = category.DATA_ROOT_PATH
	CACHE_TYPE      CacheType = CACHE_TARS
)

func init() {
	CACHE_TYPE = CACHE_CSV
}

func Today() string {
	now := time.Now()
	return now.Format(category.CACHE_DATE)
}

// CorrectDate 矫正日期, 统一格式
func CorrectDate(date string) string {
	dt, err := internal.ParseTime(date)
	if err != nil {
		return Today()
	}
	date = dt.Format(category.CACHE_DATE)
	return date
}

// FileExist 路径是否存在
func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

// CheckFilepath
//
//	检查filename 文件路径, 如果不存在就创建
func CheckFilepath(filename string) error {
	path := filepath.Dir(filename)
	dir, err := os.Stat(path)
	if err == nil {
		// 已存在
		return nil
	}
	if os.IsExist(err) {
		// 已存在
		return nil
	}
	// 不存在, 创建
	err = os.MkdirAll(path, category.CACHE_DIR_MODE)
	if err != nil {
		return err
	}
	dir, err = os.Stat(path)
	if err != nil {
		return err
	}
	if dir.IsDir() {
		return nil
	}
	return os.ErrNotExist
}

// GetInfoPath 证券信息路径
func GetInfoPath() string {
	return CACHE_ROOT_PATH + "/" + INFO_PATH
}

// GetDayPath 历史数据-日线缓存路径
func GetDayPath() string {
	return CACHE_ROOT_PATH + "/" + DAY_PATH
}

// GetTickPath tick数据路径
func GetTickPath() string {
	return CACHE_ROOT_PATH + "/" + TICK_PATH
}

// GetZxgFile 自选股文件路径
func GetZxgFile() string {
	return CACHE_ROOT_PATH + "/zxg.csv"
}

// GetXdxrPath 除权除息文件存储路径
func GetXdxrPath() string {
	return CACHE_ROOT_PATH + "/" + XDXR_PATH
}

// GetBkPath 板块路径
func GetBkPath() string {
	return CACHE_ROOT_PATH + "/" + BLOCK_PATH
}

// GetSnapshotPath 获取快照路径
func GetSnapshotPath() string {
	return CACHE_ROOT_PATH + "/" + SNAPSHOT_PATH
}
