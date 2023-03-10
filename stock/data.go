package stock

import (
	"gitee.com/quant1x/data/cache"
	"gitee.com/quant1x/data/internal/tdx"
	"gitee.com/quant1x/pandas"
	"gitee.com/quant1x/pandas/stat"
)

// KLine 加载K线
func KLine(code string) pandas.DataFrame {
	return tdx.GetCacheKLine(code, true)
}

// Tick 加载一个时间范围内的tick缓存数据
func Tick(code string, dates []string) pandas.DataFrame {
	buyVolumes := []stat.DType{}
	sellVolumes := []stat.DType{}
	buyAmounts := []stat.DType{}
	sellAmounts := []stat.DType{}
	inflowVolumes := []stat.DType{}
	inflowAmounts := []stat.DType{}
	buySpeeds := []stat.DType{}
	sellSpeeds := []stat.DType{}
	buyFront := stat.DType(0)
	sellFront := stat.DType(0)
	for _, date := range dates {
		buyVolume := stat.DType(0)
		sellVolume := stat.DType(0)
		buyAmount := stat.DType(0)
		sellAmount := stat.DType(0)
		tmp := TickByDate(code, date)
		if tmp.Nrow() == 0 {
			// 数据有缺失跳过
			return pandas.DataFrame{}
		}
		if tmp.Nrow() > 0 {
			for i := 0; i < tmp.Nrow(); i++ {
				m := tmp.IndexOf(i)
				t := stat.AnyToInt32(m["buyorsell"])
				p := stat.AnyToFloat64(m["price"])
				v := stat.AnyToFloat64(m["vol"])
				if t == 1 {
					// 卖出
					sellVolume += v
					sellAmount += v * p * 100
				} else {
					buyVolume += v
					buyAmount += v * p * 100
				}
			}
		}
		buyVolumes = append(buyVolumes, buyVolume)
		sellVolumes = append(sellVolumes, sellVolume)
		buyAmounts = append(buyAmounts, buyAmount)
		sellAmounts = append(sellAmounts, sellAmount)

		buyInflow := buyVolume - buyFront
		buyFront = buyVolume
		sellInflow := sellVolume - sellFront
		sellFront = sellVolume

		inflowVolumes = append(inflowVolumes, buyVolume-sellVolume)
		inflowAmounts = append(inflowAmounts, buyAmount-sellAmount)

		buySpeeds = append(buySpeeds, buyInflow)
		sellSpeeds = append(sellSpeeds, sellInflow)
	}
	dt := pandas.NewSeries(stat.SERIES_TYPE_STRING, "date", dates)
	bv := pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "bv", buyVolumes)
	sv := pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "sv", sellVolumes)
	ba := pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "ba", buyAmounts)
	sa := pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "sa", sellAmounts)
	iv := pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "iv", inflowVolumes)
	ia := pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "ia", inflowAmounts)
	bs := pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "bs", buySpeeds)
	ss := pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "ss", sellSpeeds)
	df := pandas.NewDataFrame(dt, bv, sv, ba, sa, iv, ia, bs, ss)
	return df
}

func TickByDate(code string, date string) pandas.DataFrame {
	date = cache.CorrectDate(date)
	var df pandas.DataFrame
	filename := cache.TickFilename(code, date)
	if !cache.FileExist(filename) {
		return df
	}
	df = pandas.ReadCSV(filename)
	return df
}

func BlockList() pandas.DataFrame {
	bkListFile := cache.BlockFilename()
	df := pandas.ReadCSV(bkListFile)
	return df
}
