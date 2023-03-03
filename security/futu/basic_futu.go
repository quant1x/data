package futu

type QotMarket int32

const (
	QotMarket_QotMarket_Unknown       QotMarket = 0  //未知市场
	QotMarket_QotMarket_HK_Security   QotMarket = 1  //香港市场
	QotMarket_QotMarket_HK_Future     QotMarket = 2  //港期货(已废弃，使用QotMarket_HK_Security即可)
	QotMarket_QotMarket_US_Security   QotMarket = 11 //美国市场
	QotMarket_QotMarket_CNSH_Security QotMarket = 21 //沪股市场
	QotMarket_QotMarket_CNSZ_Security QotMarket = 22 //深股市场
	QotMarket_QotMarket_SG_Security   QotMarket = 31 //新加坡市场
	QotMarket_QotMarket_JP_Security   QotMarket = 41 //日本市场
)

// Enum value maps for QotMarket.
var (
	QotMarket_name = map[int32]string{
		0:  "QotMarket_Unknown",
		1:  "QotMarket_HK_Security",
		2:  "QotMarket_HK_Future",
		11: "QotMarket_US_Security",
		21: "QotMarket_CNSH_Security",
		22: "QotMarket_CNSZ_Security",
		31: "QotMarket_SG_Security",
		41: "QotMarket_JP_Security",
	}
	QotMarket_value = map[string]int32{
		"QotMarket_Unknown":       0,
		"QotMarket_HK_Security":   1,
		"QotMarket_HK_Future":     2,
		"QotMarket_US_Security":   11,
		"QotMarket_CNSH_Security": 21,
		"QotMarket_CNSZ_Security": 22,
		"QotMarket_SG_Security":   31,
		"QotMarket_JP_Security":   41,
	}
)

type SecurityType int32

const (
	SecurityType_SecurityType_Unknown  SecurityType = 0  //未知
	SecurityType_SecurityType_Bond     SecurityType = 1  //债券
	SecurityType_SecurityType_Bwrt     SecurityType = 2  //一揽子权证
	SecurityType_SecurityType_Eqty     SecurityType = 3  //正股
	SecurityType_SecurityType_Trust    SecurityType = 4  //信托,基金
	SecurityType_SecurityType_Warrant  SecurityType = 5  //窝轮
	SecurityType_SecurityType_Index    SecurityType = 6  //指数
	SecurityType_SecurityType_Plate    SecurityType = 7  //板块
	SecurityType_SecurityType_Drvt     SecurityType = 8  //期权
	SecurityType_SecurityType_PlateSet SecurityType = 9  //板块集
	SecurityType_SecurityType_Future   SecurityType = 10 //期货
)

// Enum value maps for SecurityType.
var (
	SecurityType_name = map[int32]string{
		0:  "SecurityType_Unknown",
		1:  "SecurityType_Bond",
		2:  "SecurityType_Bwrt",
		3:  "SecurityType_Eqty",
		4:  "SecurityType_Trust",
		5:  "SecurityType_Warrant",
		6:  "SecurityType_Index",
		7:  "SecurityType_Plate",
		8:  "SecurityType_Drvt",
		9:  "SecurityType_PlateSet",
		10: "SecurityType_Future",
	}
	SecurityType_value = map[string]int32{
		"SecurityType_Unknown":  0,
		"SecurityType_Bond":     1,
		"SecurityType_Bwrt":     2,
		"SecurityType_Eqty":     3,
		"SecurityType_Trust":    4,
		"SecurityType_Warrant":  5,
		"SecurityType_Index":    6,
		"SecurityType_Plate":    7,
		"SecurityType_Drvt":     8,
		"SecurityType_PlateSet": 9,
		"SecurityType_Future":   10,
	}

	StockType = map[string]int32{
		"SecurityType_Unknown":  0,
		"SecurityType_Bond":     1,
		"SecurityType_Bwrt":     2,
		"STOCK":                 int32(SecurityType_SecurityType_Eqty),
		"SecurityType_Trust":    4,
		"SecurityType_Warrant":  5,
		"INDEX":                 int32(SecurityType_SecurityType_Index),
		"SecurityType_Plate":    7,
		"SecurityType_Drvt":     8,
		"SecurityType_PlateSet": 9,
		"SecurityType_Future":   10,
	}
)
