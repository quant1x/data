package date

import (
	"golang.org/x/exp/slices"
	"time"
)

const (
	kAMBegin          = "09:30"
	kAMEnd            = "11:30"
	kPMBegin          = "13:00"
	kPMEnd            = "15:00"
	kTimeMinute       = "15:04"
	BEGIN_A_AM_HOUR   = 9  // A股开市-时
	BEGIN_A_AM_MINUTE = 30 // A股开市-分
	END_A_AM_HOUR     = 11 // A股休市-时
	END_A_AM_MINUTE   = 30 // A股休市-分
	BEGIN_A_PM_HOUR   = 13 // A股开市-时
	BEGIN_A_PM_MINUTE = 0  // A股开市-分
	END_A_PM_HOUR     = 15 // A股休市-时
	END_A_PM_MINUTE   = 0  // A股休市-分
)

type TimeRange struct {
	Begin time.Time
	End   time.Time
}

var (
	cnTimeRange     []TimeRange // 交易时间范围
	trAMBegin       time.Time   // 上午开盘时间
	trAMEnd         time.Time
	trPMBegin       time.Time
	trPMEnd         time.Time
	CN_FULL_MINUTES = 0 // A股全天交易的分钟数
)

func init() {
	now := time.Now()
	trAMBegin = time.Date(now.Year(), now.Month(), now.Day(), BEGIN_A_AM_HOUR, BEGIN_A_AM_MINUTE, 0, 0, time.Local)
	trAMEnd = time.Date(now.Year(), now.Month(), now.Day(), END_A_AM_HOUR, END_A_AM_MINUTE, 0, 0, time.Local)
	tr_am := TimeRange{
		Begin: trAMBegin,
		End:   trAMEnd,
	}
	cnTimeRange = append(cnTimeRange, tr_am)

	trPMBegin = time.Date(now.Year(), now.Month(), now.Day(), BEGIN_A_PM_HOUR, BEGIN_A_PM_MINUTE, 0, 0, time.Local)
	trPMEnd = time.Date(now.Year(), now.Month(), now.Day(), END_A_PM_HOUR, END_A_PM_MINUTE, 0, 0, time.Local)
	tr_pm := TimeRange{
		Begin: trPMBegin,
		End:   trPMEnd,
	}
	_minutes := 0
	cnTimeRange = append(cnTimeRange, tr_pm)
	for _, v := range cnTimeRange {
		_minutes += int(v.End.Sub(v.Begin).Minutes())
	}
	CN_FULL_MINUTES = _minutes
}

func fixMinute(m time.Time) time.Time {
	return time.Date(m.Year(), m.Month(), m.Day(), m.Hour(), m.Minute(), 0, 0, time.Local)
}

// Minutes 分钟数
func Minutes(date ...string) int {
	lastDay := LastTradeDate()
	today := IndexToday()
	if len(date) > 0 {
		today = fixTradeDate(date[0])
	}
	if today != lastDay {
		return CN_FULL_MINUTES
	}
	tm := time.Now()
	//tm, _ = utils.ParseTime("2023-04-11 09:29:00")
	//tm, _ = utils.ParseTime("2023-04-11 09:30:00")
	//tm, _ = utils.ParseTime("2023-04-11 09:31:00")
	//tm, _ = utils.ParseTime("2023-04-11 11:31:00")
	//tm, _ = utils.ParseTime("2023-04-11 12:59:00")
	//tm, _ = utils.ParseTime("2023-04-11 13:00:00")
	//tm, _ = utils.ParseTime("2023-04-11 13:01:00")
	//tm, _ = utils.ParseTime("2023-04-11 14:59:00")
	//tm, _ = utils.ParseTime("2023-04-11 15:01:00")
	tm = fixMinute(tm)

	tr := slices.Clone(cnTimeRange)
	//tr = stat.Reverse(tr)
	var last time.Time
	for _, v := range tr {
		if tm.Before(v.Begin) {
			last = v.Begin
			break
		}
		if tm.After(v.End) {
			last = v.End
			continue
		}
		//if !tm.After(v.Begin) {
		//	last = v.Begin
		//	break
		//}
		//if !tm.Before(v.End) {
		//	last = v.End
		//	continue
		//}
		last = tm
		break
	}

	m := int(last.Sub(trAMBegin).Minutes())
	if !last.Before(trPMBegin) {
		m -= int(trPMBegin.Sub(trAMEnd).Minutes())
	}
	return m
}
