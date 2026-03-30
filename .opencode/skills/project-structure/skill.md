---
name: project-structure
description: Understand project structure, directories, and checks configuration
compatibility: opencode
---

## What I do

Explain the project's structure, directories, and checks configuration.

**Directories:**
- `.schemas/` — JSON Schema files for validating config/data. Generated schemas go here so editors and validators have a single known location. Cleared by `just clean --deep`.
- `.testignore` — file patterns excluded from triggering test re-runs in watch mode (`just test --watch`). Uses the same format as `.gitignore`. Add paths for files whose changes should not re-run tests (documentation, editor config, etc.).

**Checks:**

Checks are defined in YAML and run by `checksy`.

*Doctor checks* (`doctor.checksy.yaml`) — validate the dev environment (tools installed, env vars set, etc.). Run on container creation and with `just doctor`.

To add a check:
```yaml
rules:
  - name: description of check
    check: 'shell command that exits 0 on success'
    severity: error  # error | warning | debug
    fix: 'optional command to auto-fix'
```

*Static checks* (`static.checksy.yaml`) — validate the codebase (lint, format, etc.). Run in CI and with `just static`.

To add a check:
```yaml
rules:
  - name: description of check
    check: 'shell command that exits 0 on success'
    severity: error
    fix: 'optional auto-fix command'
```

**Devcontainer:**

Open in VS Code and choose **Reopen in Container**. The container runs `just doctor --fix` on creation, `just doctor` on each start, and includes `just`, `checksy`, and AI coding tools.

## When to use me

Use this when you need to understand project directories, add or modify checks, or work with the devcontainer configuration.