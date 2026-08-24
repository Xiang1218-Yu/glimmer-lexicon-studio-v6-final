# Bug 复现说明

## Bug 是什么

服务运行一段时间后，同样的词条再次创建会得到不同的模块证据顺序，审校页面出现不稳定的规则提示。

## 如何触发

连续创建相同内容的记录，比较每次返回证据列表中的模块顺序是否保持一致。

## 根因

根因位于 `internal/core/engine.go` 的 `Engine.Create`，并涉及 `internal/core/concept.go`、`internal/core/term.go` 的模块执行链路：创建流程改变了共享模块顺序，导致跨请求结果不一致。

## 运行指令

```bash
go test ./internal/core -v -count=1 -run '^TestBug024ModuleOrderStable$'
```

## 错误信息

第一个模块在测试运行中变成了 `term`，没有保持稳定的初始顺序。

## 错误堆栈

```text
=== RUN   TestBug024ModuleOrderStable
    bug024_module_order_stable_test.go:19: unexpected first module on run 0: term
--- FAIL: TestBug024ModuleOrderStable (0.00s)
FAIL
FAIL	glimmer-lexicon-studio/internal/core	0.497s
```
