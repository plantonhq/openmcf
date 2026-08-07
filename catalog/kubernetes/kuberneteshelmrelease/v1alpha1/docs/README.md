# Kubernetes Helm Release: Research Documentation

## Introduction

Start with the boundary, because it defines the component: **if the catalog has a first-class component for what you are deploying, use it — not this.** Typed components validate their configuration before anything reaches a cluster, export composable outputs, and teach their trade-offs field by field. A generic chart install does none of that: the chart's own values surface is the configuration contract, and Kubernetes learns about a bad value only when the chart's rendered resources misbehave.

**KubernetesHelmRelease** is the catalog's sole intentional passthrough. It exists for exactly one situation: the chart you need has no catalog component. In that situation, it installs the upstream chart as a *real* Helm release — not a client-side render of the chart's templates. The distinction matters and is the component's central design decision, so the rest of this document keeps returning to it.

Three structural facts define the component:

- **The chart's values surface is the contract.** The spec cannot validate `replicaCount` or `ingress.hosts[0].host` — those keys belong to the chart. What the spec *can* validate, it does: chart identity is complete and pinned, override layers have defined precedence, and mutually exclusive lifecycle flags are rejected before deploy.
- **Both engines create a real Helm release.** Hooks run, the release secret is written to the namespace, history is recorded, and `helm list` / `helm status` see the release exactly as if the Helm CLI had installed it.
- **The values model is Helm's own.** A values file (`values_yaml`) plus `--set`-style overrides (`set`, `set_string`, `set_sensitive`), merged with fixed precedence identically on both engines.

## Evolution and Historical Context

### Helm 2, Tiller, and the move client-side

Helm 2 (2016) installed charts through Tiller, an in-cluster server component with broad permissions that quickly became the canonical Kubernetes security objection. Helm 3 (2019) removed Tiller entirely: the client renders and applies, and release state moved from ConfigMaps in Tiller's namespace to Secrets in the release's own namespace. That release secret is now the durable record of "a Helm release exists here" — the thing `helm list` enumerates, `helm rollback` consults, and any tool claiming to install "a Helm release" must produce.

### Release vs. render: the fault line in IaC Helm support

Every IaC tool that touches Helm faces the same fork:

- **Render-only**: run `helm template` semantics client-side, apply the resulting manifests as individual resources. The IaC engine tracks each resource, which gives fine-grained diffs — but no release secret is written, chart hooks never run, and `helm list` shows nothing. Charts that rely on hooks (database migrations, admission-webhook certificate bootstrapping, pre-delete cleanup) break silently.
- **Real release**: drive Helm's install/upgrade code path. Hooks run, history accrues, Helm tooling works — the release is a first-class Helm object that happens to be managed by IaC.

Pulumi ships both (`helm.v3.Chart` renders; `helm.v3.Release` installs); Terraform's `helm_release` has always been a real release. Which side of the fork a tool sits on is the single most consequential fact about it, and it is routinely discovered in production when a hook does not fire.

### OCI registries as chart distribution

Helm 3.8 (2022) made OCI registry support generally available, and chart publishing has been migrating from index-file HTTPS repositories to OCI registries (GHCR, ECR, ACR, Artifact Registry) since. The two forms differ mechanically — an HTTPS repo serves an `index.yaml` and the chart name is looked up in it; an OCI chart is pulled directly as `<registry-path>/<chart>:<version>` — but conceptually a chart consumer should not care. Any current chart-install surface must accept both.

### The `--set` coercion footgun

Helm's `--set` parser (the `strvals` package) coerces as it parses: `true` and `false` become booleans, digit strings become numbers, `null` deletes a key. This is convenient right up until someone sets `image.tag=1.30` and the tag arrives at the template as the float `1.3`. Helm's own answer is `--set-string` (parse the path, keep the value a literal string), and any faithful modeling of Helm's values surface must carry both forms — plus `--set`'s deliberate use of `null` to *remove* a default the chart ships with.

## The Semantics in Detail

### The four values layers

Helm merges values from lowest to highest precedence: the chart's built-in defaults, then values files, then `--set`-family flags. The component models this exactly:

1. **`values_yaml`** — a YAML document, the equivalent of a values file passed with `-f`. Full expressiveness: nested maps, lists, numbers, booleans. Applied first.
2. **`set`** — dotted-path overrides with Helm's `--set` coercion (`"true"` → bool, digits → number, `"null"` deletes the key). Paths support list indexing (`ingress.hosts[0].host`).
3. **`set_string`** — same paths, values always kept as literal strings. The footgun-avoider for version-like tags.
4. **`set_sensitive`** — literal strings like `set_string`, but marked secret: kept out of rendered plans and state where each engine supports it. Highest precedence.

Both engines merge these layers in exactly this order and hand Helm one final values map, so the same manifest installs byte-identical releases on either engine.

### Release identity and the release secret

The release name defaults to the resource's `metadata.name`, overridable with `release_name` (validated to Helm's 53-character release-name limit). The namespace holds the release secret; `create_namespace` decides whether the component creates (and labels, and eventually deletes) that namespace or expects it to exist.

### Lifecycle knobs

The spec exposes Helm's install/upgrade behavior surface: `atomic` (failed deploys roll back completely), `cleanup_on_fail` (failed upgrades delete newly created resources), `skip_await` (return once the release is recorded, without waiting for readiness), `wait_for_jobs`, `timeout_seconds`, `skip_crds`, `dependency_update`, `max_history`, `replace`, `force_update`, `reuse_values` / `reset_values` (upgrade-time values handling), `disable_webhooks` (chart hooks), `disable_openapi_validation`, `take_ownership` (adopt resources an earlier tool created), and a free-form `description`.

Two pairings are contradictions, and the spec rejects them at validation rather than letting Helm fail later: `atomic` with `skip_await` (atomic must wait to know whether to roll back) and `reuse_values` with `reset_values` (one keeps the previous release's values, the other discards them).

### Version pinning is mandatory

`version` is required. An unpinned install resolves "whatever the repo serves today," which makes the manifest non-reproducible — and reproducibility is the point of declaring the release as a manifest at all. This is the spec disagreeing, deliberately, with `helm install`'s own default behavior.

## Deployment Methods Landscape

### Level 0: Helm CLI

```shell
helm repo add podinfo https://stefanprodan.github.io/podinfo
helm install podinfo podinfo/podinfo --version 6.9.2 -n podinfo --create-namespace \
  --set replicaCount=2
```

**Pros:** the reference implementation; everything works.

**Cons:** imperative — the install exists only in the cluster and the shell history; no drift detection, no review, no reproducibility unless every flag is scripted and versioned by hand.

**Verdict:** right for exploration; the thing declarative management replaces.

### Level 1: GitOps Helm controllers (Flux HelmRelease, Argo CD)

Flux's HelmRelease CRD drives real Helm installs from Git. Argo CD, notably, sits on the *render* side of the fork: it templates charts and applies the output, with hook emulation that differs from Helm's semantics.

**Pros:** declarative, continuous reconciliation, Git as the source of truth.

**Cons:** requires operating the GitOps control plane itself; secrets management and multi-cluster credentials become that platform's problem.

**Verdict:** excellent where a GitOps platform is already the deployment substrate; a heavy prerequisite where it is not.

### Level 2: Terraform

```hcl
resource "helm_release" "podinfo" {
  name       = "podinfo"
  repository = "https://stefanprodan.github.io/podinfo"
  chart      = "podinfo"
  version    = "6.9.2"
  namespace  = "podinfo"

  values = [file("values.yaml")]

  set = [{ name = "replicaCount", value = "2" }]
}
```

**Pros:** a real Helm release with full IaC lifecycle; `set_sensitive` masks individual secret values in plan output; provider 3.x supports `--take-ownership` for adopting existing resources.

**Cons:** values are strings-in-HCL; the chart's contract is still unvalidated until Helm renders.

**Verdict:** production-grade; the semantics this component's Terraform module wraps directly.

### Level 3: Pulumi

```go
helmv3.NewRelease(ctx, "podinfo", &helmv3.ReleaseArgs{
    Chart:   pulumi.String("podinfo"),
    Version: pulumi.String("6.9.2"),
    RepositoryOpts: &helmv3.RepositoryOptsArgs{
        Repo: pulumi.String("https://stefanprodan.github.io/podinfo"),
    },
    Values: pulumi.Map{"replicaCount": pulumi.Int(2)},
})
```

**Pros:** a real Helm release (`helm.v3.Release`, not the render-only `helm.v3.Chart`); values as native language maps; secrets machinery for sensitive values.

**Cons:** same unvalidated chart contract; the Release resource's diff is coarser than per-resource tracking.

**Verdict:** production-grade; the semantics this component's Pulumi module wraps directly.

### Comparative Analysis

| Aspect | Helm CLI | Flux/Argo | Terraform | Pulumi | Planton |
|--------|----------|-----------|-----------|--------|---------|
| Real Helm release (hooks, secret, `helm list`) | Yes | Flux yes; Argo renders | Yes | Release yes; Chart renders | Yes, both engines |
| Declarative + reviewable | No | Yes | Yes | Yes | Yes |
| Version pinning enforced | No | No | No | No | Yes (required field) |
| Contradictory-flag checks before deploy | No | No | No | No | Yes (CEL) |
| set / set-string / set-sensitive layers | Flags | Partial | Yes | Manual | Yes, typed, fixed precedence |
| Engine choice per deploy | N/A | N/A | TF only | Pulumi only | Both, identical result |

## The Planton Approach

### The boundary is the design

The component's own documentation leads with when *not* to use it, and that is deliberate. A first-class catalog component always wins where one exists: it validates before deploy and composes through typed outputs. KubernetesHelmRelease is scoped to the remainder — the long tail of charts no component covers — and it makes that remainder as safe as a generic passthrough can be: identity pinned, precedence fixed, contradictions rejected, secrets marked.

### A real release on both engines

Both modules install through the real-release primitive — `helm_release` on Terraform, `helm.v3.Release` on Pulumi. The render-only `helm.v3.Chart` was deliberately rejected: it silently skips hooks and leaves nothing for Helm tooling to manage. After a deploy, `helm list -n <namespace>` shows the release, `helm status` shows its state and description, and `helm rollback` works, whichever engine installed it.

### One merge, two engines

The values layers merge with the documented precedence — `values_yaml`, then `set`, then `set_string`, then `set_sensitive` — on both engines, through different mechanics that reach the same result:

- The **Pulumi module** merges module-side: it parses `values_yaml`, then applies each override entry with Helm's own `strvals` parser (`ParseInto` for `set`, `ParseIntoString` for the string layers), in sorted-key order, and hands the Release one final map.
- The **Terraform module** rides the provider's native mechanisms: `values = [values_yaml]`, `set` entries as type-`auto` attributes, `set_string` entries as type-`string`, `set_sensitive` as its own attribute — which the provider merges in exactly the same order, also via Helm's `strvals`.

Because both paths run Helm's own parser on override entries, dotted-path syntax and type coercion behave identically. Map iteration is lexical on both sides, keeping even same-path collisions deterministic.

### Namespace ownership

`create_namespace: true` creates the namespace as an explicit module-owned resource stamped with the standard Planton governance labels — not via `helm_release`'s own create-namespace flag, which would create an unlabeled namespace. Both engines stamp identical labels; the release itself carries a dependency on the namespace so ordering is guaranteed.

### Secrets, two grains

`repository_password` and `set_sensitive` are marked sensitive in the spec. The Terraform provider masks `set_sensitive` entries individually in plans and state. The Pulumi module cannot mask individual keys inside one merged values map, so when any `set_sensitive` entry is present it marks the *whole* merged values map secret in state — coarser than Terraform's per-entry masking, but it errs on the safe side. This is a documented behavioral difference, not a parity exception: the release Helm installs is identical.

### One parity exception, failed loudly

`take_ownership` (Helm's `--take-ownership`: skip existing-resource conflict checks and adopt matching resources into the release — the migration knob for pointing a release at resources an earlier tool created) is currently honored by the **Terraform provisioner only**. The Pulumi engine's pinned pulumi-kubernetes SDK predates the flag, and a set field must never be silently dropped — so the Pulumi module rejects a spec with `take_ownership: true` loudly at deploy, with an error that routes the user to the Terraform provisioner or to dropping the flag. The exception dissolves at the next pulumi-kubernetes SDK upgrade.

### Defaults resolved module-side

The optional knobs with Helm defaults — `timeout_seconds` (300) and `max_history` (10) — are resolved in each module's locals so both engines send identical values whether or not the spec set the fields. The deployed release never depends on which engine applied it.

## Implementation Landscape

### Pulumi Module Architecture

The Pulumi module (`iac/pulumi/module/`) follows the standard Planton pattern:

- **`main.go`**: Orchestrates provider init, the optional namespace, the `helm.v3.Release` (with the `take_ownership` rejection and the chart-reference resolution: OCI repos join `<repo>/<chart>`; HTTPS repos ride repository opts), and output export
- **`values.go`**: The module-side values merge — `values_yaml` parsed, then `set` / `set_string` / `set_sensitive` applied via Helm's `strvals` in sorted-key order
- **`locals.go`**: Release-name resolution, namespace resolution, governance labels, and the Helm defaults for `timeout_seconds` / `max_history`
- **`namespace.go`**: The conditional, labeled namespace
- **`outputs.go`**: Exports namespace, release name, chart version, appVersion, status, and revision — read from the Release's recorded status, not echoed from the spec

### Terraform Module Architecture

The Terraform module (`iac/tf/`) mirrors the Pulumi logic:

- **`variables.tf`**: Mirrors `spec.proto` fields as Terraform variables
- **`locals.tf`**: The same release-name / namespace / labels / defaults resolution
- **`main.tf`**: The conditional, labeled `kubernetes_namespace_v1`
- **`helm_release.tf`**: The `helm_release` resource — values list plus `set` / `set_sensitive` attributes, every lifecycle knob mapped 1:1 with the Pulumi module's arguments, including `take_ownership`
- **`outputs.tf`**: The same six outputs, read from the release's recorded metadata

### Resource Count

The component creates at most **two Kubernetes-level resources**: the optional namespace and the Helm release itself. Everything the chart creates belongs to Helm — the modules deliberately never reach into the chart's own resources (no label injection, no per-resource management).

## Production Best Practices

### Chart identity

1. **Pin the version — the spec makes you.** Then treat version bumps as reviewed changes: read the chart's release notes; a chart upgrade is an application upgrade
2. **Prefer the chart's canonical repository.** Mirrors and re-publishers lag; for OCI, the reference is `<repo>/<chart>:<version>` — verify it with `helm show chart` before committing the manifest

### Values discipline

1. **Keep non-secret configuration in `values_yaml`** where reviewers see structure, and reserve `set` / `set_string` for the handful of targeted overrides that vary per environment
2. **Use `set_string` for anything version-like.** `image.tag: "1.30"` through `set` arrives as the number 1.3 — the classic Helm coercion incident
3. **Secrets go in `set_sensitive`, never `values_yaml`.** The values document is plain text in the manifest and in state; the sensitive layer exists so secrets are marked at the source

### Lifecycle

1. **`atomic: true` and `cleanup_on_fail: true` for production releases** — a failed upgrade that rolls back completely is an incident; one that strands half a release is an outage
2. **Leave awaiting on** (the default) unless the deploy is deliberately fire-and-forget; `skip_await` trades readiness knowledge for speed and cannot be combined with `atomic`
3. **`take_ownership` is a migration tool, not a steady state** — use it to adopt resources an earlier tool created, on the Terraform provisioner, then drop the flag

### Operational awareness

1. **The release is a real Helm release** — `helm list`, `helm status`, `helm history`, `helm rollback` all work. Use `helm status` first when a deploy misbehaves; the `description` field surfaces there
2. **Chart errors surface at render or at readiness, not at validation.** The spec validates its own contract (identity, precedence, flag sanity); a typo'd chart value is between you and the chart

## Conclusion

KubernetesHelmRelease is a deliberately bounded component: the catalog's one intentional passthrough, for charts no first-class component covers, and never the recommended path where one exists. Within that boundary it is uncompromising — a real Helm release on both engines (hooks, history, `helm list`), Helm's exact four-layer values model with fixed precedence and identical merge behavior everywhere, required version pinning, contradictory flags rejected before deploy, secrets marked at the source, and the single engine gap (`take_ownership` on Pulumi's pinned SDK) converted from silent omission into a loud, routed failure.

## References

- [Helm Documentation](https://helm.sh/docs/)
- [Helm Values Files and `--set`](https://helm.sh/docs/chart_template_guide/values_files/)
- [Helm `--set` Value Parsing (strvals)](https://helm.sh/docs/intro/using_helm/#the-format-and-limitations-of---set)
- [Helm OCI Registry Support](https://helm.sh/docs/topics/registries/)
- [Helm Chart Hooks](https://helm.sh/docs/topics/charts_hooks/)
- [Terraform helm_release Resource](https://registry.terraform.io/providers/hashicorp/helm/latest/docs/resources/release)
- [Pulumi Kubernetes helm.v3.Release](https://www.pulumi.com/registry/packages/kubernetes/api-docs/helm/v3/release/)
- [Pulumi: Chart vs. Release](https://www.pulumi.com/registry/packages/kubernetes/how-to-guides/choosing-the-right-helm-resource-for-your-use-case/)
- [podinfo (the presets' demo chart)](https://github.com/stefanprodan/podinfo)
