---
title: 'Config: load, merge defaults, validate nara.yaml structure'
priority: p0
status: todo
ready: false
blocked_by:
  - 01KN05VKBF0000000000000000
tags:
  - config
  - errors
---

Implement `version`, `paths`, `meta`, `schemas`, `resolution` per PRD §5. Validate required fields, path syntax, glob-friendly schema sources. Errors: **config error** category, include file path.
