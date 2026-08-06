# Catalog-Research Eval Results

Append-only ledger of catalog-research eval runs. Each run answers the
questions in `eval_questions.yaml` under the protocol in
`_rules/docs/evaluate-planton-catalog-research.mdc`: the same questions,
answered by fresh-context agents twice -- once restricted to `planton
explain` (the CLI arm), once restricted to the reference pack's files (the
pack arm) -- grading correctness and counting the operations each side
needed. Newest run at the top. Never edit past runs.

## 2026-08-06 — introduction run

- **Catalog state**: commit `0ca2ec6a6` (all 15 questions' freshness checks green before the run).
- **Runners**: one fresh-context agent per question per arm (30 runs); arm restrictions prompt-enforced, every operation log audited, no boundary violations found.
- **Headline**: the pack arm answered **15/15 correct in 30 operations**; the CLI arm answered **13/15 (2 partial, 0 wrong, 0 DNF) in 66 operations** — 2.2x the operations for less complete answers. On the reverse-reference class alone (q07+q08) the CLI needed 28 operations and still returned one incomplete answer (hit the 15-op cap having found 8 of 11 inbound consumers); the pack needed 5 operations and returned both complete, each cross-verified against two independent pack surfaces.

| Q | Class | CLI grade | CLI ops | Pack grade | Pack ops |
|---|---|---|---|---|---|
| q01 | kind-discovery | correct | 2 | correct | 1 |
| q02 | kind-discovery | correct | 2 | correct | 3 |
| q03 | required-fields | correct | 2 | correct | 1 |
| q04 | enum-values | correct | 1 | correct | 1 |
| q05 | outputs | correct | 4 | correct | 2 |
| q06 | fk-forward | correct | 1 | correct | 1 |
| q07 | reverse-references | correct | 13 | correct | 3 |
| q08 | graph | partial (8/11 consumers, hit cap) | 15 | correct (11/11) | 2 |
| q09 | alias-substitution | correct | 3 | correct | 3 |
| q10 | wisdom | correct | 5 | correct | 2 |
| q11 | wisdom | correct | 4 | correct | 4 |
| q12 | wisdom | partial (right conclusion, wrong basis: reported "the CLI shows no diagrams" — it cannot see diagram semantics) | 4 | correct | 2 |
| q13 | grammar | correct | 5 | correct | 1 |
| q14 | cross-provider | correct | 2 | correct | 1 |
| q15 | wisdom | correct | 3 | correct | 3 |
| **Total** | | **13 correct, 2 partial** | **66** | **15 correct** | **30** |

- **Median ops per question**: CLI 3, pack 2. Wall-clock (qualitative): parallel batches containing the CLI arm's reverse-reference questions ran visibly longer than every pack batch; per-question latency tracked operation count.
- **Where each side stands**: single-kind fact lookups are near-parity (the CLI drill-down is genuinely good, and judgment pushed into proto docs at the source surfaces in both arms). The pack's decisive wins are inbound references and whole-catalog wiring (Referenced By + the graph: the CLI can only probe candidates and cannot know when it is done), catalog-level judgment homes (patterns, the substitution workflow), and platform semantics no schema dump carries (diagram-edge behavior).
- **Findings (recorded, not fixed here)**: no pack format defects surfaced — every question routed to its intended surface via the skill's workflow. One CLI-side observation: its terminal renderer displays outputs in camelCase while valueFrom fieldPaths are canonically snake_case, which cost the CLI arm extra verification operations on q05 (the pack renders the canonical spelling directly).
