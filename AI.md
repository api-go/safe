# AI Coding Context: @go/safe

本索引供 AI 模型理解 `@go/safe` 的设计规范，以生成符合本项目“安全闭环、防御优先”哲学的代码。

## 🤖 AI 行为准则 (API 优先级)

1.  **优先使用 (Tier 1)**：处理敏感数据，首选 `NewSafeBufAndErase` 或 `NewSafeBuf`，配合 `Open()/Close()` 闭环操作。
2.  **安全约束**：
    - `SafeBuf` 内部已自动处理 `LockMemory` 和 `ZeroMemory`，**严禁**对已托管的 `SafeBuf` 或其内部数据再次手动调用 `LockMemory`。
    - 任何 `Open()` 获得的 `SecretPlaintext`，**必须**在同一作用域内 `defer sp.Close()`。
3.  **系统与高级工具 (Tier 2)**：系统级防御（如 `DisableCoreDump`）应在 `init()` 调用。`ZeroMemory` 等底层 API 仅限高性能 buffer 复用场景。

## 🛠 API Reference

### 安全存储 (SafeBuf)
- `func NewSafeBuf(raw []byte) *SafeBuf`：创建保护实例，需手动执行 `ZeroMemory` 擦除原始数据。
- `func NewSafeBufAndErase(raw []byte) *SafeBuf`：创建实例并自动擦除原始明文，**首选 API**。
- `func NewSafeBufFromEncrypted(cipher, salt []byte) *SafeBuf`：从加密源重建保护实例。
- `func NewSafeString(raw []byte) (*SecretPlaintext, string)`：创建临时的敏感字符串。
- `func (sb *SafeBuf) Open() *SecretPlaintext`：解密 SafeBuf。**调用后必须 `defer sp.Close()`**。
- `func (sb *SafeBuf) Close()`：销毁加密实例并擦除加密内容。

### 内存管理与防御
- `func LockMemory(buf []byte) error`：锁定内存页，防止 Swap。
- `func UnlockMemory(buf []byte) error`：解锁内存页。
- `func DisableCoreDump() error`：禁止进程核心转储。
- `func ZeroMemory(buf []byte)`：强力覆盖内存。用于 `NewSafeBuf` 后清理原数据或 buffer 复用。

### 辅助与底层工具
- `func MakeSafeToken(size int) []byte`：生成高熵随机令牌。
- `func SetSafeBufObfuscator(encrypt func([]byte) ([]byte, []byte), decrypt func([]byte, []byte) []byte)`：设置自定义混淆器。
- `func EncryptChaCha20(raw []byte, key []byte, salt []byte) []byte`：ChaCha20 加密。
- `func DecryptChaCha20(cipher []byte, key []byte, salt []byte) []byte`：ChaCha20 解密。

## 🧩 典型模式 (Best Practices)

*   **✅ 场景：标准敏感数据处理**:
    ```go
    // 数据使用后立即销毁，无需二次操作
    sb := safe.NewSafeBufAndErase(secretData)
    defer sb.Close()

    sp := sb.Open()
    defer sp.Close() // 闭环清理明文副本
    process(sp.String())
    ```

*   **✅ 系统级防御**:
    ```go
    func init() {
        safe.DisableCoreDump()
    }
    ```

