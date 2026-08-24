# Bug 复现说明

## Bug 是什么

列表接口返回记录后，前端更新其中一条的正文，下一次读取服务端内容也被一起改写，造成词条正文被无意覆盖。

## 如何触发

调用列表接口取得记录，修改返回记录的 payload，再次读取服务端记录并比较正文内容。

## 根因

根因位于 `internal/core/engine.go` 的 `Engine.List`，并涉及 `internal/core/concept.go`、`internal/core/term.go` 的记录复制链路：列表返回对象没有与服务端记录完全隔离，调用方修改结果后改变了共享状态。

## 运行指令

```bash
go test ./internal/core -v -count=1 -run '^TestBug022ListRecordIsolation$'
```

## 错误信息

修改列表结果中的 payload 后，服务端记录正文也发生变化。

## 错误堆栈

```text
=== RUN   TestBug022ListRecordIsolation
    bug022_list_record_isolation_test.go:21: list result mutation changed stored record
--- FAIL: TestBug022ListRecordIsolation (0.00s)
FAIL
FAIL	glimmer-lexicon-studio/internal/core	0.482s
```
