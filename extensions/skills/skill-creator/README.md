> [!NOTE]
> This README is based on [openclaw/skill-creator](https://github.com/openclaw/openclaw/blob/main/skills/skill-creator/SKILL.md) (Apache 2.0), get the ZH version from [here](./README.zh.md).

# skill-creator

> A modular skill that guides agents to create, edit, improve, or audit AgentSkills with structured workflows, progressive disclosure design principles, and bundled resource patterns.<br>
> Based on [openclaw/skill-creator](https://github.com/openclaw/openclaw/blob/main/skills/skill-creator/SKILL.md) by openclaw [`v2026.3.8`](https://github.com/openclaw/openclaw/releases/tag/v2026.3.8) ([`3caab92`](https://github.com/openclaw/openclaw/commit/3caab92)), licensed under Apache 2.0. Modifications by author.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Skill Structure](#skill-structure)
- [Creation Process](#creation-process)
- [License](#license)

## Features

### Create New Skills from Scratch

Guides the agent through a 6-step process: understanding use cases with concrete examples, planning reusable resources, initializing the skill directory, writing SKILL.md, packaging, and iterating based on real usage.

### Audit and Improve Existing Skills

Reviews existing SKILL.md files against the AgentSkills spec — checking frontmatter completeness, directory structure, content conciseness, and proper progressive disclosure patterns.

### Progressive Disclosure Architecture

Enforces a three-level context loading system to minimize context window usage:

| Level | Content | When Loaded |
|-------|---------|-------------|
| Metadata | `name` + `description` | Always |
| SKILL.md body | Instructions and guidance | After skill triggers |
| Bundled resources | `scripts/`, `references/`, `assets/` | On demand |

### Bundled Resource Patterns

Provides clear guidance on when and how to use each resource type — deterministic scripts for repeated code, reference files for domain knowledge, and assets for output templates.

## Installation

Place this skill in the Claude Code skills directory:

```bash
~/.claude/skills/skill-creator/
```

Directory structure:

```
skill-creator/
├── SKILL.md              # Skill definition file
├── README.md
└── README.zh.md
```

## Usage

Trigger phrases (any of the following):

```
create a skill
author a skill
tidy up a skill
improve this skill
review the skill
clean up the skill
audit the skill
```

## Skill Structure

```
skill-name/
├── SKILL.md (required)
│   ├── YAML frontmatter — name + description
│   └── Markdown instructions
└── Bundled Resources (optional)
    ├── scripts/      — Executable code (Python/Bash/etc.)
    ├── references/   — Documentation loaded into context as needed
    └── assets/       — Files used in output (templates, icons, fonts)
```

## Creation Process

| Step | Action |
|------|--------|
| 1 | Understand the skill with concrete examples |
| 2 | Plan reusable contents (scripts, references, assets) |
| 3 | Initialize the skill directory (`init_skill.py`) |
| 4 | Write SKILL.md and implement resources |
| 5 | Package the skill (`package_skill.py`) |
| 6 | Iterate based on real usage |

> **Note**: `init_skill.py` and `package_skill.py` are tools from the openclaw toolchain and must be obtained from the original repository.

## License

Based on [openclaw/openclaw](https://github.com/openclaw/openclaw), original license: [Apache License 2.0](https://github.com/openclaw/openclaw/blob/main/LICENSE).
