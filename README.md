# pkskill - AI能力进化系统

> 一个自研自进的AI能力进化系统，通过会话数据管理、技能生成优化和能力涌现来提升AI性能。

## 概述

pkskill包含4个子skill，实现完整的数据流：

1. **pk1_ingest** - 会话数据转存

2. **pk2_tool** - 数据增强与总结

3. **pk3_lab** - Skill实验室

4. **pk4_emerge** - Skill涌现

### 数据流

```
用户输入 → pk1转存 → pk2总结 → pk3生成skill → pk4涌现调用
              ↓              ↓               ↓
         pk数据库       增强数据库        进化池
              ↓
    pk1.go自动检测路径
    生成pk_path_config.json
```

## 跨平台路径支持

pk1.exe 内置跨平台路径检测功能，运行时会自动：

* 检测操作系统（Windows/Linux/macOS）

* 检测用户名（环境变量优先级：USERNAME → USER → LOGNAME）

* 检测主目录（HOME → USERPROFILE → HOMEDRIVE+HOMEPATH）

* 生成 `pk_path_config.json` 供其他组件使用

### 路径占位符

| 占位符          | 说明         | Windows示例         | Linux示例        |
| :----------- | :--------- | :---------------- | :------------- |
| `{PK_PATH}`  | pkskill根目录 | D:/pkskill        | /opt/pkskill   |
| `{HOME}`     | 用户主目录      | C:/Users/username | /home/username |
| `{USERNAME}` | 用户名        | username          | username       |

### OpenCode路径占位符

| 占位符                 | 说明                 | Windows示例                          | Linux示例                         |
| :------------------ | :----------------- | :--------------------------------- | :------------------------------ |
| `{OPENCODE_CONFIG}` | OpenCode配置根目录      | C:/Users/username/.config/opencode | /home/username/.config/opencode |
| `{OPENCODE_SKILLS}` | OpenCode全局skills目录 | {OPENCODE_CONFIG}/skills           | {OPENCODE_CONFIG}/skills        |

### 平台示例

#### Windows

```bash
PK_PATH=D:/pkskill
HOME=C:/Users/username
pk1_output=D:/pkskill/pk1_output
pk_path_config=D:/pkskill/pk1_output/pk_path_config.json
```

#### Linux

```bash
PK_PATH=/opt/pkskill
HOME=/home/username
pk1_output=/opt/pkskill/pk1_output
pk_path_config=/opt/pkskill/pk1_output/pk_path_config.json
```

#### macOS

```bash
PK_PATH=/Users/username/pkskill
HOME=/Users/username
pk1_output=/Users/username/pkskill/pk1_output
pk_path_config=/Users/username/pkskill/pk1_output/pk_path_config.json
```

## 安装

windows：将pkskill内部的pkskill目录剪切到OpenCode全局配置目录，然后重启OpenCode

```bash
C:\Users\{USERNAME}\.config\opencode\skills\
```

### 目录结构（全局配置）

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

## 快速开始（在opencode对话框输入命令）

### 1. 转存会话数据

```bash
pk1              # 默认（自动检测）opencode本地会话源路径转存到目标路径pkskill/pk1_output
pk1++{OPENCODE_STORAGE}         # 指定源路径
pk1++{OPENCODE_STORAGE}++{PK_PATH}  # 指定源路径和目标路径
```

### 2. 增强数据库

```bash
pk2                      # 默认pk1_output源和pk2目标路径pk2_enriched
pk2++{PK_PATH}              # 指定源数据库路径
```

### 3. 生成Skill

```bash
pk3                     # 默认pk2源和pk3目标路径
pk3++{PK_PATH}              # 指定从哪个数据库路径生成或者识别导入外部skill
pk3++go                   #触发内部进化池skill融合进化
pk3++gap                   #刷新skill表现排名表
pk3++see                   #生成表现报告
```

### 4. 激活涌现

```bash
pk4                           # 默认
pk4++sense                      #自动感知（推荐）
pk4++external                  # 只用外部skill
pk4++pk3                      # 只用pk3进化池
pk4++external++pk3            # 外部skill > pk3进化池
pk4++superpowers++pk3        # superpowers > pk3进化池
pk4++external++superpowers  # 外部skill > superpowers
```

## 子Skill详解

### pk1_ingest - 会话数据转存

**pk1.go 内置跨平台路径检测**：

* OpenCode存储: `{HOME}/.local/share/opencode/storage`

* 输出目录: `{PK_PATH}/pk1_output`

### pk2_tool_v3 - 数据增强

**跨平台支持**: pk1.exe生成pk_path_config.json，pk2读取使用

**pk2路径**:

* 输入: `{PK_PATH}/pk1_output` (只读.md，跳过.json)

* 输出: `{PK_PATH}/pk2_enriched/enriched`

* 排名表: `{PK_PATH}/pk2_enriched/ranking.json`

**核心特性**:

* ✅ pk1内置跨平台路径检测

* ✅ v2版摘要结构（6个正面模式字段）

* ✅ 无意义会话判断（7种低价值类型）

* ✅ 排名表ranking.json + 特权池

* ✅ pk3标签联动（高排名skill特权）

* ✅ 只读MD + 分层读取（Token优化82%）

* ✅ 全局问题扫描 + 根因分析

**版本演进**:

| 版本          | Token/会话 | 核心特性           |
| ----------- | -------- | -------------- |
| 原版pk2       | ~34,000  | 基础增强           |
| pk2_tool_v2 | ~6,000   | 只读MD           |
| pk2_tool_v3 | ~6,000   | v2摘要+排名表+pk3标签 |

### pk3_lab - Skill实验室

**路径来源**: pk1.exe生成pk_path_config.json

* 数据库: `{PK_PATH}/pk1_output`

* 进化池: `{OPENCODE_SKILLS}/pkskill/pk3_lab/evolution_pool`

### pk3++<path> 智能识别

```
pk3++{PK_PATH}
       ↓
检测路径下的子文件夹里是否存在SKILL.md
       ↓
├─ 存在SKILL.md → 导入外部skill → 放到 external/{skill_name}/
└─ 不存在SKILL.md → 指定数据库路径 → 扫描数据生成skill
```

### pk3++go - 融合进化

触发全局配置中pk3_lab/evolution_pool/fusion/目录下的融合进化流程。

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

## pk4 - Skill涌现

### 命令格式

| 格式             | 说明                         |
| :------------- | :------------------------- |
| `pk4`          | 默认优先级：**外部skill > pk3进化池** |
| `pk4++<order>` | 自定义优先级顺序                   |
| `pk4++sense`   | 感知模式：检测上下文，建议可能需要的skill    |

### pk4++sense - 感知模式(推荐）

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

| 对比项  | pk4（优先级模式）       | pk4++sense（感知模式） |
| :--- | :--------------- | :--------------- |
| 强制程度 | 高                | 低                |
| 调用方式 | 按优先级强制调用         | 建议，AI自己决定        |
| 记录方式 | 强制记录call_records | 可选记录             |
| 氛围营造 | 必须用              | 可以用              |

### pk4++<order> 优先级格式

```
pk4++external                  # 只用外部skill
pk4++pk3                      # 只用pk3进化池
pk4++external++pk3            # 外部skill > pk3进化池
pk4++superpowers++pk3        # superpowers > pk3进化池
pk4++external++superpowers  # 外部skill > superpowers
```

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

| 特性         | 规则           |
| :--------- | :----------- |
| pk4调用优先级   | ✅ **最高**（默认） |
| 参与融合进化     | ❌ 不参与        |
| 参与排名       | ✅ 参与         |
| pk3++go融合  | ❌ 不会被融合      |
| pk3++gap刷新 | ✅ 参与排名计算     |

### 外部Skill排名指标

```
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

### pk4_emerge - Skill涌现

**路径来源**: pk1.exe生成pk_path_config.json

* 进化池: `{OPENCODE_SKILLS}/pkskill/pk3_lab/evolution_pool`

## 常见问题

### Q1: pk1.go如何检测平台？

A1: 例如windows：pk1.exe内置跨平台检测：

* Windows (win32)

* Linux

* macOS (darwin)

### Q2: 如何自定义PK_PATH？

A2: 设置环境变量或配置文件：

```bash
export PK_PATH=/custom/path/pkskill
```

### Q3: 支持哪些路径占位符？

A3:

* `{PK_PATH}` - pkskill根目录

* `{HOME}` - 用户主目录

* `{USERNAME}` - 用户名

