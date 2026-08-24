# KubernetesManifest Guide

The judgment this guide carries: raw YAML is the catalog's bottom step —
below first-class kinds, below verified alternatives, below even a chart
install. Run the search-first workflow in the
[catalog guide](../../_docs/GUIDE.md) before declaring anything here.

## The two-hatch decision

Upstream ships a Helm chart? Use KubernetesHelmRelease — you get the
chart's hooks, values model, and release lifecycle. What you hold is
plain YAML documents (a vendor's install manifest, a CRD bundle, an
exotic custom resource) — that is what this kind is for, applied exactly
as written with no mutation (the namespace-anchoring rules are on
[reference.md](v1alpha1/reference.md)).

## What raw YAML costs the architecture

Nothing this kind applies is referenceable: no validated spec, no
exported outputs, no `valueFrom` edges in or out. Anything the rest of
the architecture must wire to — an endpoint, a credential — has to be
communicated outside the platform's reference system, which also means
deploy ordering against dependents is yours to manage.

## Namespace ownership

`spec.namespace` is the anchor namespace and a required foreign key
targeting KubernetesNamespace; `createNamespace: true` follows the same
ownership contract as every workload kind — owned in IaC state, deleted
with the resource. The judgment and wiring:
[namespace-ownership pattern](../../_patterns/namespace-ownership.md).

## On the diagram

One node for the whole manifest, however many documents it contains. A
multi-document install that matters to the architecture is worth
splitting into typed kinds where they exist, so the pieces become visible,
referenceable nodes.

## Pairs well with

- KubernetesNamespace — the anchor namespace's owner (pattern above).
- KubernetesHelmRelease — the chart-shaped sibling escape hatch (decision
  above).
