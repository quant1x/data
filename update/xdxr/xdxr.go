package main

import (
	"flag"
	"fmt"
	"gitee.com/quant1x/data/category"
	"gitee.com/quant1x/data/internal/tdx"
	"gitee.com/quant1x/data/security"
	"gitee.com/quant1x/data/update/cross"
	"github.com/mymmsc/gox/cron"
	"github.com/mymmsc/gox/logger"
	"github.com/mymmsc/gox/progressbar"
	"os"
	"os/signal"
	"time"
)

var (
	// source int //数据源
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
	fullCodes := category.GetCodeList()
	count := len(fullCodes)
	bar := progressbar.NewBar(0, "执行[更新除权除息数据]", count)
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
		tdx.GetXdxrInfo(code)
	}
	logger.Info("任务执行完毕.", time.Now())
}
