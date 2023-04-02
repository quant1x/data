package main

import (
	"flag"
	"fmt"
	"gitee.com/quant1x/data/cache"
	"gitee.com/quant1x/data/category"
	"gitee.com/quant1x/data/internal"
	"gitee.com/quant1x/data/internal/tdx"
	"gitee.com/quant1x/data/security"
	"gitee.com/quant1x/data/stock"
	"github.com/mymmsc/gox/cron"
	"github.com/mymmsc/gox/logger"
	"github.com/mymmsc/gox/progressbar"
	"github.com/mymmsc/gox/signal"
	"os"
	"time"
)

var (
	MinVersion string // 版本号
)

// 更新tick数据
func main() {
	//创建监听退出chan
	c := signal.Notify()
	var (
		cronConfig string // 定时脚本
		cronTrue   bool   // 是否启用应用内定时器
		startDate  string // 开始日期
		version    bool   // 显示版本号
	)
	flag.StringVar(&cronConfig, "cron_config", "0 0 17 * * ?", "pull code data cron")
	flag.BoolVar(&cronTrue, "cron_true", false, "use crontab in application")
	flag.StringVar(&startDate, "start", cache.TickStartDate, "tick start date")
	flag.BoolVar(&version, "version", false, "print version")
	flag.Parse()

	if version {
		fmt.Println(MinVersion)
		os.Exit(0)
	}

	cache.UpdateTickStartDate(startDate)
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
		fullCodes = stock.GetCodeList()
	}
	fmt.Printf("实时更新指数及个股, 共计[%d]\n", len(fullCodes))

	count := len(fullCodes)
	bar := progressbar.NewBar(1, "执行[更新历史tick数据]", count)
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
		e := pullData(code, internal.UnixTime(listTimestamp))
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
func pullData(fc string, listTime time.Time) int {
	tdx.GetTickAll(fc)
	_ = listTime
	return 0
}
