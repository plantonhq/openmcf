# Tekton Operator

The lifecycle manager for Tekton: it installs, upgrades and removes the
Tekton components from one declaration, so the cluster's CI engine is
managed like any other reconciled resource — never by hand-applying
release manifests.

## Highlights

- **The manager, cleanly separated** — this resource installs the
  operator alone; what Tekton actually runs is declared with a
  KubernetesTekton resource the operator reconciles. Nothing fights
  over the configuration.
- **Official distribution, pinned** — the single-file release manifest
  at an exact tag; reproducible installs, no moving `latest`.
- **Clean teardown by design** — the destroy ordering that used to hang
  IaC tools on stranded finalizers is structural now: the declaration
  tears down while the operator still runs, then the operator leaves.
- **Air-gap ready** — typed image overrides for both Deployments plus
  pull-secret wiring for private mirrors.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
