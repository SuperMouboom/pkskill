---
description: Ingest OpenCode conversation data to pk database.
---

# pk1 - 会话数据转存

> **参见：** [pkskill/SKILL.md](../SKILL.md) （通用规则、路径占位符、命令格式）

## 命令格式

| 格式 | 说明 |
|------|------|
| `pk1` | 使用默认路径（pk1.exe自动检测） |
| `pk1++<source>` | 指定OpenCode存储路径 |
| `pk1++<source>++<target>` | 指定源路径和目标路径 |

**注意：** pk1使用**++分隔符**，不是空格或--flag格式。

## pk1++<path> 格式示例

```bash
# 指定OpenCode存储路径
pk1++{HOME}/.local/share/opencode/storage

# 指定源路径和目标路径
pk1++{HOME}/.local/share/opencode/storage++{PK_PATH}
pk1++{OPENCODE_STORAGE}++{PK_PATH}
```

## 路径说明

pk1.exe自动检测所有路径：
- OpenCode存储：`{HOME}/.local/share/opencode/storage`
- pk1输出：`{PK_PATH}/pk1_output`
- pk_path_config：`{PK_PATH}/pk1_output/pk_path_config.json`
