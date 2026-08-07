# KubernetesTlsRoute -- Research Documentation

This document explains why `KubernetesTlsRoute` exists, how it maps the upstream
Gateway API `TLSRoute` into the Planton component model, and the reasoning behind
its spec, validation, and IaC. It is the source-of-truth companion to the
user-facing `README.md` (getting started) and `catalog-page.md` (catalog listing).

## Table of Contents

1. Introduction
2. The Gateway API Role-Oriented Model
3. Where TLSRoute Sits
4. Anatomy of TLSRouteSpec
5. The SNI / Backend Model
6. Validation Rules
7. Why TLSRoute Is a First-Class Planton Component
8. Design Notes
9. Standard Channel Status
10. Composing in Infra Charts
11. Controller Landscape
12. Common Pitfalls
13. Conclusion
14. References

## Introduction

The Gateway API is the successor to the Kubernetes Ingress API: a standards-based,
role-oriented, expressive specification for north-south (and increasingly
east-west) traffic. `TLSRoute` routes TLS connections by their SNI hostname and
forwards them unmodified to a backend -- a layer-4 "passthrough" route. The
backend, not the Gateway, terminates TLS, so the encrypted byte stream is carried
end to end. This is the right tool when a service must hold its own certificate or
do its own mTLS, and the Gateway should route but never decrypt.

`KubernetesTlsRoute` brings that resource into Planton as a first-class deployment
component modeling the complete upstream standard-channel surface, so customers
never have to fall back to a raw `KubernetesManifest` to express TLS passthrough
routing. It is a sibling of `KubernetesHttpRoute` / `KubernetesGrpcRoute` and
shares their envelope, conventions, and IaC structure.

## The Gateway API Role-Oriented Model

The Gateway API deliberately splits responsibilities across personas:

- **Infrastructure provider** owns the `GatewayClass` (which controller, e.g.
  Istio or Envoy Gateway, implements Gateways).
- **Cluster operator** owns `Gateway` objects (listeners, ports, TLS, addresses).
- **Application developer** owns `*Route` objects (`HTTPRoute`, `TLSRoute`, ...)
  that attach to a Gateway and define routing for their app.

`TLSRoute` is in the application-developer lane. A route attaches to a Gateway via
`parentRefs`; the Gateway's TLS listener decides which route kinds and namespaces
may attach (via `allowedRoutes`). A TLSRoute must attach to a listener of protocol
`TLS` -- `Passthrough` mode for core support, `Terminate` mode for the extended
`TLSRouteTermination` feature.

## Where TLSRoute Sits

```
KubernetesGatewayApiCrds        (install CRDs + controller)
        |
KubernetesGatewayClass          (controller class)
        |
KubernetesGateway               (TLS listener, Passthrough mode)
        |
KubernetesTlsRoute              (this component: SNI-based passthrough)
        |
backend Services (terminate TLS themselves)
```

## Anatomy of TLSRouteSpec

The Planton spec flattens the upstream `TLSRouteSpec` after the standard
namespaced envelope (`namespace`):

- `parent_refs` -- the Gateways this route attaches to (max 32).
- `hostnames` -- the SNI hostnames that select the route (**required**, 1 to
  1024; wildcard prefix allowed; IPs not allowed). Upstream raised the bound
  from 16 to 1024 in the v1.6 release.
- `rules` -- **exactly one** routing rule (`min_items: 1`, `max_items: 1`). Each
  rule has only:
  - an optional `name`;
  - `backend_refs` (1 to 16) -- weighted backends.

There are no matches and no filters: a TLS passthrough route has no layer-7
visibility to match on or transform.

Reference names are foreign keys: `parent_refs[].name` defaults to
`KubernetesGateway` and `backend_refs[].name` defaults to `KubernetesService`,
both as `StringValueOrRef`. In infra charts these wire with `valueFrom` and give
the resource graph real dependency edges; when the referent is not
Planton-managed, the literal name is passed with `value:`.

## The SNI / Backend Model

- **Hostnames** are matched against the SNI attribute of the TLS ClientHello, not
  an HTTP Host header (the Gateway cannot read HTTP -- the stream is encrypted). A
  leading `*.` is a single-label wildcard suffix match. Per RFC 6066, SNI names
  may not be IP addresses.
- **Backend ref** (the shared `KubernetesGatewayApiBackendRef`):
  group/kind/name/namespace/port/weight. Because TLS routes have no per-backend
  filters, they reuse the canonical shared backend reference directly rather than
  defining a per-route backend ref (the filter-carrying HTTP/GRPC routes flatten
  their own).

## Validation Rules

Every upstream `XValidation` and kubebuilder marker is translated to
`buf.validate`:

- `hostnames`: required (`min_items: 1`), `max_items: 1024` (the upstream v1.6
  bound); per-item RFC 1123 wildcard-prefix pattern; a list-level CEL rejecting
  IPv4 literals (the SNI no-IP rule from RFC 6066).
- `rules`: `min_items: 1`, `max_items: 1` (upstream caps a TLSRoute at one rule).
- `backend_refs`: `min_items: 1`, `max_items: 16`; the shared backend ref carries
  the upstream group/kind patterns and port/weight bounds.
- `name` (rule): SectionName pattern (lowercase RFC 1123 subdomain, 1-253).

The upstream length/format constraints on foreign-key `name` fields are enforced
by the Kubernetes API server at apply time; the referenced object's own creation
already validated them.

### A note on `isIp`

Upstream expresses the SNI no-IP rule as `self.all(h, !isIP(h))`. Planton's
build-time `buf lint` CEL environment does not register protovalidate's `isIp()`
format function (only the runtime validator does), so the rule is translated as an
equivalent IPv4 dotted-quad regex. IPv6 literals contain `:` and are already
rejected by the per-item hostname pattern, so the IPv4 guard fully covers the
reachable input space. This is documented inline in `spec.proto` so a future agent
does not "restore" `isIp` and break the lint gate.

## Why TLSRoute Is a First-Class Planton Component

Without it, TLS passthrough routing forces customers back to raw
`KubernetesManifest` YAML: no proto validation, no typed SDKs, no FK wiring, no
InfraChart composability, no UI wizards. `KubernetesTlsRoute` closes that gap for
the TLS-passthrough slice of the ingress story, alongside `KubernetesHttpRoute`
and `KubernetesGrpcRoute`.

## Design Notes

- **Standard channel, `v1`.** TLSRoute is standard-channel since Gateway API v1.5
  and is served as `gateway.networking.k8s.io/v1` (it was experimental
  `v1alpha2` / `v1alpha3` in earlier releases). The typed crd2pulumi resource
  emits the `v1` apiVersion.
- **Foreign-key `parent_refs` / `backend_refs` names.** Each reference's `name`
  is a `StringValueOrRef` (Gateway and Service foreign keys respectively):
  `valueFrom` wiring creates real dependency edges in infra charts, while
  `value:` carries literal names for referents outside Planton.
- **Reuse the shared backend ref.** TLS routes have no per-backend filters, so
  they consume the canonical `KubernetesGatewayApiBackendRef` directly.
- **String-typed enum aliases stay strings.** Upstream models its enums as string
  type aliases with validation markers; the proto mirrors that with `optional
  string` plus pattern/CEL checks, so CEL comparisons and controller-specific
  values keep working, and no Planton defaults are baked into upstream fields.
- **Typed crd2pulumi IaC.** `gatewayv1.NewTLSRoute`; the spec mapping is split
  across `parent_refs.go` and `rules.go` (no matches/filters).
- **Terraform applies through `kubectl_manifest`** (alekc/kubectl) with
  server-side apply: plannable before the Gateway API CRDs exist, so an infra
  chart can deploy CRDs, Gateway, and routes in one run.

## Standard Channel Status

TLSRoute is served in the **standard** channel as `v1` (standard since Gateway
API v1.5). The standard-channel `TLSRouteSpec` has no `useDefaultGateways` field
(that experimental field is stripped from the standard CRD), so there is nothing
experimental to exclude here -- the spec is the full standard-channel surface.

## Composing in Infra Charts

`KubernetesTlsRoute` is designed as a LEGO block for Infra Charts. Every neighbor
reference is a `StringValueOrRef` foreign key:

- `namespace` (default kind `KubernetesNamespace`),
- `parent_refs[].name` (default kind `KubernetesGateway`),
- `backend_refs[].name` (default kind `KubernetesService`).

Wiring them with `valueFrom` creates real dependency edges, so the platform
builds the DAG automatically -- no manual relationship declarations:

```yaml
metadata:
  name: "{{ values.env }}-secure-route"
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
      sectionName: tls
  hostnames:
    - "secure.{{ values.domain }}"
  rules:
    - backendRefs:
        - name:
            valueFrom:
              kind: KubernetesService
              name: "{{ values.service_name }}"
              fieldPath: status.outputs.service_name
          port: 8443
```

The full ingress stack composes as:

```
KubernetesCertManager            (runs_on cluster)
  -> KubernetesClusterIssuer     (cert_manager_namespace: valueFrom CertManager)
    -> KubernetesCertificate     (issuerRef: valueFrom ClusterIssuer)
      -> (Kubernetes Secret)     (Certificate.status.outputs.secret_name)
        -> KubernetesGateway     (certificate_refs[].name: valueFrom Certificate)
          -> KubernetesTlsRoute  (parent_refs[].name: valueFrom Gateway;
                                  backend_refs[].name: valueFrom Service)
```

(For pure TLS passthrough the backend holds its own certificate, so the
cert-manager prefix is optional; it applies when the Gateway terminates TLS in the
extended `Terminate` mode.)

## Controller Landscape

TLSRoute is implemented by the major Gateway API controllers that support
layer-4 routing (Istio, Envoy Gateway, and others). The proto carries no
controller-specific behavior; upstream defaults are documented in comments and
left to the controller rather than baked into the spec. Deploying the route never
blocks on a controller: the per-parent Accepted/ResolvedRefs conditions appear
when a Gateway implementation reconciles the resource, observed via kubectl
rather than awaited at apply time.

## Common Pitfalls

- **Listener must be `TLS` protocol.** A TLSRoute attached to an `HTTP`/`HTTPS`
  listener will not attach. The parent listener must use protocol `TLS`
  (`Passthrough` for core support).
- **No HTTP-level matching.** TLSRoutes route only by SNI hostname -- there is no
  path, header, or method matching. Use `KubernetesHttpRoute` when the Gateway
  terminates TLS and you need layer-7 routing.
- **Exactly one rule.** Upstream caps a TLSRoute at one rule; express traffic
  splitting through multiple weighted `backendRefs` in that single rule.
- **Very large hostname sets.** The 1024-hostname bound matches upstream v1.6;
  validate apiserver/etcd/controller behavior with representative manifests
  before relying on very large sets in production.
- **Literal `value:` names create no dependency edge.** In infra charts, prefer
  `valueFrom` for Planton-managed Gateways and Services so ordering is
  guaranteed.

## Conclusion

`KubernetesTlsRoute` completes the TLS-passthrough slice of the Planton Gateway
API ingress layer, modeling the complete upstream standard-channel surface,
mirroring its sibling routes in every convention and fully accounting for
infra-chart composability.

## References

- [Gateway API TLSRoute](https://gateway-api.sigs.k8s.io/api-types/tlsroute/)
- Upstream types: `kubernetes-sigs/gateway-api` `apis/v1/tlsroute_types.go` (v1.6.1)
- Sibling components: `KubernetesHttpRoute`, `KubernetesGrpcRoute`, `KubernetesTcpRoute`, `KubernetesUdpRoute`
