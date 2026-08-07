# KubernetesGrpcRoute

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `kubernetes.planton.dev/v1alpha1`

KubernetesGrpcRouteSpec defines a Kubernetes Gateway API GRPCRoute: a
namespaced set of rules that match gRPC requests (by hostname, service/method,
or header), optionally transform them with filters, and forward them to
backend Services. Routes attach to a Gateway through parent_refs.

100% fidelity with the upstream Gateway API v1.6.1 GRPCRouteSpec
(kubernetes-sigs/gateway-api apis/v1/grpcroute_types.go), standard channel.
Upstream spec fields are flattened after the Planton namespaced envelope
(namespace). Experimental fields are intentionally excluded
because they are absent from the standard-channel CRD and the typed Pulumi
resource Planton provisions with, so they would have no deployable target:
  - CommonRouteSpec.useDefaultGateways
  - GRPCRouteRule.sessionPersistence

## Example

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesGrpcRoute
metadata:
  name: test-grpc-route
spec:
  namespace:
    value: test-namespace
  parentRefs:
    - name:
        value: test-gateway
      sectionName: https
  hostnames:
    - api.example.com
  rules:
    - name: greeter
      matches:
        - method:
            type: Exact
            service: helloworld.Greeter
            method: SayHello
          headers:
            - type: Exact
              name: x-tenant
              value: acme
      filters:
        - type: RequestHeaderModifier
          requestHeaderModifier:
            set:
              - name: x-forwarded-tier
                value: production
        - type: RequestMirror
          requestMirror:
            backendRef:
              name:
                value: audit-svc
              port: 9000
            fraction:
              numerator: 1
              denominator: 10
      backendRefs:
        - name:
            value: greeter
          port: 9000
          weight: 90
        - name:
            value: greeter-canary
          port: 9000
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
| `spec.hostnames` | `[]string` |  |  |  |
| `spec.rules` | `[]KubernetesGrpcRouteRule` | yes |  |  |
| `spec.rules[].name` | `string` | yes |  |  |
| `spec.rules[].matches` | `[]KubernetesGrpcRouteMatch` |  |  |  |
| `spec.rules[].matches[].method` | `KubernetesGrpcRouteMethodMatch` |  |  |  |
| `spec.rules[].matches[].method.type` | `string` |  |  |  |
| `spec.rules[].matches[].method.service` | `string` |  |  |  |
| `spec.rules[].matches[].method.method` | `string` |  |  |  |
| `spec.rules[].matches[].headers` | `[]KubernetesGrpcRouteHeaderMatch` |  |  |  |
| `spec.rules[].matches[].headers[].type` | `string` |  |  |  |
| `spec.rules[].matches[].headers[].name` | `string` | yes |  |  |
| `spec.rules[].matches[].headers[].value` | `string` | yes |  |  |
| `spec.rules[].filters` | `[]KubernetesGrpcRouteFilter` |  |  |  |
| `spec.rules[].filters[].type` | `string` | yes |  |  |
| `spec.rules[].filters[].requestHeaderModifier` | `KubernetesGrpcRouteHeaderFilter` |  |  |  |
| `spec.rules[].filters[].requestHeaderModifier.set` | `[]KubernetesGrpcRouteHeader` |  |  |  |
| `spec.rules[].filters[].requestHeaderModifier.set[].name` | `string` | yes |  |  |
| `spec.rules[].filters[].requestHeaderModifier.set[].value` | `string` | yes |  |  |
| `spec.rules[].filters[].requestHeaderModifier.add` | `[]KubernetesGrpcRouteHeader` |  |  |  |
| `spec.rules[].filters[].requestHeaderModifier.add[].name` | `string` | yes |  |  |
| `spec.rules[].filters[].requestHeaderModifier.add[].value` | `string` | yes |  |  |
| `spec.rules[].filters[].requestHeaderModifier.remove` | `[]string` |  |  |  |
| `spec.rules[].filters[].responseHeaderModifier` | `KubernetesGrpcRouteHeaderFilter` |  |  |  |
| `spec.rules[].filters[].responseHeaderModifier.set` | `[]KubernetesGrpcRouteHeader` |  |  |  |
| `spec.rules[].filters[].responseHeaderModifier.set[].name` | `string` | yes |  |  |
| `spec.rules[].filters[].responseHeaderModifier.set[].value` | `string` | yes |  |  |
| `spec.rules[].filters[].responseHeaderModifier.add` | `[]KubernetesGrpcRouteHeader` |  |  |  |
| `spec.rules[].filters[].responseHeaderModifier.add[].name` | `string` | yes |  |  |
| `spec.rules[].filters[].responseHeaderModifier.add[].value` | `string` | yes |  |  |
| `spec.rules[].filters[].responseHeaderModifier.remove` | `[]string` |  |  |  |
| `spec.rules[].filters[].requestMirror` | `KubernetesGrpcRouteRequestMirrorFilter` |  |  |  |
| `spec.rules[].filters[].requestMirror.backendRef` | `KubernetesGatewayApiBackendObjectReference` | yes |  |  |
| `spec.rules[].filters[].requestMirror.backendRef.group` | `string` |  |  |  |
| `spec.rules[].filters[].requestMirror.backendRef.kind` | `string` | yes |  |  |
| `spec.rules[].filters[].requestMirror.backendRef.name` | `string \| valueFrom` | yes |  | KubernetesService (`status.outputs.service_name`) |
| `spec.rules[].filters[].requestMirror.backendRef.namespace` | `string` |  |  |  |
| `spec.rules[].filters[].requestMirror.backendRef.port` | `int32` |  |  |  |
| `spec.rules[].filters[].requestMirror.percent` | `int32` |  |  |  |
| `spec.rules[].filters[].requestMirror.fraction` | `KubernetesGatewayApiFraction` |  |  |  |
| `spec.rules[].filters[].requestMirror.fraction.numerator` | `int32` |  |  |  |
| `spec.rules[].filters[].requestMirror.fraction.denominator` | `int32` |  |  |  |
| `spec.rules[].filters[].extensionRef` | `KubernetesGatewayApiLocalObjectReference` |  |  |  |
| `spec.rules[].filters[].extensionRef.group` | `string` | yes |  |  |
| `spec.rules[].filters[].extensionRef.kind` | `string` | yes |  |  |
| `spec.rules[].filters[].extensionRef.name` | `string` | yes |  |  |
| `spec.rules[].backendRefs` | `[]KubernetesGrpcRouteBackendRef` |  |  |  |
| `spec.rules[].backendRefs[].group` | `string` |  |  |  |
| `spec.rules[].backendRefs[].kind` | `string` | yes |  |  |
| `spec.rules[].backendRefs[].name` | `string \| valueFrom` | yes |  | KubernetesService (`status.outputs.service_name`) |
| `spec.rules[].backendRefs[].namespace` | `string` |  |  |  |
| `spec.rules[].backendRefs[].port` | `int32` |  |  |  |
| `spec.rules[].backendRefs[].weight` | `int32` |  |  |  |
| `spec.rules[].backendRefs[].filters` | `[]KubernetesGrpcRouteFilter` |  |  |  |
| `spec.rules[].backendRefs[].filters[].type` | `string` | yes |  |  |
| `spec.rules[].backendRefs[].filters[].requestHeaderModifier` | `KubernetesGrpcRouteHeaderFilter` |  |  |  |
| `spec.rules[].backendRefs[].filters[].requestHeaderModifier.set` | `[]KubernetesGrpcRouteHeader` |  |  |  |
| `spec.rules[].backendRefs[].filters[].requestHeaderModifier.set[].name` | `string` | yes |  |  |
| `spec.rules[].backendRefs[].filters[].requestHeaderModifier.set[].value` | `string` | yes |  |  |
| `spec.rules[].backendRefs[].filters[].requestHeaderModifier.add` | `[]KubernetesGrpcRouteHeader` |  |  |  |
| `spec.rules[].backendRefs[].filters[].requestHeaderModifier.add[].name` | `string` | yes |  |  |
| `spec.rules[].backendRefs[].filters[].requestHeaderModifier.add[].value` | `string` | yes |  |  |
| `spec.rules[].backendRefs[].filters[].requestHeaderModifier.remove` | `[]string` |  |  |  |
| `spec.rules[].backendRefs[].filters[].responseHeaderModifier` | `KubernetesGrpcRouteHeaderFilter` |  |  |  |
| `spec.rules[].backendRefs[].filters[].responseHeaderModifier.set` | `[]KubernetesGrpcRouteHeader` |  |  |  |
| `spec.rules[].backendRefs[].filters[].responseHeaderModifier.set[].name` | `string` | yes |  |  |
| `spec.rules[].backendRefs[].filters[].responseHeaderModifier.set[].value` | `string` | yes |  |  |
| `spec.rules[].backendRefs[].filters[].responseHeaderModifier.add` | `[]KubernetesGrpcRouteHeader` |  |  |  |
| `spec.rules[].backendRefs[].filters[].responseHeaderModifier.add[].name` | `string` | yes |  |  |
| `spec.rules[].backendRefs[].filters[].responseHeaderModifier.add[].value` | `string` | yes |  |  |
| `spec.rules[].backendRefs[].filters[].responseHeaderModifier.remove` | `[]string` |  |  |  |
| `spec.rules[].backendRefs[].filters[].requestMirror` | `KubernetesGrpcRouteRequestMirrorFilter` |  |  |  |
| `spec.rules[].backendRefs[].filters[].requestMirror.backendRef` | `KubernetesGatewayApiBackendObjectReference` | yes |  |  |
| `spec.rules[].backendRefs[].filters[].requestMirror.backendRef.group` | `string` |  |  |  |
| `spec.rules[].backendRefs[].filters[].requestMirror.backendRef.kind` | `string` | yes |  |  |
| `spec.rules[].backendRefs[].filters[].requestMirror.backendRef.name` | `string \| valueFrom` | yes |  | KubernetesService (`status.outputs.service_name`) |
| `spec.rules[].backendRefs[].filters[].requestMirror.backendRef.namespace` | `string` |  |  |  |
| `spec.rules[].backendRefs[].filters[].requestMirror.backendRef.port` | `int32` |  |  |  |
| `spec.rules[].backendRefs[].filters[].requestMirror.percent` | `int32` |  |  |  |
| `spec.rules[].backendRefs[].filters[].requestMirror.fraction` | `KubernetesGatewayApiFraction` |  |  |  |
| `spec.rules[].backendRefs[].filters[].requestMirror.fraction.numerator` | `int32` |  |  |  |
| `spec.rules[].backendRefs[].filters[].requestMirror.fraction.denominator` | `int32` |  |  |  |
| `spec.rules[].backendRefs[].filters[].extensionRef` | `KubernetesGatewayApiLocalObjectReference` |  |  |  |
| `spec.rules[].backendRefs[].filters[].extensionRef.group` | `string` | yes |  |  |
| `spec.rules[].backendRefs[].filters[].extensionRef.kind` | `string` | yes |  |  |
| `spec.rules[].backendRefs[].filters[].extensionRef.name` | `string` | yes |  |  |

## Field Details

### spec.namespace

`string | valueFrom` · required

Namespace in which the GRPCRoute is created. Backends in other namespaces,
and Gateways in other namespaces, are subject to the usual same-namespace /
ReferenceGrant rules.

- references: KubernetesNamespace (`spec.name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesNamespace, name: <that resource's name>, fieldPath: spec.name}} -- a bare string does not parse

### spec.parentRefs

`[]KubernetesGatewayApiParentReference`

References to the parent resources (usually Gateways) this route attaches
to. Each parent must allow attachment from GRPCRoutes of this namespace.
Each reference's name defaults to a KubernetesGateway foreign key — wire
it with valueFrom in an infra chart and the route deploys after its
Gateway.

Flattened from the upstream CommonRouteSpec; the experimental
useDefaultGateways field is excluded (standard channel only).

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

`[]string`

Hostnames matched against the authority (Host) pseudo-header to select this
route. A leading wildcard label (for example "*.example.com") is a suffix
match. When empty, the route matches any hostname permitted by its parent
listeners. Upstream support level: Core.

Each value is an RFC 1123 hostname (IPs not allowed). Pattern (wildcard
prefix allowed): ^(\*\.)?[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$

- rule: {"repeated":{"maxItems":"16","items":{"string":{"minLen":"1","maxLen":"253","pattern":"^(\\*\\.)?[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$"}}}}

### spec.rules

`[]KubernetesGrpcRouteRule` · required

Routing rules: each rule matches requests, optionally applies filters, and
forwards to one or more backends. At least one rule is required. Upstream
allows up to 16 rules (and up to 64 matches per rule).

Note: upstream also carries an aggregate "total matches across all rules
must be < 128" XValidation expressed as a 16-way unrolled CEL sum. It is
intentionally not translated here (consistent with KubernetesHttpRoute): it
is a controller-enforced aggregate limit, the per-field rules<=16 /
matches<=64 bounds are enforced, and the unrolled expression is
unmaintainable. The experimental rule-name-uniqueness XValidation is also
excluded (standard channel only).

- rule: {"repeated":{"minItems":"1","maxItems":"16"}}
- rule: the RequestHeaderModifier filter may be specified at most once per rule
- rule: the ResponseHeaderModifier filter may be specified at most once per rule

### spec.rules[].name

`string` · required · optional (explicit presence)

Name of the route rule; must be unique within the route if set. Upstream
support level: Extended.

Upstream SectionName constraints: 1-253 chars, lowercase RFC 1123 subdomain.
Pattern: ^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$

- rule: {"string":{"minLen":"1","maxLen":"253","pattern":"^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$"}}

### spec.rules[].matches

`[]KubernetesGrpcRouteMatch`

Conditions for matching requests. A request matches the rule if ANY one of
the matches is satisfied (matches are ORed; conditions within a match are
ANDed). When empty, the implementation matches every gRPC request.

- rule: {"repeated":{"maxItems":"64"}}
- rule: each header match name must be unique within a match

### spec.rules[].matches[].method

`KubernetesGrpcRouteMethodMatch`

gRPC service/method matcher. When unset, all services and methods match.

- rule: one or both of 'service' or 'method' must be specified
- rule: for Exact matches, service may only contain valid characters (matching ^(?i)\.?[a-z_][a-z_0-9]*(\.[a-z_][a-z_0-9]*)*$)
- rule: for Exact matches, method may only contain valid characters (matching ^[A-Za-z_][A-Za-z_0-9]*$)

### spec.rules[].matches[].method.type

`string` · optional (explicit presence)

How to match against the service and/or method. "Exact" (default) has Core
support; "RegularExpression" is implementation-specific.

Upstream default: "Exact". Closed enum.

- rule: method match type must be either 'Exact' or 'RegularExpression'

### spec.rules[].matches[].method.service

`string` · optional (explicit presence)

Value of the gRPC service to match against (for example "helloworld.Greeter").
When empty or omitted, matches any service.

- rule: {"string":{"maxLen":"1024"}}

### spec.rules[].matches[].method.method

`string` · optional (explicit presence)

Value of the gRPC method to match against (for example "SayHello"). When
empty or omitted, matches any method.

- rule: {"string":{"maxLen":"1024"}}

### spec.rules[].matches[].headers

`[]KubernetesGrpcRouteHeaderMatch`

Request header matchers. A request must match all listed headers.

- rule: {"repeated":{"maxItems":"16"}}

### spec.rules[].matches[].headers[].type

`string` · optional (explicit presence)

How the header value is compared. "Exact" (default) has Core support;
"RegularExpression" is implementation-specific.

Upstream default: "Exact". Closed enum.

- rule: header match type must be either 'Exact' or 'RegularExpression'

### spec.rules[].matches[].headers[].name

`string` · required

Name of the gRPC header to match (case-insensitive).

Upstream HeaderName pattern: ^[A-Za-z0-9!#$%&'*+\-.^_`|~]+$ (1-256 chars).

- rule: {"required":true,"string":{"minLen":"1","maxLen":"256","pattern":"^[A-Za-z0-9!#$%&'*+\\-.^_`|~]+$"}}

### spec.rules[].matches[].headers[].value

`string` · required

Value the header is matched against.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"4096"}}

### spec.rules[].filters

`[]KubernetesGrpcRouteFilter`

Filters applied to requests matching this rule (header modification,
request mirror, or an implementation-specific extension).

- rule: {"repeated":{"maxItems":"16"}}
- rule: request_header_modifier must be set when type is 'RequestHeaderModifier' and must be unset otherwise
- rule: response_header_modifier must be set when type is 'ResponseHeaderModifier' and must be unset otherwise
- rule: request_mirror must be set when type is 'RequestMirror' and must be unset otherwise
- rule: extension_ref must be set when type is 'ExtensionRef' and must be unset otherwise

### spec.rules[].filters[].type

`string` · required

Which filter to apply. Selects exactly one configuration field below.

Closed enum (external standard exception -- matches Gateway API filter type
constants).

- rule: filter type must be one of 'RequestHeaderModifier', 'ResponseHeaderModifier', 'RequestMirror', or 'ExtensionRef'
- rule: {"required":true}

### spec.rules[].filters[].requestHeaderModifier

`KubernetesGrpcRouteHeaderFilter`

Modifies request headers before forwarding. Set when type is
"RequestHeaderModifier". Upstream support level: Core.

- rule: each header name in 'set' must be unique
- rule: each header name in 'add' must be unique

### spec.rules[].filters[].requestHeaderModifier.set

`[]KubernetesGrpcRouteHeader`

Headers to overwrite (set) on the message, replacing any existing values.

- rule: {"repeated":{"maxItems":"16"}}

### spec.rules[].filters[].requestHeaderModifier.set[].name

`string` · required

Name of the header (case-insensitive).

Upstream HeaderName pattern: ^[A-Za-z0-9!#$%&'*+\-.^_`|~]+$ (1-256 chars).

- rule: {"required":true,"string":{"minLen":"1","maxLen":"256","pattern":"^[A-Za-z0-9!#$%&'*+\\-.^_`|~]+$"}}

### spec.rules[].filters[].requestHeaderModifier.set[].value

`string` · required

Value of the header.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"4096"}}

### spec.rules[].filters[].requestHeaderModifier.add

`[]KubernetesGrpcRouteHeader`

Headers to append, adding to any existing values for the same name.

- rule: {"repeated":{"maxItems":"16"}}

### spec.rules[].filters[].requestHeaderModifier.add[].name

`string` · required

Name of the header (case-insensitive).

Upstream HeaderName pattern: ^[A-Za-z0-9!#$%&'*+\-.^_`|~]+$ (1-256 chars).

- rule: {"required":true,"string":{"minLen":"1","maxLen":"256","pattern":"^[A-Za-z0-9!#$%&'*+\\-.^_`|~]+$"}}

### spec.rules[].filters[].requestHeaderModifier.add[].value

`string` · required

Value of the header.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"4096"}}

### spec.rules[].filters[].requestHeaderModifier.remove

`[]string`

Header names to remove from the message.

- rule: {"repeated":{"maxItems":"16","unique":true}}

### spec.rules[].filters[].responseHeaderModifier

`KubernetesGrpcRouteHeaderFilter`

Modifies response headers before returning to the client. Set when type is
"ResponseHeaderModifier". Upstream support level: Extended.

- rule: each header name in 'set' must be unique
- rule: each header name in 'add' must be unique

### spec.rules[].filters[].responseHeaderModifier.set

`[]KubernetesGrpcRouteHeader`

Headers to overwrite (set) on the message, replacing any existing values.

- rule: {"repeated":{"maxItems":"16"}}

### spec.rules[].filters[].responseHeaderModifier.set[].name

`string` · required

Name of the header (case-insensitive).

Upstream HeaderName pattern: ^[A-Za-z0-9!#$%&'*+\-.^_`|~]+$ (1-256 chars).

- rule: {"required":true,"string":{"minLen":"1","maxLen":"256","pattern":"^[A-Za-z0-9!#$%&'*+\\-.^_`|~]+$"}}

### spec.rules[].filters[].responseHeaderModifier.set[].value

`string` · required

Value of the header.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"4096"}}

### spec.rules[].filters[].responseHeaderModifier.add

`[]KubernetesGrpcRouteHeader`

Headers to append, adding to any existing values for the same name.

- rule: {"repeated":{"maxItems":"16"}}

### spec.rules[].filters[].responseHeaderModifier.add[].name

`string` · required

Name of the header (case-insensitive).

Upstream HeaderName pattern: ^[A-Za-z0-9!#$%&'*+\-.^_`|~]+$ (1-256 chars).

- rule: {"required":true,"string":{"minLen":"1","maxLen":"256","pattern":"^[A-Za-z0-9!#$%&'*+\\-.^_`|~]+$"}}

### spec.rules[].filters[].responseHeaderModifier.add[].value

`string` · required

Value of the header.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"4096"}}

### spec.rules[].filters[].responseHeaderModifier.remove

`[]string`

Header names to remove from the message.

- rule: {"repeated":{"maxItems":"16","unique":true}}

### spec.rules[].filters[].requestMirror

`KubernetesGrpcRouteRequestMirrorFilter`

Mirrors the request to an additional backend (responses ignored). Set when
type is "RequestMirror". Upstream support level: Extended.

- rule: specify either percent or fraction for the mirror sampling rate, but not both

### spec.rules[].filters[].requestMirror.backendRef

`KubernetesGatewayApiBackendObjectReference` · required

Backend that mirrored requests are sent to. Required. Upstream support
level: Extended for Kubernetes Service. The reference's name defaults to
a KubernetesService foreign key.

- rule: {"required":true}
- rule: port must be specified when referencing a core API Service (group is empty and kind is 'Service')

### spec.rules[].filters[].requestMirror.backendRef.group

`string` · optional (explicit presence)

Group of the referent. Empty string infers the core API group.

Upstream default: "" (core API group)
Group pattern: empty or an RFC 1123 subdomain (max 253).

- rule: {"string":{"maxLen":"253","pattern":"^$|^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$"}}

### spec.rules[].filters[].requestMirror.backendRef.kind

`string` · required · optional (explicit presence)

Kind of the referent.

Upstream default: "Service"
Kind pattern: 1-63 chars, ^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$

- rule: {"string":{"minLen":"1","maxLen":"63","pattern":"^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$"}}

### spec.rules[].filters[].requestMirror.backendRef.name

`string | valueFrom` · required

Name of the referent. Defaults to a KubernetesService foreign key: in an
infra chart, wire it with valueFrom against the backend Service and the
route deploys after its backend. When the backend is not a
Planton-managed Service (a custom backend kind, or a Service created
outside Planton), pass the literal name with `value:`.

- references: KubernetesService (`status.outputs.service_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesService, name: <that resource's name>, fieldPath: status.outputs.service_name}} -- a bare string does not parse

### spec.rules[].filters[].requestMirror.backendRef.namespace

`string` · optional (explicit presence)

Namespace of the backend. When unspecified, the local namespace is inferred.
Cross-namespace references require a ReferenceGrant.

### spec.rules[].filters[].requestMirror.backendRef.port

`int32` · optional (explicit presence)

Destination port number. Required when the referent is a core Service.

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.rules[].filters[].requestMirror.percent

`int32` · optional (explicit presence)

Percentage of requests to mirror (0-100). When neither percent nor fraction
is set, 100% are mirrored.

- rule: {"int32":{"lte":100,"gte":0}}

### spec.rules[].filters[].requestMirror.fraction

`KubernetesGatewayApiFraction`

Fraction of requests to mirror. When neither percent nor fraction is set,
100% are mirrored.

- rule: numerator must be less than or equal to denominator

### spec.rules[].filters[].requestMirror.fraction.numerator

`int32`

Numerator of the fraction. Must be non-negative.

- rule: {"int32":{"gte":0}}

### spec.rules[].filters[].requestMirror.fraction.denominator

`int32` · optional (explicit presence)

Denominator of the fraction.

Upstream default: 100

- rule: {"int32":{"gte":1}}

### spec.rules[].filters[].extensionRef

`KubernetesGatewayApiLocalObjectReference`

Implementation-specific extension filter. Set when type is "ExtensionRef".
Must not be used for core or extended filters. Upstream support:
Implementation-specific.

### spec.rules[].filters[].extensionRef.group

`string` · required · optional (explicit presence)

Group of the referent (e.g., "gateway.networking.k8s.io").
Empty string infers the core API group.

Upstream requires the KEY to be present, but its Group type explicitly
allows the empty value -- so this is a proto3 `optional` string with a
presence `required` rule: it must be SET (and is therefore emitted to
the rendered CR, whose CRD rejects a missing key) but may be empty.
The `optional` is what keeps the projection faithful: protojson omits
unset proto3 scalars, so a non-optional empty-string group would be
dropped from the manifest and rejected by the API server.

Group pattern: empty or an RFC 1123 subdomain (max 253).

- rule: {"required":true,"string":{"maxLen":"253","pattern":"^$|^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$"}}

### spec.rules[].filters[].extensionRef.kind

`string` · required

Kind of the referent (e.g., "HTTPRoute" or "Service").

Upstream models Kind as a required value type.
Kind pattern: 1-63 chars, ^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$

- rule: {"required":true,"string":{"minLen":"1","maxLen":"63","pattern":"^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$"}}

### spec.rules[].filters[].extensionRef.name

`string` · required

Name of the referent.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"253"}}

### spec.rules[].backendRefs

`[]KubernetesGrpcRouteBackendRef`

Backends matching requests are forwarded to. When multiple backends are
given, traffic is split by weight. Upstream support level: Core for
Kubernetes Service.

- rule: {"repeated":{"maxItems":"16"}}

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
Planton-managed Service, pass the literal name with `value:`.

- references: KubernetesService (`status.outputs.service_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesService, name: <that resource's name>, fieldPath: status.outputs.service_name}} -- a bare string does not parse

### spec.rules[].backendRefs[].namespace

`string` · optional (explicit presence)

Namespace of the backend. When unspecified, the route's namespace is
inferred. Cross-namespace references require a ReferenceGrant.

### spec.rules[].backendRefs[].port

`int32` · optional (explicit presence)

Destination port number. Required when the referent is a core Service.

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.rules[].backendRefs[].weight

`int32` · optional (explicit presence)

Relative weight for traffic splitting, computed as weight/(sum of weights).

Upstream default: 1.

- rule: {"int32":{"lte":1000000,"gte":0}}

### spec.rules[].backendRefs[].filters

`[]KubernetesGrpcRouteFilter`

Filters applied only when forwarding to this backend. Upstream support
level: Implementation-specific (prefer rule-level filters for portability).

- rule: {"repeated":{"maxItems":"16"}}
- rule: request_header_modifier must be set when type is 'RequestHeaderModifier' and must be unset otherwise
- rule: response_header_modifier must be set when type is 'ResponseHeaderModifier' and must be unset otherwise
- rule: request_mirror must be set when type is 'RequestMirror' and must be unset otherwise
- rule: extension_ref must be set when type is 'ExtensionRef' and must be unset otherwise

### spec.rules[].backendRefs[].filters[].type

`string` · required

Which filter to apply. Selects exactly one configuration field below.

Closed enum (external standard exception -- matches Gateway API filter type
constants).

- rule: filter type must be one of 'RequestHeaderModifier', 'ResponseHeaderModifier', 'RequestMirror', or 'ExtensionRef'
- rule: {"required":true}

### spec.rules[].backendRefs[].filters[].requestHeaderModifier

`KubernetesGrpcRouteHeaderFilter`

Modifies request headers before forwarding. Set when type is
"RequestHeaderModifier". Upstream support level: Core.

- rule: each header name in 'set' must be unique
- rule: each header name in 'add' must be unique

### spec.rules[].backendRefs[].filters[].requestHeaderModifier.set

`[]KubernetesGrpcRouteHeader`

Headers to overwrite (set) on the message, replacing any existing values.

- rule: {"repeated":{"maxItems":"16"}}

### spec.rules[].backendRefs[].filters[].requestHeaderModifier.set[].name

`string` · required

Name of the header (case-insensitive).

Upstream HeaderName pattern: ^[A-Za-z0-9!#$%&'*+\-.^_`|~]+$ (1-256 chars).

- rule: {"required":true,"string":{"minLen":"1","maxLen":"256","pattern":"^[A-Za-z0-9!#$%&'*+\\-.^_`|~]+$"}}

### spec.rules[].backendRefs[].filters[].requestHeaderModifier.set[].value

`string` · required

Value of the header.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"4096"}}

### spec.rules[].backendRefs[].filters[].requestHeaderModifier.add

`[]KubernetesGrpcRouteHeader`

Headers to append, adding to any existing values for the same name.

- rule: {"repeated":{"maxItems":"16"}}

### spec.rules[].backendRefs[].filters[].requestHeaderModifier.add[].name

`string` · required

Name of the header (case-insensitive).

Upstream HeaderName pattern: ^[A-Za-z0-9!#$%&'*+\-.^_`|~]+$ (1-256 chars).

- rule: {"required":true,"string":{"minLen":"1","maxLen":"256","pattern":"^[A-Za-z0-9!#$%&'*+\\-.^_`|~]+$"}}

### spec.rules[].backendRefs[].filters[].requestHeaderModifier.add[].value

`string` · required

Value of the header.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"4096"}}

### spec.rules[].backendRefs[].filters[].requestHeaderModifier.remove

`[]string`

Header names to remove from the message.

- rule: {"repeated":{"maxItems":"16","unique":true}}

### spec.rules[].backendRefs[].filters[].responseHeaderModifier

`KubernetesGrpcRouteHeaderFilter`

Modifies response headers before returning to the client. Set when type is
"ResponseHeaderModifier". Upstream support level: Extended.

- rule: each header name in 'set' must be unique
- rule: each header name in 'add' must be unique

### spec.rules[].backendRefs[].filters[].responseHeaderModifier.set

`[]KubernetesGrpcRouteHeader`

Headers to overwrite (set) on the message, replacing any existing values.

- rule: {"repeated":{"maxItems":"16"}}

### spec.rules[].backendRefs[].filters[].responseHeaderModifier.set[].name

`string` · required

Name of the header (case-insensitive).

Upstream HeaderName pattern: ^[A-Za-z0-9!#$%&'*+\-.^_`|~]+$ (1-256 chars).

- rule: {"required":true,"string":{"minLen":"1","maxLen":"256","pattern":"^[A-Za-z0-9!#$%&'*+\\-.^_`|~]+$"}}

### spec.rules[].backendRefs[].filters[].responseHeaderModifier.set[].value

`string` · required

Value of the header.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"4096"}}

### spec.rules[].backendRefs[].filters[].responseHeaderModifier.add

`[]KubernetesGrpcRouteHeader`

Headers to append, adding to any existing values for the same name.

- rule: {"repeated":{"maxItems":"16"}}

### spec.rules[].backendRefs[].filters[].responseHeaderModifier.add[].name

`string` · required

Name of the header (case-insensitive).

Upstream HeaderName pattern: ^[A-Za-z0-9!#$%&'*+\-.^_`|~]+$ (1-256 chars).

- rule: {"required":true,"string":{"minLen":"1","maxLen":"256","pattern":"^[A-Za-z0-9!#$%&'*+\\-.^_`|~]+$"}}

### spec.rules[].backendRefs[].filters[].responseHeaderModifier.add[].value

`string` · required

Value of the header.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"4096"}}

### spec.rules[].backendRefs[].filters[].responseHeaderModifier.remove

`[]string`

Header names to remove from the message.

- rule: {"repeated":{"maxItems":"16","unique":true}}

### spec.rules[].backendRefs[].filters[].requestMirror

`KubernetesGrpcRouteRequestMirrorFilter`

Mirrors the request to an additional backend (responses ignored). Set when
type is "RequestMirror". Upstream support level: Extended.

- rule: specify either percent or fraction for the mirror sampling rate, but not both

### spec.rules[].backendRefs[].filters[].requestMirror.backendRef

`KubernetesGatewayApiBackendObjectReference` · required

Backend that mirrored requests are sent to. Required. Upstream support
level: Extended for Kubernetes Service. The reference's name defaults to
a KubernetesService foreign key.

- rule: {"required":true}
- rule: port must be specified when referencing a core API Service (group is empty and kind is 'Service')

### spec.rules[].backendRefs[].filters[].requestMirror.backendRef.group

`string` · optional (explicit presence)

Group of the referent. Empty string infers the core API group.

Upstream default: "" (core API group)
Group pattern: empty or an RFC 1123 subdomain (max 253).

- rule: {"string":{"maxLen":"253","pattern":"^$|^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$"}}

### spec.rules[].backendRefs[].filters[].requestMirror.backendRef.kind

`string` · required · optional (explicit presence)

Kind of the referent.

Upstream default: "Service"
Kind pattern: 1-63 chars, ^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$

- rule: {"string":{"minLen":"1","maxLen":"63","pattern":"^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$"}}

### spec.rules[].backendRefs[].filters[].requestMirror.backendRef.name

`string | valueFrom` · required

Name of the referent. Defaults to a KubernetesService foreign key: in an
infra chart, wire it with valueFrom against the backend Service and the
route deploys after its backend. When the backend is not a
Planton-managed Service (a custom backend kind, or a Service created
outside Planton), pass the literal name with `value:`.

- references: KubernetesService (`status.outputs.service_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesService, name: <that resource's name>, fieldPath: status.outputs.service_name}} -- a bare string does not parse

### spec.rules[].backendRefs[].filters[].requestMirror.backendRef.namespace

`string` · optional (explicit presence)

Namespace of the backend. When unspecified, the local namespace is inferred.
Cross-namespace references require a ReferenceGrant.

### spec.rules[].backendRefs[].filters[].requestMirror.backendRef.port

`int32` · optional (explicit presence)

Destination port number. Required when the referent is a core Service.

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.rules[].backendRefs[].filters[].requestMirror.percent

`int32` · optional (explicit presence)

Percentage of requests to mirror (0-100). When neither percent nor fraction
is set, 100% are mirrored.

- rule: {"int32":{"lte":100,"gte":0}}

### spec.rules[].backendRefs[].filters[].requestMirror.fraction

`KubernetesGatewayApiFraction`

Fraction of requests to mirror. When neither percent nor fraction is set,
100% are mirrored.

- rule: numerator must be less than or equal to denominator

### spec.rules[].backendRefs[].filters[].requestMirror.fraction.numerator

`int32`

Numerator of the fraction. Must be non-negative.

- rule: {"int32":{"gte":0}}

### spec.rules[].backendRefs[].filters[].requestMirror.fraction.denominator

`int32` · optional (explicit presence)

Denominator of the fraction.

Upstream default: 100

- rule: {"int32":{"gte":1}}

### spec.rules[].backendRefs[].filters[].extensionRef

`KubernetesGatewayApiLocalObjectReference`

Implementation-specific extension filter. Set when type is "ExtensionRef".
Must not be used for core or extended filters. Upstream support:
Implementation-specific.

### spec.rules[].backendRefs[].filters[].extensionRef.group

`string` · required · optional (explicit presence)

Group of the referent (e.g., "gateway.networking.k8s.io").
Empty string infers the core API group.

Upstream requires the KEY to be present, but its Group type explicitly
allows the empty value -- so this is a proto3 `optional` string with a
presence `required` rule: it must be SET (and is therefore emitted to
the rendered CR, whose CRD rejects a missing key) but may be empty.
The `optional` is what keeps the projection faithful: protojson omits
unset proto3 scalars, so a non-optional empty-string group would be
dropped from the manifest and rejected by the API server.

Group pattern: empty or an RFC 1123 subdomain (max 253).

- rule: {"required":true,"string":{"maxLen":"253","pattern":"^$|^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$"}}

### spec.rules[].backendRefs[].filters[].extensionRef.kind

`string` · required

Kind of the referent (e.g., "HTTPRoute" or "Service").

Upstream models Kind as a required value type.
Kind pattern: 1-63 chars, ^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$

- rule: {"required":true,"string":{"minLen":"1","maxLen":"63","pattern":"^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$"}}

### spec.rules[].backendRefs[].filters[].extensionRef.name

`string` · required

Name of the referent.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"253"}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: KubernetesGrpcRoute, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.route_name` | `string` | Name of the created GRPCRoute (equals metadata.name). In InfraCharts this orders the route after the Gateway and backends it references. |
| `status.outputs.namespace` | `string` | Namespace the GRPCRoute was created in (the resolved spec.namespace). Cross-namespace parent and backend references from this route are subject to ReferenceGrant rules relative to this value. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.namespace` | KubernetesNamespace | `spec.name` |
| `spec.parentRefs[].name` | KubernetesGateway | `status.outputs.gateway_name` |
| `spec.rules[].filters[].requestMirror.backendRef.name` | KubernetesService | `status.outputs.service_name` |
| `spec.rules[].backendRefs[].name` | KubernetesService | `status.outputs.service_name` |
| `spec.rules[].backendRefs[].filters[].requestMirror.backendRef.name` | KubernetesService | `status.outputs.service_name` |

## See Also

- [Overview](../README.md)
