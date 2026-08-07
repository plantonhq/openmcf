# KubernetesListenerSet -- Research Documentation

Comprehensive background on the Kubernetes Gateway API `ListenerSet` resource,
the upstream specification this component mirrors, and the reasoning behind the
Planton `KubernetesListenerSet` component. It is the source-of-truth companion
to the user-facing `README.md` (getting started) and `catalog-page.md` (catalog
listing).

## Table of Contents

1. Introduction
2. The Shared-Gateway Delegation Problem
3. Where ListenerSet Sits
4. Anatomy of ListenerSetSpec
5. Listener Merging Semantics
6. Trust Boundaries: ReferenceGrants and Secrets
7. Validation Rules
8. Why ListenerSet Is a First-Class Planton Component
9. Design Notes
10. Composing in Infra Charts
11. Controller Landscape
12. Common Pitfalls
13. Conclusion
14. References

## Introduction

The Gateway API is the successor to the Kubernetes Ingress API: a standards-based,
role-oriented, expressive specification for north-south (and increasingly
east-west) traffic. `ListenerSet` extends the model's delegation story: where a
`Gateway` owns the traffic entry point and `*Route` objects delegate routing to
application teams, a ListenerSet delegates *listener* ownership -- ports,
hostnames, TLS certificates, and route-attachment policy -- to teams who share a
central Gateway but must not edit it.

`KubernetesListenerSet` brings that resource into Planton as a first-class
deployment component modeling the complete upstream standard-channel surface
(ListenerSet is standard since Gateway API v1.5, served as
`gateway.networking.k8s.io/v1`), so customers never have to fall back to a raw
`KubernetesManifest` to express delegated listeners.

## The Shared-Gateway Delegation Problem

A multi-team cluster typically runs one expensive traffic entry point (a cloud
load balancer programmed by one Gateway) shared by many teams. Without
ListenerSet, every new team hostname or certificate means editing the central
Gateway -- a single, contended, high-blast-radius object owned by the platform
team. ListenerSet splits that object: the platform team keeps the Gateway
(addresses, infrastructure, baseline listeners) and each team maintains its own
listeners in its own namespace, holding its own TLS Secrets. Attachment is
opt-in and auditable:

- The Gateway must declare `allowed_listeners` (by default **no** ListenerSet
  may attach).
- The ListenerSet names its parent through `parent_ref` and must live in a
  namespace the Gateway trusts (`allowed_listeners.namespaces`: `Same`, `All`,
  or a label `Selector`).

## Where ListenerSet Sits

```
KubernetesGatewayApiCrds        (install CRDs, v1.5.0+)
        |
KubernetesGatewayClass          (controller class)
        |
KubernetesGateway               (central entry point; allowed_listeners opt-in)
        |  parent_ref
KubernetesListenerSet           (this component: per-team listeners + TLS)
        |  parentRefs (kind: ListenerSet)
HTTPRoute / GRPCRoute / ...     (team routes attach to the merged listeners)
```

## Anatomy of ListenerSetSpec

The Planton spec flattens the upstream `ListenerSetSpec` after the standard
namespaced envelope (`namespace`):

| Planton field | Upstream | Notes |
|---------------|----------|-------|
| `parent_ref` | `ListenerSetSpec.parentRef` | The Gateway attached to; a whole-Gateway reference (no `sectionName`/`port`, unlike the route-side `ParentReference`). `name` is a foreign key defaulting to `KubernetesGateway`. |
| `listeners` | `ListenerSetSpec.listeners` | 1-64 listener entries, the same building blocks as Gateway listeners. |

Each listener entry (`KubernetesListenerSetListener`, upstream `ListenerEntry`)
carries `name`, `hostname`, `port`, `protocol`, `tls`, and `allowed_routes` --
the exact shape of a Gateway listener. The `tls` and `allowed_routes` fields
reuse the shared Gateway API messages (`KubernetesGatewayApiListenerTlsConfig`,
`KubernetesGatewayApiAllowedRoutes`) that `KubernetesGateway` uses, because
upstream uses one type for both kinds -- sharing the messages makes drift
between the two components structurally impossible.

TLS `certificate_refs[].name` is a `StringValueOrRef` foreign key defaulting to
`KubernetesSecret` -- the Secret is typically produced by a
`KubernetesCertificate` (the cert-manager seam), so `valueFrom` against the
Certificate's `status.outputs.secret_name` wires the listener to the issued
certificate with a real dependency edge.

## Listener Merging Semantics

The parent Gateway merges listeners from itself and all attached ListenerSets:

- **Precedence:** the Gateway's own listeners first, then ListenerSets by
  creation time, then alphabetical namespace/name order.
- **Name scoping:** listener names must be unique only *within* their
  ListenerSet -- two ListenerSets (or a ListenerSet and the Gateway) may reuse a
  name. Routes disambiguate by attaching to the specific ListenerSet
  (`parentRef` with `kind: ListenerSet` and optional `sectionName`).
- **Route attachment default:** a ListenerSet listener with no `allowed_routes`
  accepts routes from the ListenerSet's *own* namespace (not the Gateway's).

## Trust Boundaries: ReferenceGrants and Secrets

ListenerSet keeps namespaces as hard trust boundaries:

- **ReferenceGrants are not inherited** in either direction between a Gateway
  and its ListenerSets. A grant that authorizes the Gateway to read a Secret
  does not authorize any attached ListenerSet, and vice versa.
- A ListenerSet may reference TLS Secrets **in its own namespace without any
  ReferenceGrant** -- which is exactly the delegation point: each team keeps its
  certificate material in its own namespace, no cross-namespace trust needed.
- A cross-namespace `certificate_refs` entry still needs a ReferenceGrant in
  the Secret's namespace, exactly as for a Gateway.

## Validation Rules

Every upstream `XValidation` and kubebuilder marker is translated to
`buf.validate`, mirroring the Gateway listener rules:

- Listener names unique within the ListenerSet.
- The (port, protocol, hostname) combination unique across listeners.
- `tls` must not be set when protocol is HTTP, TCP, or UDP.
- An HTTPS listener may only `Terminate` (mode unset or `Terminate`).
- A TLS listener must declare its `tls.mode` (`Terminate` or `Passthrough`).
- `hostname` must not be set for TCP/UDP listeners.
- A `Terminate` listener must provide `certificate_refs` or `options` (enforced
  on the shared TLS config message).
- `listeners`: `min_items: 1`, `max_items: 64`; `port` 1-65535; `protocol` is
  the upstream open-set pattern (custom domain-prefixed protocols stay valid).

## Why ListenerSet Is a First-Class Planton Component

Without it, shared-gateway tenancy has two bad options: give every team write
access to the central Gateway, or funnel every listener change through the
platform team. Raw `KubernetesManifest` YAML restores neither proto validation
nor FK wiring nor InfraChart composability. As a typed component, ListenerSet
gives each team a self-contained, validated, DAG-ordered unit: namespace +
certificate + ListenerSet + routes.

## Design Notes

- **Standard channel, `v1`.** ListenerSet is standard-channel since Gateway API
  v1.5 and served as `gateway.networking.k8s.io/v1`; the component tracks
  v1.6.1.
- **Shared listener building blocks.** The `tls` and `allowed_routes` messages
  are the same shared Gateway API types `KubernetesGateway` uses -- one upstream
  type serves both kinds, so sharing prevents drift.
- **Dedicated parent reference type.** `parent_ref` uses
  `KubernetesGatewayApiParentGatewayReference` (upstream
  `ParentGatewayReference`), not the route-side `ParentReference`: a
  ListenerSet always attaches to the Gateway as a whole, so there is no
  `section_name` or `port`.
- **Reference names are foreign keys.** `parent_ref.name` defaults to
  `KubernetesGateway` and `certificate_refs[].name` to `KubernetesSecret`, both
  as `StringValueOrRef`: `valueFrom` wiring creates real dependency edges in
  infra charts, while `value:` carries literal names for referents outside
  Planton.
- **No baked-in Planton defaults for upstream fields.** Upstream defaults
  (`tls.mode=Terminate`, route `from` defaulting to the ListenerSet's own
  namespace) are documented in comments and left to the controller/CRD.
- **String-typed enum aliases stay strings.** Open-set values (`protocol`) are
  validated with the upstream regex; closed enums (`tls.mode`, namespace
  `from`) with CEL membership checks -- exact upstream casing, so CEL string
  comparisons and controller-specific values keep working.
- **Typed crd2pulumi IaC.** The Pulumi module uses `gatewayv1.NewListenerSet`;
  the listener mapping lives in `listeners.go`. The Terraform module applies
  through `kubectl_manifest` (alekc/kubectl) with server-side apply, plannable
  before the Gateway API CRDs exist.
- **Outputs expose the chain.** Besides `listener_set_name` and `namespace`,
  the component exports `gateway_name` (the resolved parent) so downstream
  resources can follow the chain to the Gateway without re-resolving the
  reference. Per-listener Accepted/Programmed conditions and the Gateway's
  AttachedListenerSets count are controller-managed and observed via kubectl,
  not stored in outputs.

## Composing in Infra Charts

The delegation pattern splits cleanly across charts. The platform chart owns
the Gateway with the opt-in:

```yaml
# platform chart: the shared Gateway opts in to same-namespace ListenerSets
spec:
  allowed_listeners:
    namespaces:
      from: Selector
      selector:
        matchLabels:
          gateway-access: "true"
```

Each team chart owns a ListenerSet wired entirely through foreign keys --
`namespace` (-> `KubernetesNamespace`), `parent_ref.name` (->
`KubernetesGateway.status.outputs.gateway_name`), and
`certificate_refs[].name` (-> `KubernetesSecret`, typically via a
`KubernetesCertificate`'s `status.outputs.secret_name`). `valueFrom` wiring
creates real dependency edges, so the platform builds the DAG automatically:

```yaml
metadata:
  name: "{{ values.team }}-listeners"
spec:
  namespace:
    valueFrom:
      kind: KubernetesNamespace
      name: "{{ values.team }}-ns"
      fieldPath: spec.name
  parentRef:
    name:
      valueFrom:
        kind: KubernetesGateway
        name: shared-gateway
        fieldPath: status.outputs.gateway_name
    namespace: gateway-ns
  listeners:
    - name: team-https
      hostname: "{{ values.team }}.example.com"
      port: 443
      protocol: HTTPS
      tls:
        mode: Terminate
        certificateRefs:
          - name:
              valueFrom:
                kind: KubernetesCertificate
                name: "{{ values.team }}-cert"
                fieldPath: status.outputs.secret_name
```

Team routes then attach to the merged listeners (parentRef `kind: ListenerSet`
plus `sectionName: team-https`). Full delegation stack:

```
KubernetesGateway (platform; allowed_listeners opt-in)
   -> KubernetesListenerSet (team; parent_ref valueFrom Gateway)
        -> KubernetesHttpRoute (team; parentRef kind ListenerSet)
```

## Controller Landscape

ListenerSet is implemented by Gateway API controllers that support merged
listeners (Envoy Gateway, Istio, and others); support is newer than for the
core Gateway/route kinds, so confirm your controller's version. The proto
carries no controller-specific behavior. Deploying the ListenerSet never blocks
on a controller: per-listener Accepted/Programmed conditions appear when the
Gateway implementation reconciles the merge, observed via kubectl rather than
awaited at apply time.

## Common Pitfalls

- **The Gateway must opt in first.** By default a Gateway allows no ListenerSet
  attachment; without `allowed_listeners` covering the ListenerSet's namespace,
  the attachment is rejected.
- **ReferenceGrants do not flow through.** A grant for the Gateway does not
  authorize its ListenerSets (or vice versa); each object needs its own grants
  for cross-namespace references.
- **Certificates belong in the ListenerSet's namespace.** Same-namespace
  Secrets need no grant -- keeping each team's certificate next to its
  ListenerSet is the intended (and simplest) pattern.
- **Setting `tls` on an HTTP/TCP/UDP listener.** TLS only applies to HTTPS/TLS
  listeners; the spec rejects it otherwise.
- **Name collisions across ListenerSets are legal.** Listener names are scoped
  to their ListenerSet; if a route must target a specific listener, attach it
  to that ListenerSet (`kind: ListenerSet` + `sectionName`), not the Gateway.
- **Merge order matters for conflicts.** The Gateway's listeners win, then
  earlier-created ListenerSets, then alphabetical order -- a conflicting (port,
  protocol, hostname) entry in a later ListenerSet loses.

## Conclusion

`KubernetesListenerSet` completes the shared-gateway delegation story of the
Planton Gateway API layer: the platform team keeps one Gateway, each tenant team
keeps its own validated, FK-wired, DAG-ordered listener bundle -- certificates
included -- without any team ever editing the central entry point.

## References

- [Gateway API ListenerSet](https://gateway-api.sigs.k8s.io/api-types/listenerset/)
- [Gateway API: Gateway](https://gateway-api.sigs.k8s.io/api-types/gateway/)
- Upstream types: `kubernetes-sigs/gateway-api` `apis/v1/listenerset_types.go` (v1.6.1)
- Sibling components: `KubernetesGateway`, `KubernetesHttpRoute`, `KubernetesReferenceGrant`
