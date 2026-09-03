# Planton Helm Charts

Helm charts for running [Planton](https://planton.ai) on your own Kubernetes
cluster. All charts are published as OCI artifacts to GitHub Container
Registry and are publicly pullable — no registry login, no repo to add:

```
oci://ghcr.io/plantonhq/charts/<chart>
```

> Not to be confused with the repository's [`charts/`](../charts) directory —
> that is the **InfraChart catalog**, Planton's own blueprint format for
> deploying infrastructure *through* Planton. The charts here deploy Planton
> itself.

## Which chart do I want?

| Chart | Use it when |
|---|---|
| [`planton`](planton) | **Start here.** The batteries-included install: one `helm install` deploys the operator AND a `PlantonPlatform` resource with proven defaults. A few minutes later the full platform is running — console, API, sign-in, an in-cluster deploy runner — reachable over a single `kubectl port-forward` command the install prints for you. |
| [`planton-operator`](planton-operator) | You want the operator alone and will manage the `PlantonPlatform` resource yourself (GitOps, custom manifests). The `planton` chart composes this one — same operator either way. |
| [`planton-runner`](planton-runner) | Deploy a standalone Planton Runner into a remote cluster with a runner token — the runner enrolls itself with the control plane on first boot and receives its own identity. Normally driven by `planton runner deploy` rather than installed by hand. |

## The one-command install

```bash
helm install planton oci://ghcr.io/plantonhq/charts/planton \
  --namespace planton \
  --create-namespace
```

No values are required on any cluster. The install output tells you what to
do next: watch the platform reach `Ready`, run the printed port-forward
command, and open the console — the first visitor becomes the administrator.

### Per-cloud values files

To publish Planton at your own URL and give the runner a cloud identity, use
the recipe values files shipped beside the `planton` chart — they work
straight from a raw GitHub URL:

```bash
helm install planton oci://ghcr.io/plantonhq/charts/planton \
  --namespace planton --create-namespace \
  --values https://raw.githubusercontent.com/plantonhq/planton/main/helm/planton/values.rke2.yaml
```

| File | For |
|---|---|
| [`values.eks.yaml`](planton/values.eks.yaml) | Amazon EKS |
| [`values.gke.yaml`](planton/values.gke.yaml) | Google GKE |
| [`values.aks.yaml`](planton/values.aks.yaml) | Azure AKS |
| [`values.doks.yaml`](planton/values.doks.yaml) | DigitalOcean Kubernetes |
| [`values.rke2.yaml`](planton/values.rke2.yaml) | RKE2 / Rancher |

Only facts true of every such cluster are active in these files (for example,
RKE2's bundled nginx ingress class); everything specific to your account —
hostnames, IAM roles, certificates — is a precisely annotated comment, one
uncomment away.

## Rules worth knowing

- **One Planton operator per cluster.** The operator watches every namespace,
  so a second installation would fight the first over the same platforms —
  and it refuses to start, with the remedy in its log. Already running the
  operator? Install the umbrella with `--set planton-operator.enabled=false`.
- **Helm never upgrades CRDs.** When a chart release notes a schema change,
  apply the refreshed CRD before `helm upgrade`:
  `kubectl apply -f https://raw.githubusercontent.com/plantonhq/planton/main/helm/planton-operator/crds/planton.ai_plantonplatforms.yaml`
  (and likewise for `planton.ai_plantonidentityproviders.yaml`).
- **Your data outlives an uninstall.** `helm uninstall` removes the platform
  workloads but deliberately leaves the data services' volumes; delete the
  namespace to reclaim them — and always do so before reinstalling into the
  same namespace (generated credentials die with the release, so a reinstall
  against surviving volumes would not authenticate).

## Versioning and publishing

Each chart owns its semantic `version` in `Chart.yaml`; `appVersion` pins the
image line the chart deploys. Two version lines meet here: the operator is
built from this repository's `operator/` directory and released on its own
line (`ghcr.io/plantonhq/planton/operator:<tag>`, pinned by
`planton-operator`'s `appVersion`), while the platform images the operator
runs (control plane, runner, console) are built in the platform repository at
the version a `PlantonPlatform`'s `spec.version` names, which is what the
umbrella chart's `appVersion` records. Charts publish
automatically on every change pushed to `main`
([`release.helm.yaml`](../.github/workflows/release.helm.yaml)): lint →
package → push to the OCI registry. Published versions are immutable — bump
the chart `version` to release a change — and new chart packages are made
publicly pullable as part of the same workflow.
