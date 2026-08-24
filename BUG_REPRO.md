# Bug 复现说明

## Bug 是什么

先校验信息不足的词条，再校验完整内容时，第二次结果仍携带第一次请求留下的旧提醒，页面会误报问题。

## 如何触发

连续调用校验入口，第一次传入信息不足的内容，第二次传入完整内容，再比较第二次返回的提醒列表。

## 根因

根因涉及 `internal/core/engine.go` 的 `Engine.ValidatePayload`，以及 `internal/core/concept.go`、`internal/core/term.go` 的校验链路：前一次请求使用的 warnings slice 被后续校验复用，造成跨请求状态污染。

## 运行指令

```bash
go test ./internal/core -v -count=1 -run '^TestBug016ValidationWarningsIsolation$'
```

## 错误信息

第二次校验仍返回了第一次请求遗留的多条提醒。

## 错误堆栈

```text
=== RUN   TestBug016ValidationWarningsIsolation
    bug016_validation_warnings_isolation_test.go:22: warnings leaked across requests: []string{"term expects richer context", "variant expects richer context", "context expects richer context", "proposal expects richer context", "review expects richer context", "relation expects richer context", "synonym expects richer context", "antonym expects richer context", "forbidden expects richer context", "preferred expects richer context", "deprecated expects richer context", "replacement expects richer context", "locale expects richer context", "script expects richer context", "grammar expects richer context", "gender expects richer context", "number expects richer context", "tone expects richer context", "brand expects richer context", "product expects richer context", "feature expects richer context", "interface expects richer context", "error expects richer context", "command expects richer context", "measurement expects richer context", "unit expects richer context", "date expects richer context", "numberformat expects richer context", "punctuation expects richer context", "capitalization expects richer context", "abbreviation expects richer context", "acronym expects richer context", "token expects richer context", "phrase expects richer context", "example expects richer context", "source expects richer context", "translator expects richer context", "reviewer expects richer context", "comment expects richer context", "decision expects richer context", "status expects richer context", "version expects richer context", "release expects richer context", "bundle expects richer context", "manifest expects richer context", "checksum expects richer context", "importer expects richer context", "exporter expects richer context", "memory expects richer context", "match expects richer context", "ambiguity expects richer context", "exception expects richer context", "expiry expects richer context", "notification expects richer context", "history expects richer context", "search expects richer context", "access expects richer context", "workspace expects richer context", "language expects richer context", "scriptmap expects richer context", "quality expects richer context"}
--- FAIL: TestBug016ValidationWarningsIsolation (0.00s)
FAIL
FAIL	glimmer-lexicon-studio/internal/core	0.723s
```
