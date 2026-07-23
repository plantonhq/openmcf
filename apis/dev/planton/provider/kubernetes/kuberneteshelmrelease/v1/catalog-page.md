# Kubernetes Helm Release

Installs an upstream Helm chart as a real Helm release through a single declarative manifest — hooks run, release history is recorded, and `helm list` sees the release exactly as if the Helm CLI had installed it.

**Check the catalog first.** If a first-class component exists for what you're deploying, use it: typed components validate their configuration before deploy and export composable outputs, while a generic chart install trusts you to get the chart's values right. KubernetesHelmRelease is the catalog's sole intentional passthrough, for charts no component covers — not the recommended path where one exists.

## What Gets Created

When you deploy a KubernetesHelmRelease resource, Planton provisions:

- **Helm Release** — a real release installed from the chart (HTTP(S) repository or OCI registry) at the pinned version, with your merged values. Everything the chart creates belongs to Helm; the module never reaches into the chart's own resources
- **Namespace** (optional) — when `create_namespace` is true, the target namespace is created with standard Planton governance labels and deleted with the resource

The chart's values surface is the configuration contract: values pass through to the chart unvalidated, and a typo'd value surfaces when the chart renders or its resources misbehave — not at apply.

## Prerequisites

- **Kubernetes credentials** configured via environment variables or Planton provider config
- **A reachable chart source** — the HTTPS repository or OCI registry must be accessible from where the deploy runs; private sources need the credential fields

## Quick Start

Create a file `helm-release.yaml` — podinfo 6.9.2 from its HTTPS repository:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesHelmRelease
metadata:
  name: podinfo
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesHelmRelease.podinfo
spec:
  namespace:
    value: podinfo
  create_namespace: true
  repo: https://stefanprodan.github.io/podinfo
  chart: podinfo
  version: 6.9.2
  values_yaml: |
    replicaCount: 2
```

Deploy:

```shell
planton apply -f helm-release.yaml
```

The release is a real Helm release — `helm list -n podinfo`, `helm status podinfo -n podinfo`, and `helm rollback` all work on it.

## The Values Model

Four layers, merged in this order on both engines (later wins), exactly as Helm defines them:

| Layer | Helm equivalent | Semantics |
|-------|-----------------|-----------|
| `values_yaml` | values file (`-f`) | A full YAML document: nested maps, lists, numbers, booleans. Applied first. No secrets here |
| `set` | `--set` | Dotted-path overrides with Helm's coercion: `"true"` → boolean, digits → number, `"null"` deletes the key |
| `set_string` | `--set-string` | Same paths, values stay literal strings — use for version-like tags (`"1.30"` would otherwise become the number 1.3) |
| `set_sensitive` | `--set-string`, secret | Literal strings, highest precedence, kept out of rendered plans and state where each engine supports it |

Both engines hand Helm one final merged map — the same manifest installs byte-identical releases on either engine.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `spec.namespace` | `StringValueOrRef` | Namespace to install into — a literal (`namespace: {value: my-ns}`) or a reference to a KubernetesNamespace resource. | required |
| `spec.repo` | `string` | Chart repository: an HTTP(S) Helm repository URL or an `oci://` registry reference. For OCI, the chart is pulled as `<repo>/<chart>:<version>`. | must match `^(https?\|oci)://` |
| `spec.chart` | `string` | Chart name within the repository (e.g. `podinfo`, `cert-manager`). | required |
| `spec.version` | `string` | Exact chart version (e.g. `6.9.2`). Required — unpinned installs are not reproducible. | required |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `spec.create_namespace` | `bool` | `false` | Create the namespace (with Planton governance labels) before installing; delete it with the resource. When false, the namespace must exist. |
| `spec.release_name` | `string` | `metadata.name` | Overrides the Helm release name. Lowercase alphanumerics and hyphens, at most 53 characters (Helm's limit). |
| `spec.values_yaml` | `string` | `""` | Chart values as a YAML document (see values model). |
| `spec.set` / `spec.set_string` / `spec.set_sensitive` | `map<string, string>` | `{}` | The three override layers (see values model). |
| `spec.repository_username` / `spec.repository_password` | `string` | `""` | Credentials for a private repository (HTTP basic auth) or OCI registry login. Must be set together — the spec rejects one without the other. Password is sensitive. |
| `spec.atomic` | `bool` | `false` | A failed install/upgrade rolls everything back — never a half-deployed release. Cannot be combined with `skip_await`. |
| `spec.cleanup_on_fail` | `bool` | `false` | A failed upgrade deletes the resources it newly created (lighter than `atomic`). |
| `spec.skip_await` | `bool` | `false` | Return once Helm records the release, without waiting for readiness. Default is to wait. |
| `spec.wait_for_jobs` | `bool` | `false` | When awaiting, also wait for chart-created Jobs to complete. |
| `spec.timeout_seconds` | `int32` | `300` | Maximum seconds per Kubernetes operation (and for readiness when awaiting). |
| `spec.skip_crds` | `bool` | `false` | Do not install CRDs from the chart's `crds/` directory. |
| `spec.dependency_update` | `bool` | `false` | Run `helm dependency update` before installing. Most published charts bundle dependencies and don't need this. |
| `spec.max_history` | `int32` | `10` | Release revisions kept for rollback (0 = unlimited). |
| `spec.replace` | `bool` | `false` | Re-use a release name left behind by a failed/uninstalled release. Helm marks this unsafe in production. |
| `spec.force_update` | `bool` | `false` | Force delete-and-recreate when a field can't be patched in place. Disruptive. |
| `spec.reuse_values` / `spec.reset_values` | `bool` | `false` | Upgrade-time values handling: keep the last release's values vs. reset to chart defaults. Mutually exclusive. |
| `spec.disable_webhooks` | `bool` | `false` | Chart hooks (pre/post install, upgrade, delete) do not run. |
| `spec.disable_openapi_validation` | `bool` | `false` | Skip validating rendered manifests against the cluster's schema. |
| `spec.take_ownership` | `bool` | `false` | Adopt existing resources into the release (Helm `--take-ownership`) — the migration knob. **Terraform provisioner only for now**; the Pulumi engine rejects it loudly rather than silently ignoring it. |
| `spec.description` | `string` | `""` | Free-form note stored with the release (visible in `helm status`). |

## Examples

### OCI Registry Chart

The same chart pulled from an OCI registry — only the `repo` scheme changes:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesHelmRelease
metadata:
  name: podinfo-oci
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesHelmRelease.podinfo-oci
spec:
  namespace:
    value: podinfo
  create_namespace: true
  repo: oci://ghcr.io/stefanprodan/charts
  chart: podinfo
  version: 6.9.2
  set:
    replicaCount: "2"
  set_string:
    image.tag: "6.9.2"
```

Note the split: `replicaCount` rides `set` (coerced to a number, which the chart expects), while `image.tag` rides `set_string` (stays a literal string — a tag like `"1.30"` under `set` would become the number 1.3).

### Private Repository With Secret Values

Credentials for the chart pull, a secret chart value kept out of plans and state, and production lifecycle knobs:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesHelmRelease
metadata:
  name: private-chart
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesHelmRelease.private-chart
spec:
  namespace:
    value: private-chart
  create_namespace: true
  repo: https://charts.example.com
  chart: my-app
  version: 1.4.2
  repository_username: ci-bot
  repository_password: <your-repository-password>
  set_sensitive:
    auth.apiKey: <your-api-key>
  atomic: true
  cleanup_on_fail: true
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs` — what `helm list` and `helm status` would show, exported identically by both engines:

| Output | Type | Description |
|--------|------|-------------|
| `namespace` | `string` | The namespace the release is installed in |
| `releaseName` | `string` | The Helm release name (`helm list` NAME column) |
| `version` | `string` | The installed chart version |
| `appVersion` | `string` | The chart's appVersion — the upstream application version the chart packages |
| `status` | `string` | The release status as Helm records it (e.g. `deployed`) |
| `revision` | `int32` | The release revision number (1 on install, incremented by upgrades/rollbacks) |

## Related Components

- [KubernetesNamespace](/docs/catalog/kubernetes/kubernetesnamespace) — manage the target namespace as its own resource and reference it from `spec.namespace` instead of using `create_namespace`
- [KubernetesManifest](/docs/catalog/kubernetes/kubernetesmanifest) — the passthrough for raw manifests, when what you have is YAML rather than a chart
