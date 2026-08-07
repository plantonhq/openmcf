# Kubernetes Helm Release

## When NOT to Use This

**A first-class catalog component always wins.** Typed components validate their configuration before anything reaches a cluster, export composable outputs other resources can reference, and teach their trade-offs field by field. A generic chart install does none of that — the chart's values surface is the contract, and a typo'd value is discovered when the chart's resources misbehave, not at validation.

**KubernetesHelmRelease is the catalog's sole intentional passthrough.** Reach for it only when the catalog has no component for the chart you need. It is never the recommended path where a first-class component exists.

## Overview

**KubernetesHelmRelease** installs an upstream Helm chart as a *real* Helm release: on both engines, hooks run, the release secret is written, history is recorded, and `helm list` shows the release exactly as if the Helm CLI had installed it. This is not a client-side template render — charts that rely on hooks (migrations, certificate bootstrapping, pre-delete cleanup) behave exactly as Helm intends.

The spec covers chart identity (HTTP(S) repositories and OCI registries, private-repo credentials, a required pinned version), Helm's full values model (a values file plus three `--set`-style override layers with fixed precedence), and the release lifecycle surface (`atomic`, `cleanup_on_fail`, awaiting, timeouts, history, CRD handling, upgrade values behavior, resource adoption).

**Key value over raw `helm install`:**

- **Declarative and reviewable**: the release is a manifest — versioned, diffed, reproducible
- **Version pinning is enforced**: `version` is required; an unpinned "latest" install is not reproducible, and reproducibility is the point of declaring the release
- **Contradictions rejected before deploy**: `atomic` + `skip_await` and `reuse_values` + `reset_values` fail validation with messages that say which flag to drop
- **Secrets marked at the source**: `set_sensitive` and `repository_password` are sensitive fields, kept out of rendered plans and state where each engine supports it
- **Dual IaC support**: Pulumi and Terraform modules that merge values identically and install byte-identical releases (one documented exception, below)

## The Values Model

The four layers are Helm's own, applied in this order on both engines — later layers win:

1. **`values_yaml`** — a YAML document, the equivalent of a Helm values file passed with `-f`. Full expressiveness: nested maps, lists, numbers, booleans. Do not put secrets here.
2. **`set`** — dotted-path overrides (e.g. `image.tag`, `ingress.hosts[0].host`) with Helm's `--set` coercion: `"true"`/`"false"` become booleans, digits become numbers, `"null"` removes the key.
3. **`set_string`** — same paths, but values always stay literal strings. Use this for version-like tags: `image.tag: "1.30"` through `set` would arrive as the number 1.3.
4. **`set_sensitive`** — literal strings like `set_string`, marked secret, highest precedence. For credentials and other sensitive chart values.

Both engines run Helm's own `--set` parser on the override entries and hand Helm one final merged map, so the same manifest installs byte-identical releases on either engine.

## Essential Configuration Fields

### Required

- **`spec.namespace`**: The namespace to install into — a literal name (`namespace: {value: my-ns}`) or a reference to a KubernetesNamespace resource
- **`spec.repo`**: The chart repository — an HTTP(S) Helm repository URL (`https://stefanprodan.github.io/podinfo`) or an OCI registry reference (`oci://ghcr.io/stefanprodan/charts`). For OCI, the chart is pulled as `<repo>/<chart>:<version>`
- **`spec.chart`**: The chart name within the repository (e.g. `podinfo`, `cert-manager`)
- **`spec.version`**: The exact chart version (e.g. `6.9.2`). Required — unpinned installs are not reproducible

### Common

- **`spec.create_namespace`**: When true, the namespace is created (with standard Planton governance labels) before the install and deleted with the resource; when false it must already exist
- **`spec.release_name`**: Overrides the Helm release name (defaults to `metadata.name`); lowercase alphanumerics and hyphens, at most 53 characters (Helm's limit)
- **`spec.values_yaml` / `spec.set` / `spec.set_string` / `spec.set_sensitive`**: The values layers (see above)
- **`spec.repository_username` / `spec.repository_password`**: Credentials for a private chart repository (HTTP basic auth or OCI registry login). Must be set together — the spec rejects one without the other

### Lifecycle

- **`spec.atomic`**: A failed install or upgrade rolls everything back and purges new resources — never a half-deployed release. Implies waiting for readiness, so it cannot be combined with `skip_await`
- **`spec.cleanup_on_fail`**: A failed upgrade deletes the resources that upgrade newly created (lighter than `atomic`)
- **`spec.skip_await`**: Return as soon as Helm records the release, without waiting for resource readiness. Both engines default to waiting
- **`spec.wait_for_jobs`**: When awaiting, also wait for chart-created Jobs to complete
- **`spec.timeout_seconds`**: Per-operation timeout (Helm's default 300)
- **`spec.skip_crds`**: Skip installing CRDs from the chart's `crds/` directory
- **`spec.dependency_update`**: Run `helm dependency update` before installing (rarely needed — most published charts bundle dependencies)
- **`spec.max_history`**: Release revisions kept for rollback (Helm's default 10; 0 = unlimited)
- **`spec.replace`**: Re-use a release name left behind by a failed/uninstalled release — Helm marks this unsafe in production
- **`spec.force_update`**: Force delete-and-recreate when a field cannot be patched in place — disruptive
- **`spec.reuse_values` / `spec.reset_values`**: Upgrade-time values handling (keep the last release's values vs. reset to chart defaults). Mutually exclusive — the spec rejects both together
- **`spec.disable_webhooks`**: Chart hooks do not run
- **`spec.disable_openapi_validation`**: Skip schema validation of rendered manifests against the cluster
- **`spec.take_ownership`**: Adopt existing resources into the release (Helm `--take-ownership`) — the migration knob. **Terraform provisioner only for now**; the Pulumi engine's pinned SDK predates the flag and rejects it loudly rather than silently ignoring it
- **`spec.description`**: Free-form note stored with the release (visible in `helm status`)

## Stack Outputs

After deployment, the following outputs are available in `status.outputs` — what `helm list` and `helm status` would show, exported identically by both engines:

- **`namespace`**: The namespace the release is installed in
- **`release_name`**: The Helm release name (`helm list` NAME column)
- **`version`**: The installed chart version
- **`app_version`**: The chart's appVersion — the upstream application version the chart packages
- **`status`**: The release status as Helm records it (e.g. `deployed`)
- **`revision`**: The release revision number (1 on install, incremented by upgrades/rollbacks)

## Quick Start

Create a file `helm-release.yaml` — podinfo from its HTTPS repository:

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesHelmRelease
metadata:
  name: podinfo
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

Then verify with Helm itself — the release is real:

```shell
helm list -n podinfo
helm status podinfo -n podinfo
```

## How It Works

This component includes both **Pulumi** (Go) and **Terraform** (HCL) modules that:

1. Optionally create the release namespace as an explicit, labeled module-owned resource (never via Helm's own create-namespace flag, which would create it unlabeled)
2. Resolve the release name (`release_name` or `metadata.name`) and the Helm defaults for `timeout_seconds` (300) and `max_history` (10), so both engines send identical values whether or not the spec set them
3. Resolve the chart reference: OCI repositories join `<repo>/<chart>`; HTTPS repositories pass the repo URL alongside the bare chart name. Private-repo credentials ride the repository options in both forms
4. Merge the values layers with the documented precedence — the Pulumi module merges module-side using Helm's own `strvals` parser; the Terraform module uses `helm_release`'s native `values` + `set`/`set_sensitive` attributes, which the provider merges in exactly the same order
5. Install the chart as a real Helm release (`helm.v3.Release` / `helm_release`) with every lifecycle knob mapped 1:1 between the engines
6. Export the release's observable handles from what Helm actually recorded — not echoed from the spec

**One parity exception**: `take_ownership` is honored by the Terraform provisioner only. The Pulumi engine's pinned pulumi-kubernetes SDK predates the flag; rather than silently dropping a set field, the Pulumi module fails the deploy with an error that routes you to the Terraform provisioner (or to dropping the flag).

**One behavioral difference (not a parity exception)**: Terraform masks each `set_sensitive` entry individually in plans and state; Pulumi marks the whole merged values map secret when any `set_sensitive` entry is present — coarser, but safe, and the installed release is identical.

## When to Use

Use **KubernetesHelmRelease** when — and only when — the catalog has no first-class component for the chart you need:

- Third-party or vendor charts outside the catalog's coverage
- Internal charts your organization publishes to a private repository or OCI registry
- Migrating Helm-managed workloads into declarative management before (or instead of) modeling them as typed components

**Do NOT use** when:

- A first-class catalog component covers the workload — it always wins (validation before deploy, composable outputs, documented trade-offs)
- You have raw manifests rather than a chart — that is KubernetesManifest's job
- You want to render a chart without creating a release — this component deliberately does not do that; hooks and Helm tooling are the point

## Prerequisites

- **Kubernetes Cluster**: Access to a Kubernetes cluster
- **Credentials**: Kubernetes cluster credentials (kubeconfig)
- **Chart availability**: The repository or registry must be reachable from where the deploy runs; private repos need `repository_username` / `repository_password`

## Best Practices

1. **Pin versions and treat bumps as reviewed changes**: the spec forces the pin; the chart's release notes are part of the review — a chart upgrade is an application upgrade
2. **Structure values deliberately**: non-secret configuration in `values_yaml` where reviewers see structure; `set`/`set_string` for the targeted per-environment overrides; secrets only in `set_sensitive`
3. **Use `set_string` for anything version-like**: `"1.30"` through `set` becomes the number 1.3 — the classic Helm coercion incident
4. **Enable `atomic` and `cleanup_on_fail` in production**: a failed upgrade that rolls back completely is an incident; one that strands half a release is an outage
5. **Reach for Helm tooling when debugging**: `helm status`, `helm history`, and `helm rollback` all work on releases this component installs — that is what a real release buys you

## References

- [Helm Documentation](https://helm.sh/docs/)
- [Helm Values Files and `--set`](https://helm.sh/docs/chart_template_guide/values_files/)
- [Helm OCI Registry Support](https://helm.sh/docs/topics/registries/)
- [Helm Chart Hooks](https://helm.sh/docs/topics/charts_hooks/)
- [Terraform helm_release Resource](https://registry.terraform.io/providers/hashicorp/helm/latest/docs/resources/release)
- [Pulumi Kubernetes helm.v3.Release](https://www.pulumi.com/registry/packages/kubernetes/api-docs/helm/v3/release/)

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
