# KubernetesKyverno Guide

The judgment this guide carries: Kyverno's superpower — policies as plain
Kubernetes resources that can also mutate and generate — comes with a
lifecycle quirk worth knowing before you operate it: its admission
webhooks are registered at runtime, and force-deleting the release can
strand them.

## When Kyverno over Gatekeeper

The comparison home is the
[Gatekeeper guide](../../kubernetesgatekeeper/v1alpha1/GUIDE.md); the short
form: pick Kyverno for policies-as-plain-YAML (no Rego to learn) and for
capabilities beyond validation — mutating resources on admission,
GENERATING companion resources (a default NetworkPolicy per new
namespace is the classic), and scheduled cleanup. Run one engine, not
both.

## The runtime-webhook lifecycle

Unlike Gatekeeper, the chart templates NO webhook configurations — the
admission controller registers and updates them at runtime per installed
policies. The modules add cleanup beyond the chart's own (the chart's
delete-webhooks helper at the pinned release deletes the wrong API), so
a normal destroy is clean — but a FORCE-deleted release can leave
`kyverno-*` webhook configurations behind, which then block admissions
for a controller that no longer exists (the reference page carries the
recovery detail). Tear down through the platform, not by deleting the
namespace by hand.

## Engine here, policies beside it

This installs four controllers + CRDs. ClusterPolicy/Policy resources
are authored separately (KubernetesManifest or GitOps) — the
engine-plus-declarations shape. An installed Kyverno with no policies
does nothing.

## Once per cluster; the infra exception

Shared chart, own namespace
([namespace-ownership pattern](../../../_patterns/namespace-ownership.md)
sole-tenant case;
[operator-prerequisite pattern](../../../_patterns/operator-prerequisite.md)
for the CRD seam).

## Pairs well with

- KubernetesManifest — carries the ClusterPolicy/Policy resources today.
- KubernetesGatekeeper — the alternative lane (comparison in its guide).
