# Bug Reproduction

## 包的性质

当前 test_model_fix 保存的是被测模型修复后的结果源码，不是初始含 Bug 源码。要复现原始缺陷，必须检出下面固定的 parent SHA；不要在当前修复结果源码上期待重新出现修复前失败。生成系统使用的可信验证补丁和完整验证日志仅在本地留存，不提交到结果分支。

## 问题现象

报关员给一票资料齐全的货物办理放行时，案件先显示“正在签章”，随后签章服务报错；接口返回失败后案件一直停在这个中间状态，再次提交即使签章服务恢复也只会收到状态冲突。请修复失败后的状态交接：保留可识别的签章错误，把案件恢复到可复核状态，并允许后续成功签章进入 released，同时不能重复追加放行备注。

## 含 Bug 版本

- 仓库：VanceMichael/go-label-17
- 仓库地址：https://github.com/VanceMichael/go-label-17.git
- parent SHA：8b7aec6c18e2e36d779d21150e7b3e7b4c741dae

## 复现步骤

```bash
git clone -- https://github.com/VanceMichael/go-label-17.git bug-repro
cd bug-repro
git checkout --detach 8b7aec6c18e2e36d779d21150e7b3e7b4c741dae
go test ./internal/customs -run ^TestFailedReleaseSigningRestoresReviewableCase$ -count=1
```

## 双架构完整错误信息

### linux/amd64

- 容器内复现预期退出码：1
- 容器内复现实际退出码：1

stdout：

```text
$ go test ./internal/customs -run ^TestFailedReleaseSigningRestoresReviewableCase$ -count=1
--- FAIL: TestFailedReleaseSigningRestoresReviewableCase (0.00s)
    release_test.go:35: status after failed signing = releasing, want pending
FAIL
FAIL	github.com/VanceMichael/go-base-airbridge/internal/customs	0.030s
FAIL

```

stderr：

```text
(empty)
```

### linux/arm64

- 容器内复现预期退出码：1
- 容器内复现实际退出码：1

stdout：

```text
$ go test ./internal/customs -run ^TestFailedReleaseSigningRestoresReviewableCase$ -count=1
--- FAIL: TestFailedReleaseSigningRestoresReviewableCase (0.00s)
    release_test.go:35: status after failed signing = releasing, want pending
FAIL
FAIL	github.com/VanceMichael/go-base-airbridge/internal/customs	0.002s
FAIL

```

stderr：

```text
(empty)
```

## 通过条件

资料齐全的 pending 海关案件进入签章后，如果签章服务返回错误，该错误仍能被 errors.Is 识别，案件状态恢复为 pending 且没有放行备注；服务恢复后对同一案件重试应成功转为 released，并且只追加一条签章备注。TestFailedReleaseSigningRestoresReviewableCase 必须由红转绿，customs 与 domain 回归、全仓 go test ./... 和 go build ./... 均通过，不得吞掉签章失败、放宽状态断言或绕过真实重试路径。
