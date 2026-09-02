# One skill, to spec — the service skill folds into `planton`, the dot taxonomy lands, and agentskills.io conformance becomes a gate

## What changed

- **The assistant has ONE skill again.** `planton-service` folded into
  `planton`: service delivery is now a first-class domain beside
  infrastructure in a single SKILL.md with a single doctrine — one "Know
  your instruments" ladder, one invariants block, one prohibitions list
  (never-approve-gates promoted to a shared hard boundary). The separate
  skill duplicated the assistant's doctrine while NOTHING consumed it: the
  local seed and the hosted publish share one slug list that never included
  it, and no agent referenced it — so the merge is a pure restructure with
  zero consumer migration, taken while that was still true.
- **Every reference is renamed into a dot-separated domain taxonomy** —
  `infra.*` (chart craft and lifecycle), `cloud.*` (provider judgment),
  `catalog.*` (the research layer's doors), `craft.*` (working method and
  the person), `service.*` (service delivery) — flat files, per the Agent
  Skills specification's own keep-references-one-level-deep guidance. Git
  tracks all 37 as renames; the 5 cross-citations between references moved
  with them.
- **The Agent Skills specification (agentskills.io) is now the tree's
  binding standard, enforced by the gate** the PR lint and the release
  packager share (`pkg/skills/defspack`): the spec's name grammar; the
  1,024-character description cap — all three descriptions were over it
  (1,856 / 2,768 / 1,762 characters on the surface every agent loads at
  startup), all now within spec; frontmatter restricted to the spec's field
  set; a 500-line SKILL.md ceiling (the spec's guidance adopted as law —
  `planton` was 684 lines; the worked example and the workspace-posture
  choreography moved into references, which load on demand); and a
  dot-taxonomy filename grammar.
- **A directory under `references/` now refuses loudly.** Previously the
  loader silently skipped directories, so a nested reference file would
  ship in NO release while PASSING the bidirectional-citation lint — the
  file list that check runs over had already dropped it. The refusal names
  the directory and the flat contract; the silent-drop class is dead.
- **New reference: `service.offline-deploy.md`** — the complete journey
  from "set up CI/CD on GitHub for my repo" to a watched live deploy:
  authoring offline-clean kustomize trees (zero `$var`/`$secret`; runtime
  secrets as provider-native references composed from the secret snippets
  surface, never written from memory), the offline verb's
  report-then-deploy contract with its exit codes and node-addressed
  `--set` overrides, state-backend truths for laptops vs ephemeral
  runners, the forced-offline local verification road taught honestly
  (read the REPORT, not the exit code — 2 covers both refused and
  unapproved — and local preflight does not prove the RUNNER's
  credentials), the `gh`-driven setup choreography with a
  verify-everything-before-"ready" checklist and the offered finale
  (`gh run watch`), the one-line upgrade story to connected, and the
  refusal classes with their fixes.
- **The assistant's instructions frame both domains** of the one skill;
  **`skills/README.md` is rewritten as the complete authoring doctrine**
  (the one-skill rule with its admission test — the catalog skill's
  machine-assembled-plus-consumed-by-reference exception — the spec limits
  the gate enforces, the taxonomy and how to choose segment depth, the
  reference authoring bar, the compat-floor raising discipline); the lint
  workflow's header tells the same story.

## Why

One assistant with two doctrine documents is a drift machine: the two
"Know your instruments" ladders had already diverged in shape, and every
future session would have widened the gap. The specification limits exist
because the always-loaded surfaces are paid for on every conversation —
and they were being exceeded threefold. Both problems are now structural
impossibilities: the gate refuses spec drift on every pull request with
the same validator the release lane packages with, and there is exactly
one doctrine to drift from — none.

## Verification

- `go test ./pkg/skills/...` green: 16 violation classes each proven
  caught against a planted-defect fixture (including the six new spec
  classes and the nested-directory refusal), the unmutated fixture proven
  clean, and the committed-tree gate passing over the restructured tree.
- `go run ./pkg/skills/defspack` (the lint workflow's own second step):
  "validated 2 skill(s), 1 agent(s), and 1 automation(s)".
- The gate's before state, captured: the four violations it reported
  against the pre-cleanup tree (three description caps, one body ceiling)
  are exactly the measurements that drove the cleanup.
- `planton` SKILL.md: 500 lines, description 1,010 characters;
  `multi-cloud-catalog` description 1,003 characters — all within spec.
