# Changelog
All notable changes to this project will be documented in this file.

## [Unreleased]

## [0.8.16] - 2023-03-16
### Changed
- 更新通达信工具库.
- 更新tdxzs*文件, 后续的处理可以手动更新.quant1x/bk/*.cfg.
- 更新tdx工具库版本.
- 增加windows arm64编译脚本.
- 编译时控制台输出版本号.
- !4 #I6MFH8 每天同步板块数据文件 * 调整快照数据工具 * 调整通达信工具类版本.
- 调整板块源文件名.
- 增加涨跌停板限制幅度函数.
- 板块代码去重.
- 批量获取快照数据, 返回快照结构的数组.
- 调整BlockList返回的code字段, 要加前缀.
- !3 #I6MFH4 增加板块的快照数据 * 增加windows打包脚本 * 调整快照的package * 增加实时行情的快照数据 * 计划调整板块列表的字段类型, 目前板块代码是int64 * 新增获取最后一个交易日的函数 * 修订实时数据函数内部变量名, 使变量名贴合意图 * 增加mac arm64的编译脚本 * 增加mac amd64的编译脚本 * 增加函数的注释 * 增加独立的mac打包脚本.
- 调整五档行情数据更新K线时, 需要区分指数和个股的字段不同.
- 区分个股和指数的K线数据, 自动修正本地缓存.
- 修订版本号.

## [0.8.15] - 2023-03-10
### Changed
- 删掉无用的代码.

## [0.8.14] - 2023-03-10
### Changed
- 修订板块代码.
- 修正板块代码.
- 修订板块代码.
- 修订板块代码.

## [0.8.13] - 2023-03-10
### Changed
- 修订板块代码.

## [0.8.12] - 2023-03-10
### Changed
- 修订版本代码.
- 修订板块代码.

## [0.8.11] - 2023-03-10
### Changed
- !2 #I6LR5W 增加板块数据 * add block data.
- Tick数据不全跳过.
- 调整版本.

## [0.8.10] - 2023-03-08
### Changed
- 修订依赖版本.

## [0.8.9] - 2023-03-08
### Changed
- 更新pandas以支持arm框架.
- 增加测试代码.

## [0.8.8] - 2023-03-05
### Changed
- 修复windows创建路径失败的bug.

## [0.8.7] - 2023-03-05
### Changed
- 修订交易日历的bug.

## [0.8.6] - 2023-03-05
### Changed
- 增加版本号.

## [0.8.5] - 2023-03-05
### Changed
- 梳理部分代码.

## [0.8.4] - 2023-03-05
### Changed
- 更新版本.
- 更新版本.

## [0.8.3] - 2023-03-05
### Changed
- 升级pandas版本.

## [0.8.2] - 2023-03-03
### Changed
- Tick返回增加日期.

## [0.8.1] - 2023-03-03
### Changed
- !1 #I6J74Z 增加成交数据按天统计接口 * 增加成交数据按天统计接口 * 修订部分代码.

## [0.8.0] - 2023-03-03
### Changed
- 删除废弃的代码.

## [0.7.22] - 2023-03-03
### Changed
- 重新构建 quant1x.data.

[Unreleased]: https://gitee.com/quant1x/data/compare/v1.0.5...HEAD
[1.0.5]: https://gitee.com/quant1x/data/compare/v1.0.4...v1.0.5
[1.0.4]: https://gitee.com/quant1x/data/compare/v1.0.3...v1.0.4
[1.0.3]: https://gitee.com/quant1x/data/compare/v1.0.2...v1.0.3
[1.0.2]: https://gitee.com/quant1x/data/compare/v1.0.1...v1.0.2
[1.0.1]: https://gitee.com/quant1x/data/compare/v1.0.0...v1.0.1
[1.0.0]: https://gitee.com/quant1x/data/compare/v0.9.28...v1.0.0
[0.9.28]: https://gitee.com/quant1x/data/compare/v0.9.27...v0.9.28
[0.9.27]: https://gitee.com/quant1x/data/compare/v0.9.26...v0.9.27
[0.9.26]: https://gitee.com/quant1x/data/compare/v0.9.25...v0.9.26
[0.9.25]: https://gitee.com/quant1x/data/compare/v0.9.24...v0.9.25
[0.9.24]: https://gitee.com/quant1x/data/compare/v0.9.23...v0.9.24
[0.9.23]: https://gitee.com/quant1x/data/compare/v0.9.22...v0.9.23
[0.9.22]: https://gitee.com/quant1x/data/compare/v0.9.21...v0.9.22
[0.9.21]: https://gitee.com/quant1x/data/compare/v0.9.20...v0.9.21
[0.9.20]: https://gitee.com/quant1x/data/compare/v0.9.19...v0.9.20
[0.9.19]: https://gitee.com/quant1x/data/compare/v0.9.18...v0.9.19
[0.9.18]: https://gitee.com/quant1x/data/compare/v0.9.17...v0.9.18
[0.9.17]: https://gitee.com/quant1x/data/compare/v0.9.16...v0.9.17
[0.9.16]: https://gitee.com/quant1x/data/compare/v0.9.15...v0.9.16
[0.9.15]: https://gitee.com/quant1x/data/compare/v0.9.14...v0.9.15
[0.9.14]: https://gitee.com/quant1x/data/compare/v0.9.13...v0.9.14
[0.9.13]: https://gitee.com/quant1x/data/compare/v0.9.12...v0.9.13
[0.9.12]: https://gitee.com/quant1x/data/compare/v0.9.11...v0.9.12
[0.9.11]: https://gitee.com/quant1x/data/compare/v0.9.10...v0.9.11
[0.9.10]: https://gitee.com/quant1x/data/compare/v0.9.9...v0.9.10
[0.9.9]: https://gitee.com/quant1x/data/compare/v0.9.8...v0.9.9
[0.9.8]: https://gitee.com/quant1x/data/compare/v0.9.7...v0.9.8
[0.9.7]: https://gitee.com/quant1x/data/compare/v0.9.6...v0.9.7
[0.9.6]: https://gitee.com/quant1x/data/compare/v0.9.5...v0.9.6
[0.9.5]: https://gitee.com/quant1x/data/compare/v0.9.4...v0.9.5
[0.9.4]: https://gitee.com/quant1x/data/compare/v0.9.3...v0.9.4
[0.9.3]: https://gitee.com/quant1x/data/compare/v0.9.2...v0.9.3
[0.9.2]: https://gitee.com/quant1x/data/compare/v0.9.1...v0.9.2
[0.9.1]: https://gitee.com/quant1x/data/compare/v0.9.0...v0.9.1
[0.9.0]: https://gitee.com/quant1x/data/compare/v0.8.16...v0.9.0
[0.8.16]: https://gitee.com/quant1x/data/compare/v0.8.15...v0.8.16
[0.8.15]: https://gitee.com/quant1x/data/compare/v0.8.14...v0.8.15
[0.8.14]: https://gitee.com/quant1x/data/compare/v0.8.13...v0.8.14
[0.8.13]: https://gitee.com/quant1x/data/compare/v0.8.12...v0.8.13
[0.8.12]: https://gitee.com/quant1x/data/compare/v0.8.11...v0.8.12
[0.8.11]: https://gitee.com/quant1x/data/compare/v0.8.10...v0.8.11
[0.8.10]: https://gitee.com/quant1x/data/compare/v0.8.9...v0.8.10
[0.8.9]: https://gitee.com/quant1x/data/compare/v0.8.8...v0.8.9
[0.8.8]: https://gitee.com/quant1x/data/compare/v0.8.7...v0.8.8
[0.8.7]: https://gitee.com/quant1x/data/compare/v0.8.6...v0.8.7
[0.8.6]: https://gitee.com/quant1x/data/compare/v0.8.5...v0.8.6
[0.8.5]: https://gitee.com/quant1x/data/compare/v0.8.4...v0.8.5
[0.8.4]: https://gitee.com/quant1x/data/compare/v0.8.3...v0.8.4
[0.8.3]: https://gitee.com/quant1x/data/compare/v0.8.2...v0.8.3
[0.8.2]: https://gitee.com/quant1x/data/compare/v0.8.1...v0.8.2
[0.8.1]: https://gitee.com/quant1x/data/compare/v0.8.0...v0.8.1
[0.8.0]: https://gitee.com/quant1x/data/compare/v0.7.22...v0.8.0
[0.7.22]: https://gitee.com/quant1x/data/releases/tag/v0.7.22
