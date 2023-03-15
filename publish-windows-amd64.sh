#!/bin/sh
# 获取当前路径, 用于返回
p0=`pwd`
# 获取脚本所在路径, 防止后续操作在非项目路径
p1=$(cd $(dirname $0);pwd)

BIN=~/runtime/bin
PREFIX=stock-update

version=$(git describe --tags `git rev-list --tags --max-count=1`)
version=${version:1}
echo "version: ${version}"
# windows amd64
GOOS=windows
GOARCH=amd64
apps=("kline" "realtime" "xdxr" "zxg" "tick" "snapshot")
for app in ${apps[@]}
do
  echo "正在编译应用:${app} => $BIN/$PREFIX-${app}..."
  env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-X 'main.MinVersion=${version}'" -o $BIN/$PREFIX-${app} gitee.com/quant1x/data/update/${app}
  echo "正在编译应用:${app} => $BIN/$PREFIX-${app}...OK"
done
cd $p0