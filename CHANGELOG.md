# Changelog: @go/safe

## v1.5.2 (2026-06-21)
- **依赖对齐**:
  - 升级 `jsmod` 依赖至 v1.5.3，`cast` 依赖至 v1.5.3，`rand` 依赖至 v1.5.3。

## v1.5.1 (2026-06-11)
- **版本对齐**: 基础设施全局对齐 v1.5.1。

## [v1.0.6] - 2026-05-09

### Added
- **安全拼接 (Concat)**：新增 `Concat` 工具，支持对多个敏感部分（SafeBuf, SecretPlaintext, string, []byte）进行零残留安全拼接，返回持久加密的 `SafeBuf`。
- **安全编码 (Base64/Hex/UrlEncode)**：新增多种编码工具，支持接受 `...any` 变长参数。内置隐式安全拼接能力（无需额外调用 `Concat`），直接返回持久加密的 `SafeBuf` 对象，脱离对 GC Finalizer 的强依赖，安全且可重复使用。

### Improved
- **生命周期增强**：优化了 `SecretPlaintext` 的安全性，增强了 `String()` 方法的鲁棒性。

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
