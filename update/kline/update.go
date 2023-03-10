package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"gitee.com/quant1x/data/cache"
	"gitee.com/quant1x/data/category"
	"gitee.com/quant1x/data/internal"
	"gitee.com/quant1x/data/internal/dfcf"
	"gitee.com/quant1x/data/internal/tdx"
	"gitee.com/quant1x/data/security"
	"gitee.com/quant1x/data/stock"
	"gitee.com/quant1x/data/update/cross"
	"github.com/mymmsc/gox/api"
	"github.com/mymmsc/gox/cron"
	"github.com/mymmsc/gox/logger"
	"github.com/mymmsc/gox/progressbar"
	"os"
	"os/signal"
	"time"
)

var (
	source     int    //数据源
	MinVersion string // 版本号
)

// 更新日线数据工具
func main() {
	//创建监听退出chan
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, cross.StopSignals...)
	var (
		//path       string // 数据路径
		//logPath    string // 日志输出路径
		cronConfig string // 定时脚本
		cronTrue   bool   // 是否启用应用内定时器
		version    bool   // 显示版本号
	)
	flag.StringVar(&cronConfig, "cron_config", "0 0 17 * * ?", "pull code data cron")
	flag.BoolVar(&cronTrue, "cron_true", false, "use crontab in application")
	flag.IntVar(&source, "source", 0, "data source, default from tdx,1-dfcf")
	flag.BoolVar(&version, "version", false, "print version")
	flag.Parse()

	if version {
		fmt.Println(MinVersion)
		os.Exit(0)
	}

	if !cronTrue {
		handleCodeData()
	} else {
		crontab := cron.New(cron.WithSeconds()) //精确到秒
		// 添加定时任务,
		_, _ = crontab.AddFunc(cronConfig, handleCodeData)
		//启动定时器
		crontab.Start()
		select {
		case sig := <-c:
			{
				logger.Info("进程结束:", sig)
				os.Exit(1)
			}
		}
	}
}

func handleCodeData() {
	logger.Info("任务开始启动...")
	fullCodes := stock.GetCodeList()
	count := len(fullCodes)
	bar := progressbar.NewBar(0, "执行[更新历史数据]", count)
	for _, code := range fullCodes {
		basicInfo, err := security.GetBasicInfo(code)
		if err == security.ErrCodeNotExist {
			// 证券代码不存在
			bar.Add(1)
			continue
		}
		if err != nil {
			// 其它错误 输出错误信息
			logger.Errorf("%s => %v", code, err)
			bar.Add(1)
			continue
		}
		if basicInfo.Delisting {
			logger.Errorf("%s => 已退市", code)
			bar.Add(1)
			continue
		}
		bar.Add(1)
		listTimestamp := int64(basicInfo.ListTimestamp)
		var e = 0
		if source == 1 {
			e = pullData_em(code, internal.UnixTime(listTimestamp))
		} else {
			e = pullData_tdx(code, internal.UnixTime(listTimestamp))
		}
		if e&category.D_ENEED == 0 {
			sleep()
		}
	}
	logger.Info("任务执行完毕.", time.Now())
}

func sleep() {
	// 休眠2秒
	//time.Sleep(time.Second * 2)
}

// 拉取数据
func pullData_em(fc string, listTime time.Time) int {
	ks, err := dfcf.A(fc)
	if err != nil {
		_ = fmt.Errorf("error :%v", err.Error())
		return category.D_ENET
	}
	ToCSV(fc, ks)
	return 0
}

func pullData_tdx(fc string, listTime time.Time) int {
	df := tdx.GetKLineAll(fc)
	filename := cache.KLineFilename(fc)
	_ = df.WriteCSV(filename)
	return 0
}

func ToCSV(code string, ks []dfcf.KLine) {
	filename := cache.GetKLineFilename(code)
	count := len(ks)
	wrote := 0
	//fw, _ := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, category.CACHE_FILE_MODE)
	fw, _ := os.OpenFile(filename, category.CACHE_UPDATE, category.CACHE_FILE_MODE)
	_writer := csv.NewWriter(fw)
	if count > 0 {
		var cskHead dfcf.KLine
		err := cskHead.Init(_writer)
		if err != nil {
			logger.Errorf("code[%s]: 写日线文件, failed: ", code, err)
			return
		}
	}
	for j := 0; j < count; j++ {
		kl := ks[j]
		err := kl.WriteCSV(_writer)
		if err != nil {
			logger.Errorf("code[%s]: 写日线文件, failed: ", code, err)
			return
		}
		wrote += 1
	}
	_writer.Flush()
	api.CloseQuietly(fw)
	if wrote > 0 {
		logger.Infof("code[%s]: 写日线文件, SUCCESS", code)
	}
}
