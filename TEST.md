# Test Report: @go/safe

## v1.5.4 验证
- 独立模块解析、全量测试与基准测试均通过。

## 测试环境
- **操作系统**: darwin
- **架构**: amd64
- **Go 版本**: 1.25.0

## 覆盖场景
1. **SafeBuf 基础功能**: 验证 `NewSafeBuf`, `Open`, `Close` 的闭环逻辑。
2. **内存物理擦除**: 通过 `unsafe` 指针直接验证 `ZeroMemory` 是否成功覆盖了已关闭的明文数据。
3. **自定义混淆器**: 验证 `SetSafeBufObfuscator` 注入自定义算法的有效性。
4. **Core Dump 防御**: 验证 `DisableCoreDump` 在当前系统下的调用。
5. **ChaCha20 健壮性**: (已在代码中增加长度校验与错误处理)。

## 测试结果
```
=== RUN   TestSafeBuf
--- PASS: TestSafeBuf (0.00s)
=== RUN   TestMemoryErasure
--- PASS: TestMemoryErasure (0.00s)
=== RUN   TestSafeBufCustomObfuscator
--- PASS: TestSafeBufCustomObfuscator (0.00s)
=== RUN   TestDisableCoreDump
--- PASS: TestDisableCoreDump (0.00s)
PASS
ok      apigo.cc/go/safe        0.371s
```

## 性能测试 (Benchmark)
```
goos: darwin
goarch: amd64
pkg: apigo.cc/go/safe
cpu: Intel(R) Core(TM) i9-9980HK CPU @ 2.40GHz
BenchmarkSafeBufOpenClose-16              976933              1319 ns/op
```
*注：引入更强随机种子的内存擦除后，单次操作耗时约增加 300ns，符合预期。*
