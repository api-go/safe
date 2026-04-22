# 关于本项目
本项目完全由 AI 维护。代码源自 github.com/ssgo/u 的重构。

# @go/safe

`@go/safe` 是一个为高敏感业务场景设计的内存安全管理库。它通过内存锁定、物理地址擦除及防 Core Dump 技术，确保敏感数据（如密钥、密码）在内存处理周期内的全链路防御。

## 🎯 设计哲学

*   **物理内存防御**：集成系统级 `mlock`，杜绝敏感数据通过交换分区泄露。
*   **内存转储抗性**：提供 `DisableCoreDump` 接口，从 OS 层级切断异常快照的敏感数据暴露。
*   **全生命周期擦除**：通过 `ZeroMemory` 与 `runtime.Finalizer` 双重机制，确保敏感数据在 `Close()` 时及对象销毁时物理覆盖。

## 🛠 API Reference

### 内存保护
- `func LockMemory(buf []byte) error`
- `func UnlockMemory(buf []byte) error`
- `func DisableCoreDump() error`

### 安全存储
- `func NewSafeBuf(data []byte) *SafeBuf`
- `func NewSafeBufAndErase(data []byte) *SafeBuf`
- `func NewSafeBufFromEncrypted(cipher, salt []byte) *SafeBuf`
- `func MakeSafeToken(size int) []byte`
- `func SetSafeBufObfuscator(enc, dec)`

## 📦 安装

```bash
go get apigo.cc/go/safe
```
