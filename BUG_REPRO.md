# Bug 复现说明

## Bug 是什么

服务对象已经创建但引擎尚未注入时，健康检查返回正常，部署系统会误以为实例可以接收流量。

## 如何触发

创建尚未注入引擎的服务对象，请求 `/readyz`，检查状态和 ready 字段。

## 根因

根因涉及 `internal/api/server.go`、`internal/api/modules.go`、`internal/api/middleware.go` 的 `GET /readyz` 调用链：nil engine 被当成可用服务，健康检查没有把未初始化状态映射为未就绪响应。

## 运行指令

```bash
go test ./internal/api -v -count=1 -run '^TestBug029NilEngineReadiness$'
```

## 错误信息

未注入引擎时仍返回 200 和 `{"ready":true}`，部署系统因此会错误放行流量。

## 错误堆栈

```text
=== RUN   TestBug029NilEngineReadiness
2026/08/24 12:11:43 INFO request service=glimmer-lexicon-studio method=GET path=/readyz duration=220.958µs
    bug029_nil_engine_readiness_test.go:17: expected not-ready response, got 200 body={"ready":true}
--- FAIL: TestBug029NilEngineReadiness (0.00s)
FAIL
FAIL	glimmer-lexicon-studio/internal/api	0.641s
```
