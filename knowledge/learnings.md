# Project Learnings

Structured knowledge base for patterns and anti-patterns discovered across features.
Inspired by CASS Memory System concepts: confidence tracking, staleness decay, and cross-feature learning.

**How to use this file:**
- `/speckit.specify` and `/speckit.plan` check this file before generation for relevant anti-patterns and patterns
- `/speckit.implement` appends new learnings after each implementation session
- `/speckit.analyze` validates staleness of entries (Detection Pass G)
- Entries not validated within 90 days are flagged as potentially stale
- Patterns validated across 3+ features are candidates for constitution promotion

---

## Patterns (validated approaches)

<!-- Add validated patterns here. Each pattern should include:
- **confidence**: high | medium | low
- **last_validated**: YYYY-MM-DD
- **validated_count**: N (number of features that confirmed this pattern)
- **context**: Description of the pattern and when to apply it
- **features**: [list of feature IDs where this pattern was used]
-->

_No patterns recorded yet. Patterns will be captured after implementation sessions._

---

## Anti-Patterns (failed approaches)

<!-- Add failed approaches here. Each anti-pattern should include:
- **type**: anti-pattern
- **discovered**: YYYY-MM-DD
- **context**: What was attempted and why it failed
- **features**: [list of feature IDs where this was discovered]
- **severity**: critical | high | medium | low
-->

_No anti-patterns recorded yet. Anti-patterns will be captured after implementation sessions._

## See Also
- [anti patterns failed approaches](code-conventions/code-style/anti-patterns-failed-approaches.md) <!-- rel:strong -->
- [patterns validated approaches](patterns/patterns-validated-approaches.md) <!-- rel:strong -->
- [common mistakes and anti patterns](guides/troubleshooting/common-mistakes-and-anti-patterns.md) <!-- rel:related -->
### Code-conventions


- **[high]** This is a **public GitHub repository**. Never commit: _(source: CLAUDE.md)_
- **[high]** Each package must be importable independently. Do NOT create imports between sibling packages (e.g. `client` importing `middleware`, `database` importing `metrics` directly). Dependencies flow ONE WAY: utility packages → infrastructure packages → integration packages. _(source: CLAUDE.md)_
- **[high]** All metric constructors (`NewPerformanceMetric`, `NewHttpServerMetric`, etc.) accept a `prometheus.Registerer` parameter. **Never** use: _(source: CLAUDE.md)_
- **[high]** The `testy` package provides integration test helpers (real DB + Redis connections). Never import it from non-test (production) code. It should only appear in `_test.go` files. _(source: CLAUDE.md)_
