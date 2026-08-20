# planton

The batteries-included Planton install: one `helm install` deploys the
[Planton operator](../planton-operator) **and** a `PlantonPlatform` resource
with proven defaults. A few minutes later the full platform is running --
console, API, sign-in, and an in-cluster runner -- reachable over a single
`kubectl port-forward` command that the install prints for you.

```bash
helm install planton oci://ghcr.io/plantonhq/charts/planton \
  --namespace planton \
  --create-namespace
```

No values are required. The first person to open the console becomes the
administrator (the setup page asks for their email plus a setup code read
from the cluster).

Prefer to manage the `PlantonPlatform` resource yourself (GitOps, custom
manifests)? Install the [`planton-operator`](../planton-operator) chart
instead -- this chart simply composes it.

## Values

```yaml
# Passed through to the bundled planton-operator subchart (any key from that
# chart's values.yaml works here):
planton-operator:
  enabled: true          # set false to use an operator you already run

platform:
  name: planton          # name of the PlantonPlatform resource
  spec: {}               # the PlantonPlatform spec, passed through VERBATIM
```

`platform.spec` is not curated by this chart: every field the
`PlantonPlatform` CRD supports works here, today and in future operator
versions. `spec.version` defaults to this chart's `appVersion`. See the
[operator chart's README](../planton-operator/README.md) for the full
resource reference (exposure ladder, identity, runner, components).

## Per-cloud values files

The base install needs no values file on any cloud. The files below are
recipes for the graduation step -- publishing Planton at your own URL and
giving the runner a cloud identity -- with only universally-true facts active
(e.g. RKE2's bundled nginx ingress class) and everything account-specific as
annotated comments:

| File | For |
|---|---|
| [`values.eks.yaml`](values.eks.yaml) | Amazon EKS (ALB/nginx ingress, IRSA runner identity, Secrets Manager backend) |
| [`values.gke.yaml`](values.gke.yaml) | Google GKE (built-in `gce` ingress active, Workload Identity runner) |
| [`values.aks.yaml`](values.aks.yaml) | Azure AKS (application routing add-on, Entra Workload ID runner) |
| [`values.doks.yaml`](values.doks.yaml) | DigitalOcean Kubernetes (nginx ingress, credentials Secret runner) |
| [`values.rke2.yaml`](values.rke2.yaml) | RKE2 / Rancher (bundled nginx ingress active) |

Use them straight from GitHub:

```bash
helm install planton oci://ghcr.io/plantonhq/charts/planton \
  --namespace planton --create-namespace \
  --values https://raw.githubusercontent.com/plantonhq/planton/main/helm/planton/values.rke2.yaml
```

## Storage

Planton's data services (PostgreSQL, Valkey, the runner's
IaC-state volume, and any optional components you enable) all request
persistent volumes. Two prerequisites decide whether an install works out of
the box on your cluster:

1. **A StorageClass Planton can provision from.** Either your cluster has a
   default StorageClass whose storage driver is actually installed, or you
   pin one explicitly. A default class that merely *exists* is not enough --
   on some managed clusters (fresh EKS is the famous case) the default class
   points at a CSI driver nobody installed, and every volume request hangs.
2. **Volume sizes at or above your backend's minimum.** Some enterprise
   backends enforce large minimum volume sizes (some NAS systems refuse
   anything under 800Gi). Planton's defaults (1-10Gi per volume) are sized
   for ordinary clusters; on such backends set sizes explicitly.

Both are one block in the values -- the class and one size for every volume,
with per-component overrides when volumes should differ:

```yaml
platform:
  spec:
    storage:
      storageClassName: your-class   # every volume; omit to use the cluster default
      size: 800Gi                    # every volume; omit for per-component defaults
    # Per-component settings win over the storage block:
    # database:
    #   postgresql: {storageSize: 2Ti, storageClassName: fast-ssd}
    #   redis:      {storageSize: 1Gi}
    # runner:
    #   storageSize: 2Gi
    #   storageClassName: your-class
    # vault:
    #   storageSize: 2Gi
    #   storageClassName: your-class
```

Set storage before installing: Kubernetes fixes a volume's class and size at
creation, so changing them on a running install requires recreating the
volume (scale the component down, delete its PersistentVolumeClaim, let the
operator recreate it -- destroying that volume's data). The install preflight
checks the default class can provision, and if a volume ever sticks, the
`PlantonPlatform` component status names the exact problem and fix:

```bash
kubectl get plantonplatform planton -n planton -o jsonpath='{.status.components}' | jq
```

## Secrets

Every install ships with a bundled secrets manager (OpenBAO, the open-source
Vault fork), deployed and bootstrapped automatically: the operator
initializes it, unseals it after every restart, and registers it as the
organization's default secret store -- pasting a cloud credential into a
connection works on a fresh install with zero configuration. The unseal keys
and root token live in the `<name>-openbao-init` Secret next to the platform
(the Secret's annotation explains itself); teams that prefer to hold their
own keys set `spec.vault.initMode: manual` and run the unseal ceremony
themselves.

Preferring a cloud secret store is a layered choice, not an opt-out: declare
it and it wins (see `values.eks.yaml` for AWS Secrets Manager with pod
identity -- no stored keys at all). The vault still runs underneath, serving
the platform's own signing and encryption keys.

Opting out entirely is supported but deliberate:

```yaml
platform:
  spec:
    vault:
      enabled: false   # loses pasted-credential storage and keyless connections
```

## One operator per cluster

Run exactly one Planton operator per cluster. If this cluster already runs
the standalone `planton-operator` chart, install the umbrella with
`--set planton-operator.enabled=false` so it only creates the
`PlantonPlatform` resource. A second operator refuses to start and its log
names the conflict and both ways out -- two operators would fight over the
same platforms.

## Upgrades

`helm upgrade` rolls the operator and the platform version together. One
Helm limitation to know: Helm never upgrades CRDs. When a new chart version
carries CRD schema changes, apply the refreshed CRD first:

```bash
kubectl apply -f https://raw.githubusercontent.com/plantonhq/planton/main/helm/planton-operator/crds/planton.ai_plantonplatforms.yaml
```

## Teardown

```bash
helm uninstall planton -n planton
kubectl delete namespace planton
```

`helm uninstall` removes the `PlantonPlatform` resource, the operator, and
the platform's own workloads (control plane, console, front door, identity,
runner). The PostgreSQL database goes WITH the platform: CloudNativePG owns
the database's volumes and credentials as one unit, so removing the platform
removes the database cleanly -- credentials can never survive apart from the
volumes they unlock. The remaining data-service volumes (the Valkey cache,
Temporal's chart-rendered claims) linger until the namespace is deleted,
which is why deleting the namespace is part of teardown. To REINSTALL,
always delete the namespace first: a new install against leftover volumes
would mint new credentials that the old data refuses. The CRD is
cluster-scoped; remove it last, after every PlantonPlatform in the cluster
is gone:

```bash
kubectl delete crd plantonplatforms.planton.ai
```

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
