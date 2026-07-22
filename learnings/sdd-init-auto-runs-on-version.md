---
title: sdd-init Runs Silently on Invocation via --version
date: 2026-07-22
pr: TBD
area: [tooling, kb-bootstrap]
files: []
symptom: Running `sdd-tools --version` in a fresh repo executed `sdd-init` and copied all skill/command files before the bootstrap was ready to run it explicitly
tags: [sdd-tools, sdd-init, kb-bootstrap, tooling, idempotent]
summary: sdd-tools --version auto-runs sdd-init; this is idempotent so it's fine, but be aware that sdd-init may have already run before you invoke it explicitly in step 1.
---

## The pattern

When checking if `sdd-tools` was installed by calling `sdd-tools --version`, the tool executed `sdd-init` automatically (treating `--version` as the init trigger). This populated `.claude/`, `.specify/`, `.mcp.json` etc. before step 1's explicit `sdd-init` call.

## The rule

`sdd-init` is idempotent — it skips already-existing files. The bootstrap proceeds normally. Do not skip step 1's `sdd-init` call, but be aware the files may already exist when you get there.

## Detection

Step 1's `sdd-init` output shows all files as `skip (exists): ...` rather than `copied: ...`.

## Related

- None
EOF
