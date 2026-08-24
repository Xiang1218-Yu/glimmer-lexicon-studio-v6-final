# Bug 复现说明

## Bug 是什么

创建记录后，证据列表的模块顺序会在不同请求之间变化，审校页面因此把规则说明显示在错误的位置。

## 如何触发

连续创建多条记录并比较每条记录的证据模块顺序，观察后续记录是否仍与首次创建保持一致。

## 根因

根因位于 `internal/core/engine.go` 的 `Engine.Create`，并经过 `internal/core/concept.go`、`internal/core/term.go` 的模块链路：一次创建改变了共享模块 slice 的顺序，后续记录因此得到不同的证据排列。

## 运行指令

```bash
go test ./internal/core -v -count=1 -run '^TestBug017ModuleEvidenceOrder$'
```

## 错误信息

第二次记录的证据模块顺序从预期的首个模块开始变成了 `term`、`variant`，说明共享顺序被改写。

## 错误堆栈

```text
=== RUN   TestBug017ModuleEvidenceOrder
    bug017_module_evidence_order_test.go:19: module order changed: first=term second=variant
--- FAIL: TestBug017ModuleEvidenceOrder (0.00s)
FAIL
FAIL	glimmer-lexicon-studio/internal/core	0.466s
```
