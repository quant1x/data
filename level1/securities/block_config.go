package securities

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/quant1x/data/cache"
	"github.com/quant1x/gox/api"
	"github.com/quant1x/gox/text/encoding"
)

const (
	BLK_ZIP_FILENAME = "zhb.zip"
	BLK_ZS_FILENAME  = "tdxzs.cfg"
	BLK_ZS3_FILENAME = "tdxzs3.cfg"
)

var (
	need_blk_files = []string{
		BLK_ZS_FILENAME,
		BLK_ZS3_FILENAME,
	}
)

// 加载板块和板块名称对应
func loadIndexBlockInfos() []BlockInfo {
	bks := need_blk_files
	bis := []BlockInfo{}
	tmpMap := map[string]BlockInfo{}
	for _, v := range bks {
		bi := getBlockInfoFromConfig(v)
		if len(bi) == 0 {
			continue
		}
		for _, info := range bi {
			if bv, ok := tmpMap[info.Code]; !ok {
				bis = append(bis, info)
				tmpMap[info.Code] = info
			} else {
				_ = bv
			}
		}
	}
	return bis
}

func getBlockInfoFromConfig(name string) []BlockInfo {
	cacheFilename := cache.GetBlockPath() + "/" + name
	if !api.FileExist(cacheFilename) {
		// 如果文件不存在, 导出内嵌资源
		err := export(cacheFilename, name)
		if err != nil {
			return nil
		}
	}
	file, err := os.Open(cacheFilename)
	if err != nil {
		return nil
	}
	defer api.CloseQuietly(file)
	reader := bufio.NewReader(file)
	// 按行处理txt
	decoder := encoding.NewDecoder("GBK")
	var blocks = []BlockInfo{}
	for {
		data, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		line := decoder.ConvertString(string(data))
		arr := strings.Split(line, "|")
		bk := BlockInfo{
			Name:  arr[0],
			Code:  arr[1],
			Type:  int(api.ParseInt(arr[2])),
			Block: arr[5],
		}
		blocks = append(blocks, bk)
	}
	return blocks
}
