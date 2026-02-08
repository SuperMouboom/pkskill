---
name: pkskill-pk4-emerge
description: Activate skill emergence with priority-based invocation.
requires: pkskill/SKILL.md
---

# pk4_emerge - Skill涌现

> **强制：** 修改pk4前必须阅读 [pkskill/SKILL.md](../SKILL.md) （通用规则、路径占位符、命令格式、外部Skill机制）

## 核心定位

pk4营造**all-in-skill**的涌现氛围，确保AI在每个对话中自动识别并使用相应Skill。

## 命令格式

| 格式 | 说明 |
|------|------|
| `pk4` | 默认优先级：**外部skill > pk3进化池** |
| `pk4++<order>` | 自定义优先级顺序 |
| `pk4++sense` | 感知模式：检测上下文，建议可能需要的skill |

## pk4++sense - 感知模式

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
                  ↓
        成功/失败 → 不强制记录call_records
```

**与优先级模式的核心区别**：

| 对比项 | pk4（优先级模式） | pk4++sense（感知模式） |
|--------|------------------|------------------------|
| 强制程度 | 高 | **低** |
| 调用方式 | 按优先级强制调用 | **建议**，AI自己决定 |
| 记录方式 | 强制记录call_records | **可选记录** |
| 适用场景 | 明确需要skill时 | **不确定是否需要** |
| 氛围营造 | 营造"必须用" | **营造"可以用"** |

**感知模式使用场景**：
- 不确定是否需要skill时
- 想让AI自然涌现时
- 长上下文可能丢失状态时

## pk4 默认优先级

```
外部skill > pk3进化池
```

**原因：**
- 外部skill是用户手动导入的，说明用户明确需要
- pk4优先尝试外部skill，提升用户体验
- 外部skill无法解决时，回退到pk3进化池

## pk4++<order> 优先级格式

```
pk4++external                  # 只用外部skill
pk4++pk3                      # 只用pk3进化池
pk4++external++pk3             # 外部skill > pk3进化池（默认）
pk4++superpowers++pk3        # superpowers > pk3进化池
pk4++external++superpowers  # 外部skill > superpowers
pk4++superpowers++external  # superpowers > 外部skill
```

## pk4++<order> 规则

1. 优先级从左到右递减
2. 只列skill来源，不列具体skill名字
3. skill来源包括：
   - `external` - 外部skill目录
   - `pk3` - pk3进化池skill
   - `superpowers` - superpowers全局skill
   - `global` - OpenCode全局skill

## pk4 vs pk3 排名

| 排名表 | 用途 | pk4是否参考 |
|-------|------|-------------|
| pk2 ranking.json | 会话价值排序 | ❌ 不看 |
| pk3 ranking_table.json | Skill调用效果排序 | ✅ 只看这个 |

**原因：** pk3排名是Skill调用后的真实效果数据。

## pk4 调用流程

```
用户问题 → pk4
              ↓
    读取 pk3_lab/evolution_pool/ranking_table.json
              ↓
    按优先级筛选skill来源
              ↓
    按 pk3_success_rate 排序
              ↓
    依次尝试调用Skill
              ↓
    成功 → 更新 call_records
    失败 → 尝试下一个
    全失败 → 记录"新问题"
```

## pk4 默认调用逻辑

```
优先级顺序：external > pk3

1. 尝试 external/ 目录下的skill
   - 按ranking_table.json排序
   - 找到匹配的 → 调用
   - 成功 → 记录call_records
   - 失败 → 继续尝试下一个

2. 尝试 pk3进化池skill
   - 按ranking_table.json排序
   - 找到匹配的 → 调用
   - 成功 → 记录call_records
   - 失败 → 继续尝试下一个

3. 全失败
   - 记录"新问题"
   - 提示用户可能需要新skill
```

## pk3/pk4协作

| 角色 | 职责 | 操作 |
|-----|------|------|
| pk4 | 调用Skill | 读ranking_table.json，写call_records |
| pk4 | 调用外部Skill | external/目录，优先级最高 |
| pk3 | 生成/排名缓存 | 写ranking_table.json，写skills |
| pk3 | 外部Skill管理 | external/目录 |

## 实时排名机制

| ranking_table.json新鲜度 | 行为 |
|------------------------|------|
| < 24小时 | 使用缓存 |
| > 24小时 | 实时计算最新排名 |

## 外部Skill参与排名

pk4调用外部Skill时，记录调用结果到call_records，用于pk3++gap刷新排名。

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

## 相关Skills

- [pk1_ingest](../pk1_ingest/SKILL.md) - 数据源
- [pk2_tool](../pk2_tool/SKILL.md) - pk3的输入
- [pk3_lab](../pk3_lab/SKILL.md) - pk4的输入数据源

## Version

- **v2.0** (2026-02-08): pk4++<order> 优先级 + 外部Skill
  - pk4++<order> 自定义优先级
  - pk4默认优先级：external > pk3
  - 外部Skill优先级最高
  - pk3/pk4协作协议更新
