# Kubernetes Tekton Operator — design notes

## Grain

One resource = the Tekton Operator, installed from the official
single-file release manifest at the pinned tag (the in-repo Helm chart
is version "devel" and unpublished — not a distribution). Exactly one
install per cluster; the namespace is the manifest's fixed
`tekton-operator`. The manifest applies per document, keyed by each
document's composed identity (`apiVersion//kind//name[//namespace]`) —
the multi-GVK release-manifest bundle class.

## The single-owner patch

The release ships `AUTOINSTALL_COMPONENTS: "true"` in its
`tekton-config-defaults` ConfigMap, which makes the operator create a
default `TektonConfig` named `config` at startup. Both modules ALWAYS
patch it to `"false"`: the `KubernetesTekton` declaration renders
exactly that object, and two managers co-writing it fight through
server-side apply field ownership. The patch is a design invariant,
not a knob — the E2E verifier asserts NO TektonConfig exists after
installing the operator alone.

## Typed overrides

Image overrides for the two Deployments (the operator Deployment's two
containers share one image upstream), per-Deployment resources, and
pod scheduling — patched onto the manifest's own documents; every
other document applies verbatim (faithful distribution). Pull-secret
names join the image overrides' own entries, deduplicated.

## The destroy ordering contract

Every document deletes with the resource, INCLUDING the 14
`operator.tekton.dev` CRDs — destroying the operator cascade-deletes
any TektonConfig. The safe order is structural: `KubernetesTekton`
destroys first (its deletion blocks on the operator processing the
InstallerSet finalizers), then this operator. The registry prerequisite
edge encodes the forward ordering; the reverse teardown rides it.

## No cert-manager prerequisite

Unlike fail-closed-webhook operators, the Tekton operator registers its
webhook configurations at runtime (knative-style) and manages its own
serving cert Secret — the bundle ships no
Mutating/ValidatingWebhookConfiguration documents.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
