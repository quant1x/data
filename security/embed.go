package security

import (
	"embed"
	"fmt"
	"io/fs"
)

var (
	// ResourcesPath 资源路径
	ResourcesPath = "resources"
)

//go:embed resources/*
var resources embed.FS

// OpenEmbed 打开嵌入式文件
func OpenEmbed(name string) (fs.File, error) {
	filename := fmt.Sprintf("%s/%s", ResourcesPath, name)
	reader, err := resources.Open(filename)
	if err != nil {
		return nil, err
	}
	return reader, nil
}
