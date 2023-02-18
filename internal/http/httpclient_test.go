package http

import (
	"fmt"
	"testing"
)

func TestRequest(t *testing.T) {
	url := "https://baidu.com"
	data, err := Get(url)
	if err != nil {
		return
	}
	fmt.Println(data)
}
