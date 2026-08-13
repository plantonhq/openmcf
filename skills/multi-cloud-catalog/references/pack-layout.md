# Pack Layout and Where to Find It

The reference pack is one directory tree. This skill ships self-contained:
the pack travels inside the skill itself, so wherever this file is mounted,
the pack sits beside it in `components/`. You never need a repo checkout,
a network fetch, or a marker hunt to start researching.

## The shape

```
<pack-root>/
├── _docs/
│   ├── reference-index.md        # root index: provider table with kind,
│   │                             # example, and guide counts (live numbers)
│   ├── reference-commons.md      # the manifest grammar every component
│   │                             # shares + the pack's search grammar
│   ├── reference-graph.yaml      # every foreign-key edge in the catalog
│   └── GUIDE.md                  # catalog-level wisdom: the substitution
│                                 # workflow and verified alternatives
├── _patterns/                    # architecture patterns (README.md lists
│                                 # them; one .md per pattern, each with
│                                 # validated manifests)
└── <provider>/                   # aws/, gcp/, azure/, kubernetes/, ...
    ├── reference-index.md        # every kind in the provider, one line each
    └── <kind>/
        ├── GUIDE.md              # authored judgment -- present only where
        │                         # wisdom has been written
        └── <api-version>/
            └── reference.md      # the component's complete reference page
```

## Finding the pack, in order

Every probe below stays inside the filesystem your session already grants
-- the working tree, an attached workspace, or a mount your tools handed
you. Never search the wider machine for pack files (home directories --
`~/.stigmer` included -- tool install paths, other checkouts that happen
to exist on the host): a pack that is not reachable inside that boundary
is simply not reachable -- take the fallback below.

1. **Your own skill mount** -- the pack root is `components/` in the same
   directory as this skill's `SKILL.md`. This is the common case: the
   skill and its pack are one artifact, published and versioned together,
   so the pack is always exactly as fresh as the skill teaching you to
   read it.
2. **A Planton open-source repo checkout** -- the pack root is `catalog/`
   under the repo root. Prefer it over the mount only when you are working
   IN the repo (contributing, or reviewing unreleased changes): a checkout
   may carry pages newer than any published skill.
3. **The release artifact** -- every Planton release also publishes the
   pack as one zip:

   ```
   https://downloads.planton.dev/releases/<version>/content/reference-pack.zip
   ```

   Entries carry repo-relative paths, so after extracting into an empty
   directory the pack root is `<dir>/catalog/`. Useful for tools that want
   the pack without the skill; as a mounted agent you should never need it.

## When no pack is reachable

Two per-component fallbacks, in order of preference:

1. **`planton explain <Kind>`** (when the CLI exists where you run) prints
   the same schema-derived facts for one component at a time -- fields,
   validation rules, outputs, outbound references -- fully offline.
2. **Fetch the component's pages into your workspace** (when you have
   network but no CLI): the generated reference page and its neighbors live
   at stable public paths --

   ```
   https://raw.githubusercontent.com/plantonhq/planton/main/catalog/<provider>/<kind>/<api-version>/reference.md
   ```

   Download into the workspace you already own and read there -- the fetch
   stays inside your filesystem boundary, and the page carries the same
   facts the pack ships. The provider index
   (`catalog/<provider>/reference-index.md`) and the commons page fetch the
   same way when discovery or grammar questions arise. Prefer fetching only
   the pages the answer needs over mirroring the tree.

Be honest about what per-component fallbacks cannot give you: full-text
capability search, the inbound `Referenced By` view, the catalog-wide graph,
and the authored guides and patterns (the fetched index partially covers
"what exists"). When an answer would have needed one of those, say so
instead of approximating.

## Freshness

The pack is regenerated from the schemas and shipped whole: within one
skill mount, one checkout, or one release zip, every page, index, and the
graph are mutually consistent. Never combine files from two different pack
versions in one answer -- resolve one pack root and stay in it. Your mount's
pack matches your skill's version by construction; only a repo checkout can
legitimately be newer.
