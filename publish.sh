#!/bin/sh
# 获取当前路径, 用于返回
p0=`pwd`
# 获取脚本所在路径, 防止后续操作在非项目路径
p1=$(cd $(dirname $0);pwd)

# windows amd64
env GOOS=windows GOARCH=amd64 go build -o bin/kline-win-amd64.exe gitee.com/quant1x/data/update/kline
# darwin amd64
env GOOS=darwin GOARCH=amd64 go build -o bin/kline-mac-amd64 gitee.com/quant1x/data/update/kline
cd $p0