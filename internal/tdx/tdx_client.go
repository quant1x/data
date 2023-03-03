package tdx

import "gitee.com/quant1x/gotdx/quotes"

var (
	stdApi *quotes.StdApi = nil
)

func prepare() *quotes.StdApi {
	if stdApi == nil {
		api_, err := quotes.NewStdApi()
		if err != nil {
			return nil
		}
		stdApi = api_
	}
	return stdApi
}
