# Research Recipes

The concrete command moves for catalog research. Every recipe runs from
`<pack-root>` (see `references/pack-layout.md`). They work because the pack
renders a fixed heading vocabulary on every page -- documented in
`reference-commons.md` and pinned by the catalog's own contract tests, so
these patterns hold across every component.

## Find components by name or capability

```
rg -il "kafka" -g 'reference*.md' -g 'GUIDE.md' -g 'patterns/*.md' .
```

The `-g` globs scope the search to pack files -- in a repo checkout the same
directories also hold protos, generated code, and IaC modules, and an
unscoped search drowns the answer in those. Then shortlist through the
provider's `reference-index.md` -- each row is
`[Kind](path) | purpose | example? | guide?`. Full-text search matters more
than it looks: compatible alternatives document the well-known names they
substitute for in their own pages, so searching the name the user said finds
the alternative even when no component carries that name. When that happens,
follow the substitution workflow in the catalog root's `GUIDE.md` -- propose
openly, never silently.

## One component's required fields

Read `## Example` first -- it is a validated manifest, the fastest picture
of the component's real shape. Then the spec table:

```
rg "^## Spec Fields" -A 40 <page>
```

Columns: `Path | Type | Required | Default | References`. A `References`
cell names the kind (and field path) that field points at by default.

## One field, fully

```
rg "^### spec.replication" -A 12 <page>
```

Each field has its own `### <dotted path>` block under `## Field Details`:
docs, validation rules, and allowed enum values. Enum lists sit BELOW the
field's doc prose -- when a list looks truncated, widen `-A` rather than
concluding the values are missing.

## What a component exports

```
rg "^## Outputs" -A 20 <page>
```

Output paths are spelled snake_case (`status.outputs.vpc_id`) -- that is the
canonical `fieldPath` spelling, while spec YAML keys are camelCase; the
asymmetry is explained once in `reference-commons.md`.

## Wiring: both directions

- Outbound (what K's fields can point at): K's `## References` table.
- Inbound (who can point at K): K's `## Referenced By` table --
  `Kind | Field | Reads`.
- Catalog-wide, without opening pages:

```
rg 'to: "KubernetesValkey"' reference-graph.yaml      # every field that can target it
rg -A 3 'from: "AwsEcsService"' reference-graph.yaml  # everything it can reference
```

## Where judgment has been written

A page that has authored wisdom links it in its head:

```
rg -l '^\*\*Guide\*\*:' kubernetes/            # every guided kind in a provider
```

The per-provider index carries the same signal as its Guide column, and
`patterns/` (its `README.md` is the list) holds the multi-component recipes
-- each pattern declares the kinds it composes and embeds validated
manifests.

## Craft notes, learned the hard way

- Prefer anchored heading greps (`rg "^## Outputs" -A 20`) and then reading
  the section over regexing markdown table cells -- cell patterns are
  whitespace-sensitive and silently miss rows.
- `-l` before content: when a search may hit many pages, list files first
  (`rg -il <term> .`), pick from the index, then read one page deeply.
- The `## Example` manifest is validated against the schema at generation
  time -- trust it as a starting shape, then adjust against `## Spec Fields`
  and `## Validation Rules`.
