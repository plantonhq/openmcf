# Pack Layout and Where to Find It

The reference pack is one directory tree. The pack root is the directory
that contains `reference-commons.md` -- that filename exists exactly once in
the pack, which makes it the reliable root marker.

## The shape

```
<pack-root>/
├── reference-index.md            # root index: provider table with kind,
│                                 # example, and guide counts (live numbers)
├── reference-commons.md          # the manifest grammar every component
│                                 # shares + the pack's search grammar
├── reference-graph.yaml          # every foreign-key edge in the catalog
├── GUIDE.md                      # catalog-level wisdom: the substitution
│                                 # workflow and verified alternatives
├── patterns/                     # architecture patterns (README.md lists
│                                 # them; one .md per pattern, each with
│                                 # validated manifests)
└── <provider>/                   # aws/, gcp/, azure/, kubernetes/, ...
    ├── reference-index.md        # every kind in the provider, one line each
    └── <kind>/<api-version>/
        ├── reference.md          # the component's complete reference page
        └── GUIDE.md              # authored judgment -- present only where
                                  # wisdom has been written
```

## Finding the pack, in order

Every probe below stays inside the filesystem your session already grants
-- the working tree, an attached workspace, or a mount your tools handed
you. Never search the wider machine for pack files (home directories, tool
install paths, other checkouts that happen to exist on the host): a pack
that is not reachable inside that boundary is simply not reachable -- take
the fallback below.

1. **A Planton open-source repo checkout** -- the pack root is
   `catalog/` under the repo root. This is the common case
   for coding agents working in or beside the repo.
2. **The release artifact** -- every Planton release publishes the pack as
   one zip:

   ```
   https://downloads.planton.dev/releases/<version>/content/reference-pack.zip
   ```

   Entries carry repo-relative paths, so after extracting into an empty
   directory the pack root is
   `<dir>/catalog/`.
3. **A vendored copy**, when the tree you are working in embeds one. Do not
   assume it exists -- probe for the root marker and use it only on a hit:

   ```
   rg --files -g 'reference-commons.md' <tree>
   ```

   Its directory is the pack root.

## When no pack is reachable

`planton explain <Kind>` prints the same schema-derived facts for one
component at a time -- fields, validation rules, outputs, outbound
references -- and works fully offline. Use it as the per-component fallback,
and be honest about what it cannot give you: the indexes ("what exists"),
full-text capability search, the inbound `Referenced By` view, the
catalog-wide graph, and the authored guides and patterns. When an answer
would have needed one of those, say so instead of approximating.

## Freshness

The pack is regenerated from the schemas and shipped whole: within one
checkout or one release zip, every page, index, and the graph are mutually
consistent. Never combine files from two different pack versions in one
answer -- resolve one pack root and stay in it.
