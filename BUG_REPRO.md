# Bug 复现说明

## Bug 是什么

提交重复编号的记录时，服务返回内部故障状态，客户端不断重试而不是提示编号已存在，监控也把正常校验失败统计成服务器事故。

## 如何触发

先创建一条记录，再提交相同编号，检查第二次请求的 HTTP 状态和错误 body。

## 根因

根因涉及 `internal/api/records.go`、`internal/api/codec.go`、`internal/core/engine.go` 的 `POST /v1/records` 调用链：重复编号的业务错误跨层映射为内部故障状态，调用方无法区分可修复请求和服务崩溃。

## 运行指令

```bash
go test ./internal/api -v -count=1 -run '^TestBug030CreateErrorStatus$'
```

## 错误信息

重复编号请求返回 500 和 `{"error":"record already exists"}`，预期应是客户端错误状态。

## 错误堆栈

```text
=== RUN   TestBug030CreateErrorStatus
2026/08/24 12:11:44 INFO request method=POST path=/v1/records duration=668.167µs
2026/08/24 12:11:44 INFO request method=POST path=/v1/records duration=6.583µs
    bug030_create_error_status_test.go:21: expected client error for duplicate record, got 500 body={"error":"record already exists"}
--- FAIL: TestBug030CreateErrorStatus (0.00s)
FAIL
FAIL	glimmer-lexicon-studio/internal/api	0.601s
```
