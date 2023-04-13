# Changelog
All notable changes to this project will be documented in this file.

## [Unreleased]

## [1.0.15] - 2023-04-13
### Changed
- 调整tick的告警信息日志.

## [1.0.14] - 2023-04-13
### Changed
- !15 #I6V85X 调整实时数据的更新策略 * 调整实时数据的更新策略, 非交易时段不更新 * 梳理K线新增计算方法 * 对齐分钟数.

## [1.0.13] - 2023-04-12
### Changed
- 字段对齐.

## [1.0.12] - 2023-04-12
### Changed
- !14 #I6V0A6 K线增加涨幅 * 日线周线增加涨跌幅字段 * 日K线增加涨跌幅字段.

## [1.0.11] - 2023-04-12
### Changed
- !13 #I6UL6N 日K线增加量比 * 日K线增加量比 * 增加注释 * 删除无用的代码 * 增加计算某个交易日的分钟数.
- 更新pandas版本号.

## [1.0.10] - 2023-04-11
### Changed
- 删除无用的代码.

## [1.0.9] - 2023-04-10
### Changed
- 删除无用的代码.

## [1.0.8] - 2023-04-10
### Changed
- 删除废弃的代码.

## [1.0.7] - 2023-04-10
### Changed
- !12 #I6U890 修复tick数据的bug * 增加对空文件的判断.
- 修订量比函数的注释, 过去5个交易日(不包含当日).
- 删除废弃的代码.
- Update changelog.
- Add change log.
- 删除英文版本README.
- Update changelog.
- Update changelog.
- 修改changelog.

## [1.0.6] - 2023-04-09
### Changed
- Update changelog.

## [1.0.5] - 2023-04-06
### Changed
- 强化交易日历校验.
- 强化交易日历校验.
- 备注未整合进日历.
- 增加交易日历校验函数.

## [1.0.4] - 2023-04-06
### Changed
- !11 #I6T5GE 调整板块代码和类型 * 恢复风格板块, 过滤在策略工具中自助实现 * 忽略风格板块 * 新增通过板块代码获取板块类型名称.

## [1.0.3] - 2023-04-05
### Changed
- !10 #I6T4B8 修复创业板涨停板的错误 * 创业板涨跌停限制应该是20%.

## [1.0.2] - 2023-04-05
### Changed
- !9 #I6T43C 用最新的通达信股票列表更新富途静态的证券名称 * 更新股票名称.

## [1.0.0] - 2023-04-05
### Changed
- 补全字段.
- 字段改名.
- !8 #I6T1IE K线增加内外盘数据 * 去掉早期附加成交量拆分的做法 * K线弃用废弃的函数 * 实时数据增加内外盘数据 * 日K线增加成交量和成交金额的买卖方向 * 增加退市和摘牌标志 * 增加st的过滤 * tick函数处理中补全代码前缀 * 增加获取tick文件名, 不创建路径 * 修复证券代码中深圳市场缺少前缀的bug * 去掉无效的引用 * 修订slice唯一的处理方式.

## [0.9.28] - 2023-04-05
### Changed
- 如果没有date序列, 直接返回.

## [0.9.27] - 2023-04-04
### Changed
- !7 #I6T08K K线附加成交量内外盘数据 * K线附加成交量细节.

## [0.9.26] - 2023-04-04
### Changed
- !6 #I6SS99 补全行业板块数据 * 增加行业指数.

## [0.9.24] - 2023-04-04
### Changed
- 去掉风格指数.

## [0.9.23] - 2023-04-03
### Changed
- 增加周K线函数.
- 增加是否前复权的处理.
- 测试周线.
- 优化代码.

## [0.9.22] - 2023-04-03
### Changed
- 预留开发任务.
- 调整js的工具包路径.
- 调整退出信号的监控方式.
- 去掉无用的被动注释.
- 调整缓存类型.

## [0.9.21] - 2023-03-24
### Changed
- !5 #I6OXJ9 增加证券代码和证券名称的映射关系 * 增加当日证券代码和证券名称的映射关系.

## [0.9.20] - 2023-03-24
### Changed
- 修订pandas版本.

## [0.9.19] - 2023-03-20
### Changed
- 修订交易日历的处理方式.
- 更新依赖库版本号.
- 更新依赖库版本号.
- 修订交易日期的csv文件名.

## [0.9.18] - 2023-03-19
### Changed
- 拆分源代码.

## [0.9.17] - 2023-03-18
### Changed
- 更新依赖库.

## [0.9.16] - 2023-03-18
### Changed
- 调整日志级别为INFO.
- 更新依赖库版本号.
- 增加计算量比的函数.

## [0.9.15] - 2023-03-16
### Changed
- 微调脚本代码段落.

## [0.9.14] - 2023-03-16
### Changed
- 更新通达信工具库.
- 更新tdxzs*文件, 后续的处理可以手动更新.quant1x/bk/*.cfg.
- 更新tdx工具库版本.
- 增加windows arm64编译脚本.
- 编译时控制台输出版本号.

## [0.9.13] - 2023-03-13
### Changed
- !4 #I6MFH8 每天同步板块数据文件 * 调整快照数据工具 * 调整通达信工具类版本.
- 调整板块源文件名.

## [0.9.12] - 2023-03-13
### Changed
- 增加涨跌停板限制幅度函数.

## [0.9.11] - 2023-03-13
### Changed
- 板块代码去重.

## [0.9.10] - 2023-03-13
### Changed
- 批量获取快照数据, 返回快照结构的数组.

## [0.9.9] - 2023-03-13
### Changed
- 调整BlockList返回的code字段, 要加前缀.

## [0.9.8] - 2023-03-13
### Changed
- !3 #I6MFH4 增加板块的快照数据 * 增加windows打包脚本 * 调整快照的package * 增加实时行情的快照数据 * 计划调整板块列表的字段类型, 目前板块代码是int64 * 新增获取最后一个交易日的函数 * 修订实时数据函数内部变量名, 使变量名贴合意图 * 增加mac arm64的编译脚本 * 增加mac amd64的编译脚本 * 增加函数的注释 * 增加独立的mac打包脚本.

## [0.9.7] - 2023-03-11
### Changed
- 调整五档行情数据更新K线时, 需要区分指数和个股的字段不同.
- 区分个股和指数的K线数据, 自动修正本地缓存.
- 修订版本号.

## [0.9.6] - 2023-03-10
### Changed
- 删掉无用的代码.

## [0.9.5] - 2023-03-10
### Changed
- 修订板块代码.

## [0.9.4] - 2023-03-10
### Changed
- 修正板块代码.
- 修订板块代码.
- 修订板块代码.
- 修订版本代码.
- 修订板块代码.

## [0.9.3] - 2023-03-10
### Changed
- 修订板块代码.

## [0.9.2] - 2023-03-10
### Changed
- !2 #I6LR5W 增加板块数据 * add block data.

## [0.9.1] - 2023-03-09
### Changed
- Tick数据不全跳过.

## [0.9.0] - 2023-03-09
### Changed
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

[Unreleased]: https://gitee.com/quant1x/data/compare/v1.0.15...HEAD
[1.0.15]: https://gitee.com/quant1x/data/compare/v1.0.14...v1.0.15
[1.0.14]: https://gitee.com/quant1x/data/compare/v1.0.13...v1.0.14
[1.0.13]: https://gitee.com/quant1x/data/compare/v1.0.12...v1.0.13
[1.0.12]: https://gitee.com/quant1x/data/compare/v1.0.11...v1.0.12
[1.0.11]: https://gitee.com/quant1x/data/compare/v1.0.10...v1.0.11
[1.0.10]: https://gitee.com/quant1x/data/compare/v1.0.9...v1.0.10
[1.0.9]: https://gitee.com/quant1x/data/compare/v1.0.8...v1.0.9
[1.0.8]: https://gitee.com/quant1x/data/compare/v1.0.7...v1.0.8
[1.0.7]: https://gitee.com/quant1x/data/compare/v1.0.6...v1.0.7
[1.0.6]: https://gitee.com/quant1x/data/compare/v1.0.5...v1.0.6
[1.0.5]: https://gitee.com/quant1x/data/compare/v1.0.4...v1.0.5
[1.0.4]: https://gitee.com/quant1x/data/compare/v1.0.3...v1.0.4
[1.0.3]: https://gitee.com/quant1x/data/compare/v1.0.2...v1.0.3
[1.0.2]: https://gitee.com/quant1x/data/compare/v1.0.0...v1.0.2
[1.0.0]: https://gitee.com/quant1x/data/compare/v0.9.28...v1.0.0
[0.9.28]: https://gitee.com/quant1x/data/compare/v0.9.27...v0.9.28
[0.9.27]: https://gitee.com/quant1x/data/compare/v0.9.26...v0.9.27
[0.9.26]: https://gitee.com/quant1x/data/compare/v0.9.24...v0.9.26
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
[0.9.0]: https://gitee.com/quant1x/data/compare/v0.8.10...v0.9.0
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
