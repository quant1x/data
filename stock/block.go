package stock

import (
	"gitee.com/quant1x/data/cache"
	"gitee.com/quant1x/pandas"
	"gitee.com/quant1x/pandas/stat"
	"github.com/mymmsc/gox/api"
)

// BlockList 获取板块列表
func BlockList() pandas.DataFrame {
	bkListFile := cache.BlockFilename()
	df := pandas.ReadCSV(bkListFile)
	codes := df.Col("code").Strings()
	names := df.Col("name").Strings()
	types := df.Col("type").Ints()
	for i, v := range codes {
		if api.StartsWith(v, []string{"88"}) {
			codes[i] = "sh" + v
		}
	}
	oc := pandas.NewSeries(stat.SERIES_TYPE_STRING, "code", codes)
	on := pandas.NewSeries(stat.SERIES_TYPE_STRING, "name", names)
	ot := pandas.NewSeries(stat.SERIES_TYPE_INT32, "type", types)
	df = pandas.NewDataFrame(oc, on, ot)
	return df
}
