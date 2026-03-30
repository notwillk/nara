---
title: 'Schema discovery: glob CUE sources, one schema per file, unique schema names'
priority: p0
status: todo
ready: false
blocked_by:
  - 01KN05VZHT0000000000000000
tags:
  - schema
  - cue
---

Load `./schemas/*.cue` (and config-driven globs). Filename stem = schema name; fail on duplicates. Surface **schema error** with file path.
