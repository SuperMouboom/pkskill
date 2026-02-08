---
description: Enrich pk database with AI-generated summaries.
---

# pk2 - 数据增强

**参见：** [pkskill/SKILL.md#pk2-tool](../SKILL.md#pk2-tool)

## 命令格式

| 格式 | 说明 |
|------|------|
| `pk2` | 使用默认路径（pk1_output） |
| `pk2++<path>` | 指定pk1数据库路径 |

**注意：** pk2使用**++分隔符**，不是--flag格式。

## 路径说明

pk2从pk1_output读取，增强后输出到pk2_enriched：
- 输入：`{PK_PATH}/pk1_output`
- 输出：`{PK_PATH}/pk2_enriched/enriched`
