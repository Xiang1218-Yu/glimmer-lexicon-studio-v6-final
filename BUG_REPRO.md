# Bug 复现说明

## Bug 是什么

客户端提交格式错误的记录请求后，连接迟迟不释放，连续重试会消耗服务端连接资源。

## 如何触发

向记录创建接口提交无法解码的请求体，观察解码失败路径是否关闭请求体。

## 根因

根因涉及 `internal/api/codec.go`、`internal/api/records.go`、`internal/api/middleware.go` 的 `POST /v1/records` 调用链：请求体只在成功路径释放，解码错误路径没有完成资源生命周期。

## 运行指令

```bash
go test ./internal/api -v -count=1 -run '^TestBug026RequestBodyClosedOnDecodeError$'
```

## 错误信息

解码失败后请求体仍未关闭，客户端持续重试时会累积连接资源。

## 错误堆栈

```text
=== RUN   TestBug026RequestBodyClosedOnDecodeError
2026/08/24 12:11:41 INFO request method=POST path=/v1/records duration=128.875µs
    bug026_request_body_closed_on_decode_error_test.go:25: request body was not closed on the decoder error path
--- FAIL: TestBug026RequestBodyClosedOnDecodeError (0.00s)
FAIL
FAIL	glimmer-lexicon-studio/internal/api	0.526s
```
