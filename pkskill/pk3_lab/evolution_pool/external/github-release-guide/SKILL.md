---
name: github-release-guide
description: Professional GitHub project release automation skill - helps AI automate complete release workflow including version management, changelog, git operations, and GitHub release creation.
---

# GitHub项目发布指南

> **强制：** 修改本skill前必须阅读 [pkskill/SKILL.md](../pkskill/SKILL.md) （通用规则、路径占位符、命令格式、外部Skill机制）

## 核心定位

github-release-guide是**GitHub项目自动化发布专家**，帮助AI按照用户描述自动化完成专业级的GitHub项目发布流程。

```
用户描述 → AI解析需求 → 执行发布流程 → 完成GitHub Release
```

## 核心功能

1. **智能需求解析** - 理解用户发布需求，识别发布类型
2. **版本号管理** - 自动遵循语义化版本规范(SemVer)
3. **变更日志生成** - 自动从git历史生成CHANGELOG
4. **Git操作** - 自动执行commit、tag、push等操作
5. **GitHub Release** - 自动创建GitHub Release页面
6. **发布检查** - 自动检查发布前准备状态

## 使用场景

- 新版本发布（major/minor/patch）
- 预发布版本（alpha/beta/rc）
- 修复版本发布
- 功能版本发布
- 一次性发布多个仓库

## 使用方法

### 基础发布

```bash
# 标准发布流程
github-release --type patch --message "fix: 修复登录bug"

# 主要版本发布
github-release --type major --message "feat: 全新UI设计"

# 次要版本发布  
github-release --type minor --message "feat: 新增用户配置文件功能"
```

### 高级发布

```bash
# 预发布版本
github-release --type alpha --version 1.0.0-alpha.1
github-release --type beta --version 1.0.0-beta.1
github-release --type rc --version 1.0.0-rc.1

# 自定义版本号
github-release --version 2.1.0 --message "feat: 重大更新"

# 跳过某些步骤
github-release --type patch --skip-tests --force
```

### 交互式发布

```bash
# 引导式发布流程
github-release --interactive

# 只生成变更日志
github-release --changelog-only

# 只创建Git标签
github-release --tag-only
```

## 发布流程

### 第一阶段：需求解析

```
用户输入 → 解析发布类型 → 提取版本信息 → 生成发布计划
```

**发布类型识别：**
- `patch` - 补丁发布（bug修复，小改动）
- `minor` - 次要版本（新功能，向后兼容）
- `major` - 主要版本（重大变更，可能不兼容）
- `alpha` - 内部测试版
- `beta` - 公开测试版
- `rc` - 发布候选版
- `custom` - 自定义版本

**版本号规则（SemVer）：**
```
MAJOR.MINOR.PATCH[-PRERELEASE][+BUILD]
例子：
- 1.0.0 正式版
- 1.0.1 补丁版
- 1.1.0 次要版本
- 2.0.0 主要版本
- 1.0.0-alpha.1 预发布版
- 1.0.0+build.123 构建号
```

### 第二阶段：预发布检查

```
检查项 → 状态 → 操作
```

**必须检查项：**
- [ ] Git工作区清洁（无未提交更改）
- [ ] 远程仓库连接正常
- [ ] 有推送权限
- [ ] 所有测试通过
- [ ] 代码符合规范

**建议检查项：**
- [ ] 版本号无冲突
- [ ] 没有未合并的PR
- [ ] 文档已更新
- [ ] 依赖项已更新

### 第三阶段：变更准备

```
git diff → 变更收集 → CHANGELOG生成 → 版本号更新
```

**变更收集规则：**
```
feat: → 新功能（MINOR版本）
fix: → Bug修复（PATCH版本）
docs: → 文档更新（PATCH版本）
style: → 代码格式（不改变语义）
refactor: → 重构（MINOR或PATCH）
perf: → 性能优化（PATCH版本）
test: → 测试相关
chore: → 构建工具相关
break: → 破坏性变更（MAJOR版本）
```

**CHANGELOG模板：**

```markdown
# Changelog

## [版本号] - 日期

### Features
- 新功能描述

### Bug Fixes
- 修复内容

### Breaking Changes
- 破坏性变更

### Other
- 其他变更
```

### 第四阶段：Git操作

```
git add . → git commit -m "release: 版本号" → git tag -a v版本号 -m "Release 版本号" → git push && git push --tags
```

**Commit消息规范：**
```
feat: 新功能（MINOR）
fix: 修复bug（PATCH）
docs: 文档更新（PATCH）
style: 代码格式（无版本号变更）
refactor: 重构（MINOR/PATCH）
perf: 性能优化（PATCH）
test: 测试（无版本号变更）
chore: 构建（无版本号变更）

带scope的格式：
feat(core): 新增核心功能
fix(auth): 修复登录bug
```

### 第五阶段：GitHub Release

```
GitHub API → 创建Release → 上传资产 → 发布通知
```

**Release页面要素：**
- 标题：版本号或发布名称
- 描述：变更日志内容
- 标签：选择合适的标签
- 预发布：如果是alpha/beta/rc版本
- 资产：上传编译产物、压缩包等

**GitHub API命令示例：**
```bash
# 使用gh CLI
gh release create v1.0.0 --notes "Release notes here" --title "Version 1.0.0"

# 使用git tag
git tag -a v1.0.0 -m "Release 1.0.0"
git push origin v1.0.0
```

## 命令详细说明

### github-release 命令参数

| 参数 | 说明 | 示例 |
|------|------|------|
| `--type` | 发布类型 | `--type patch` |
| `--version` | 自定义版本号 | `--version 1.2.0` |
| `--message` | 发布消息 | `--message "fix: bug"` |
| `--changelog` | 自定义changelog | `--changelog CHANGELOG.md` |
| `--skip-tests` | 跳过测试 | `--skip-tests` |
| `--skip-changelog` | 跳过changelog生成 | `--skip-changelog` |
| `--force` | 强制发布 | `--force` |
| `--interactive` | 交互模式 | `--interactive` |
| `--dry-run` | 演练模式 | `--dry-run` |
| `--assets` | 上传资产 | `--assets "*.zip,*.exe"` |
| `--draft` | 草稿发布 | `--draft` |
| `--prerelease` | 预发布 | `--prerelease` |
| `--repo` | 指定仓库 | `--repo owner/repo` |

### 发布类型与版本号自动递增

| 当前版本 | patch | minor | major |
|---------|-------|-------|-------|
| 1.0.0 | 1.0.1 | 1.1.0 | 2.0.0 |
| 1.2.3 | 1.2.4 | 1.3.0 | 2.0.0 |
| 2.0.0-rc.1 | 2.0.0-rc.2 | 2.1.0 | 3.0.0 |

## 发布检查清单

### 发布前检查

```bash
# 1. 检查工作区状态
git status

# 2. 检查远程连接
git remote -v

# 3. 检查权限
gh repo view --json owner,name

# 4. 运行测试
npm test  # 或其他测试命令

# 5. 检查代码规范
npm run lint  # 或其他lint命令
```

### 版本号检查

```bash
# 查看当前版本
cat package.json | grep version

# 查看最近的tag
git describe --abbrev=0 --tags

# 查看版本差异
git log --oneline $(git describe --abbrev=0 --tags)..HEAD
```

### GitHub Release检查

```bash
# 查看已存在的releases
gh release list

# 查看特定release
gh release view v1.0.0

# 删除错误的release（谨慎使用）
gh release delete v1.0.0 --yes
```

## 高级用法

### 多仓库同时发布

```bash
# 发布列表中的所有仓库
github-release --multi-repo repos.txt

# repos.txt格式：
# owner/repo1
# owner/repo2
# owner/repo3
```

### 定时发布

```bash
# 使用cron定时发布
0 9 * * 1 github-release --type minor --repo owner/repo --force

# 每周一发布一次次要版本
```

### 发布回滚

```bash
# 回滚到上一个版本
github-release --rollback

# 回滚到指定版本
github-release --rollback --version 1.0.0
```

## 最佳实践

### 1. 发布前准备

- [ ] 确保代码已全面测试
- [ ] 更新文档和README
- [ ] 检查依赖项安全性
- [ ] 清理不必要的文件
- [ ] 确认变更日志准确

### 2. 版本号管理

- 严格遵循SemVer规范
- 预发布版本要有明确标识
- 不要跳过版本号
- 记录版本号变更原因

### 3. GitHub Release页面

- 使用清晰的标题
- 详细描述变更内容
- 添加升级指南（如果需要）
- 链接相关Issue和PR
- 上传必要的资产文件

### 4. 沟通与通知

- 在Release中说明变更
- 通知相关团队成员
- 更新内部文档
- 考虑用户通知方式

## 错误处理

### 常见错误

| 错误 | 原因 | 解决方案 |
|------|------|----------|
| push rejected | 权限不足 | 检查token权限 |
| tag already exists | 标签重复 | 删除旧标签或使用新版本 |
| merge conflict | 代码冲突 | 先pull最新代码 |
| tests failed | 测试失败 | 修复测试后再发布 |
| version conflict | 版本号冲突 | 使用新版本号 |

### 错误恢复

```bash
# 取消本地commit
git reset --soft HEAD~1

# 删除本地tag
git tag -d v1.0.0

# 删除远程tag
git push origin :refs/tags/v1.0.0

# 取消已发布的release
gh release delete v1.0.0 --yes
```

## 与其他工具集成

### CI/CD集成

```yaml
# GitHub Actions示例
name: Release
on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Build
        run: npm build
      - name: Release
        run: |
          gh release create ${{ github.ref_name }} \
            --title "Release ${{ github.ref_name }}" \
            --notes "$(git log --oneline --reverse --format='- %s' ${{ github.ref_name }}^..${{ github.ref_name }})"
```

### npm发布集成

```bash
# npm publish后再创建GitHub Release
npm publish && github-release --type patch --repo owner/repo
```

## 相关Skills

- [pk1_ingest](../pk1_ingest/SKILL.md) - 数据摄入
- [pk2_tool](../pk2_tool/SKILL.md) - pk2智能总结
- [pk3_lab](../pk3_lab/SKILL.md) - pk3实验室
- [pk4_emerge](../pk4_emerge/SKILL.md) - pk4涌现

## Version

- **v1.0** (2026-02-08): 初始版本
  - 基础发布流程（patch/minor/major）
  - 语义化版本管理
  - 自动CHANGELOG生成
  - Git操作自动化
  - GitHub Release创建
  - 预发布版本支持（alpha/beta/rc）
  - 交互式发布模式
  - 演练模式支持
  - 错误处理与恢复
  - CI/CD集成示例
