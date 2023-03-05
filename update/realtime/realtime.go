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
	source   = 0  // 数据源
	batchMax = 10 // 批量最大100
)

// 更新日线数据工具
func main() {
	//创建监听退出chan
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, cross.StopSignals...)
	var (
		cronConfig string // 定时脚本
		cronTrue   bool   // 是否启用应用内定时器
	)
	flag.StringVar(&cronConfig, "cron_config", "0 */2 9-15 * * ?", "pull code data cron")
	flag.BoolVar(&cronTrue, "cron_true", false, "use crontab in application")
	flag.IntVar(&source, "source", 0, "data source, default from tdx,1-dfcf")
	flag.Parse()
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
	fullCodes := tdx.GetZxgList()
	fmt.Printf("自选股, 共计[%d]\n", len(fullCodes))
	if len(fullCodes) == 0 {
		fmt.Printf("没有指定自选股, 全量更新\n")
		fullCodes = category.GetCodeList()
	}
	fmt.Printf("实时更新指数及个股, 共计[%d]\n", len(fullCodes))
	count := len(fullCodes)
	bar := progressbar.NewBar(0, "执行[实时更新日线数据]", count)
	total := len(fullCodes)
	for start := 0; start < total; start += batchMax {
		codes := []string{}
		length := total - start
		if length >= batchMax {
			length = batchMax
		}
		for i := 0; i < length; i++ {
			code := fullCodes[start+i]
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
			codes = append(codes, code)
		}
		logger.Infof("%+v", codes)
		if len(codes) == 0 {
			continue
		}
		var e = 0
		if source == 1 {
			e = emRealTime(codes)
		} else {
			e = tdxRealTime(codes)
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
//
//	Deprecated: 告警, 东方财富的即时行情未实现
func emRealTime(codes []string) int {
	//TODO:未实现
	return 0
}

func tdxRealTime(codes []string) int {
	tdx.BatchRealtime(codes)
	return 0
}
