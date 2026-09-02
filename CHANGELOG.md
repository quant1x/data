# Changelog
All notable changes to this project will be documented in this file.

## [Unreleased]

## [0.2.5] - 2026-09-02
### Changed
- fix: slice bounds out of range
- release version 0.2.5

## [0.2.4] - 2026-08-26
### Changed
- Merge branch 'master' of https://gitee.com/quant1x/data
- release v0.2.4

## [0.2.3] - 2026-08-26
### Changed
- 更新go版本到1.27.0
- 更新go版本到1.27.0
- release v0.2.3

## [0.2.2] - 2026-08-24
### Changed
- bug fix: slice out of range
- release v0.2.2

## [0.2.1] - 2026-07-23
### Changed
- 修复交易日范围越界问题并补充回归测试
- release v0.2.1

## [0.2.0] - 2026-06-12
### Changed
- 调整git仓库地址
- update changelog

## [0.1.18] - 2025-12-29
### Changed
- rust: 重构代码结构，与go匹配
- 将cache目录提升到项目顶层
- git忽略vscode的项目配置
- go: 删除废弃的cache测试代码
- go: 基础配置信息
- Merge branch 'fix-pkg' into 0.1.x
- update changelog

## [0.1.17] - 2025-12-07
### Changed
- 删除根路径的go模块配置
- update changelog

## [0.1.16] - 2025-12-07
### Changed
- 增加演示add函数
- update changelog

## [0.1.15] - 2025-12-07
### Changed
- update changelog

## [0.1.14] - 2025-12-07
### Changed
- update changelog

## [0.1.13] - 2025-12-07
### Changed
- 删除git 子目录的尝试的代码
- 恢复gitee.com仓库
- update changelog

## [0.1.12] - 2025-12-07
### Changed
- 修改index.html，去掉body
- update changelog

## [0.1.11] - 2025-12-07
### Changed
- 补充注释
- go: 尝试将git子目录作为go module
- update changelog
- 新增vscode的launch配置
- 调整测试代码

## [0.1.10] - 2025-11-03
### Changed
- 如果没有数据不更新缓存文件的修改时间
- update changelog

## [0.1.9] - 2025-10-30
### Changed
- 删除废弃的单元测试
- 将证券名称扩展到16个字节
- update changelog

## [0.1.8] - 2025-10-30
### Changed
- 通过协议下载zhb.zip文件，并指定解压需要的板块文件
- update changelog

## [0.1.7] - 2025-10-29
### Changed
- 添加北证50指数
- update changelog

## [0.1.6] - 2025-10-29
### Changed
- 去掉K上datetime字段上的毫秒
- update changelog

## [0.1.5] - 2025-10-29
### Changed
- 补充除权除息字段说明
- update changelog

## [0.1.4] - 2025-10-29
### Changed
- 删除废弃的代码
- update changelog

## [0.1.3] - 2025-10-29
### Changed
- 删除废弃的协议文档
- 合并字符串功能函数
- 解压缩功能源文件名改为bytes
- 日期时间类功能源文件改为datetime
- 删除废弃的命令字
- 整理协议包, 删除废弃的代码
- update changelog

## [0.1.2] - 2025-10-29
### Changed
- 新增公开的command函数
- update changelog

## [0.1.1] - 2025-10-29
### Changed
- 更新num依赖库版本到0.7.15
- update changelog

## [0.1.0] - 2025-10-29
### Changed
- go代码改内部功能为utils
- update changelog

## [0.0.8] - 2025-10-28
### Changed
- 修订缓存的证券列表, 改由最新的接口更新
- update changelog

## [0.0.7] - 2025-10-28
### Changed
- 更新pkg版本到0.9.0, 剔除对javescript的依赖
- update changelog

## [0.0.6] - 2025-10-28
### Changed
- 新增exchange代码
- update changelog

## [0.0.5] - 2025-10-28
### Changed
- 屏蔽默认端口
- update changelog

## [0.0.4] - 2025-10-28
### Changed
- 删除使用域名的服务器地址
- 删除废弃的字段
- update changelog

## [0.0.3] - 2025-10-28
### Changed
- 新增新的证券列表接口
- 新增rust实现level1接口
- update changelog

## [0.0.2] - 2025-10-14
### Changed
- 修订readme
- update changelog

## [0.0.1] - 2025-10-13
### Changed
- Initial commit
- 初始化level1协议集合
- 新增README文档
- 支持go1.25
- update changelog


[Unreleased]: https://gitee.com/quant1x/data.git/compare/v0.2.5...HEAD
[0.2.5]: https://gitee.com/quant1x/data.git/compare/v0.2.4...v0.2.5
[0.2.4]: https://gitee.com/quant1x/data.git/compare/v0.2.3...v0.2.4
[0.2.3]: https://gitee.com/quant1x/data.git/compare/v0.2.2...v0.2.3
[0.2.2]: https://gitee.com/quant1x/data.git/compare/v0.2.1...v0.2.2
[0.2.1]: https://gitee.com/quant1x/data.git/compare/v0.2.0...v0.2.1
[0.2.0]: https://gitee.com/quant1x/data.git/compare/v0.1.18...v0.2.0
[0.1.18]: https://gitee.com/quant1x/data.git/compare/v0.1.17...v0.1.18
[0.1.17]: https://gitee.com/quant1x/data.git/compare/v0.1.16...v0.1.17
[0.1.16]: https://gitee.com/quant1x/data.git/compare/v0.1.15...v0.1.16
[0.1.15]: https://gitee.com/quant1x/data.git/compare/v0.1.14...v0.1.15
[0.1.14]: https://gitee.com/quant1x/data.git/compare/v0.1.13...v0.1.14
[0.1.13]: https://gitee.com/quant1x/data.git/compare/v0.1.12...v0.1.13
[0.1.12]: https://gitee.com/quant1x/data.git/compare/v0.1.11...v0.1.12
[0.1.11]: https://gitee.com/quant1x/data.git/compare/v0.1.10...v0.1.11
[0.1.10]: https://gitee.com/quant1x/data.git/compare/v0.1.9...v0.1.10
[0.1.9]: https://gitee.com/quant1x/data.git/compare/v0.1.8...v0.1.9
[0.1.8]: https://gitee.com/quant1x/data.git/compare/v0.1.7...v0.1.8
[0.1.7]: https://gitee.com/quant1x/data.git/compare/v0.1.6...v0.1.7
[0.1.6]: https://gitee.com/quant1x/data.git/compare/v0.1.5...v0.1.6
[0.1.5]: https://gitee.com/quant1x/data.git/compare/v0.1.4...v0.1.5
[0.1.4]: https://gitee.com/quant1x/data.git/compare/v0.1.3...v0.1.4
[0.1.3]: https://gitee.com/quant1x/data.git/compare/v0.1.2...v0.1.3
[0.1.2]: https://gitee.com/quant1x/data.git/compare/v0.1.1...v0.1.2
[0.1.1]: https://gitee.com/quant1x/data.git/compare/v0.1.0...v0.1.1
[0.1.0]: https://gitee.com/quant1x/data.git/compare/v0.0.8...v0.1.0
[0.0.8]: https://gitee.com/quant1x/data.git/compare/v0.0.7...v0.0.8
[0.0.7]: https://gitee.com/quant1x/data.git/compare/v0.0.6...v0.0.7
[0.0.6]: https://gitee.com/quant1x/data.git/compare/v0.0.5...v0.0.6
[0.0.5]: https://gitee.com/quant1x/data.git/compare/v0.0.4...v0.0.5
[0.0.4]: https://gitee.com/quant1x/data.git/compare/v0.0.3...v0.0.4
[0.0.3]: https://gitee.com/quant1x/data.git/compare/v0.0.2...v0.0.3
[0.0.2]: https://gitee.com/quant1x/data.git/compare/v0.0.1...v0.0.2

[0.0.1]: https://gitee.com/quant1x/data.git/releases/tag/v0.0.1
