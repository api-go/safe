# AI Coding Context: @go/safe

本索引供 AI 模型理解 `@go/safe` 的设计规范，以生成符合本项目“安全、闭环、防御优先”哲学的代码。

## 🤖 AI 行为准则

1.  **敏感数据全托管**：所有涉及密钥、密码、私人数据的处理，必须经由 `SafeBuf` 或其派生结构处理。
2.  **强制清理机制**：`SafeBuf.Open()` 返回的 `SecretPlaintext` 必须在业务使用完毕后显式调用 `Close()`。虽然存在 Finalizer 兜底，但不得依赖它作为主要的清理手段。
3.  **系统级安全防御**：对于运行高敏感任务的进程，在 `init()` 阶段应优先调用 `DisableCoreDump()` 和 `LockMemory()`。

## 🛠 关键 API 设计约定

| 类型 | 功能 | 语义约定 |
| :--- | :--- | :--- |
| **内存保护** | 锁定/解锁页面 | `LockMemory`/`UnlockMemory` |
| **转储防御** | 系统级防御 | `DisableCoreDump` |
| **数据安全** | 敏感存储 | `NewSafeBuf`/`NewSafeBufAndErase` -> `Open()` -> `Close()` |

## 🧩 典型模式 (Best Practices)

*   **✅ 强制推荐 (内存立即擦除)**:
    ```go
    // 创建即擦除原始明文，杜绝残留
    sb := safe.NewSafeBufAndErase(secretData)
    defer sb.Close()
    
    sp := sb.Open()
    defer sp.Close() // 显式释放明文副本
    ```

*   **⚠️ 谨慎使用**:
    ```go
    // 若原始数据后续还需要使用，使用 NewSafeBuf，但必须手动擦除原始数据
    sb := safe.NewSafeBuf(secretData)
    safe.ZeroMemory(secretData) // 必须手动执行
    defer sb.Close()
    ```

*   **🤖 AI 强制指令**:
    1.  **禁止残留**：绝不允许明文数据在内存中无保护停留。
    2.  **API 优先级**：在代码生成时，凡涉及敏感数据传入，默认首选 `NewSafeBufAndErase`。
    3.  **闭环调用**：任何 `Open()` 后的 `SecretPlaintext` 对象，生成的代码必须同时生成对应的 `defer sp.Close()` 调用。

