# Planton Operator Helm Chart

Deploy the [Planton Operator](https://github.com/plantonhq/planton) to any Kubernetes cluster.

The Planton Operator watches for `PlantonPlatform` custom resources and deploys the full
Planton platform stack: PostgreSQL, Valkey, OpenFGA, Temporal, the control plane
monolith, and the web console. A single YAML file is all it takes to go from an empty cluster
to a running Planton instance.

Want the operator AND the platform in one command? Use the
[`planton` umbrella chart](../planton) instead -- it composes this chart and adds the
`PlantonPlatform` resource with proven defaults. Run exactly one Planton operator per
cluster: a second installation refuses to start and its log explains why.

## Prerequisites

- Kubernetes 1.24+
- Helm 3.x

## Installation

The chart is published to GitHub Container Registry as an OCI artifact:

```bash
helm install planton-operator oci://ghcr.io/plantonhq/charts/planton-operator \
  --namespace planton-operator-system \
  --create-namespace
```

Pin a specific chart version with `--version <x.y.z>`. The operator image the chart
deploys is set by the chart's `appVersion` (override with `--set image.tag=<tag>`).

After the operator is running, create a `PlantonPlatform` resource to deploy the platform:

```yaml
apiVersion: planton.ai/v1
kind: PlantonPlatform
metadata:
  name: planton
spec:
  version: v1.0.0
```

Apply it with `kubectl apply -f` and watch progress:

```bash
kubectl get plantonplatform -w
```

### Publishing Planton at a URL

External access is a friction ladder in `spec.ingress` -- each rung is one field
up from the last, and every rung is a working deployment:

```yaml
spec:
  ingress:
    enabled: true                # rung 1: URL auto-derived from the ingress
                                 #   controller's address (magic DNS), plain HTTP
    hostname: planton.corp.com   # rung 2: your hostname, plain HTTP
    ingressClassName: nginx      # optional; omit to use the cluster default
    tls:                         # rung 3/4: HTTPS
      secretName: planton-tls    #   EITHER a kubernetes.io/tls Secret you bring
      # issuer:                  #   OR a cert-manager issuer
      #   name: lets-encrypt
      #   kind: ClusterIssuer
```

One hostname serves the web console and the API the browser calls. The platform
reports its URL in `status.consoleUrl` (the `URL` column of
`kubectl get plantonplatform`), and the ingress component's status explains any
misconfiguration in plain language (missing class, missing TLS secret,
cert-manager not installed).

## Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `image.repository` | Container image repository | `ghcr.io/plantonhq/planton/operator` |
| `image.tag` | Container image tag | Chart `appVersion` |
| `image.pullPolicy` | Image pull policy | `IfNotPresent` |
| `replicaCount` | Number of operator replicas | `1` |
| `leaderElection.enabled` | Enable leader election for HA | `true` |
| `healthProbe.port` | Health check endpoint port | `8081` |
| `resources.requests.cpu` | CPU request | `10m` |
| `resources.requests.memory` | Memory request | `256Mi` |
| `resources.limits.cpu` | CPU limit | `500m` |
| `resources.limits.memory` | Memory limit | `512Mi` |
| `serviceAccount.create` | Create a ServiceAccount | `true` |
| `serviceAccount.annotations` | ServiceAccount annotations | `{}` |
| `serviceAccount.name` | Override ServiceAccount name | `""` |
| `nodeSelector` | Node selector for scheduling | `{}` |
| `tolerations` | Tolerations for scheduling | `[]` |
| `affinity` | Affinity rules for scheduling | `{}` |

## Architecture

The operator runs as a single Deployment and watches for `PlantonPlatform` resources
in any namespace. When a resource is created, the operator:

1. Deploys the data layer (PostgreSQL via CloudNativePG, Valkey)
2. Deploys supporting services (OpenFGA with authorization model, Temporal with schema)
3. Deploys the application layer (control plane monolith, web console)

Each component is reconciled independently with explicit dependency tracking.
The operator reports per-component status and an aggregate `Ready` condition:

```
$ kubectl get plantonplatform
NAME      PHASE   VERSION   AGE
planton   Ready   v1.0.0    5m
```

## CRD Management

The operator's CRDs (`PlantonPlatform` and `PlantonIdentityProvider`) are installed
automatically during `helm install`. Helm places CRDs in a special `crds/` directory that
has the following behavior:

- CRDs are installed **before** any templates
- CRDs are **not deleted** on `helm uninstall` (protects existing custom resources)
- CRDs are **not upgraded** on `helm upgrade`

**Source of truth:** the files in `crds/` and the manager's permissions in
`rbac/manager-role.yaml` are controller-gen output, written by
`make -C operator manifests` in this repository from the operator's Go types and RBAC
markers. CI regenerates and diffs them on every change, so they are never edited by hand:
change the Go source, regenerate, and the chart follows in the same commit.

To upgrade the CRDs after a chart version bump:

```bash
kubectl apply -f https://raw.githubusercontent.com/plantonhq/planton/main/helm/planton-operator/crds/planton.ai_plantonplatforms.yaml
kubectl apply -f https://raw.githubusercontent.com/plantonhq/planton/main/helm/planton-operator/crds/planton.ai_plantonidentityproviders.yaml
```

## Uninstallation

```bash
# Remove the operator (CRD and custom resources are preserved)
helm uninstall planton-operator -n planton-operator-system

# To also remove the CRD (destroys all PlantonPlatform resources):
kubectl delete crd plantonplatforms.planton.ai
```

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
