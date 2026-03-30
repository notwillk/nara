---
title: 'Resolver: implement ref forms and resolution order'
priority: p1
status: todo
ready: false
blocked_by:
  - 01KN05VZHT0000000000000000
  - 01KN060VX80000000000000000
tags:
  - resolver
  - config
---

Support sibling, relative, root (`/`), alias (`~alias/...`) per PRD §6. Resolution order: `~alias` → `/` → `./`/`../` → bare. Fail with **resolution error** (path + field).
