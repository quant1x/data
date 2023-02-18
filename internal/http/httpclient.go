package http

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	httpclient "github.com/guonaihong/gout"
	"github.com/guonaihong/gout/dataflow"
	"github.com/guonaihong/gout/middler"
	"github.com/mymmsc/gox/api"
	"github.com/mymmsc/gox/logger"
	"io"
	"net/http"
	URL "net/url"
	"strings"
)

type Header = httpclient.H
type Query = httpclient.H

// ResponseHeader 用于解析 服务端 返回的http header
type ResponseHeader struct {
	ContentEncoding string `header:"Content-Encoding"`
}

// ResponseMiddler response拦截器修改示例
type ResponseMiddler struct{}

func (d *ResponseMiddler) ModifyResponse(response *http.Response) error {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	contentEncoding := response.Header.Get("Content-Encoding")
	var reader io.ReadCloser = nil
	if len(contentEncoding) > 0 {
		contentEncoding = strings.ToLower(contentEncoding)

		switch contentEncoding {
		case "gzip":
			reader, err = gzip.NewReader(bytes.NewBuffer(body))
			if err != nil {
				logger.Error(err)
				reader = nil
			}
		case "deflate":
			reader = flate.NewReader(bytes.NewReader(body))
		}
	}
	if reader != nil {
		defer api.CloseQuietly(reader)
		body, err = io.ReadAll(reader)
		if err != nil {
			return err
		}
	}
	response.Body = io.NopCloser(bytes.NewReader(body))
	return nil
}

func responseMiddler() *ResponseMiddler {
	return &ResponseMiddler{}
}

type Method string

const (
	get  Method = "GET"
	post Method = "POST"
)

// Get HTTP GET请求
func Get(url string, params ...any) ([]byte, error) {
	return httpRequest(url, get, params...)
}

func httpRequest(url string, method Method, q ...any) ([]byte, error) {
	u, err := URL.Parse(url)
	if err != nil {
		return nil, err
	}
	requestHeader := make(map[string]string)
	requestHeader["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"
	requestHeader["Accept-Encoding"] = "gzip, deflate"
	requestHeader["Accept-Language"] = "zh-CN,zh;q=0.9,en;q=0.8"
	requestHeader["Cache-Control"] = "no-cache"
	requestHeader["Connection"] = "keep-alive"
	requestHeader["Host"] = u.Host
	requestHeader["Pragma"] = "no-cache"
	requestHeader["Upgrade-Insecure-Requests"] = "1"
	requestHeader["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.84 Safari/537.36"

	responseHeader := ResponseHeader{}
	body := []byte{}
	var flow *dataflow.DataFlow
	client := httpclient.New()
	if post == method {
		flow = client.POST(url)
	} else {
		flow = client.GET(url)
	}
	err = flow.SetHost(u.Host).
		SetHeader(requestHeader).
		SetQuery(q...).
		BindHeader(&responseHeader).
		BindBody(&body).
		ResponseUse(
			responseMiddler(),
			middler.WithResponseMiddlerFunc(func(response *http.Response) error {
				return nil
			})).
		Do()
	if err != nil {
		return nil, err
	}
	return body, nil
}
