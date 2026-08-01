# Harbor

CNCF-graduated container registry: store, sign, and scan OCI artifacts
from the official chart, on a data plane you compose. Admin credentials
are generated per install — the chart's public defaults never ship —
and north-south exposure is a separate composed resource pointed at the
exported front-door Service.

## Highlights

- **Composition-first data plane** — PostgreSQL, Redis, and artifact
  storage each have an evaluation arm and an external arm; a
  `KubernetesPostgres`, `KubernetesValkey`, and `KubernetesSeaweedFs`
  pair naturally for production.
- **Secret-native credentials** — every credential site is an
  `existingSecret` reference (or a module-materialized Secret); the
  publicly documented `Harbor12345` / `changeit` / `not-a-secure-key`
  defaults never reach the cluster.
- **Exposure composes** — ClusterIP / NodePort / LoadBalancer front
  door plus a required `externalUrl` (load-bearing for OCI auth);
  ingress and Gateway API ride the catalog's exposure kinds.
- **Trivy first-class** — vulnerability scanning on by default, with
  honest air-gap knobs (`skip_update`, `offline_scan`) for clusters
  that cannot reach the public vulnerability DB.
- **Name budget enforced** — `metadata.name` capped at 39 characters
  so the chart's longest derived object name stays inside Kubernetes'
  63-character limit.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
