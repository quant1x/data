package category

import (
	"github.com/mymmsc/gox/logger"
	"github.com/mymmsc/gox/util/homedir"
	"os"
)

const (
	// CACHE_DIR_MODE 目录权限
	CACHE_DIR_MODE os.FileMode = 0755
	// CACHE_FILE_MODE 文件权限
	CACHE_FILE_MODE os.FileMode = 0644

	// DEBUG 调试开关
	DEBUG = false
)

var (
	// DATA_ROOT_PATH 数据根路径
	DATA_ROOT_PATH = "~/.quant1x"
	// KLINE_PATH 日线数据文件路径
	KLINE_PATH = DATA_ROOT_PATH + "/day"
	// LOG_ROOT_PATH 日志路径
	LOG_ROOT_PATH = DATA_ROOT_PATH + "/logs"
)

func init() {
	rootPath, err := homedir.Expand(DATA_ROOT_PATH)
	if err != nil {
		panic(err)
	}
	DATA_ROOT_PATH = rootPath

	// 创建根路径
	if err := os.MkdirAll(DATA_ROOT_PATH, CACHE_DIR_MODE); err != nil {
		panic(err)
	}

	// 创建日志路径
	logsPath, err := homedir.Expand(LOG_ROOT_PATH)
	if err != nil {
		panic(err)
	}
	LOG_ROOT_PATH = logsPath
	if err := os.MkdirAll(LOG_ROOT_PATH, CACHE_DIR_MODE); err != nil {
		panic(err)
	}

	logger.SetLogPath(LOG_ROOT_PATH)

	klinePath, err := homedir.Expand(KLINE_PATH)
	if err != nil {
		panic(err)
	}
	KLINE_PATH = klinePath
	if err := os.MkdirAll(KLINE_PATH, CACHE_DIR_MODE); err != nil {
		panic(err)
	}
}
