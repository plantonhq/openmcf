# KubernetesOpenSearch Guide

The judgment this guide carries: when a user asks for Elasticsearch —
or "Elasticsearch plus Kibana" — this one kind is the catalog's answer,
and the two things a proposal must get right are the operator
prerequisite and not shipping the demo credentials past development.

Substitutes for: Elasticsearch (Apache-2.0 fork, drop-in for the
Elasticsearch APIs at the 7.10 fork line; existing clients connect
unchanged). Dashboards — the Kibana-role console — is a section of this
same spec, not a separate component.

## Answering "give me Elasticsearch"

Propose this kind and say what you did: OpenSearch is the catalog's
Elasticsearch-compatible engine, existing clients connect unchanged, and
enabling `dashboards` covers the Kibana half of the ask without another
component. The substitution workflow is the
[catalog guide](../../_docs/GUIDE.md)'s first law.

## The operator, and the credentials trap

KubernetesOpenSearchOperator is the registry prerequisite — include it in
the shared-cluster chart, once. Then read the security section of
[reference.md](v1alpha1/reference.md) before proposing anything user-facing: the
bootstrapped admin Secret carries the image's well-known demo
credentials, and the API serves the demo security config even with no
`security` block declared. For anything beyond a private dev cluster,
bring your own security config or rotate the admin password immediately
— the reference page carries both paths.

## Namespace ownership

`spec.namespace` is a required foreign key targeting KubernetesNamespace.
A search cluster commonly shares its namespace with the pipeline that
feeds it — the multi-tenant case where `createNamespace: true` is wrong;
[namespace-ownership pattern](../../_patterns/namespace-ownership.md).

## On the diagram

Cluster (with Dashboards inside it) and operator render as separate
nodes — the operator in the shared-cluster layer, the cluster in its
application environment. Like every operator-backed kind, no edge
connects them (the
[operator-prerequisite pattern](../../_patterns/operator-prerequisite.md));
verify the operator node exists.

## Pairs well with

- KubernetesOpenSearchOperator — required, once per cluster, shared
  chart.
- KubernetesNamespace — the namespace owner (pattern above).
- Exposure kinds (KubernetesIngress, Gateway API routes) when Dashboards
  or the API must be reachable from outside — composed, never embedded.
