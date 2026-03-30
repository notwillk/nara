---
title: 'Entity I/O: parse yaml/yml/json, extract id, schema, ref, meta keys'
priority: p1
status: todo
ready: false
blocked_by:
  - 01KN05VZHT0000000000000000
tags:
  - graph
---

Use `gopkg.in/yaml.v3` + JSON decode. Respect `meta` renames from config. Preserve locations for errors (line/column where feasible).
