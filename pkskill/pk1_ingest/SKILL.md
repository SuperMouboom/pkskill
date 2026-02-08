---
name: pkskill-pk1
description: Ingest OpenCode conversation data to pk database. pk1.exe auto-detects OS, username, and paths.
requires: pkskill/SKILL.md
---

# pk1 - OpenCode会话数据转存

> **强制：** 修改pk1前必须阅读 [pkskill/SKILL.md](../SKILL.md)

## 独立运作信息

### pk1核心功能

pk1从OpenCode存储目录提取**完整**会话数据到pk数据库：

1. **扫描** `{OPENCODE_STORAGE}/message/` 获取会话列表
2. **提取** `{OPENCODE_STORAGE}/part/` 获取完整内容
3. **关联** `{OPENCODE_STORAGE}/session_diff/` 获取代码差异
4. **过滤** 空消息
5. **格式化** Markdown
6. **输出** 到 `{PK_PATH}/pk1_output/`
7. **记录** SHA256 hash用于增量更新
8. **生成** `pk_path_config.json` 供pk2/pk3/pk4使用

### 数据源结构

```
{OPENCODE_STORAGE}/
├── message/{sessionID}/              # 消息元数据
├── part/{messageID}/                  # 完整消息内容
├── session_diff/{sessionID}.json     # 代码差异
└── session/global/{sessionID}.json   # 会话上下文
```

### 输出结构

```
{PK_PATH}/
├── pk1_output/
│   ├── conv_{会话ID}_{时间戳}.json
│   ├── conv_{会话ID}_{时间戳}.md
│   └── pk_path_config.json           # pk1生成的路径配置
└── pk1_ingest_metadata.json
```

### 独立CLI使用

```bash
# 使用默认路径（自动检测）
{PK_PATH}/pk1.exe

# 指定源路径
{PK_PATH}/pk1.exe "{HOME}/.local/share/opencode/storage"

# 指定源路径和目标路径
{PK_PATH}/pk1.exe "{HOME}/.local/share/opencode/storage" "{PK_PATH}"
```

## 相关Skills

- [pk2_tool](../pk2_tool/SKILL.md) - 数据增强（pk1的输出）
- [pk3_lab](../pk3_lab/SKILL.md) - Skill实验室
- [pk4_emerge](../pk4_emerge/SKILL.md) - Skill涌现

## Version

- **v4.4** (2026-02-08): 精简文档
  - 引用主SKILL.md通用规则
  - 只保留pk1独立运作信息
