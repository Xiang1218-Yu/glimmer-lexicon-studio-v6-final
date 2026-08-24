# Bug 复现说明

## Bug 是什么

推进阶段失败时接口返回错误，但记录版本已经增加，客户端重试会误以为发生过一次成功变更。

## 如何触发

记录推进到不被接受的阶段后读取版本号，比较失败调用前后的版本是否保持不变。

## 根因

根因位于 `internal/core/engine.go` 的 `Engine.Advance`，并涉及 `internal/core/review.go`、`internal/core/reviewer.go` 的阶段错误处理链路：阶段错误返回前先修改了版本，导致失败路径与客户端重试状态不一致。

## 运行指令

```bash
go test ./internal/core -v -count=1 -run '^TestBug020AdvanceErrorVersion$'
```

## 错误信息

阶段推进失败后版本从 1 变为 2，错误路径没有保持记录版本不变。

## 错误堆栈

```text
=== RUN   TestBug020AdvanceErrorVersion
    bug020_advance_error_version_test.go:20: error path changed version: before=1 after=2
--- FAIL: TestBug020AdvanceErrorVersion (0.00s)
FAIL
FAIL	glimmer-lexicon-studio/internal/core	0.468s
```
