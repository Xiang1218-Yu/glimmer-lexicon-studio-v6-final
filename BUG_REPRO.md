# Bug 复现说明

## Bug 是什么

监控快照把已经推进过的词条继续计入草稿阶段，运营据此判断流程积压时会得到错误结论。

## 如何触发

创建词条并推进到其他阶段，再调用 `Engine.Snapshot`，比较快照中的阶段计数和词条当前实际阶段。

## 根因

根因涉及 `internal/core/engine.go` 的 `Engine.Snapshot`，以及 `internal/core/history.go`、`internal/core/search.go` 的状态读取链路：快照只按记录总量增加 `stage_draft`，没有根据每条记录的当前阶段建立对应计数。

## 运行指令

```bash
go test ./internal/core -v -count=1 -run '^TestBug015SnapshotStageCounts$'
```

## 错误信息

快照仍报告 `stage_draft=1`，没有反映词条已经推进到其他阶段。

## 错误堆栈

```text
=== RUN   TestBug015SnapshotStageCounts
    bug015_snapshot_stage_counts_test.go:19: snapshot counts are stale: map[string]int{"modules":62, "records":1, "stage_draft":1}
--- FAIL: TestBug015SnapshotStageCounts (0.00s)
FAIL
FAIL	glimmer-lexicon-studio/internal/core	0.224s
```
