# Architecture Patterns

This folder holds authored architecture patterns: named compositions of two
or more catalog kinds, with validated manifest wiring and the trade-offs
behind each choice. This is NOT a provider directory -- it sits beside the
catalog's generated indexes because patterns span providers.

Patterns are the catalog's judgment layer for anything bigger than one kind:
where a kind's `GUIDE.md` teaches wisdom about that kind alone, a pattern
teaches how kinds compose -- and what each composition looks like on the
architecture diagram the platform renders.

## Reading a pattern

One file per pattern. Frontmatter declares the kinds the pattern composes,
so tooling (and you) can find every pattern touching a kind:

```
rg -l "KubernetesNamespace" patterns/     # every pattern composing a kind
```

Each pattern states the problem, the composition choices with their
consequences (state ownership, diagram shape, customization surface), and
complete manifests you can adapt.

## Writing a pattern

Patterns are contributed through pull requests -- the file you edit here is
exactly the file agents read, in the repo, the release archive, and the
packaged skill alike. Embedded complete manifests must validate against
their schemas, and declared kind names must resolve: catalog checks enforce
both, so a pattern can age gracefully but can never be silently
schema-wrong. Authoring standards live in the repository's contribution
docs.
