# Test Report: @go/safe

## 📋 测试概览
- **测试时间**: 2026-04-22
- **测试环境**: darwin/amd64
- **Go 版本**: 1.25.0

## ✅ 功能测试 (Functional Tests)
| 场景 | 状态 | 描述 |
| :--- | :--- | :--- |
| `TestSafeBuf` | PASS | 验证加解密一致性，及 `Close()` 后的数据清理。 |
| `TestMemoryErasure` | PASS | 验证内存物理擦除，确保明文内存被覆盖。 |
| `TestSafeBufCustomObfuscator` | PASS | 验证自定义加解密钩子的注入逻辑。 |
| `TestDisableCoreDump` | PASS | 验证系统级 CoreDump 防御接口调用。 |

## 🛡️ 鲁棒性防御 (Robustness)
- **内存锁定 (Mlock)**: 针对 Unix/Windows 实现了系统级内存锁定，防止敏感数据 Swapping 到硬盘。
- **Core Dump 防御**: 从操作系统层级禁止核心转储，防御异常导致的数据泄露。
- **全生命周期销毁**: 结合 `ZeroMemory` 与 `runtime.Finalizer`，确保敏感数据随对象销毁自动物理擦除。

## ⚡ 性能基准 (Benchmarks)
| 函数 | 平均耗时 | 性能分析 |
| :--- | :--- | :--- |
| `BenchmarkSafeBufOpenClose` | **1095 ns/op** | 纳秒级加解密处理，对业务性能影响极小。 |
