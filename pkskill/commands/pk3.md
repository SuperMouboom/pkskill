---
description: Analyze pk database and generate new skills.
---

# pk3 - Skill实验室

> **参见：** [pkskill/SKILL.md](../SKILL.md) （通用规则、路径占位符、命令格式、外部Skill机制）

## 命令格式

| 格式 | 说明 |
|------|------|
| `pk3` | 默认（pk2数据库） |
| `pk3++<path>` | 智能识别：导入skill 或 指定数据库 |
| `pk3++go` | 触发融合进化 |
| `pk3++gap` | 刷新排名表 |
| `pk3++see` | 生成表现报告 |

**注意：** pk3使用**++分隔符**，不是空格或--flag格式。

## pk3++<path> 智能识别

```
pk3++{PK_PATH}/external_skill/
       ↓
检测路径下是否存在SKILL.md
       ↓
├─ 存在SKILL.md → 导入外部skill → 放到 external/{skill_name}/
└─ 不存在SKILL.md → 指定数据库路径 → 扫描数据生成skill
```

## pk3++go - 融合进化

触发pk3_lab/evolution_pool/fusion/目录下的融合进化流程。

**融合进化规则：**
- 只融合pk3生成的skill
- 不融合外部skill
- 融合结果放在 fusion/ 目录

## pk3++gap - 刷新排名表

刷新pk3_ranking_table.json，计算所有skill的最新成功率。

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

## 路径说明

pk3从pk2_enriched读取，生成skill到evolution_pool：
- 输入：`{PK_PATH}/pk2_enriched`
- pk3生成skill：`{OPENCODE_SKILLS}/pkskill/pk3_lab/evolution_pool/skills/`
- 外部skill：`{OPENCODE_SKILLS}/pkskill/pk3_lab/evolution_pool/external/{skill_name}/`
