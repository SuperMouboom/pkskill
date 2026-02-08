---
name: pkskill
description: Use when you want to build an AI capability evolution system with data ingestion, enrichment, skill generation, and emergence features.
---

# pkskill - AI能力进化系统

## Overview

pkskill是一个自研自进的AI能力进化系统，通过4个子skill实现：
- **pk1**: 会话数据转存与数据库建立（完整内容 + Markdown）
- **pk2**: 数据增强与智能总结
- **pk3**: Skill自动生成与进化
- **pk4**: Skill涌现与优先级调用

---

## 🔒 修改门禁规则

> **强制：** AI试图修改pkskill任意部分前，必须先阅读此文档（SKILL.md）

**修改前检查清单**：
```
□ 是否涉及路径？ → 检查占位符 {PK_PATH}, {HOME}, {USERNAME}
□ 是否涉及命令？ → 检查命令格式表
□ 是否修改配置文件？ → 移除硬编码
□ 是否只改一个文件？ → 检查关联文件
```

**违反后果**：AI必须重新阅读规则后重试。

---

# 系统契约

> **重要：** AI修改pkskill前必须阅读此部分

## 核心规则

### 路径规则

| 规则 | 说明 |
|------|------|
| **占位符** | 所有路径使用 `{PK_PATH}`, `{HOME}`, `{USERNAME}` |
| **禁止硬编码** | 不允许出现 `C:\Users\GGBOM`, `D:/pkskill` 等具体路径 |
| **禁止bat脚本** | 不创建 pk1.bat, pk2.bat 等，所有命令通过 skill 调用 |
| **自动检测** | pk1.exe 自动检测操作系统、用户名、主目录，无需手动配置 |

### 命令规则（++分隔符）

| 类型 | Skill | 说明 |
|------|-------|------|
| 可执行 | pk1 | pk1.exe（需要编译） |
| AI行为准则 | pk2 | AI直接加载 skill 执行 |
| AI行为准则 | pk3 | AI直接加载 skill 执行 |
| AI行为准则 | pk4 | AI直接加载 skill 执行 |

**pk1 可执行命令**：
```
pk1++<source>           # 指定OpenCode存储路径
pk1++<source>++<target> # 指定源路径和目标路径
```

### 语言选择规则

所有操作前必须询问用户语言偏好：

```
请选择会话语言 / Please select conversation language:

[1] 中文 (Chinese)
[2] English (英文)

输入编号或直接输入语言 / Enter number or language:
```

**规则**：
- 默认选项：中文 / English
- 选择后，所有问答和输出使用该语言
- 代码/技术术语保持英文（如 skills_used, tools_used）
- 子 SKILL 引用此规则，无需重复

### 文件权威顺序

1. **pk1.go** - 代码真相，具体实现
2. **config.json** - 配置默认值
3. **SKILL.md** - 文档
4. **commands/*.md** - 命令引用

### 修改规则

当用户要求修改时：

1. **改路径规则** → 更新此文档（SKILL.md）
2. **改命令格式** → 更新此文档 + commands/*.md
3. **改config.json** → 检查是否有硬编码
4. **改子SKILL.md** → 检查是否要更新关联文档
5. **用户要改命令** → 直接列出命令格式表，不反问用户
6. **创建/删除文件** → 先判断在"文件权威顺序"中的位置
7. **修改任意部分前** → 先读 SKILL.md 开头建立全局认知

### AI检查清单

修改前快速检查：

```
- [ ] 是否涉及路径？ → 检查占位符使用
- [ ] 是否涉及命令？ → 检查命令格式表
- [ ] 是否修改了配置文件？ → 移除硬编码
- [ ] 是否只改了一个文件？ → 检查关联文件
- [ ] 是否要创建/删除文件？ → 判断在文件权威顺序中的位置
- [ ] 是否修改任意部分？ → 先读 SKILL.md 开头
```

---

## pk1 - 会话数据转存

### 命令格式

| 格式 | 说明 |
|------|------|
| `pk1` | 使用默认路径（pk1.exe自动检测） |
| `pk1++<source>` | 指定OpenCode存储路径 |
| `pk1++<source>++<target>` | 指定源路径和目标路径 |

### pk1++<path> 格式

```
pk1++{HOME}/.local/share/opencode/storage
pk1++{OPENCODE_STORAGE}++{PK_PATH}
```

---

## pk2 - 数据增强

### 命令格式

| 格式 | 说明 |
|------|------|
| `pk2` | 使用默认路径（pk1_output） |
| `pk2++<path>` | 指定pk1数据库路径 |

### pk2++<path> 格式

```
pk2++{PK_PATH}/pk1_output
```

---

## pk3 - Skill实验室

### 命令格式

| 格式 | 说明 |
|------|------|
| `pk3` | 默认（pk2数据库） |
| `pk3++<path>` | 智能识别：导入skill 或 指定数据库 |
| `pk3++go` | 触发融合进化 |
| `pk3++gap` | 刷新排名表 |
| `pk3++see` | 生成表现报告 |

### pk3++<path> 智能识别

```
pk3++{PK_PATH}/external_skill/
       ↓
检测路径下是否存在SKILL.md
       ↓
├─ 存在SKILL.md → 导入外部skill → 放到 external/{skill_name}/
└─ 不存在SKILL.md → 指定数据库路径 → 扫描数据生成skill
```

### pk3++go - 融合进化

触发pk3_lab/evolution_pool/fusion/目录下的融合进化流程。

### pk3++gap - 刷新排名表

刷新pk3_ranking_table.json，计算所有skill的最新成功率。

### pk3++see - 生成表现报告

生成进化池skill和外部skill的表现详细报告：

```
报告内容：
1. 进化池skill表现统计
2. 外部skill表现统计
3. 对比分析
4. 外部skill无法解决的问题列表
5. 进化池skill解决这些问题的情况
6. 改进建议
```

---

## pk4 - Skill涌现

### 命令格式

| 格式 | 说明 |
|------|------|
| `pk4` | 默认优先级：**外部skill > pk3进化池** |
| `pk4++<order>` | 自定义优先级顺序 |
| `pk4++sense` | 感知模式：检测上下文，建议可能需要的skill |

### pk4++sense - 感知模式

pk4++sense是**非强制**的感知-建议模式：

```
用户问题 → pk4++sense
                  ↓
        检测当前对话上下文
                  ↓
        扫描可用skill（文件名语义）
                  ↓
        建议："我注意到这个问题可能需要 xxx skill"
                  ↓
        AI自己决定是否调用
```

**与pk4优先级模式的核心区别**：

| 对比项 | pk4（优先级模式） | pk4++sense（感知模式） |
|--------|------------------|------------------------|
| 强制程度 | 高 | 低 |
| 调用方式 | 按优先级强制调用 | 建议，AI自己决定 |
| 记录方式 | 强制记录call_records | 可选记录 |
| 氛围营造 | 必须用 | 可以用 |

### pk4++<order> 优先级格式

```
pk4++external                  # 只用外部skill
pk4++pk3                      # 只用pk3进化池
pk4++external++pk3            # 外部skill > pk3进化池
pk4++superpowers++pk3        # superpowers > pk3进化池
pk4++external++superpowers  # 外部skill > superpowers
```

---

## 外部Skill机制

### 外部Skill定义

外部Skill是用户手动导入的第三方Skill，不参与pk3融合进化，但参与pk4调用和排名。

### 外部Skill存放位置

```
{OPENCODE_SKILLS}/pkskill/pk3_lab/evolution_pool/external/
└── {skill_name}/
    └── SKILL.md
```

### 外部Skill规则

| 特性 | 规则 |
|------|------|
| pk4调用优先级 | ✅ **最高**（默认） |
| 参与融合进化 | ❌ 不参与 |
| 参与排名 | ✅ 参与 |
| pk3++go融合 | ❌ 不会被融合 |
| pk3++gap刷新 | ✅ 参与排名计算 |

### 外部Skill排名指标

```json
{
  "skill_name": "external-skill-xxx",
  "source": "external",
  "total_calls": 100,
  "total_success": 80,
  "success_rate": 0.8,
  "unsolved_problems": [
    {
      "problem": "问题描述",
      "count": 5,
      "solved_by_evolution_pool": 3,
      "solved_rate": 0.6
    }
  ]
}
```

### pk3++see 报告中的外部Skill分析

外部Skill的独特价值在于：记录哪些问题是外部Skill无法解决的，以及这些问题是否被进化池Skill解决。

```
报告格式：
1. 外部Skill统计（成功率、调用次数）
2. 外部Skill无法解决的问题列表
3. 进化池Skill解决这些问题的情况
4. 对比分析：外部Skill vs 进化池Skill
5. 改进建议：哪些问题需要进化池关注
```

---

## pk3_lab 目录结构

```
{OPENCODE_SKILLS}/pkskill/pk3_lab/evolution_pool/
├── skills/                      # pk3生成的skill
│   └── pkskill-*.md
├── external/                    # 外部skill（新）
│   └── {skill_name}/
│       └── SKILL.md
├── fusion/                     # 融合进化
│   └── pkskill-*-fused-v*.md
├── ranking_table.json         # 排名表
└── call_records/             # 调用记录
    └── pkskill-*.json
```

---

## 跨平台路径支持

pk1.exe 内置跨平台路径检测功能，运行时会自动：
- 检测操作系统（Windows/Linux/macOS）
- 检测用户名（环境变量优先级：USERNAME → USER → LOGNAME）
- 检测主目录（HOME → USERPROFILE → HOMEDRIVE+HOMEPATH）
- 生成 `pk_path_config.json` 供其他组件使用

### 路径占位符

| 占位符 | 说明 | Windows示例 | Linux示例 |
|--------|------|-------------|-----------|
| `{PK_PATH}` | pkskill根目录 | D:/pkskill | /opt/pkskill |
| `{HOME}` | 用户主目录 | C:/Users/username | /home/username |
| `{USERNAME}` | 用户名 | username | username |

### OpenCode路径占位符

| 占位符 | 说明 | Windows示例 | Linux示例 |
|--------|------|-------------|-----------|
| `{OPENCODE_CONFIG}` | OpenCode配置根目录 | C:/Users/username/.config/opencode | /home/username/.config/opencode |
| `{OPENCODE_SKILLS}` | OpenCode全局skills目录 | {OPENCODE_CONFIG}/skills | {OPENCODE_CONFIG}/skills |

### PK_PATH检测优先级

1. 环境变量 `PK_PATH`
2. 配置文件 `~/.config/pkskill/config.json`
3. 平台默认：
   - Windows: `D:/pkskill`
   - Linux: `/opt/pkskill`
   - macOS: `/Users/{USERNAME}/pkskill`

---

## Directory Structure

```
pkskill/
├── SKILL.md                           # 主入口 + 系统契约
├── README.md                          # 使用说明
├── config.json                        # 默认配置
├── commands/                          # 命令引用
│   ├── pk1.md
│   ├── pk2.md
│   ├── pk3.md
│   └── pk4.md
├── pk1_ingest/
│   └── SKILL.md                       # pk1入口
├── pk2_tool/
│   └── SKILL.md                       # pk2智能总结入口
├── pk3_lab/
│   └── SKILL.md                       # pk3入口
└── pk4_emerge/
    └── SKILL.md                       # pk4入口
```

---

## Related Skills

- [pk1_ingest](./pk1_ingest/SKILL.md) - 数据转存
- [pk2_tool](./pk2_tool/SKILL.md) - 数据增强
- [pk3_lab](./pk3_lab/SKILL.md) - Skill实验室
- [pk4_emerge](./pk4_emerge/SKILL.md) - Skill涌现

---

## Version

- **v2.0** (2026-02-08): 外部Skill机制 + 命令格式统一
  - pk1/pk2/pk3/pk4 全部使用++分隔符
  - pk3++<path> 智能识别（导入skill 或 指定数据库）
  - pk3++go / pk3++gap / pk3++see 新命令
  - pk4默认优先级：外部skill > pk3进化池
  - 外部Skill机制（不参与融合，参与排名，pk4调用优先级最高）
  - pk3++see 表现报告
