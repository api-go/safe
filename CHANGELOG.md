# Changelog: @go/safe

## [v1.0.0] - 2026-04-22

### Added
- **内存锁定机制**：实现 `LockMemory`/`UnlockMemory` (支持 Unix/Windows 条件编译)，防止敏感数据被交换至磁盘。
- **Core Dump 防御**：新增 `DisableCoreDump`，切断异常进程的核心转储文件暴露路径。
- **安全缓存 (SafeBuf)**：实现动态加密内存缓冲，提供基于 `ChaCha20` 的存储保护，并支持自定义混淆器注入。
- **内存物理擦除**：提供 `ZeroMemory` 强力覆盖机制，配合 `runtime.Finalizer` 实现全链路数据自动销毁。
