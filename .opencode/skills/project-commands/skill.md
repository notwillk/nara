---
name: project-commands
description: Use common project commands via just
compatibility: opencode
---

## What I do

Run common project tasks using [`just`](https://just.systems). Run `just help` for a summary.

**Build:**
- `just build` — compile the project
- `just build --watch` — compile with watch mode

**Test:**
- `just test` — run tests
- `just test --watch` — run tests with watch mode

**Clean:**
- `just clean` — remove build artifacts
- `just clean --deep` — also removes generated files in `.schemas/`

**Format:**
- `just format` — format code in-place
- `just format --check` — check formatting without modifying

**Doctor:**
- `just doctor` — check environment health
- `just doctor --fix` — auto-fix environment issues

**Static:**
- `just static` — run static checks including format
- `just static --fix` — format and auto-fix

## When to use me

Use this when you need to run build, test, format, or other common project tasks.