$id: string
$schema: "task"
title: string
status: "todo" | "done"
dependsOn?: [..._]
metadata?: {
  channel?: string
}
