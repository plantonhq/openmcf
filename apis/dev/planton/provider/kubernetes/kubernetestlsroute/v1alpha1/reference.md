# KubernetesTlsRoute

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `kubernetes.planton.dev/v1alpha1`

KubernetesTlsRouteSpec defines a Kubernetes Gateway API TLSRoute: a namespaced
set of rules that match TLS connections by their SNI hostname and forward them,
unmodified, to backend Services (TLS passthrough). The route attaches to a
Gateway listener of protocol TLS through parent_refs.

100% fidelity with the upstream Gateway API v1.6.1 TLSRouteSpec
(kubernetes-sigs/gateway-api apis/v1/tlsroute_types.go), standard channel.
TLSRoute graduated to the standard channel and is served as
gateway.networking.k8s.io/v1 (it was experimental v1alpha2/v1alpha3 in earlier
releases). Upstream spec fields are flattened after the Planton namespaced
envelope (namespace).

TLSRoute is a layer-4 (connection) route: it has no path/header/method matches
and no filters. Each rule simply forwards to one or more backends, so it reuses
the shared KubernetesGatewayApiBackendRef directly rather than defining its own
per-route backend ref.

## Example

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesTlsRoute
metadata:
  name: test-tls-route
spec:
  namespace:
    value: test-namespace
  parentRefs:
    - name:
        value: test-gateway
      sectionName: tls-passthrough
  hostnames:
    - secure.example.com
    - "*.tenants.example.com"
  rules:
    - name: passthrough
      backendRefs:
        - name:
            value: secure-backend
          port: 8443
          weight: 100
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.namespace` | `string \| valueFrom` | yes |  | KubernetesNamespace (`spec.name`) |
| `spec.parentRefs` | `[]KubernetesGatewayApiParentReference` |  |  |  |
| `spec.parentRefs[].group` | `string` |  |  |  |
| `spec.parentRefs[].kind` | `string` | yes |  |  |
| `spec.parentRefs[].namespace` | `string` |  |  |  |
| `spec.parentRefs[].name` | `string \| valueFrom` | yes |  | KubernetesGateway (`status.outputs.gateway_name`) |
| `spec.parentRefs[].sectionName` | `string` | yes |  |  |
| `spec.parentRefs[].port` | `int32` |  |  |  |
| `spec.hostnames` | `[]string` | yes |  |  |
| `spec.rules` | `[]KubernetesTlsRouteRule` | yes |  |  |
| `spec.rules[].name` | `string` | yes |  |  |
| `spec.rules[].backendRefs` | `[]KubernetesGatewayApiBackendRef` | yes |  |  |
| `spec.rules[].backendRefs[].group` | `string` |  |  |  |
| `spec.rules[].backendRefs[].kind` | `string` | yes |  |  |
| `spec.rules[].backendRefs[].name` | `string \| valueFrom` | yes |  | KubernetesService (`status.outputs.service_name`) |
| `spec.rules[].backendRefs[].namespace` | `string` |  |  |  |
| `spec.rules[].backendRefs[].port` | `int32` |  |  |  |
| `spec.rules[].backendRefs[].weight` | `int32` |  |  |  |

## Field Details

### spec.namespace

`string | valueFrom` · required

Namespace in which the TLSRoute is created. Backends in other namespaces,
and Gateways in other namespaces, are subject to the usual same-namespace /
ReferenceGrant rules.

- references: KubernetesNamespace (`spec.name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesNamespace, name: <that resource's name>, fieldPath: spec.name}} -- a bare string does not parse

### spec.parentRefs

`[]KubernetesGatewayApiParentReference`

References to the parent resources (usually Gateways) this route attaches
to. Each parent must have a TLS listener that allows attachment from
TLSRoutes of this namespace. Each reference's name defaults to a
KubernetesGateway foreign key — wire it with valueFrom in an infra chart
and the route deploys after its Gateway.

Flattened from the upstream CommonRouteSpec.

- rule: {"repeated":{"maxItems":"32"}}

### spec.parentRefs[].group

`string` · optional (explicit presence)

Group of the referent.
When unspecified, "gateway.networking.k8s.io" is inferred.
Set to "" (empty string) for the core API group.

Upstream default: "gateway.networking.k8s.io"
Group pattern: empty or an RFC 1123 subdomain (max 253).

- rule: {"string":{"maxLen":"253","pattern":"^$|^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$"}}

### spec.parentRefs[].kind

`string` · required · optional (explicit presence)

Kind of the referent.
Core parent kinds: Gateway (Gateway conformance), Service (Mesh conformance).

Upstream default: "Gateway"
Kind pattern: 1-63 chars, ^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$

- rule: {"string":{"minLen":"1","maxLen":"63","pattern":"^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$"}}

### spec.parentRefs[].namespace

`string` · optional (explicit presence)

Namespace of the referent. When unspecified, refers to the local
namespace of the Route. Cross-namespace references require a
ReferenceGrant in the target namespace.

### spec.parentRefs[].name

`string | valueFrom` · required

Name of the referent. Defaults to a KubernetesGateway foreign key: in an
infra chart, wire it with valueFrom against the Gateway resource and the
route deploys after (and follows renames of) its Gateway. When the parent
is not a Planton-managed Gateway (a ListenerSet, a mesh Service, or an
externally created Gateway), pass the literal name with `value:`.

- references: KubernetesGateway (`status.outputs.gateway_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesGateway, name: <that resource's name>, fieldPath: status.outputs.gateway_name}} -- a bare string does not parse

### spec.parentRefs[].sectionName

`string` · required · optional (explicit presence)

Name of a section within the target resource (e.g., a Gateway Listener
name). When unspecified, references the entire resource.

SectionName pattern: 1-253 chars, RFC 1123 subdomain
^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$

- rule: {"string":{"minLen":"1","maxLen":"253","pattern":"^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$"}}

### spec.parentRefs[].port

`int32` · optional (explicit presence)

Network port this Route targets on the parent resource.

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.hostnames

`[]string` · required

SNI hostnames matched against the SNI attribute of the TLS ClientHello during
the handshake. A leading wildcard label (for example "*.example.com") is a
suffix match. At least one hostname is required (upstream requires hostnames
for TLSRoute, unlike HTTP/GRPC routes where it is optional). Upstream support
level: Core.

Each value is an RFC 1123 hostname; per RFC 6066, SNI hostnames may not be IP
addresses. The per-item pattern enforces the RFC 1123 + wildcard-first-label
shape; the list-level CEL below additionally rejects IPv4 literals. This
no-IP rule is specific to TLS SNI (RFC 6066) and is genuinely absent from the
upstream HTTP/GRPC hostname rules -- so translating it here is fidelity to the
TLSRoute spec, not divergence from the sibling routes.

Upstream expresses this as `!isIP(h)`. Planton's buf-lint CEL environment does
not register protovalidate's isIp() format function, so we use an equivalent
IPv4 dotted-quad regex: any value containing ':' (IPv6) is already rejected by
the per-item hostname pattern above, so an IPv4 guard fully covers the
reachable input space.

- rule: SNI hostnames cannot be IP addresses (RFC 6066)
- rule: {"repeated":{"minItems":"1","maxItems":"1024","items":{"string":{"minLen":"1","maxLen":"253","pattern":"^(\\*\\.)?[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$"}}}}

### spec.rules

`[]KubernetesTlsRouteRule` · required

Routing rules. Each rule forwards matching connections to one or more
backends. Upstream allows exactly one rule for a TLSRoute (max_items 1): a
TLS passthrough route has no layer-7 matching, so a second rule would be
ambiguous. At least one rule is required.

- rule: {"repeated":{"minItems":"1","maxItems":"1"}}

### spec.rules[].name

`string` · required · optional (explicit presence)

Name of the route rule; must be unique within the route if set. Upstream
support level: Extended.

Upstream SectionName constraints: 1-253 chars, lowercase RFC 1123 subdomain.
Pattern: ^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$

- rule: {"string":{"minLen":"1","maxLen":"253","pattern":"^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$"}}

### spec.rules[].backendRefs

`[]KubernetesGatewayApiBackendRef` · required

Backends matching connections are forwarded to. When multiple backends are
given, traffic is split by weight. At least one backend is required.
Upstream support level: Core for Kubernetes Service. Each backend's name
defaults to a KubernetesService foreign key — wire it with valueFrom in an
infra chart and the route deploys after its backend.

- rule: {"repeated":{"minItems":"1","maxItems":"16"}}

### spec.rules[].backendRefs[].group

`string` · optional (explicit presence)

Group of the referent. Empty string infers the core API group.

Upstream default: "" (core API group).
Group pattern: empty or an RFC 1123 subdomain (max 253).

- rule: {"string":{"maxLen":"253","pattern":"^$|^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$"}}

### spec.rules[].backendRefs[].kind

`string` · required · optional (explicit presence)

Kind of the referent.

Upstream default: "Service".
Kind pattern: 1-63 chars, ^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$

- rule: {"string":{"minLen":"1","maxLen":"63","pattern":"^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$"}}

### spec.rules[].backendRefs[].name

`string | valueFrom` · required

Name of the referent. Defaults to a KubernetesService foreign key: in an
infra chart, wire it with valueFrom against the backend Service and the
route deploys after its backend. When the backend is not a
Planton-managed Service (a custom backend kind, or a Service created
outside Planton), pass the literal name with `value:`.

- references: KubernetesService (`status.outputs.service_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesService, name: <that resource's name>, fieldPath: status.outputs.service_name}} -- a bare string does not parse

### spec.rules[].backendRefs[].namespace

`string` · optional (explicit presence)

Namespace of the backend. When unspecified, the local namespace is inferred.
Cross-namespace references require a ReferenceGrant.

### spec.rules[].backendRefs[].port

`int32` · optional (explicit presence)

Destination port number. Required when the referent is a core Service.

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.rules[].backendRefs[].weight

`int32` · optional (explicit presence)

Proportion of requests forwarded to this backend, computed as
weight/(sum of all weights). Unspecified defaults to 1.

Upstream default: 1

- rule: {"int32":{"lte":1000000,"gte":0}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: KubernetesTlsRoute, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.route_name` | `string` | Name of the created TLSRoute (equals metadata.name). In InfraCharts this orders the route after the Gateway and backends it references. |
| `status.outputs.namespace` | `string` | Namespace the TLSRoute was created in (the resolved spec.namespace). Cross-namespace parent and backend references from this route are subject to ReferenceGrant rules relative to this value. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.namespace` | KubernetesNamespace | `spec.name` |
| `spec.parentRefs[].name` | KubernetesGateway | `status.outputs.gateway_name` |
| `spec.rules[].backendRefs[].name` | KubernetesService | `status.outputs.service_name` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
