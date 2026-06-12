package level1

import (
	"sync"

	"github.com/quant1x/data/level1/quotes"
)

var (
	stdApi   *quotes.StdApi = nil
	apiMutex sync.Mutex
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
	apiMutex.Lock()
	defer apiMutex.Unlock()
	initStdApi()
	return stdApi
}

func ReOpen() {
	apiMutex.Lock()
	defer apiMutex.Unlock()
	if stdApi != nil {
		stdApi.Close()
	}
}
