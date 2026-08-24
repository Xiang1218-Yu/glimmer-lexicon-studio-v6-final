# Bug 复现说明

## Bug 是什么

一次创建请求在规则校验阶段报错后，列表里仍能查到空词条记录，重新提交同一个编号也会被拦住。

## 如何触发

提交一条会在规则校验阶段失败的词条，随后查询记录列表并再次提交相同编号。失败请求不应留下空记录，也不应污染后续请求。

## 根因

根因涉及 `internal/core/engine.go`、`internal/core/concept.go` 和 `internal/core/term.go` 的 `Engine.Create` 调用链：创建失败后共享记录表残留临时记录，后续请求读取到这份不完整状态并受到污染。该问题属于跨层错误传播与错误状态映射失效。

## 运行指令

```bash
go test ./internal/core -v -count=1 -run '^TestBug010CreateRollbackAfterError$'
```

## 错误信息

创建失败后仍能查到一条空记录，说明失败路径没有回滚已经写入的状态。

## 错误堆栈

```text
=== RUN   TestBug010CreateRollbackAfterError
    bug010_create_rollback_after_error_test.go:21: failed create polluted records: []core.Record{core.Record{ID:"rollback-term", Stage:"draft", Payload:"rollback-error", Score:0, Version:1, Evidence:[]core.Evidence(nil), History:[]core.Transition(nil), CreatedAt:time.Date(2026, time.August, 24, 5, 3, 25, 479184000, time.UTC), UpdatedAt:time.Date(2026, time.August 24, 5, 3, 25, 479184000, time.UTC)}}
--- FAIL: TestBug010CreateRollbackAfterError (0.00s)
FAIL
FAIL	glimmer-lexicon-studio/internal/core	0.773s
FAIL
```
