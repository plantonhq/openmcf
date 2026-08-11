---
name: multi-cloud-catalog
description: Research Planton's cloud component catalog from its file-based reference pack -- generated per-component reference pages (spec fields, validation rules, outputs, cross-component wiring in both directions, a validated example manifest), per-provider and root indexes, the catalog-wide foreign-key graph, and the authored GUIDE.md / patterns wisdom layer. Use when composing or reviewing an architecture and you need component facts (what exists for a provider, which fields a component requires, what an output is called, what can reference what), operational judgment before choosing between components, or a compatible alternative when the software the user asked for has no component of its own. Also covers giving wisdom back -- when a session learns judgment the pack does not teach, this skill's contribution workflow turns it into a reviewed pull request against the catalog's guides and patterns. The discipline is absolute -- every schema claim is read from a pack file at answer time, never recalled from model memory; when no pack is reachable, fall back to `planton explain` per component and say what that fallback cannot see. Do not use for editing schemas or generated reference pages (facts are fixed at their proto source, never on the page) and do not use for deploying resources -- this skill is the research layer, plus the loop that improves it, that other workflows build on.
---

# Multi-Cloud Catalog

Planton's deployable building blocks -- cloud components spanning AWS, GCP,
Azure, Kubernetes, and every other supported provider -- are documented as a
reference pack: plain files you read and search. One page per component,
regenerated from the schemas so facts cannot drift, plus indexes, a
catalog-wide wiring graph, and an authored wisdom layer. This skill is how
you research that pack: which file answers which question, in what order to
look, and the discipline that keeps answers grounded.

## The one law: facts are read, never recalled

Every field name, type, enum value, default, validation rule, and output
path you state must come from a pack file you read while answering -- or,
when no pack is on disk, from `planton explain` or from a component page
fetched into your workspace (`references/pack-layout.md` names the fallback
ladder). Whatever you remember about a component's schema is stale by
construction: schemas change and the pack regenerates with them. A confident
answer from memory that names one wrong field costs the user a failed
deployment; a ten-second file read costs nothing.

## Locate the pack first

Find the pack root before the first research step --
`references/pack-layout.md` has the tree shape, the places the pack lives
(repo checkout, release zip, vendored copies), and the honest fallback when
none is reachable. Everything below writes `<pack-root>` for the directory
that contains `reference-commons.md`.

## Which file answers which question

| Question | Where the answer lives |
|---|---|
| What components exist for this provider / this need? | The provider's `reference-index.md` table; full-text search for capability words when the name is unknown |
| What does component K require? | K's `reference.md`: read `## Example` first (a validated manifest), then the `## Spec Fields` table's Required column and `## Validation Rules` |
| What does one field mean and accept? | The field's own `### spec.<path>` block in `## Field Details` -- docs, rules, allowed enum values |
| What does K export after deployment? | K's `## Outputs` table |
| How do components wire together? | K's `## References` (outbound) and `## Referenced By` (inbound) tables; `reference-graph.yaml` for catalog-wide edges |
| What is the judgment call before choosing K? | `GUIDE.md` beside K's page -- the page head links it when one exists -- and `patterns/` for multi-component recipes |
| The asked-for software has no component of its own | The catalog `GUIDE.md` beside the root index owns this workflow: search the pack by the name the user said, propose the compatible alternative openly, generic mechanisms only as the last resort |
| Manifest grammar: envelope, value/valueFrom, fieldPath spelling | `reference-commons.md` -- read it once per session; it also documents the search grammar these pages share |
| This session learned judgment the pack does not teach | `references/contributing-wisdom.md` -- offer once, verify, draft, and deliver it as a pull request |

The concrete command moves for each row -- with the friction traps already
learned so you do not relearn them -- are in
`references/research-recipes.md`.

## The research workflow

1. **Commons once.** Read `reference-commons.md` at the start of catalog
   work: the manifest grammar every component shares and the fixed heading
   vocabulary that makes the whole pack greppable.
2. **Index before grep.** For "what exists" questions, open the provider's
   `reference-index.md` before reaching for full-text search -- the table
   carries every kind with a one-line purpose, and its Guide column tells
   you where judgment has been written.
3. **Page before graph.** For one component's wiring, its own References /
   Referenced By tables are complete; the graph file is for catalog-wide
   questions ("every field anywhere that can point at X").
4. **Facts, then wisdom.** Settle the schema facts first, then read the
   component's guide and any pattern that composes it -- wisdom changes
   which component you pick and what else the architecture needs, not what
   the fields are.
5. **Verify before you assert.** Before stating a field path, an enum
   value, or an output name in an answer or a manifest, confirm the exact
   spelling in the page you have open. If a guide and a generated page ever
   disagree on a fact, the generated page wins -- and the disagreement is
   worth contributing back (`references/contributing-wisdom.md`), as is any
   judgment this session had to work out that no guide teaches.

## Speak the user's language

The catalog's internal vocabulary -- kind names, `valueFrom`, `fieldPath`,
provisioner annotations -- is your building material, never the user's
curriculum. Deliver answers in the terms the user asked in; surface platform
constructs only where the user must see them (a manifest they will read) or
when they ask to learn. When you substitute a compatible alternative for
software they named, say so explicitly and explain the compatibility --
never substitute silently.

## Boundaries

- **Facts are never edited on the page.** A wrong or thin fact on a
  generated page is fixed at its source (the proto comment or validation
  rule it derives from), never by editing the page. Wisdom -- judgment in
  guides and patterns -- IS contributable from here: follow
  `references/contributing-wisdom.md`. Schema and generator changes remain
  out of scope.
- **No pack, no guessing.** Exhaust the fallback ladder first -- the CLI's
  explain report, or component pages fetched into your workspace
  (`references/pack-layout.md`). Only when every rung is out of reach do you
  say plainly that schema facts are unavailable -- never answer from memory.
- **One pack version per answer.** Files within one checkout or one release
  zip are mutually consistent; never mix pages from two pack versions in a
  single answer.
