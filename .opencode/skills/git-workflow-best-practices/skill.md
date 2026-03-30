---
name: git-workflow-best-practices description: Apply modern Git branching best practices (trunk-based / GitHub Flow) for repositories, pull requests, and releases. compatibility: opencode
---

# Git Workflow Best Practices

Use this skill when designing or modifying a repository's Git workflow. The goal is a simple, modern branching model optimized for CI/CD and frequent integration.

Follow trunk-based development / GitHub Flow principles.

## Core Rules

### 1. Main Branch

The repository has one primary branch:

```text
main
```

Requirements:

- `main` must always be releasable
- CI must pass before merge
- direct pushes should usually be disabled
- merges happen through pull requests

### 2. Short-Lived Branches

All work happens on temporary branches.

Naming:

```text
feature/<name>
fix/<name>
chore/<name>
refactor/<name>
docs/<name>
```

Example:

```text
feature/add-auth
fix/timeout-bug
refactor/config-loader
```

Rules:

- branches should live hours or days, not weeks
- branches should stay small
- merge frequently

### 3. Pull Requests

All changes go through pull requests.

PR best practices:

- small and focused
- one logical change per PR
- CI must pass
- code review before merge

Preferred merge strategy:

```text
squash merge
```

This keeps history readable.

### 4. Continuous Integration

CI must run on:

- pull requests
- merges to `main`

Typical checks:

- build
- lint
- tests
- formatting

PRs must not merge if CI fails.

### 5. Releases

Releases come directly from `main`.

Create tags:

```text
v1.2.0
v1.2.1
```

Example:

```bash
git tag v1.4.0
git push origin v1.4.0
```

Use release branches only if necessary for long-lived maintenance.

Example:

```text
release/1.4
```

Most repositories do not need release branches.

### 6. Hotfixes

Urgent fixes use the same flow:

```text
fix/critical-bug
→ PR
→ merge to main
→ tag release
```

Avoid separate `hotfix` branches unless required.

## What To Avoid

Avoid complex workflows such as:

```text
develop
release/*
hotfix/*
```

This is the classic Gitflow model, which is heavier and slower for CI/CD workflows.

Prefer:

```text
main
└─ feature/*
```

## When To Use This Skill

Use this skill when:

- creating a new repository
- defining contribution guidelines
- designing a CI/CD workflow
- simplifying an existing Gitflow workflow
- documenting a repo’s branching policy

## Example Repository Structure

```text
main
├─ feature/login
├─ fix/cache-timeout
└─ docs/update-readme
```

All merge back into:

```text
main
```

## Outcome

Repositories following this workflow should have:

- fast integration
- minimal branching complexity
- CI-enforced quality
- simple release management

