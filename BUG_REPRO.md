# Bug 复现说明

## Bug 是什么

记录列表请求超时后，后台仍继续扫描词条，接口最后返回过期列表，页面刷新会继续展示失效内容。

## 如何触发

准备一批词条，使用会被取消的 Context 调用 `Engine.List`，观察取消发生后调用是否仍返回列表结果。取消的列表请求应停止扫描并返回取消错误。

## 根因

根因涉及 `internal/core/engine.go`、`internal/core/history.go` 和 `internal/core/search.go` 的 `Engine.List` 调用链：扫描和结果返回边界没有持续检查已取消的 Context，导致取消请求仍完成并返回旧快照。该问题属于 Context 取消传播与跨层状态一致性失效。

## 运行指令

```bash
go test ./internal/core -v -count=1 -run '^TestBug008ListCancellation$'
```

## 错误信息

取消后的列表请求仍然返回了结果，说明扫描没有按请求生命周期停止。

## 错误堆栈

```text
=== RUN   TestBug008ListCancellation
    bug008_list_cancellation_test.go:25: expected list cancellation
--- FAIL: TestBug008ListCancellation (0.04s)
FAIL
FAIL	glimmer-lexicon-studio/internal/core	0.733s
FAIL
```
