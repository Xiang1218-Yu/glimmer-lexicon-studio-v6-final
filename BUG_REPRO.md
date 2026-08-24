# Bug 复现说明

## Bug 是什么

词条推进到不被接受的阶段时，接口虽然返回失败，但失败操作已经向历史记录写入了一条假的推进记录，版本状态也可能被提前改变。

## 如何触发

创建一条词条后调用 `Engine.Advance` 推进到不被对应模块接受的阶段，再检查该词条的历史长度和版本号。失败调用不应改变这两项状态。

## 根因

根因位于 `internal/core/engine.go` 的 `Engine.Advance`：阶段校验完成前就写入历史并保存记录，校验失败后没有回滚这次状态写入。`internal/core/review.go` 与 `internal/core/status.go` 参与阶段接受性校验，延长了错误状态被观察到的调用链。

## 运行指令

```bash
go test ./internal/core -v -count=1 -run '^TestBug011AdvanceFailureHistory$'
```

## 错误信息

失败推进后历史长度从 1 变为 2，说明失败操作已经污染了记录状态。

## 错误堆栈

```text
=== RUN   TestBug011AdvanceFailureHistory
    bug011_advance_failure_history_test.go:23: failed advance changed history: before=1/1 after=2/1
--- FAIL: TestBug011AdvanceFailureHistory (0.00s)
FAIL
FAIL	glimmer-lexicon-studio/internal/core	0.518s
```
