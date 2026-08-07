# KubernetesUdpRoute

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `kubernetes.planton.dev/v1alpha1`

KubernetesUdpRouteSpec defines a Kubernetes Gateway API UDPRoute: a namespaced
set of rules that forward UDP datagrams arriving on a Gateway listener to
backend Services. A UDPRoute has no matching logic at all -- traffic on the
listener's port is forwarded to the rule's backends. The route attaches to a
Gateway listener of protocol UDP through parent_refs. Typical backends are
DNS servers, syslog collectors, game servers, and other datagram protocols.

100% fidelity with the upstream Gateway API v1.6.1 UDPRouteSpec
(kubernetes-sigs/gateway-api apis/v1/udproute_types.go), standard channel.
UDPRoute graduated to GA and is served as gateway.networking.k8s.io/v1 (it
was an experimental v1alpha2 resource in earlier releases). The experimental
CommonRouteSpec field useDefaultGateways is intentionally excluded: it is
absent from the standard-channel CRD, so it would have no deployable target.

UDPRoute is a layer-4 route: it has no hostnames, no matches, and no filters
(UDP has no connection or request structure to match on). Each rule simply
forwards to one or more backends, so it reuses the shared
KubernetesGatewayApiBackendRef directly rather than defining its own
per-route backend ref.

## Example

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesUdpRoute
metadata:
  name: test-udp-route
spec:
  namespace:
    value: test-namespace
  parentRefs:
    - name:
        value: test-gateway
      sectionName: dns
      port: 53
  rules:
    - name: dns-forward
      backendRefs:
        - name:
            value: dns-backend
          port: 53
          weight: 90
        - name:
            value: dns-backend-canary
          port: 53
          weight: 10
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
| `spec.rules` | `[]KubernetesUdpRouteRule` | yes |  |  |
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

Namespace in which the UDPRoute is created. Backends in other namespaces,
and Gateways in other namespaces, are subject to the usual same-namespace /
ReferenceGrant rules.

- references: KubernetesNamespace (`spec.name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesNamespace, name: <that resource's name>, fieldPath: spec.name}} -- a bare string does not parse

### spec.parentRefs

`[]KubernetesGatewayApiParentReference`

References to the parent resources (usually Gateways) this route attaches
to. Each parent must have a UDP listener that allows attachment from
UDPRoutes of this namespace. Each reference's name defaults to a
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

### spec.rules

`[]KubernetesUdpRouteRule` · required

Routing rules. Each rule forwards matching datagrams to one or more
backends. At least one rule is required; upstream allows up to 16.

- rule: {"repeated":{"minItems":"1","maxItems":"16"}}

### spec.rules[].name

`string` · required · optional (explicit presence)

Name of the route rule; must be unique within the route if set. Upstream
support level: Extended.

Upstream SectionName constraints: 1-253 chars, lowercase RFC 1123 subdomain.
Pattern: ^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$

- rule: {"string":{"minLen":"1","maxLen":"253","pattern":"^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$"}}

### spec.rules[].backendRefs

`[]KubernetesGatewayApiBackendRef` · required

Backends matching datagrams are forwarded to. When multiple backends are
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

Reference an output from another manifest as `valueFrom: {kind: KubernetesUdpRoute, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.route_name` | `string` | Name of the created UDPRoute (equals metadata.name). In InfraCharts this orders the route after the Gateway and backends it references. |
| `status.outputs.namespace` | `string` | Namespace the UDPRoute was created in (the resolved spec.namespace). Cross-namespace parent and backend references from this route are subject to ReferenceGrant rules relative to this value. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.namespace` | KubernetesNamespace | `spec.name` |
| `spec.parentRefs[].name` | KubernetesGateway | `status.outputs.gateway_name` |
| `spec.rules[].backendRefs[].name` | KubernetesService | `status.outputs.service_name` |

## See Also

- [Overview](../README.md)
