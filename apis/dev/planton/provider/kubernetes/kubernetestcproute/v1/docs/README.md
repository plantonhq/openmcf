# KubernetesTcpRoute -- Research Documentation

This document explains why `KubernetesTcpRoute` exists, how it maps the upstream
Gateway API `TCPRoute` into the Planton component model, and the reasoning behind
its spec, validation, and IaC. It is the source-of-truth companion to the
user-facing `README.md` (getting started) and `catalog-page.md` (catalog listing).

## Table of Contents

1. Introduction
2. Standard Channel, GA
3. The Gateway API Role-Oriented Model
4. Where TCPRoute Sits
5. Anatomy of TCPRouteSpec
6. The Backend Model
7. Validation Rules
8. Why TCPRoute Is a First-Class Planton Component
9. Design Notes
10. Composing in Infra Charts
11. Controller Landscape
12. Common Pitfalls
13. Conclusion
14. References

## Introduction

The Gateway API is the successor to the Kubernetes Ingress API: a standards-based,
role-oriented, expressive specification for north-south (and increasingly
east-west) traffic. `TCPRoute` is the simplest route kind: it forwards raw TCP
connections arriving on a Gateway listener to a set of backend Services. There is
no matching of any kind -- the listener's port selects the traffic, and the route
forwards it (splitting by weight across backends). Use it to put non-HTTP TCP
services (Postgres, Redis, Kafka, a custom protocol) behind a Gateway.

`KubernetesTcpRoute` brings that resource into Planton as a first-class deployment
component modeling the complete upstream surface, so customers never have to fall
back to a raw `KubernetesManifest` to express TCP routing.

## Standard Channel, GA

TCPRoute graduated to GA in the Gateway API standard channel and is served as
`gateway.networking.k8s.io/v1` (it was an experimental `v1alpha2` resource in
earlier Gateway API releases). Consequences:

- The default standard-channel install of `KubernetesGatewayApiCrds` (Gateway API
  v1.6.1) includes the TCPRoute CRD; no experimental channel is required.
- The typed crd2pulumi resource is `gatewayv1.NewTCPRoute`, and the Terraform
  module applies apiVersion `gateway.networking.k8s.io/v1`.
- The experimental `CommonRouteSpec` field `useDefaultGateways` is excluded: it
  is absent from the standard-channel CRD, so it would have no deployable target.

## The Gateway API Role-Oriented Model

The Gateway API splits responsibilities across personas:

- **Infrastructure provider** owns the `GatewayClass` (which controller implements
  Gateways).
- **Cluster operator** owns `Gateway` objects (listeners, ports, TLS).
- **Application developer** owns `*Route` objects that attach to a Gateway.

`TCPRoute` is in the application-developer lane. A route attaches to a Gateway via
`parentRefs`; the Gateway's TCP listener decides which route kinds and namespaces
may attach (via `allowedRoutes`). A TCPRoute must attach to a listener of protocol
`TCP`.

## Where TCPRoute Sits

```
KubernetesGatewayApiCrds        (install standard-channel CRDs)
        |
KubernetesGatewayClass          (controller class)
        |
KubernetesGateway               (TCP listener)
        |
KubernetesTcpRoute              (this component: raw TCP forwarding)
        |
backend Services
```

## Anatomy of TCPRouteSpec

The Planton spec flattens the upstream `TCPRouteSpec` after the standard
namespaced envelope (`namespace`):

- `parent_refs` -- the Gateways this route attaches to (max 32). Flattened from
  the upstream `CommonRouteSpec`.
- `rules` -- 1 to 16 routing rules. Each rule has only:
  - an optional `name`;
  - `backend_refs` (1 to 16) -- weighted backends.

There are no hostnames, matches, or filters: a TCP route has no application-layer
visibility.

Reference names are foreign keys: `parent_refs[].name` defaults to
`KubernetesGateway` and `backend_refs[].name` defaults to `KubernetesService`,
both as `StringValueOrRef`. In infra charts these wire with `valueFrom` and give
the resource graph real dependency edges; when the referent is not
Planton-managed, the literal name is passed with `value:`.

## The Backend Model

A TCP route forwards to the shared `KubernetesGatewayApiBackendRef`
(group/kind/name/namespace/port/weight). Because TCP routes have no per-backend
filters, they reuse the canonical shared backend reference directly rather than
defining a per-route backend ref (the filter-carrying HTTP/GRPC routes flatten
their own). If a backend is invalid, the implementation rejects connection
attempts in proportion to the backend's weight.

## Validation Rules

Every upstream `XValidation` and kubebuilder marker is translated to
`buf.validate`:

- `parent_refs`: `max_items: 32`.
- `rules`: `min_items: 1`, `max_items: 16`.
- `backend_refs`: `min_items: 1`, `max_items: 16`; the shared backend ref carries
  the upstream group/kind patterns and port/weight bounds.
- `name` (rule): SectionName pattern (lowercase RFC 1123 subdomain, 1-253).

The upstream length/format constraints on foreign-key `name` fields are enforced
by the Kubernetes API server at apply time; the referenced object's own creation
already validated them.

The experimental rule-name-uniqueness `XValidation` is not translated (consistent
with the rest of the route family -- it is `<gateway:experimental>` even in the
GA CRD and is controller-enforced).

## Why TCPRoute Is a First-Class Planton Component

Without it, TCP routing forces customers back to raw `KubernetesManifest` YAML: no
proto validation, no typed SDKs, no FK wiring, no InfraChart composability, no UI
wizards. `KubernetesTcpRoute` closes that gap for the raw-TCP slice of the ingress
story.

## Design Notes

- **Standard channel, `v1`.** Dictated by upstream Gateway API v1.6.1; the typed
  crd2pulumi resource and the Terraform manifest both use the `v1` apiVersion.
- **`use_default_gateways` is excluded.** This experimental `CommonRouteSpec`
  field is absent from the standard-channel CRD, so including it would produce a
  field with no deployable target -- the same reasoning applied across the whole
  route family.
- **Foreign-key `parent_refs` / `backend_refs` names.** Each reference's `name`
  is a `StringValueOrRef` (Gateway and Service foreign keys respectively):
  `valueFrom` wiring creates real dependency edges in infra charts, while
  `value:` carries literal names for referents outside Planton.
- **Reuse the shared backend ref.** TCP routes have no per-backend filters, so
  they consume the canonical `KubernetesGatewayApiBackendRef` directly.
- **String-typed enum aliases stay strings.** Upstream models its enums as string
  type aliases with validation markers; the proto mirrors that with `optional
  string` plus CEL `in [...]` checks, so CEL comparisons and controller-specific
  values keep working, and no Planton defaults are baked into upstream fields.
- **Typed crd2pulumi IaC.** `gatewayv1.NewTCPRoute`; the spec mapping is split
  across `parent_refs.go` and `rules.go` (no matches/filters).
- **Terraform applies through `kubectl_manifest`** (alekc/kubectl) with
  server-side apply: plannable before the Gateway API CRDs exist, so an infra
  chart can deploy CRDs, Gateway, and routes in one run.

## Composing in Infra Charts

`KubernetesTcpRoute` is designed as a LEGO block for Infra Charts. Every neighbor
reference is a `StringValueOrRef` foreign key:

- `namespace` (default kind `KubernetesNamespace`),
- `parent_refs[].name` (default kind `KubernetesGateway`),
- `backend_refs[].name` (default kind `KubernetesService`).

Wiring them with `valueFrom` creates real dependency edges, so the platform
builds the DAG automatically -- no manual relationship declarations:

```yaml
metadata:
  name: "{{ values.env }}-postgres-route"
spec:
  namespace:
    valueFrom:
      kind: KubernetesNamespace
      name: "{{ values.env }}-ns"
      fieldPath: spec.name
  parentRefs:
    - name:
        valueFrom:
          kind: KubernetesGateway
          name: "{{ values.env }}-gateway"
          fieldPath: status.outputs.gateway_name
      sectionName: tcp
  rules:
    - backendRefs:
        - name:
            valueFrom:
              kind: KubernetesService
              name: "{{ values.service_name }}"
              fieldPath: status.outputs.service_name
          port: 5432
```

When the parent or backend is not a Planton-managed resource, pass the literal
name with `value:` -- no edge is created, matching reality.

## Controller Landscape

TCPRoute is implemented by Gateway API controllers that support layer-4 routing
(Istio, Envoy Gateway, and others). The proto carries no controller-specific
behavior; upstream defaults are documented in comments and left to the controller
rather than baked into the spec. Deploying the route never blocks on a
controller: the per-parent Accepted/ResolvedRefs conditions appear when a Gateway
implementation reconciles the resource, which is observed via kubectl rather than
awaited at apply time.

## Common Pitfalls

- **Listener must be `TCP` protocol.** A TCPRoute attached to an HTTP/HTTPS/TLS
  listener will not attach.
- **No matching.** TCP routes cannot match on hostnames, paths, headers, or
  methods -- they forward by listener port. Use TLSRoute (SNI) or HTTPRoute
  (layer 7) when you need matching.
- **Controller support varies.** Not every Gateway controller implements
  layer-4 routes; confirm yours does before attaching.
- **Literal `value:` names create no dependency edge.** In infra charts, prefer
  `valueFrom` for Planton-managed Gateways and Services so ordering is
  guaranteed.

## Conclusion

`KubernetesTcpRoute` completes the raw-TCP slice of the Planton Gateway API
ingress layer, modeling the complete upstream GA surface, mirroring its sibling
routes in every convention and fully accounting for infra-chart composability.

## References

- [Gateway API TCPRoute](https://gateway-api.sigs.k8s.io/api-types/tcproute/)
- Upstream types: `kubernetes-sigs/gateway-api` `apis/v1/tcproute_types.go` (v1.6.1)
- Sibling components: `KubernetesUdpRoute`, `KubernetesTlsRoute`, `KubernetesHttpRoute`, `KubernetesGrpcRoute`
