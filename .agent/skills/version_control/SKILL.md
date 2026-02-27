---
name: Version Control
description: 严格的代码版本管理技能，在AI完成代码修改后自动进行 commit 并 push 到 GitHub。
---

# Version Control Skill

这套技能赋予你（AI）在完成每个阶段的代码修改后，强制并自动执行版本控制的职责，以确保用户的代码安全并同步到远程 GitHub 仓库。

## ⏳ 触发条件 (When to Use)
- 当你完成了一轮代码编写、Debug、或者需求修改，**且验证无误，准备向用户报告 "任务完成" 或结束回答前**。
- 当用户主动要求你“提交一下代码”或“推送到远端”时。

## 📝 严格执行步骤 (Execution Steps)

一旦触发本技能，你【必须】自动使用 `run_command` 工具依次执行以下操作：

### 第一步：检查状态 (Check Status)
运行命令：
```sh
git status
```
*如果输出显示 `nothing to commit, working tree clean`，则说明没有改动，你可以直接向用户说明“没有检测到更改，无需提交”，然后终止后续流程。否则，继续下一步。*

### 第二步：暂存更改 (Stage Changes)
运行命令：
```sh
git add .
```

### 第三步：生成并应用提交 (Commit)
根据你刚才为用户所做的所有修改，总结出一句高度概括且专业的 `commit message`（中英文不限，如："Refactor: Extract NPC interactions to builder and controllers" 或 "Fix: 修复双重边界限制与雷达识别问题"）。
然后运行命令：
```sh
git commit -m "<生成的准确提交信息>"
```

### 第四步：推送到 GitHub (Push)
运行命令：
```sh
git push
```

## 🚨 必须遵守的红线规则 (Critical Rules)
1. **静默高效**：不要问用户“我是否可以提交”，只要你确认代码修改完成且运行没有报错，**请大胆且自觉地执行上述全套流程**（可以使用多个相关的命令调用同时或连续执行）。
2. **禁止强制推送**：如果 `git push` 失败报错（例如：远程有新的变更导致冲突，`git checkout` 错误等），【立刻停止推送】，然后向用户抛出异常报告，【绝对严禁】擅自使用 `git push -f` 覆盖远端代码！
3. **忽略多余文件**：如果在第一步 `git status` 看到了大批量的临时文件，且项目中存在 `.gitignore`，你只需正常执行 `git add .` 即可（Git 会自动忽略）。不要尝试修改忽略规则，除非用户明确要求。