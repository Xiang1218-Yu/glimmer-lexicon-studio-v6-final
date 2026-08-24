# Bug 复现说明

## Bug 是什么

审校人员推进记录后，历史显示为匿名操作，页面无法追溯实际完成阶段变更的操作者。

## 如何触发

使用明确的审校人员身份推进一条记录，再读取该记录历史并检查阶段变更的操作者字段。

## 根因

根因位于 `internal/core/engine.go` 的 `Engine.Advance`，并涉及 `internal/core/review.go`、`internal/core/reviewer.go` 的操作者传递链路：推进函数退出路径上的状态写入使用了错误的 actor 值，导致审校历史丢失操作者。

## 运行指令

```bash
go test ./internal/core -v -count=1 -run '^TestBug018HistoryActorPreserved$'
```

## 错误信息

阶段变更历史的操作者被写成 `anonymous`，没有保留调用时提供的审校人员身份。

## 错误堆栈

```text
=== RUN   TestBug018HistoryActorPreserved
    bug018_history_actor_preserved_test.go:22: actor was lost: []core.Transition{core.Transition{From:"", To:"draft", By:"tester", At:time.Date(2026, time.August, 24, 3, 28, 24, 463279000, time.UTC)}, core.Transition{From:"draft", To:"review", By:"anonymous", At:time.Date(2026, time.August, 24, 3, 28, 24, 463682000, time.UTC)}}
--- FAIL: TestBug018HistoryActorPreserved (0.00s)
FAIL
FAIL	glimmer-lexicon-studio/internal/core	0.465s
```
