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
runner). The data services (PostgreSQL, Redis, NATS, Temporal) and their
volumes deliberately do not vanish with it -- your data outlives an
accidental uninstall -- so deleting the namespace is the step that reclaims
them. To REINSTALL, always delete the namespace first: an uninstall removes
the generated credential Secrets with the release, so a new install against
the surviving volumes would mint new passwords that the old data refuses.
The CRD is cluster-scoped; remove it last, after every PlantonPlatform in
the cluster is gone:

```bash
kubectl delete crd plantonplatforms.planton.ai
```
