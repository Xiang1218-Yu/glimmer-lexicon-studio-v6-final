# Bug 复现说明

## Bug 是什么

记录请求包含服务不接受的内容时，客户端收到的反馈与普通服务故障相同，无法判断应修改请求还是重试。

## 如何触发

向记录创建接口提交会触发解码或输入边界错误的请求，检查 HTTP 响应是否保留清晰的业务错误信息。

## 根因

根因涉及 `internal/api/codec.go`、`internal/api/records.go`、`internal/api/middleware.go` 的 `POST /v1/records` 调用链：输入错误跨过解码边界后被映射成缺少具体类别的响应，业务错误信息在 HTTP 边界丢失。

## 运行指令

```bash
go test ./internal/api -v -count=1 -run '^TestBug028DecodeErrorBoundary$'
```

## 错误信息

解码器的具体错误没有传递到 HTTP 响应，调用方只能看到不清晰的错误反馈。

## 错误堆栈

```text
=== RUN   TestBug028DecodeErrorBoundary
2026/08/24 12:11:42 INFO request method=POST path=/v1/records duration=313.166µs
    bug028_decode_error_boundary_test.go:20: decoder detail was lost at the HTTP boundary
--- FAIL: TestBug028DecodeErrorBoundary (0.00s)
FAIL
FAIL	glimmer-lexicon-studio/internal/api	0.506s
```
