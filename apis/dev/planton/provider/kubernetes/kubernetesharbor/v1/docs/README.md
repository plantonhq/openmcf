# Deploying Harbor on Kubernetes: Composition Over Bundles

## Introduction

Harbor is a CNCF-graduated container registry — store, sign, and scan
OCI artifacts behind one front door. The deployment question is not
"can Helm install it" (it can, in one command). The question is which
parts of the stack you let the chart own, and which you compose from
purpose-built backends.

This component takes a firm position: the chart's in-cluster PostgreSQL
and Redis are evaluation-grade by upstream's own documentation, artifact
storage on a filesystem PVC is fine until you need multi-replica
registries, and credentials that ship as chart defaults are public
knowledge. Production means composing the catalog's own Postgres,
Valkey, and object-storage kinds, with every credential riding a Secret.

## What Planton models

- **Database** — internal XOR external. The external arm takes host and
  username as references and a password Secret whose key is the chart's
  contract key `password` — exactly the shape a `KubernetesPostgres`
  application Secret exports.
- **Cache** — internal XOR external Redis. A `KubernetesValkey` speaks
  the Redis protocol; bridge its username-keyed auth Secret to the
  chart's `REDIS_PASSWORD` contract, or declare the password and let
  the module materialize `<name>-redis-auth`.
- **Artifact storage** — filesystem XOR S3-compatible XOR GCS XOR
  Azure. S3-compatible pairs with `KubernetesSeaweedFs` (set
  `disable_redirect` for in-cluster endpoints clients cannot reach) or
  with cloud S3 under IRSA/ambient credentials (omit the credentials
  block entirely).
- **Exposure** — ClusterIP (default), NodePort, or LoadBalancer. The
  chart's ingress and Gateway API route types are deliberately not
  modeled; point a catalog exposure kind at the exported
  `expose_service`. `externalUrl` is required and load-bearing: Harbor
  embeds it in the token-service URL, so OCI push/pull fail auth when
  the dialed address disagrees.
- **Admin credential** — unset means module-generated into
  `<name>-admin-auth` (key `HARBOR_ADMIN_PASSWORD`). The chart default
  `Harbor12345` never ships.
- **Trivy** — on by default. Air-gapped clusters set `skip_update` and
  `offline_scan` and preload the vulnerability DB; with `skip_update`
  and no pre-loaded DB, scans fail closed.
- **Notary and Clair** — not modeled. Notary left Harbor upstream at
  v2.6; Clair was replaced by Trivy.

## Evaluation vs production

| Concern | Evaluation (minimal preset) | Production (composed preset) |
|---|---|---|
| Database | chart-internal PostgreSQL | `KubernetesPostgres` reference |
| Cache | chart-internal Redis | `KubernetesValkey` + bridge Secret |
| Artifacts | filesystem PVC | S3-compatible / GCS / Azure |
| Exposure | ClusterIP + port-forward | LoadBalancer or composed ingress |
| Admin password | module-generated Secret | same (or your existing Secret) |

## Destroy truth

With the default `keep_volumes_on_uninstall`, the registry and
jobservice PVCs carry `helm.sh/resource-policy: keep` and survive
uninstall for a reinstall to adopt. The internal database and Redis
volumes are StatefulSet-template PVCs Helm never deletes regardless.
Retiring an install for good means sweeping those PVCs explicitly.
Harbor installs no CRDs — destroy leaves no cluster-scoped residue.

## Name budget

`metadata.name` caps at 39 characters. The chart truncates its
fullname at 63 and then appends component suffixes; the longest
(`-jobservice-internal-tls`) is 24 characters. Both engines fail loud
when the budget is exceeded.

## Quick start

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesHarbor
metadata:
  name: registry
spec:
  namespace:
    value: harbor
  createNamespace: true
  externalUrl: http://localhost:8080
  database:
    internal: {}
  cache:
    internal: {}
  storage:
    filesystem:
      diskSize: 20Gi
```

Retrieve the admin password:

```bash
kubectl get secret registry-admin-auth -n harbor \
  -o jsonpath='{.data.HARBOR_ADMIN_PASSWORD}' | base64 -d
```

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
