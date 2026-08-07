# Helm Release

Installs an upstream Helm chart as a REAL Helm release — the catalog's sole intentional passthrough. Both engines perform an actual `helm install`: hooks run, the release secret is written, and `helm list` shows the release exactly as if installed by the Helm CLI. Reach for this only when the catalog has no first-class component for the chart you need: typed components validate their configuration before deploy, export composable outputs, and teach their trade-offs field by field — a generic chart install does none of that.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Helm Release** -- the chart fetched from `repo` (an HTTPS repository index or an `oci://` registry pull), installed at the PINNED `version` into `namespace`, named `release_name` when set (otherwise `metadata.name`)
- **Everything the chart renders** -- workloads, Services, ConfigMaps, and (unless `skip_crds`) the CRDs shipped in the chart's `crds/` directory; Helm installs chart CRDs only when absent and never upgrades or deletes them
- **Namespace** (optional) -- created with standard governance labels when `create_namespace` is true, and deleted with the resource

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** -- an active connection in the Connect module with credentials for the target cluster. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- For a **private chart source**: an organization secret holding the repository password or registry token — `repository_password` stores a `$secret/<name>` reference, resolved just-in-time at deploy.

### Cluster

- The chart repository or OCI registry must be reachable from the execution environment.
- If `create_namespace` is false, the target namespace must already exist.
- Any operators or CRDs the CHART's resources depend on (but the chart does not itself ship) must be in place first.

## Deploy

### Console

Open the deployment store, find **Helm Release**, and click **Deploy**. The creation wizard walks you through placement, the pinned chart coordinates (with OCI-vs-HTTPS teaching), private-source access, the values file, the three `--set` override layers (with reference-only pickers for sensitive values), install behavior (atomic vs skip-await), CRDs and dependencies, upgrade behavior (reuse vs reset, rollback history), and the advanced escape dials. Start from a preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesHelmRelease
metadata:
  name: podinfo
  org: acme-corp
  env: prod
spec:
  namespace:
    value: apps
  createNamespace: true
  repo: https://stefanprodan.github.io/podinfo
  chart: podinfo
  version: 6.9.2
  valuesYaml: |
    replicaCount: 2
    resources:
      requests:
        cpu: 100m
  setString:
    image.tag: "6.9.2"
  atomic: true
```

```shell
planton apply -f helm-release.yaml
```

## Key Configuration

These are the most important decisions when configuring a Helm release. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**The version is pinned on purpose** -- an unpinned "latest" install is not reproducible, and reproducibility is the point of declaring the release here. Editing the version on a deployed release IS the `helm upgrade`.

**The values model is Helm's own** -- `values_yaml` is the values file (applied first), and `set`, `set_string`, `set_sensitive` are the `--set` layers applied on top, in that order. Both engines merge the layers with identical precedence and hand Helm ONE final values map — the same manifest installs byte-identical releases on either engine.

**The coercion boundary** -- `set` values get Helm's coercion ("true"/"false" become booleans, digits become numbers, null removes the key); a version-like value that must STAY a string ("1.30") belongs in `set_string`.

**Secrets never enter the spec** -- `set_sensitive` values and `repository_password` are `$secret/<name>` references to organization secrets, resolved at deploy time and kept out of rendered plans and state where each engine supports it.

**Atomic vs skip-await** -- atomic rolls a failed install/upgrade back whole (it must WAIT for readiness to know whether to roll back), skip-await returns as soon as Helm records the release; the spec rejects the pair. `wait_for_jobs` extends the readiness wait to chart Jobs and is ignored while the wait is skipped.

**Upgrades are declarative by default** -- the values declared here are the whole story on every upgrade. `reuse_values` merges over the previous release's values instead, `reset_values` starts from chart defaults — mutually exclusive. `max_history` bounds rollback history (0 keeps everything).

**Adoption is engine-asymmetric** -- `take_ownership` (adopt matching resources an earlier tool created) is honored by the Terraform provisioner only; the Pulumi engine's pinned SDK rejects the flag loudly rather than silently ignoring it.

## Outputs and Dependencies

### What This Component Consumes

| Field | References | Purpose |
|-------|-----------|---------|
| `spec.namespace` | KubernetesNamespace (`spec.name`) | The namespace the release installs into |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `namespace` | The namespace the release is installed in | Debugging and composition |
| `release_name` | The Helm release name (as shown by `helm list`) — `spec.release_name` when set, otherwise `metadata.name` | Addressing the release's derived objects |
| `version` | The installed chart version | Auditing which chart bytes run |
| `app_version` | The chart's appVersion — the upstream application version the chart packages | Auditing which software version runs |
| `status` | The release status as Helm records it (e.g. `deployed`) | Health checks and composition gates |
| `revision` | The release revision (1 on first install, incremented by each upgrade or rollback) | Change tracking |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Public HTTPS chart** -- a pinned chart from a public repository with a values document — the everyday shape.

**OCI registry chart** -- an `oci://` reference pulled like a container image, with targeted `set`/`set_string` overrides.

**Private repository with secrets** -- repository credentials and sensitive chart values, all as organization-secret references; atomic install for rollback safety.

## Works With

- **Kubernetes Namespace** -- the placement target; reference one by name or compose it as a first-class resource.
- **First-class chart components** -- when the catalog grows a typed component for a chart you run through this kind (cert-manager, ingress-nginx, Istio, and the rest already exist), prefer it: validation, outputs, and field-level teaching come with it.
