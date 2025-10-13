package level1

import (
	"sync"

	"gitee.com/quant1x/data/level1/quotes"
)

var (
	stdApi   *quotes.StdApi = nil
	tdxMutex sync.Mutex
)

func initStdApi() {
	if stdApi == nil {
		api_, err := quotes.NewStdApi()
		if err != nil {
			return
		}
		stdApi = api_
	}
}

func GetApi() *quotes.StdApi {
	tdxMutex.Lock()
	defer tdxMutex.Unlock()
	initStdApi()
	return stdApi
}

func ReOpen() {
	tdxMutex.Lock()
	defer tdxMutex.Unlock()
	if stdApi != nil {
		stdApi.Close()
	}
}
