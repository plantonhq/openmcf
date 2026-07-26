# Tekton

Cloud-native CI/CD on your own cluster: every pipeline step is a
container, every pipeline a Kubernetes resource. This component
declares the whole installation — which pieces run, how pipelines
behave, and how completed runs get cleaned up — as one resource the
Tekton Operator keeps converged.

## Highlights

- **Profiles, honestly** — Pipelines alone (`lite`), + Triggers
  (`basic`), or the full stack with the Dashboard (`all`); Chains
  supply-chain signing rides along unless disabled.
- **The pipeline surface, typed** — execution defaults, API stability
  gates, feature flags, resolver toggles, metrics granularity, and the
  controller's replicas+buckets HA tuning.
- **One event stream** — the cluster-global CloudEvents sink streams
  every run's lifecycle to your orchestration (one URL by Tekton's
  design; fan out downstream for multi-tenant routing).
- **Cleanup as configuration** — the pruner cron (keep newest N or
  younger-than) is a first-class field, because unbounded completed
  runs are the default failure mode of busy CI clusters.
- **Destroys clean** — teardown completes while the operator still
  runs; the stranded-finalizer hang is designed out, not scripted
  around.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
