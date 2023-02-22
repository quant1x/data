数据工具使用说明
===

1. 同步自选股到缓存目录
```shell
 zxg -path ~/workspace/data/tdx
```
2. *nix 定时任务设定
```shell
*/2 9-15 * * 1-5 XXXXX/stock-realtime-csv
30 17 * * 1-5 XXXXX/stock-update-csv
```