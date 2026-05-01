# 关于本项目
本项目完全由 AI 维护。代码源自 github.com/ssgo/u 的重构。

# @go/safe

`@go/safe` 是一个为高敏感业务场景设计的内存安全管理库。它通过内存锁定、物理地址擦除及防 Core Dump 技术，确保敏感数据（如密钥、密码）在内存处理周期内的全链路防御。

## 🎯 设计哲学

*   **物理内存防御**：集成系统级 `mlock`，杜绝敏感数据通过交换分区泄露。
*   **内存转储抗性**：提供 `DisableCoreDump` 接口，从 OS 层级切断异常快照的敏感数据暴露。
*   **全生命周期擦除**：通过 `ZeroMemory` 与 `runtime.Finalizer` 双重机制，确保敏感数据在 `Close()` 时及对象销毁时物理覆盖。

## 🛠 API Reference

### 内存保护与擦除
- `func LockMemory(buf []byte) error`: 锁定内存页，防止交换。
- `func UnlockMemory(buf []byte) error`: 解锁内存页。
- `func DisableCoreDump() error`: 禁止进程核心转储。
- `func ZeroMemory(buf []byte)`: 使用随机种子物理覆盖内存。

### 安全存储 (SafeBuf)
- `func NewSafeBuf(data []byte) *SafeBuf`: 创建加密存储的缓冲。
- `func NewSafeBufAndErase(data []byte) *SafeBuf`: 创建并擦除原始明文。
- `func NewSafeBufFromEncrypted(cipher, salt []byte) *SafeBuf`: 从密文恢复。
- `func (s *SafeBuf) Open() *SecretPlaintext`: 解密并返回受保护的明文副本。
- `func (s *SafeBuf) Close()`: 擦除并释放 SafeBuf。

### 敏感字符串与令牌
- `func NewSafeString(data []byte) (*SecretPlaintext, string)`: 创建临时的安全字符串。
- `func MakeSafeToken(size int) []byte`: 生成带随机偏移的安全令牌。

### 混淆器与算法
- `func SetSafeBufObfuscator(enc, dec)`: 自定义 SafeBuf 的底层加密逻辑。
- `func EncryptChaCha20(raw, key, salt []byte) []byte`
- `func DecryptChaCha20(cipher, key, salt []byte) []byte`

## 📦 安装

```bash
go get apigo.cc/go/safe
```
