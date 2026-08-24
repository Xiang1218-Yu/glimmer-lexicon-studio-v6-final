# Bug 复现说明

## Bug 是什么

页面取消查看请求后，查询接口仍从内存返回记录，调用方会把已经无效的结果继续展示。

## 如何触发

先创建一条词条，再使用已经取消的 Context 调用 `Engine.Get`，观察调用是否仍返回记录而不是取消错误。

## 根因

根因涉及 `internal/core/engine.go` 的 `Engine.Get`，以及 `internal/core/context.go`、`internal/core/workspace.go` 的上下文传递链路：已取消的 Context 没有在查询入口和返回边界被正确传播，导致取消请求仍被当作成功查询处理。

## 运行指令

```bash
go test ./internal/core -v -count=1 -run '^TestBug012GetCanceledContext$'
```

## 错误信息

已取消的查询仍返回了一条记录，调用方因此继续处理失效结果。

## 错误堆栈

```text
=== RUN   TestBug012GetCanceledContext
    bug012_get_canceled_context_test.go:17: canceled get returned a record
--- FAIL: TestBug012GetCanceledContext (0.00s)
FAIL
FAIL	glimmer-lexicon-studio/internal/core	0.269s
```
