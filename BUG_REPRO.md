# Bug 复现说明

## Bug 是什么

管理页面读取模块列表后做本地排序，随后新建词条时引擎出现空模块崩溃，说明返回的模块集合与内部状态发生了串联。

## 如何触发

读取 `Engine.Modules` 返回的模块列表并修改其中的元素，再创建一条词条并观察创建流程。外部修改模块快照不应改变引擎内部模块，也不应触发运行时崩溃。

## 根因

根因涉及 `internal/core/engine.go`、`internal/core/concept.go` 和 `internal/core/term.go` 的 `Engine.Modules` 调用链：返回的模块 slice 与引擎内部 slice 共享底层数组，调用方修改快照后，后续创建流程读取到 nil 模块并崩溃。该问题属于并发访问下的状态污染与一致性失效。

## 运行指令

```bash
go test ./internal/core -v -count=1 -run '^TestBug004ModuleSnapshotIsolation$'
```

## 错误信息

修改返回的模块快照后，后续创建流程读取到了空模块并崩溃。

## 错误堆栈

```text
=== RUN   TestBug004ModuleSnapshotIsolation
    bug004_module_snapshot_isolation_test.go:20: mutating a returned module snapshot crashed the engine
--- FAIL: TestBug004ModuleSnapshotIsolation (0.00s)
FAIL
FAIL	glimmer-lexicon-studio/internal/core	0.700s
FAIL
```
