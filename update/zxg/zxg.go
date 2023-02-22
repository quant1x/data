package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"gitee.com/quant1x/data/cache"
	"gitee.com/quant1x/data/category"
	"gitee.com/quant1x/data/internal/tdx"
	"gitee.com/quant1x/data/security"
	"github.com/mymmsc/gox/api"
	"github.com/mymmsc/gox/logger"
	"github.com/mymmsc/gox/util/homedir"
	"io"
	"log"
	"os"
	"os/user"
	"strings"
	"syscall"
)

const (
	tdx_path = "/workspace/data/tdx"
)

func main() {
	var (
		path string // 数据路径
	)
	flag.StringVar(&path, "path", "", "通达信安装目录")
	flag.Parse()
	if len(path) == 0 {
		u, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Home dir:", u.HomeDir+tdx_path)
		path = u.HomeDir + tdx_path
	}
	filename, err := homedir.Expand(path + tdx.BlockPath + "/" + tdx.ZxgBlk)
	f, err := os.Open(filename)
	if err != nil {
		if errors.Is(err, syscall.ENOENT) {
			logger.Debugf("自选股[%s]: 通达信自选股blk文件不存在", filename)
			return
		} else {
			logger.Errorf("自选股[%s]: 通达信自选股blk文件操作失败:%v", filename, err)
			return
		}
	}

	//var data []byte
	data, err := io.ReadAll(f)
	if err != nil {
		logger.Errorf("自选股[%s]: 通达信自选股blk文件操作失败:%v", filename, err)
	}
	if len(data) == 0 {
		logger.Errorf("自选股[%s]: 通达信自选股blk空", filename)
	}

	s := string(data)
	arr := strings.Split(s, "\r\n")
	// 深圳指数(0, ‘399001’)，上海大盘 (1, ‘999999’)。
	// 数据在’ZXG.blk’中以8个字节来存放。
	fcsv, _ := os.OpenFile(cache.GetZxgFile(), category.CACHE_REPLACE, category.CACHE_FILE_MODE)
	defer api.CloseQuietly(fcsv)
	//fcsv.WriteString("\xEF\xBB\xBF")
	out := csv.NewWriter(fcsv)
	var header = []string{"market", "code", "name"}
	_ = out.Write(header)
	for _, d := range arr {
		d = strings.TrimSpace(d)
		if len(d) != 7 {
			continue
		}
		market := d[:1]
		code := d[1:]
		fmt.Printf("市场编码:%s, 证券代码:%s\t=>\t", market, code)
		fullCode := ""
		code = strings.TrimSpace(code)
		if market == "1" {
			market = "上海"
			fullCode = "sh" + code
		} else if market == "0" {
			market = "深圳"
			fullCode = "sz" + code
		} else {
			continue
		}
		fmt.Printf("市场编码:%s, 证券代码:%s\n", market, fullCode)
		info, err := security.GetBasicInfo(fullCode)
		if err != nil {
			fmt.Printf("没有找到 %s\n", fullCode)
			continue
		}
		row := []string{market, fullCode, info.Name}
		_ = out.Write(row)
		out.Flush()
	}
	out.Flush()
	api.CloseQuietly(fcsv)
}
