# Bug 复现说明

## Bug 是什么

服务尚未完成引擎初始化时，模块列表请求返回内部故障而不是稳定的未就绪响应，监控端无法区分服务未就绪和服务崩溃。

## 如何触发

在服务对象存在但引擎尚未注入时请求模块列表，观察 HTTP 状态和错误响应格式。

## 根因

根因涉及 `cmd/glimmer-lexicon-studio/main.go`、`internal/api/modules.go`、`internal/api/server.go` 和 `internal/core/engine.go` 的模块列表调用链：nil engine 被解引用后进入恢复中间件，最终被映射为 500，而不是未就绪响应。

## 运行指令

```bash
go test ./internal/api -v -count=1 -run '^TestBug025NilEngineModulesResponse$'
```

## 错误信息

nil engine 的模块列表请求返回 500，实际 body 为 `{"error":"internal server error"}`，预期应为 503。

## 错误堆栈

```text
=== RUN   TestBug025NilEngineModulesResponse
2026/08/24 12:11:40 INFO request service=glimmer-lexicon-studio method=GET path=/v1/modules duration=447.042µs
    bug025_nil_engine_modules_response_test.go:17: expected 503 for nil engine, got 500 body={"error":"internal server error"}
--- FAIL: TestBug025NilEngineModulesResponse (0.00s)
FAIL
FAIL	glimmer-lexicon-studio/internal/api	0.483s
```
