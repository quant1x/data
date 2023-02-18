package category

const (
	D_OK    = 0x00000000 // 数据正常
	D_ERROR = 0x40000000 // 数据错误
	D_ECODE = 0x00000001 // 代码错误
	D_ENET  = 0x00000002 // 网络异常
	D_EDATA = 0x00000004 // 数据错误
	D_EDISK = 0x00000008 // 写文件错误
	D_ENEED = 0x00000010 // 不需要更新
)
