---
description: Activate skill emergence with priority-based invocation.
---

# pk4 - Skill涌现

> **参见：** [pkskill/SKILL.md](../SKILL.md) （通用规则、路径占位符、命令格式、外部Skill机制）

## 命令格式

| 格式 | 说明 |
|------|------|
| `pk4` | 默认优先级：**外部skill > pk3进化池** |
| `pk4++<order>` | 自定义优先级顺序 |
| `pk4++sense` | 感知模式：检测上下文，建议可能需要的skill |

## pk4++sense 感知模式

**非强制**的感知-建议模式：

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

**与pk4优先级模式的区别**：
| pk4 | pk4++sense |
|-----|------------|
| 按优先级强制调用 | 建议，AI自己决定 |
| 强制记录调用结果 | 可选记录 |

**注意：** pk4使用**++分隔符**，不是空格或--flag格式。

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

## 路径说明

pk4从pk3_evolution_pool读取，按排名调用skill：
- 输入：`{OPENCODE_SKILLS}/pkskill/pk3_lab/evolution_pool`
- 排名表：`{OPENCODE_SKILLS}/pkskill/pk3_lab/evolution_pool/pk3_ranking_table.json`
- 外部skill：`{OPENCODE_SKILLS}/pkskill/pk3_lab/evolution_pool/external/`
