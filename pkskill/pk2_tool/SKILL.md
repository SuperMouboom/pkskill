---
name: pkskill-pk2-tool-v3
description: pk2智能总结工具v3 - v2版摘要结构、排名表、pk3标签联动、无意义会话判断。只读MD、分层读取、Token优化82%。
requires: pkskill/SKILL.md
---

# pk2_tool_v3 - AI对话数据智能增强工具

> **前置要求：** 修改pk2前必须阅读 [pkskill/SKILL.md](../SKILL.md)

## 核心定位

pk2 是**AI行为准则**，不是可执行工具。

**pk2 的本质**：
- 当用户要求总结/增强/分析 OpenCode 会话时，AI 加载此 skill
- AI 直接执行总结任务，不需要调用任何外部命令
- pk2 存在于 AI 的认知中，不是独立的二进制程序

## 执行步骤

**当用户要求总结、增强或分析 OpenCode 会话时**：
1. **语言选择** - 遵循主 SKILL.md"语言选择规则"
2. 扫描 `{PK_PATH}/pk1_output/*.md` 文件
3. 为每个会话生成 individual v2 摘要
4. 保存摘要到 `pk2_enriched/enriched/`
5. 生成 `pk2_meta.json` 记录扫描结果
6. 生成 `ranking.json` 排名表

**v2摘要语言规则**：
- skills_used, tools_used 保持英文
- 其他字段（success_patterns, best_practices 等）使用用户选择语言

## 输出结构

```
{PK_PATH}/
└── pk2_enriched/
    ├── ranking.json           # 排名表
    ├── pk2_meta.json          # 扫描记录（SHA256哈希用于增量更新）
    └── enriched/
        └── {核心主题}-{技术栈}-{评分}.md   # 每个会话的v2摘要
```

## 文件命名规则

**文件名格式**：`{核心主题}-{技术栈}-{评分}.md`

**命名要求**：
1. **核心主题**：从会话标题中提取1-3个关键词，用中划线连接
2. **技术栈**：主要使用的语言/框架/工具（如：python、rust、powershell、general）
3. **评分**：A+、A、B、C、D、F（根据价值评估体系）

**示例**：
- `brainstorming-pk2-optimization-A.md`
- `powershell-script-debugging-B.md`
- `tauri-lifetime-annotation-A.md`
- `tokenization-multilingual-B.md`

**AI感知规则**：
- AI查找pk2数据库时，**先扫描文件名**获取主题和评分信息
- 高评分（A+/A）文件优先处理
- 文件名即快速索引，无需读取完整内容即可判断相关性

**禁止行为**：
- 制作 pk2.exe、pk2.bat 等工具
- 调用外部命令执行 pk2
- 只生成汇总文件，不保存 individual v2 摘要

## v2摘要结构（输出格式）

每个增强文件遵循v2标准格式：

```markdown
<!-- AI_QUICK_LOOKUP_START -->
```json
{
  "session_id": "conv_ses_xxx",
  "model": "minimax-m2.1-free",
  "skills_used": ["brainstorming", "writing-plans"],
  "tools_used": ["read", "write", "glob"],
  "session_type": "NORMAL",
  "success_patterns": [
    {
      "pattern": "parallel_skill_loading",
      "description": "同时加载多个skill并行执行",
      "effectiveness": "显著提高效率"
    }
  ],
  "design_decisions": [
    {
      "decision": "选择A而非B",
      "reasoning": "A更适合当前场景",
      "alternative_considered": "B的缺点"
    }
  ],
  "best_practices": [
    "先理解上下文再行动",
    "利用现有skill避免重复"
  ],
  "workflow_patterns": [
    {
      "pattern": "scan_first_design_later",
      "description": "先扫描再设计",
      "applies_to": ["重构", "集成"]
    }
  ],
  "pitfalls_avoided": [
    {
      "pitfall": "过早优化",
      "description": "在性能需求明确前避免优化",
      "severity": "medium"
    }
  ],
  "errors_analyses": [
    {
      "error": "API方法名不匹配",
      "analysis": "过度信任用户输入，跳过了独立验证环节。应该先grep API定义再修改。",
      "root_cause": "认知偏差 - 假设用户提供的信息100%准确",
      "solution": "任何API修改先执行grep验证方法签名，不直接信任行号信息",
      "pattern": "input_trust_validation"
    }
  ],
  "code_strategies": {
    "language": "python",
    "patterns": ["context_manager", "generator"],
    "anti_patterns": ["global_state"]
  },
  "key_insights": "TDD显著减少bug数量",
  "learnable": true,
  "priority": "HIGH",
  "pk2_version": "v3"
}
```
<!-- AI_QUICK_LOOKUP_END -->
```

### errors_analyses 深度要求

**每个错误必须包含4个层面**：

| 字段 | 必须回答的问题 | 示例 |
|------|--------------|------|
| `error` | 具体犯了什么错？ | "API方法名不匹配" |
| `analysis` | 当时是怎么做的？当时在想什么？ | "用户给了行号，我直接相信了，没有自己验证" |
| `root_cause` | 本质问题是什么？认知偏差还是知识不足？ | "过度信任用户输入，缺少独立验证思维" |
| `solution` | 具体、可执行的规避措施？ | "先grep API定义，不直接信任行号" |
| `pattern` | 这个错误属于哪类认知模式？ | "input_trust_validation" |

**禁止**：
- 只有错误描述，没有分析
- 只有解决方案，没有原因
- 用"粗心"、"大意"解释（这是认知惰性）

**必须**：
- 挖掘到认知层面
- 抽象出可复用的思维模式
- 给出下次能直接执行的行动

## 价值评估体系

| 维度 | 权重 | 说明 |
|------|------|------|
| 成功率 | 40% | 任务完成度 |
| Token效率 | 20% | 输入输出比 |
| Skill多样性 | 15% | Skill调用数 |
| 工具多样性 | 15% | 工具使用数 |
| 代码变更量 | 10% | 合理范围 |

## 价值等级

| 等级 | 评分 | 定义 |
|------|------|------|
| A+ | 95-100 | 完整、有深度、可复用工作流 |
| A | 85-94 | 高质量，有洞察 |
| B | 70-84 | 基本完整 |
| C | 55-69 | 有缺失 |
| D | 40-54 | 不完整 |
| F | <40 | 无意义/反面教材 |

## 相关Skills

- [pk1_ingest](../pk1_ingest/SKILL.md) - 数据源（pk2的输入）
- [pk3_lab](../pk3_lab/SKILL.md) - Skill实验室
- [pk4_emerge](../pk4_emerge/SKILL.md) - Skill涌现

## Version

- **v4.9** (2026-02-08): 强化errors_analyses深度要求，每个错误必须包含error/analysis/root_cause/solution/pattern五个层面
- **v4.8** (2026-02-08): 添加文件命名规则，文件名格式 `{核心主题}-{技术栈}-{评分}.md`，AI通过文件名快速感知内容
- **v4.7** (2026-02-08): 语言规则引用主 SKILL.md
- **v4.6** (2026-02-08): 修复缺失输出问题，添加执行步骤和输出结构
- **v4.5** (2026-02-08): 禁止制作工具，AI必须直接使用pk2命令
