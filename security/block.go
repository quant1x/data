package security

import (
	"bufio"
	"gitee.com/quant1x/data/cache"
	"gitee.com/quant1x/data/category"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/pandas"
	"gitee.com/quant1x/pandas/stat"
	"github.com/mymmsc/gox/api"
	"github.com/mymmsc/gox/encoding/binary/struc"
	"github.com/mymmsc/gox/text/encoding"
	"golang.org/x/exp/slices"
	"io"
	"os"
	"strings"
)

type BlockType int

const (
	//BK_UNKNOWN BlockType = iota
	BK_HANGYE  BlockType = 2  // 行业
	BK_DIQU    BlockType = 3  // 地区
	BK_GAINIAN BlockType = 4  // 概念
	BK_FENGGE  BlockType = 5  // 风格
	BK_ZHISHU  BlockType = 6  // 指数
	BK_YJHY    BlockType = 12 // 研究行业

	BKN_HANGYE  = "行业"
	BKN_DIQU    = "地区"
	BKN_GAINIAN = "概念"
	BKN_FENGGE  = "风格"
	BKN_ZHISHU  = "指数"
	BKN_YJHY    = "研究行业"
)

var (
	kMapBlock = map[BlockType]string{
		BK_HANGYE:  BKN_HANGYE,
		BK_DIQU:    BKN_DIQU,
		BK_GAINIAN: BKN_GAINIAN,
		BK_FENGGE:  BKN_FENGGE,
		BK_ZHISHU:  BKN_ZHISHU,
		BK_YJHY:    BKN_YJHY,
	}
)

type __blockHeader struct {
	Unknown [384]byte `struc:"[384]byte,little"`
	Count   uint16    `struc:"uint16,little"`
}
type __raw_block_info struct {
	BlockName string             `struc:"[9]byte,little"`             // 板块名称
	Num       uint16             `struc:"uint16,little"`              // 个股数量
	BlockType uint16             `struc:"uint16,little"`              // 板块类型
	List      [400]__block_stock `struct:"[400]__block_stock,little"` // 个股列表
}

type __block_stock struct {
	Code string `struc:"[7]byte,little"` // 证券代码
}

type __raw_block_data struct {
	//Header blockHeader `struc:"[386]byte,little"`
	Unknown [384]byte          `struc:"[384]byte,little"`          // 头信息, 忽略
	Count   uint16             `struc:"uint16,little,sizeof=Data"` // 板块数量
	Data    []__raw_block_info `struc:"[2813]byte, little"`        // 板块数据
}

func get_block_file(blockFilename string) *__raw_block_data {
	fn := cache.GetBkPath() + "/" + blockFilename
	file, err := os.Open(fn)
	if err != nil {
		return nil
	}
	defer api.CloseQuietly(file)
	var block __raw_block_data
	err = struc.Unpack(file, &block)
	if err != nil {
		return nil
	}
	decoder := encoding.NewDecoder("GBK")
	for i, v := range block.Data {
		name := decoder.ConvertString(v.BlockName)
		block.Data[i].BlockName = strings.ReplaceAll(name, string([]byte{0x00}), "")
		for j, s := range v.List {
			block.Data[i].List[j].Code = strings.ReplaceAll(s.Code, string([]byte{0x00}), "")
		}
	}
	return &block
}

// BlockInfo 板块信息
type BlockInfo struct {
	Name  string   // 名称
	Code  string   // 代码
	Type  int      // 类型
	Count int      // 个股数量
	Block string   // 通达信板块编码
	List  []string // 代码列表
}

// HyInfo 行业板块对应
type HyInfo struct {
	Code   string // 股票代码
	Block  string // 通达信板块代码
	Block5 string // 通达信板块代码
}

func get_zs_file(name string) []BlockInfo {
	cacheFilename := cache.GetBkPath() + "/" + name
	if !cache.FileExist(cacheFilename) {
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
			Type:  int(stat.AnyToInt32(arr[2])),
			Block: arr[5],
		}
		blocks = append(blocks, bk)
	}
	return blocks
}

// 板块和板块名称对应
func get_zs_blocks() []BlockInfo {
	bks := []string{"tdxzs.cfg", "tdxzs3.cfg"}
	bis := []BlockInfo{}
	tmpMap := map[string]bool{}
	for _, v := range bks {
		bi := get_zs_file(v)
		if len(bi) == 0 {
			continue
		}
		for _, info := range bi {
			if _, ok := tmpMap[info.Code]; !ok {
				//if info.Code == "880482" {
				//	fmt.Println(info)
				//}
				bis = append(bis, info)
				tmpMap[info.Code] = true
			}
		}
	}
	return bis
}

// 获取行业板块
func get_hy_blocks() []HyInfo {
	hyfile := "tdxhy.cfg"
	name := hyfile
	cacheFilename := cache.GetBkPath() + "/" + name
	if !cache.FileExist(cacheFilename) {
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
	var hys = []HyInfo{}
	for {
		data, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		line := decoder.ConvertString(string(data))
		arr := strings.Split(line, "|")
		bc := arr[2]
		bc5 := bc
		if len(bc5) >= 5 {
			bc5 = bc5[0:5]
		}
		hy := HyInfo{
			Code:   arr[1],
			Block:  bc,
			Block5: bc5,
		}
		hys = append(hys, hy)
	}
	return hys
}

// 从行业信息中提取股票代码列表
func get_stock_list(hys []HyInfo, block string) []string {
	list := []string{}
	for _, v := range hys {
		if len(block) == 5 {
			if v.Block5 == block {
				list = append(list, v.Code)
			}
		} else {
			if v.Block == block {
				list = append(list, v.Code)
			}
		}
	}
	slices.Sort(list)
	return list
}

// 读取板块数据
func genBlockFile() {
	blockInfos := get_zs_blocks()
	//bks := []string{"block.dat", "block_gn.dat", "block_fg.dat", "block_zs.dat"}
	bks := []string{"block_gn.dat", "block_fg.dat", "block_zs.dat"}
	name2block := map[string]__raw_block_info{}
	for _, v := range bks {
		bi := get_block_file(v)
		if bi != nil {
			for _, bk := range (*bi).Data {
				name2block[bk.BlockName] = bk
			}
		}
	}
	// 行业代码映射
	code2hy := map[string]string{}
	for _, v := range blockInfos {
		if v.Name != v.Block {
			code2hy[v.Block] = v.Name
		}
	}
	// 行业板块数据
	hys := get_hy_blocks()
	for i, v := range blockInfos {
		bn := v.Name
		__info, ok := name2block[bn]
		if ok {
			list := []string{}
			for _, sc := range __info.List {
				if len(sc.Code) < 5 {
					continue
				}
				marketId, _, _ := category.DetectMarket(sc.Code)
				if marketId == category.MARKET_ID_BEIJING {
					continue
				}
				list = append(list, sc.Code)
			}
			blockInfos[i].Count = int(__info.Num)
			blockInfos[i].List = list
			continue
		}
		bc := v.Block
		stockList := get_stock_list(hys, bc)
		if len(stockList) > 0 {
			blockInfos[i].Count = len(stockList)
			blockInfos[i].List = stockList
		}
	}

	bk_code := []string{}
	bk_name := []string{}
	bk_type := []int{}
	for _, v := range blockInfos {
		if v.Count == 0 {
			continue
		}
		bk_stock := v.List
		codes := pandas.NewSeries(stat.SERIES_TYPE_STRING, "code", bk_stock)
		tmp := pandas.NewDataFrame(codes)
		_ = tmp.WriteCSV(cache.GetBkPath() + "/" + v.Code + ".csv")
		bk_code = append(bk_code, v.Code)
		bk_name = append(bk_name, v.Name)
		bk_type = append(bk_type, v.Type)
	}
	bkc := pandas.NewSeries(stat.SERIES_TYPE_STRING, "code", bk_code)
	bkn := pandas.NewSeries(stat.SERIES_TYPE_STRING, "name", bk_name)
	bkt := pandas.NewSeries(stat.SERIES_TYPE_STRING, "type", bk_type)
	dfBk := pandas.NewDataFrame(bkc, bkn, bkt)
	_ = dfBk.WriteCSV(cache.GetBkPath() + "/block.csv")
}

func init() {
	// 如果板块数据不存在, 从应用内导出
	blockFile := cache.GetBkPath() + "/block.csv"
	reCreate := false
	if !cache.FileExist(blockFile) {
		reCreate = true
	}
	if !reCreate {
		blockData := cache.GetBkPath() + "/" + quotes.BLOCK_DEFAULT
		stat_, err := os.Stat(blockData)
		if err == nil || os.IsExist(err) {
			dataModifyTime := stat_.ModTime()
			bk881266 := cache.GetBkPath() + "/" + "881266.csv"
			stat_bk, err := os.Stat(bk881266)
			if err == nil || os.IsExist(err) {
				if stat_bk.ModTime().Before(dataModifyTime) {
					reCreate = true
				}
			} else {
				reCreate = true
			}
		} else {
			reCreate = true
		}
	}
	if reCreate {
		genBlockFile()
	}
}

// BlockTypeNameByCode 通过板块类型代码获取板块类型名称
func BlockTypeNameByCode(blockCode int) (name string, ok bool) {
	blockType := BlockType(blockCode)
	bkName, found := kMapBlock[blockType]
	return bkName, found
}
