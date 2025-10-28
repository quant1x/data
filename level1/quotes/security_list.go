package quotes

import (
	"bytes"
	"encoding/binary"

	"gitee.com/quant1x/data/level1/internal"
	"gitee.com/quant1x/data/level1/proto"
)

const (
	// SECURITY_LIST_A_PER_REQUEST_MAX 单次最大获取多少条股票数据
	SECURITY_LIST_A_PER_REQUEST_MAX = 1600
)

// SecurityListAPackage 股票列表
type SecurityListAPackage struct {
	reqHeader  *StdRequestHeader
	respHeader *StdResponseHeader
	request    *SecurityListARequest
	reply      *SecurityListAReply

	contentHex string
}

type SecurityListARequest struct {
	Market  uint16
	Start   uint32
	Count   uint32
	Unknown uint32 // 未知字段, 固定0
}

type SecurityListAReply struct {
	Count uint16
	List  []Security
}

// type SecurityA struct {
// 	Code      string
// 	VolUnit   uint16
// 	Reversed1 [4]byte `dataframe:"-"`
// 	//R1           uint32
// 	//P1           float64
// 	DecimalPoint int8
// 	Name         string
// 	IgnoreString string
// 	PreClose     float64
// 	Reversed2    [4]byte `dataframe:"-"`
// 	//R2           uint32
// 	//P2           float64
// }

func NewSecurityListAPackage() *SecurityListAPackage {
	obj := new(SecurityListAPackage)
	obj.reqHeader = new(StdRequestHeader)
	obj.respHeader = new(StdResponseHeader)
	obj.request = new(SecurityListARequest)
	obj.reply = new(SecurityListAReply)

	obj.reqHeader.ZipFlag = proto.FlagNotZipped
	obj.reqHeader.SeqID = internal.SequenceId()
	obj.reqHeader.PacketType = 0x01
	obj.reqHeader.Method = proto.STD_MSG_SECURITY_LIST_A
	return obj
}

func (obj *SecurityListAPackage) SetParams(req *SecurityListARequest) {
	obj.request = req
}

func (obj *SecurityListAPackage) Serialize() ([]byte, error) {
	obj.reqHeader.PkgLen1 = 2 + 2 + 4 + 4 + 4
	obj.reqHeader.PkgLen2 = 2 + 2 + 4 + 4 + 4

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, obj.reqHeader)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.LittleEndian, obj.request)

	//b, err := hex.DecodeString(obj.contentHex)
	//buf.Write(b)

	//err = binary.Write(buf, binary.LittleEndian, uint16(len(obj.stocks)))

	return buf.Bytes(), err
}

func (obj *SecurityListAPackage) UnSerialize(header interface{}, data []byte) error {
	obj.respHeader = header.(*StdResponseHeader)

	pos := 0
	err := binary.Read(bytes.NewBuffer(data[pos:pos+2]), binary.LittleEndian, &obj.reply.Count)
	pos += 2
	for index := uint16(0); index < obj.reply.Count; index++ {
		ele := Security{}
		var code [6]byte
		_ = binary.Read(bytes.NewBuffer(data[pos:pos+6]), binary.LittleEndian, &code)
		pos += 6
		ele.Code = string(code[:])

		_ = binary.Read(bytes.NewBuffer(data[pos:pos+2]), binary.LittleEndian, &ele.VolUnit)
		pos += 2

		var name [8]byte
		_ = binary.Read(bytes.NewBuffer(data[pos:pos+8]), binary.LittleEndian, &name)
		ele.Name = internal.Utf8ToGbk(name[:])
		pos += 8

		var ignore_string [8]byte
		_ = binary.Read(bytes.NewBuffer(data[pos:pos+8]), binary.LittleEndian, &ignore_string)
		//ele.IgnoreString = internal.Utf8ToGbk(ignore_string[:])
		pos += 8

		_ = binary.Read(bytes.NewBuffer(data[pos:pos+4]), binary.LittleEndian, &ele.Reversed1)
		pos += 4

		_ = binary.Read(bytes.NewBuffer(data[pos:pos+1]), binary.LittleEndian, &ele.DecimalPoint)
		pos += 1
		var rawPreClose uint32
		_ = binary.Read(bytes.NewBuffer(data[pos:pos+4]), binary.LittleEndian, &rawPreClose)
		ele.PreClose = internal.IntToFloat64(int(rawPreClose))
		pos += 4

		_ = binary.Read(bytes.NewBuffer(data[pos:pos+4]), binary.LittleEndian, &ele.Reversed2)
		//_ = binary.Read(bytes.NewBuffer(data[pos:pos+4]), binary.LittleEndian, &ele.R2)
		//ele.P2 = getVolume(int(ele.R2))
		pos += 4

		obj.reply.List = append(obj.reply.List, ele)
	}
	return err
}

func (obj *SecurityListAPackage) Reply() interface{} {
	return obj.reply
}
