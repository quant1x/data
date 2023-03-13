package main

import (
	"flag"
	"fmt"
	"gitee.com/quant1x/data/category"
	"gitee.com/quant1x/data/internal/tdx"
	"gitee.com/quant1x/data/security"
	"gitee.com/quant1x/data/stock"
	"gitee.com/quant1x/data/update/cross"
	"gitee.com/quant1x/gotdx/quotes"
	"github.com/mymmsc/gox/api"
	"github.com/mymmsc/gox/cron"
	"github.com/mymmsc/gox/logger"
	"github.com/mymmsc/gox/progressbar"
	"os"
	"os/signal"
	"time"
)

const (
	application = "板块指数"
)

var (
	source     = 0                              // 数据源
	batchMax   = quotes.TDX_SECURITY_QUOTES_MAX // 批量最大100
	MinVersion string                           // 版本号
)

// 更新快照数据工具
func main() {
	//创建监听退出chan
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, cross.StopSignals...)
	var (
		cronConfig string // 定时脚本
		cronTrue   bool   // 是否启用应用内定时器
		version    bool   // 显示版本号
	)
	flag.StringVar(&cronConfig, "cron_config", "0 * 9-15 * * ?", "pull code data cron")
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
				logger.Info(application+"进程结束:", sig)
				os.Exit(1)
			}
		}
	}
}

func handleCodeData() {
	logger.Info(application + "任务开始启动...")
	blockList := stock.BlockList()
	fmt.Printf("板块, 共计[%d]\n", blockList.Nrow())
	if blockList.Nrow() == 0 {
		fmt.Printf("没有指定板块列表, 任务结束\n")
		return
	}
	CODE := blockList.Col("code").Strings()
	//NAME := blockList.Col("name").Strings()
	count := blockList.Nrow()
	bar := progressbar.NewBar(0, "执行[实时更新板块数据]", count)
	for start := 0; start < count; start += batchMax {
		codes := []string{}
		length := count - start
		if length >= batchMax {
			length = batchMax
		}
		for i := 0; i < length; i++ {
			code := CODE[start+i]
			if api.StartsWith(code, []string{"88"}) {
				code = "sh" + code
			}
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
		e := tdxSnapshot(codes)
		if e&category.D_OK != 0 {
			sleep()
		}
	}
	logger.Info(application+"任务执行完毕.", time.Now())
}

func sleep() {
	// 休眠2秒
	time.Sleep(time.Second * 2)
}

func tdxSnapshot(codes []string) int {
	tdx.BatchSnapShot(codes)
	return 0
}
