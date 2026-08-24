# Bug 复现说明

## Bug 是什么

审校人员推进词条阶段时页面超时，但后台仍把词条推进到新阶段，之后重新提交会看到版本跳跃。

## 如何触发

创建一条词条，使用已经取消的 Context 调用 `Engine.Advance`，再检查阶段和版本是否保持不变。取消后的推进应返回错误，且不能提交状态变化。

## 根因

根因涉及 `internal/core/engine.go`、`internal/core/review.go` 和 `internal/core/status.go` 的 `Engine.Advance` 调用链：阶段校验期间 Context 已取消，但取消状态没有阻止后续状态提交，导致阶段和版本仍然发生变化。该问题属于 Context 取消传播与跨层状态一致性失效。

## 运行指令

```bash
go test ./internal/core -v -count=1 -run '^TestBug006AdvanceCancellation$'
```

## 错误信息

取消后的推进调用没有失败，后台状态已经被推进。

## 错误堆栈

```text
=== RUN   TestBug006AdvanceCancellation
    bug006_advance_cancellation_test.go:23: expected canceled advance to fail
--- FAIL: TestBug006AdvanceCancellation (0.02s)
FAIL
FAIL	glimmer-lexicon-studio/internal/core	0.727s
FAIL
```
