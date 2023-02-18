package cache

import (
	"errors"
	"fmt"
	"gitee.com/quant1x/data/category"
	"github.com/mymmsc/gox/logger"
	"os"
	"strings"
	"syscall"
)

type FastCache struct {
	filename string
	f        *os.File
	size     int64
	Data     []byte
}

func Open(filename string) (*FastCache, error) {
	_ = CheckFilepath(filename)
	f, err := os.OpenFile(filename, os.O_RDONLY, category.CACHE_FILE_MODE)
	if err != nil {
		return nil, err
	}
	fileinfo, err := f.Stat()
	if err != nil {
		return nil, err
	}
	filelength := fileinfo.Size()
	data, err := syscall.Mmap(int(f.Fd()), 0, int(filelength), syscall.PROT_READ, syscall.MAP_SHARED)
	if nil != err {
		return nil, err
	}

	return &FastCache{
		filename: filename,
		f:        f,
		size:     filelength,
		Data:     data,
	}, nil
}

func Create(filename string, size int64) (*FastCache, error) {
	_ = CheckFilepath(filename)
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, category.CACHE_FILE_MODE)
	if nil != err {
		return nil, err
	}
	err = f.Truncate(size)
	if err != nil {
		return nil, err
	}
	data, err := syscall.Mmap(int(f.Fd()), 0, int(size), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if nil != err {
		return nil, err
	}

	return &FastCache{
		filename: filename,
		f:        f,
		size:     size,
		Data:     data,
	}, nil
}

func (fc *FastCache) Close() {
	if fc != nil && fc.f != nil {
		_ = syscall.Munmap(fc.Data)
		_ = fc.f.Close()
	}
}

// GetKLineFilename 获取缓存的文件名
func GetKLineFilename(fullCode string) string {
	fullCode = strings.TrimSpace(fullCode)
	if len(fullCode) != 7 && len(fullCode) != 8 {
		return fullCode
	}
	pos := len(fullCode) - 3
	fullCode = strings.ToLower(fullCode)
	// 组织存储路径
	filename := GetDayPath() + "/" + fullCode[0:pos] + "/" + fullCode
	if CACHE_TYPE == CACHE_CSV {
		filename += ".csv"
	}
	err := CheckFilepath(filename)
	if err != nil {
		panic(fmt.Errorf("create file %s, failed", fullCode))
	}
	return filename
}

func GetCache(fullCode string) *os.File {
	filename := GetKLineFilename(fullCode)
	file, err := os.Open(filename)
	if err != nil {
		//ENOENT (2)
		if errors.Is(err, syscall.ENOENT) {
			logger.Debugf("code[%s]: K线数据文件不存在", fullCode)
			return nil
		} else {
			logger.Errorf("code[%s]: K线数据文件操作失败:%v", fullCode, err)
			return nil
		}
	}
	return file
}
