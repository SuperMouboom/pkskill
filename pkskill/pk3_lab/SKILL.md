---
name: pkskill-pk3-lab
description: Analyze pk2 database, generate skills, manage evolution pool, and perform skill fusion.
requires: pkskill/SKILL.md
---

# pk3_lab - AI Skill实验室

> **强制：** 修改pk3前必须阅读 [pkskill/SKILL.md](../SKILL.md) （通用规则、路径占位符、命令格式、外部Skill机制）

## 核心定位

pk3是pkskill系统的**AI辅助Skill实验室**：
- 分析pk2增强数据库中的对话痛点
- 自动生成针对性skill
- 管理进化池实现skill持续进化与优化
- 支持外部Skill导入和管理

```
pk1转存 → pk2增强 → pk3分析痛点 → 生成Skill → pk4涌现
                                      ↑
                              pk3进化池管理
```

## 命令格式

| 格式 | 说明 |
|------|------|
| `pk3` | 默认（pk2数据库） |
| `pk3++<path>` | 智能识别：导入skill 或 指定数据库 |
| `pk3++go` | 触发融合进化 |
| `pk3++gap` | 刷新排名表 |
| `pk3++see` | 生成表现报告 |

## pk3++<path> 智能识别

```
输入: pk3++{PK_PATH}/external_skill/
       ↓
检测路径下是否存在SKILL.md
       ↓
├─ 存在SKILL.md → 导入外部skill → 放到 external/{skill_name}/
└─ 不存在SKILL.md → 指定数据库路径 → 扫描数据生成skill
```

## 数据库源优先级

| 优先级 | 数据源 | 说明 |
|--------|--------|------|
| 1 | pk2_enriched/enriched | v2摘要结构化数据，质量最高 |
| 2 | pk1_output | pk1原始对话，降级使用 |
| 3 | 用户指定路径 | 通过`pk3++<path>`格式 |

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

### pk3++see 报告中的外部Skill分析

外部Skill的独特价值在于：记录哪些问题是外部Skill无法解决的，以及这些问题是否被进化池Skill解决。

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

## Skill模板

```markdown
---
name: pkskill-{能力英文名}-v{版本号}
description: {一句话描述}
---

# {能力中文名}

## 核心功能

{3-5句话描述这个skill能做什么}

## 使用场景

- {场景1}
- {场景2}
- {场景3}

## 使用方法

```bash
{具体命令}
```

## 示例

```bash
{使用示例}
```
```

## pk3++go - 融合进化

触发融合进化流程，只融合pk3生成的skill。

**融合规则：**
- 只融合skills/目录下的skill
- 不融合external/目录下的外部skill
- 融合结果放在fusion/目录
- 命名格式：`pkskill-{能力}-fused-v{版本}.md`

## pk3++gap - 刷新排名表

刷新ranking_table.json，计算所有skill的最新成功率。

**排名指标：**
- pk3生成skill：成功率、调用次数
- 外部skill：成功率、调用次数、**被进化池解决**

## pk3++see - 生成表现报告

生成进化池skill和外部skill的表现详细报告。

**报告内容：**
1. 进化池skill表现统计
2. 外部skill表现统计
3. 对比分析
4. 外部skill无法解决的问题列表
5. 进化池skill解决这些问题的情况
6. 改进建议

## 进化池结构

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

## 相关Skills

- [pk1_ingest](../pk1_ingest/SKILL.md) - 数据源
- [pk2_tool](../pk2_tool/SKILL.md) - pk3的输入数据源
- [pk4_emerge](../pk4_emerge/SKILL.md) - pk3的输出用户

## Version

- **v2.0** (2026-02-08): 外部Skill机制 + pk3++go/gap/see
  - pk3++<path> 智能识别（导入skill 或 指定数据库）
  - pk3++go 融合进化（不融合外部skill）
  - pk3++gap 刷新排名表
  - pk3++see 生成表现报告
  - 外部Skill机制（不参与融合，参与排名）
