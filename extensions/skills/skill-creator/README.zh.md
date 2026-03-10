> [!NOTE]
> 此 README 基於 [openclaw/skill-creator](https://github.com/openclaw/openclaw/blob/main/skills/skill-creator/SKILL.md)（Apache License 2.0），英文版請參閱 [這裡](./README.md)。

# skill-creator

> 引導 Agent 建立、編輯、改善或審查 AgentSkill 的模組化 Skill，包含結構化工作流程、漸進式揭露設計原則與打包資源模式。<br>
> 基於 [openclaw/skill-creator](https://github.com/openclaw/openclaw/blob/main/skills/skill-creator/SKILL.md) by openclaw [`v2026.3.8`](https://github.com/openclaw/openclaw/releases/tag/v2026.3.8)（[`3caab92`](https://github.com/openclaw/openclaw/commit/3caab92)，Apache License 2.0），作者僅做部分調整。

## 目錄

- [功能特點](#功能特點)
- [安裝](#安裝)
- [使用方法](#使用方法)
- [Skill 結構](#skill-結構)
- [建立流程](#建立流程)
- [授權](#授權)

## 功能特點

### 從零建立新 Skill

引導 Agent 完成六步驟流程：透過具體範例理解使用情境、規劃可重用資源、初始化 Skill 目錄、撰寫 SKILL.md、打包，以及依實際使用回饋迭代改善。

### 審查與改善現有 Skill

依 AgentSkills 規範審查現有 SKILL.md——檢查 Frontmatter 完整性、目錄結構、內容精簡度，以及漸進式揭露模式是否正確。

### 漸進式揭露架構

強制採用三層 Context 載入系統，最小化 Context Window 使用量：

| 層級 | 內容 | 載入時機 |
|------|------|----------|
| Metadata | `name` + `description` | 始終載入 |
| SKILL.md body | 指令與指引 | Skill 觸發後 |
| 打包資源 | `scripts/`、`references/`、`assets/` | 依需求載入 |

### 打包資源模式

提供清楚的指引，說明何時及如何使用各類型資源——確定性腳本用於重複程式碼、參考文件用於領域知識、Assets 用於輸出模板。

## 安裝

將此 Skill 放置於 Claude Code 的技能目錄：

```bash
~/.claude/skills/skill-creator/
```

目錄結構：

```
skill-creator/
├── SKILL.md              # 技能定義檔
├── README.md
└── README.zh.md
```

## 使用方法

觸發短語（以下任一均可）：

```
create a skill
author a skill
tidy up a skill
improve this skill
review the skill
clean up the skill
audit the skill
```

## Skill 結構

```
skill-name/
├── SKILL.md（必要）
│   ├── YAML Frontmatter — name + description
│   └── Markdown 指令
└── 打包資源（選用）
    ├── scripts/      — 可執行程式碼（Python/Bash 等）
    ├── references/   — 依需求載入 Context 的參考文件
    └── assets/       — 輸出用資源（模板、圖示、字型）
```

## 建立流程

| 步驟 | 動作 |
|------|------|
| 1 | 透過具體範例理解 Skill 使用情境 |
| 2 | 規劃可重用內容（scripts、references、assets） |
| 3 | 初始化 Skill 目錄（`init_skill.py`） |
| 4 | 撰寫 SKILL.md 並實作資源 |
| 5 | 打包 Skill（`package_skill.py`） |
| 6 | 依實際使用回饋迭代改善 |

> **注意**：`init_skill.py` 與 `package_skill.py` 為 openclaw 工具鏈的腳本，需從原始 repo 取得。

## 授權

基於 [openclaw/openclaw](https://github.com/openclaw/openclaw)，原始授權為 [Apache License 2.0](https://github.com/openclaw/openclaw/blob/main/LICENSE)。
