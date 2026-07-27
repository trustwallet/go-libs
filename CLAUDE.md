# go-libs

## What this repo is

- **Domain**: Public shared Go utility library for Trust Wallet backend services — infrastructure helpers (HTTP, caching, DB, messaging, metrics, workers, crypto).
- **Route here**: anything about the go-libs packages themselves — `cache/redis`, `client` (HTTP + JSON-RPC), `database` (GORM/Postgres, migrations), `metrics` (Prometheus), `middleware` (Gin), `mq` (RabbitMQ/AMQP), `worker` (periodic background jobs), `testy` (integration test helpers), `crypto` (AES/ECDSA), `logging`, `gin`, `health`, `slice`, `set`, `ctask`, `nullable`, `eventer`, `httplib`. Adding, modifying, or debugging any of these packages belongs here.
- **Do not route here**: domain-specific business logic or crypto primitives (those belong in `go-primitives`); the backend microservices that consume go-libs (separate private repos); wallet-core or mobile code; blockchain-specific coin types or addresses.
- **Consumers**: Trust Wallet Go backend microservices that import `github.com/trustwallet/go-libs` (private repos, multiple teams); any public Go project that adopts these utilities.
- **Ships**: Go module `github.com/trustwallet/go-libs` — importable packages (no binaries, no CLI, no HTTP service).
- **Agent map**: understanding a package → `knowledge/features/<package>.md`; architecture + god-nodes → `knowledge/architecture/`; patterns → `knowledge/patterns/`; security → `knowledge/security/`; test helpers → `knowledge/tests/`.

## Knowledge Map

For the structured knowledge base, see [knowledge/constitution.md](knowledge/constitution.md).

- [code-conventions](knowledge/code-conventions/index.md) — Code conventions, style rules, and decision records
- [patterns](knowledge/patterns/index.md) — Coding patterns, recipes, and proven approaches

- [architecture](knowledge/architecture/index.md) — Architecture
- [features](knowledge/features/index.md) — Features
- [libs](knowledge/libs/index.md) — Libs
- [security](knowledge/security/index.md) — Security
- [tests](knowledge/tests/index.md) — Tests

- [guides](knowledge/guides/index.md) — Guides

## Learnings

This repo may keep a living archive of incident-derived rules in ~~[`learnings/`](learnings/)~~ — each file a postmortem of a real bug or a non-obvious pattern that bit once and would bite again: root cause, the rule that prevents recurrence, and tags for matching. The folder is **optional and may be absent** — create it the first time you have a learning worth saving.

**Before** investigating any bug, regression, or "weird behavior", *if a `learnings/` directory exists*:

1. Search the frontmatter directly — it's the source of truth and always present:
   - `grep -ril "<keyword>" learnings/` — matches the frontmatter `tags:`/`summary:` + body.
   - `ls learnings/ | grep -i "<keyword>"` — matches the slug-style filename.
   - Skim each match's `summary:` line to decide whether to read the full body.
2. For a topic-organized ToC (grouped by surface + a tag index), open `learnings/index.md`. It is a **generated** artifact that garden **always regenerates** from frontmatter — never hand-edit it (any edit is discarded next run). Depending on the repo it's either gitignored (a derived artifact) or committed; either way it can be stale if a learning file changed without a regen, so prefer reading the learning files' frontmatter over trusting it blindly.
3. Found a match? **Read it before forming a hypothesis** — a 30-second read can turn a 2-hour investigation into a 5-minute fix.
4. Every file has frontmatter (`title`, `date`, `area`, `files`, `symptom`, `tags`, `summary`; `pr` when tied to a specific PR). `area` drives the index's surface grouping; `tags` drive its tag index.

**After** any fix, feature, or non-trivial change — if you learned something not already obvious from the code:

1. Add a new file `learnings/<slug>.md` with the frontmatter above, then a body covering: the symptom, the root cause (the actual mechanism, not just "the bug"), why prior fixes weren't enough if applicable, the rule going forward, and any regression guards. Create the `learnings/` folder if it doesn't exist yet.
2. If the learning extends an existing entry, edit that file instead of creating a duplicate.
3. Make the new file's `area`, `tags`, and `summary` accurate — those drive both `grep` and the generated index (`area` → its surface grouping, `tags` → its tag index, `summary` → its hook). **Never hand-edit `learnings/index.md`** — it's generated and always regenerated; edit the learning file's frontmatter instead.
4. Commit the learning **in the same PR as the fix** — never as a follow-up.

The bar: would a future agent save time by reading this before touching the same surface? If yes, write it; if it would just say "read the diff," skip it. Don't ask which learnings to capture — commit every candidate that clears the bar.

## Repository Knowledge Scope

This repo's `knowledge/` covers: **code-conventions, patterns**

Topics NOT documented locally: architecture, build, ci, conventions, core-libs, decisions, design, features, git-conventions, guides, observability, product, quality, security, tests, workflows, brand, business, legal, hr, prompts, api, specs, libs, components, references

## Constraints

- This is a **public repo** — never commit secrets, tokens, private DSNs, or any internal infrastructure details.
- Consumers import packages individually — avoid creating circular imports between packages (e.g. `client` must not import `middleware`).
- All metric constructors accept a `prometheus.Registerer` — never use `prometheus.DefaultRegisterer` or `prometheus.MustRegister` with global metrics.
- The `testy` package is test-only — never import it from non-test code.

<!-- sdd-knowledge-generated -->
