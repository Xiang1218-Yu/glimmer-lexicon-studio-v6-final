# Bug 复现说明

## Bug 是什么

客户端读取词条后整理证据列表，下一次查看时发现服务内保存的证据也被改掉，页面因此出现脏数据。

## 如何触发

创建一条带证据的词条，调用 `Engine.Get` 获取记录并修改返回结果中的证据 slice，再次读取同一条记录并比较证据内容。调用方修改返回值不应污染引擎内部状态。

## 根因

根因涉及 `internal/core/engine.go`、`internal/core/concept.go` 和 `internal/core/term.go` 的 `Engine.Get` 调用链：查询结果中的证据 slice 与内存记录共享底层数组，调用方修改查询结果后，服务内记录被同步污染。该问题属于 slice 底层数组复用造成的状态污染。

## 运行指令

```bash
go test ./internal/core -v -count=1 -run '^TestBug003RecordEvidenceIsolation$'
```

## 错误信息

修改查询返回值后，服务内保存的证据内容也发生变化，说明返回对象没有与内部状态隔离。

## 错误堆栈

```text
=== RUN   TestBug003RecordEvidenceIsolation
    bug003_record_evidence_isolation_test.go:27: caller mutation changed the stored evidence slice
--- FAIL: TestBug003RecordEvidenceIsolation (0.00s)
FAIL
FAIL	glimmer-lexicon-studio/internal/core	0.739s
FAIL
```
