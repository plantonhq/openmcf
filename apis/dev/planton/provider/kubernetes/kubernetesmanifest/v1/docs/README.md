# Kubernetes Manifest: Research Documentation

## Introduction

Start with when NOT to use this component, because that decision shapes everything else: **a first-class catalog component always wins over raw YAML.** Typed components validate configuration before anything reaches a cluster, export composable outputs that other resources can reference through the deployment graph, and document their trade-offs field by field. Raw YAML does none of that — the API server is the first thing that checks it, and its outputs are whatever you go look up afterwards. **KubernetesManifest** is the catalog's escape hatch, deliberately last in line: reach for it only when no component covers what you need to apply — a vendor's install manifest, a CRD bundle, an exotic custom resource.

What the escape hatch does, it does with a strict contract. Hand it any valid Kubernetes YAML — a single document or many separated by `---`, core kinds or custom resources, even a CRD and its custom resources in the same manifest — and both IaC engines apply it to the cluster exactly as written. No injected labels, no rewritten fields, no interpretation. The only defaulting that ever happens is the namespace anchor, and it is precisely bounded:

- Documents that declare their own `metadata.namespace` **keep it**.
- Namespaced documents that declare none **land in `spec.namespace`** (the anchor).
- Cluster-scoped documents (CRDs, ClusterRoles, ...) are **applied as-is** — the anchor never distorts them.

That contract — apply-exactly-as-written plus one bounded namespace default — is the whole design. Everything below explains why it is shaped this way and how both engines uphold it identically.

## The Problem Space

### Raw manifests refuse to die

In a fully typed platform, every deployment would flow through a purpose-built component. Reality keeps producing YAML that fits no type:

1. **Vendor install manifests.** Many projects publish their installation as one raw YAML file for `kubectl apply -f` — hundreds of documents mixing CRDs, RBAC, Services, Deployments, and webhook configurations. Wrapping that in anything loses the property that matters most: the applied objects are byte-for-byte what the vendor tested.
2. **CRD bundles.** Operators arrive as CRDs plus the custom resources that configure them. The CRDs are cluster-scoped, the custom resources often namespaced, and the two have a hard ordering dependency.
3. **Exotic custom resources.** Once an operator is installed, its custom resources are the API — and no general-purpose catalog can type every operator's API.

### What `kubectl apply` cannot give you

The obvious tool for this YAML is `kubectl apply -f`, and its gaps are the reason an IaC wrapper exists at all:

- **No lifecycle.** Apply creates and updates, but nothing tracks what was applied, so deletion means remembering every object (or the fragile `--prune`). Renaming a resource in the file orphans the old object silently.
- **No CRD ordering.** A manifest containing a CRD and a custom resource of that type fails on first apply — the type does not exist yet when the custom resource is submitted. The folk remedy is running apply twice.
- **No state, no drift detection, no plan.** There is no way to see what an apply *would* change, and no record connecting cluster objects back to the file that produced them.

### The design tension

A raw-manifest component lives on a knife's edge: every convenience it adds (label injection, name rewriting, templating) breaks the byte-for-byte property that justifies its existence, while every convenience it refuses pushes work back onto the user. KubernetesManifest resolves the tension by adding exactly one convenience — the anchor namespace — and making it scope-aware so it can never distort a document that already knows where it belongs.

## The Semantics in Detail

### Namespace anchoring

`spec.namespace` is the anchor: **namespaced documents that declare no `metadata.namespace` are applied there.** Documents with an explicit namespace keep it, and cluster-scoped documents are untouched. This mirrors how `kubectl apply --namespace` behaves, so a manifest written for kubectl lands identically here.

The anchor is a `StringValueOrRef`: a literal namespace name, or a reference to a KubernetesNamespace resource so the namespace and the manifest compose in one deployment graph. With `create_namespace: true`, the component creates the anchor namespace (with the standard Planton governance labels) before the manifest applies and deletes it with the resource; with `false`, the namespace must already exist. Note the labels go on the *created namespace object only* — never into the manifest's own documents.

### The two engines, one contract

Both engines reach the same anchoring outcome by different mechanisms:

- **Pulumi**: the Kubernetes provider is constructed with `spec.namespace` as its default namespace. The provider resolves each kind's scope before defaulting, so cluster-scoped kinds are skipped client-side.
- **Terraform**: each document becomes a `kubectl_manifest` resource, and `override_namespace` is set only on documents that declare no `metadata.namespace`. Cluster-scoped documents without a namespace also receive the override, but the API server ignores `metadata.namespace` on cluster-scoped objects — so the outcome matches the Pulumi provider's scope-aware defaulting.

### CRD ordering

A CRD and its custom resources can ship in one manifest, and both engines apply them in one pass:

- **Pulumi**: the `yaml/v2` ConfigGroup is CRD-aware — it registers CRDs before creating the custom resources that use them.
- **Terraform**: `kubectl_manifest` applies server-side and needs no cluster connection at plan time, so the plan does not fail on a type that does not exist yet, and CRD+CR manifests apply in a single run.

### Await behavior (`skip_await`)

By default (`skip_await: false`), **both engines block until readiness**: Deployments, DaemonSets, and StatefulSets complete their rollout, and other kinds pass their engine's readiness checks. With `skip_await: true`, the deploy returns as soon as the API server accepts every document.

When to skip: manifests whose readiness depends on something deployed later (the classic case is a webhook configuration waiting on its service), or resources that intentionally stay not-ready at install time. Vendor install bundles are the usual customer — awaiting readiness on a bundle with internal ordering assumptions can deadlock the deploy.

Engine mechanics: Pulumi sets `SkipAwait` on the ConfigGroup; Terraform inverts the flag into `wait` and `wait_for_rollout` on each `kubectl_manifest`. One benign breadth difference exists: when awaiting, the Pulumi engine also readiness-checks non-workload kinds (Services and the like), while kubectl awaits workload rollouts only. The applied objects are identical either way.

### The applied-resource inventory

The `applied_resources` stack output lists one `"<apiVersion>/<Kind>/<name>"` entry per document, in manifest order. Both engines derive it by **parsing the input YAML** — never by reflecting over engine-side child resources — so the inventory is identical by construction, and downstream tooling can see what a manifest contains without re-parsing it.

The document-split rule is likewise identical on both engines: documents are separated on lines starting with `---`, a newline is prepended so a manifest that *starts* with `---` does not lose its first document, blank and comment-only chunks are dropped, and an invalid document fails loudly (at plan on Terraform, at apply on Pulumi) rather than being skipped.

### State identity under document reordering

On the Terraform engine, each document's state address is keyed by its full identity (`apiVersion/Kind/namespace/name`), not its position in the file. Reordering documents in the manifest therefore never churns state; two documents with the same identity collide loudly at plan time, because a duplicate document is a manifest bug, not something to apply twice.

## Deployment Methods Landscape

### Level 0: `kubectl apply -f`

The baseline everyone knows. Immediate, universal — and stateless: no deletion tracking, no drift detection, no CRD ordering (apply twice), no plan.

**Verdict:** Fine for a terminal session; not a lifecycle.

### Level 1: kubectl + GitOps (Argo CD / Flux)

Raw manifests in Git, reconciled continuously. Solves drift and deletion tracking; adds sync-wave annotations for ordering.

**Pros:** continuous reconciliation, audit trail, prune-on-delete.

**Cons:** a separate control plane to operate; ordering is annotation-driven; the manifest's lifecycle is decoupled from the rest of the infrastructure graph (the namespace it needs, the cluster it lands on).

**Verdict:** Strong for teams already invested; heavy as an entry price for "apply this one vendor file."

### Level 2: Terraform with raw manifests

The hashicorp provider's `kubernetes_manifest` gives per-object IaC lifecycle, but it needs a live cluster connection *at plan time* (to fetch schemas), which breaks plan-before-cluster workflows and makes CRD+CR-in-one-apply painful. The community `kubectl_manifest` (alekc/kubectl) fixes both: no plan-time connection, server-side apply, one resource per document.

**Verdict:** Production-grade lifecycle; the provider choice matters more than it should, and multi-document YAML needs hand-rolled splitting.

### Level 3: Pulumi `yaml/v2`

`ConfigGroup` accepts whole multi-document manifests, orders CRDs before their custom resources, and awaits readiness per kind. The closest upstream analogue to what this component does.

**Verdict:** Excellent mechanism — still needs the namespace-anchoring, output, and lifecycle conventions built around it.

### Comparative Analysis

| Aspect | kubectl | GitOps | Terraform (kubectl provider) | Pulumi yaml/v2 | Planton |
|--------|---------|--------|------------------------------|----------------|---------|
| Deletion tracking | No (`--prune` fragile) | Yes | Yes | Yes | Yes |
| CRD + CR in one apply | No (apply twice) | Via sync waves | Yes | Yes | Yes (both engines) |
| Namespace anchoring | `-n` flag, per-invocation | Per-app config | Hand-rolled | Provider default | First-class `spec.namespace` (value or ref) |
| Readiness await | `kubectl wait`, manual | Health checks | `wait`/`wait_for_rollout` | `SkipAwait` | One `skip_await` field, both engines |
| Applied-resource inventory | No | UI only | State inspection | State inspection | `applied_resources` output, engine-independent |
| Manifest mutated? | No | Sometimes (labels) | No | No | Never |

## The Planton Approach

### Boundary first

The component's own documentation leads with when not to use it, and the spec enforces nothing about manifest *content* beyond "at least one non-blank YAML document" — deliberately. Validating raw YAML would mean re-implementing the API server; the component's value is lifecycle and anchoring, not schema-checking someone else's kinds.

### A four-field spec

The entire surface is `namespace` (the anchor, value-or-ref), `create_namespace`, `manifest_yaml`, and `skip_await`. Every field maps to exactly one decision the user actually faces: where unanchored documents land, who owns the namespace lifecycle, what to apply, and whether to wait.

### Parity as a contract, not an aspiration

Both engines uphold the same externally observable behavior: identical anchoring outcomes (scope-aware provider defaulting vs. conditional `override_namespace`), identical apply mechanism (server-side apply on both), identical await defaults, identical outputs derived from the same input parsing with the same document-split rule. Where an internal mechanism differs (await breadth on non-workload kinds), it is documented as a behavioral note — the applied objects never differ.

### Outputs from the input, not the engine

Deriving `applied_resources` by parsing the input YAML rather than reflecting engine state is a deliberate inversion: it guarantees the two engines export identical inventories, and it makes the output meaningful *before* deploy (the inventory is a property of the manifest, not of the run).

## Implementation Landscape

### Pulumi Module Architecture

The Pulumi module (`iac/pulumi/module/`) follows the standard Planton pattern:

- **`main.go`**: builds the provider anchored to `spec.namespace`, conditionally creates the namespace, applies the manifest through a `yaml/v2` ConfigGroup with `SkipAwait` wired to the spec
- **`locals.go`**: resolves the anchor namespace, computes the governance labels for the created namespace, and parses the applied-resource inventory
- **`inventory.go`**: the document-split-and-parse rule that produces `applied_resources`
- **`namespace.go`**: the conditional namespace creation
- **`outputs.go`**: exports `namespace` and `applied_resources`

### Terraform Module Architecture

The Terraform module (`iac/tf/`) mirrors the contract:

- **`variables.tf`**: mirrors `spec.proto` (the anchor arrives flattened to a plain string)
- **`locals.tf`**: the same labels, the same document-split rule, identity-keyed document map, and the same inventory derivation
- **`main.tf`**: optional `kubernetes_namespace_v1`, then one `kubectl_manifest` per document with conditional `override_namespace`, server-side apply, and the await flags
- **`outputs.tf`**: exports the same two outputs

### Resource Count

One optional namespace plus exactly as many Kubernetes objects as the manifest declares. The component adds nothing of its own to the cluster.

## Production Best Practices

1. **Exhaust the catalog first.** Every document you move from raw YAML into a typed component gains validation, outputs, and documentation. Treat a KubernetesManifest resource as a to-do item pointing at a missing component.
2. **Leave first-party documents unanchored.** Omitting `metadata.namespace` in your own documents and anchoring through `spec.namespace` makes the bundle retargetable by changing one field.
3. **Paste vendor manifests verbatim.** The value of the pattern is byte-for-byte fidelity; upgrading means pasting the next release's file, and diffing releases stays meaningful.
4. **Use `skip_await: true` for install bundles with internal ordering.** Webhook configurations waiting on their services are the canonical await deadlock. For your own simple bundles, keep the default and let the deploy verify readiness.
5. **One concern per resource.** A manifest resource that mixes a vendor install with your own configuration couples their lifecycles; split them so each can change independently.
6. **Check `applied_resources` in reviews.** The inventory output is the quickest manifest-level diff: a document you didn't mean to add shows up as an extra entry.

## Conclusion

KubernetesManifest is deliberately the smallest possible component wrapped around the largest possible input: four spec fields, two outputs, and a strict apply-exactly-as-written contract with one bounded convenience (the anchor namespace) and one lifecycle knob (`skip_await`). Both engines uphold the same externally observable behavior — same anchoring, same one-pass CRD ordering, same await defaults, same input-derived inventory. It is the right tool precisely when nothing else in the catalog is, and the wrong tool the moment a typed component exists.

## References

- [Kubernetes Objects and Manifests](https://kubernetes.io/docs/concepts/overview/working-with-objects/)
- [Namespaces and Object Scope](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/)
- [Custom Resource Definitions](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)
- [Server-Side Apply](https://kubernetes.io/docs/reference/using-api/server-side-apply/)
- [Pulumi Kubernetes yaml/v2 ConfigGroup](https://www.pulumi.com/registry/packages/kubernetes/api-docs/yaml/v2/configgroup/)
- [alekc/kubectl Terraform Provider](https://registry.terraform.io/providers/alekc/kubectl/latest/docs)
- [kubectl apply and declarative management](https://kubernetes.io/docs/tasks/manage-kubernetes-objects/declarative-config/)
