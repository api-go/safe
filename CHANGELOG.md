# Changelog: @go/safe

## [v1.0.4] - 2026-05-01

### Fixed
- **Windows 兼容性**：修复 `mem_windows.go` 缺失 `unsafe` 导入导致的编译失败。

### Improved
- **命名规范**：重构短变量名（如 `sb`, `sp`）为更具描述性的名称，提升代码可读性。
- **鲁棒性增强**：为 `ChaCha20` 加解密增加长度校验与错误处理。
- **安全性提升**：改进 `ZeroMemory` 的随机种子生成机制，引入 `crypto/rand` 熵池。
- **内存锁定处理**：显式处理 `LockMemory`/`UnlockMemory` 的潜在错误。

## [v1.0.0] - 2026-04-22

### Added
- **内存锁定机制**：实现 `LockMemory`/`UnlockMemory` (支持 Unix/Windows 条件编译)，防止敏感数据被交换至磁盘。
- **Core Dump 防御**：新增 `DisableCoreDump`，切断异常进程的核心转储文件暴露路径。
- **安全缓存 (SafeBuf)**：实现动态加密内存缓冲，提供基于 `ChaCha20` 的存储保护，并支持自定义混淆器注入。
- **内存物理擦除**：提供 `ZeroMemory` 强力覆盖机制，配合 `runtime.Finalizer` 实现全链路数据自动销毁。
