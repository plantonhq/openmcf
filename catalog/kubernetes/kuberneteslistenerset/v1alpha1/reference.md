# KubernetesListenerSet

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `kubernetes.planton.dev/v1alpha1`

KubernetesListenerSetSpec defines a Kubernetes Gateway API ListenerSet: a
namespaced set of additional listeners merged into an existing Gateway.
ListenerSets let teams attach their own listeners (ports, hostnames, TLS
certificates) to a centrally managed Gateway without editing the Gateway
itself — the per-tenant/per-team delegation model for shared gateways.

The parent Gateway must explicitly opt in through its allowed_listeners
configuration (Gateways allow NO ListenerSet attachment by default). Routes
can attach to a ListenerSet by naming it as a parentRef, optionally
targeting a specific listener with sectionName.

100% fidelity with the upstream Gateway API v1.6.1 ListenerSetSpec
(kubernetes-sigs/gateway-api apis/v1/listenerset_types.go), standard
channel, served as gateway.networking.k8s.io/v1 (ListenerSet graduated to
the standard channel in the v1.5 release).

The listener entries reuse the shared Gateway API listener building blocks
(TLS config, allowed routes) — upstream uses the same types for Gateway
listeners and ListenerSet entries, so the shared messages make drift
between the two kinds structurally impossible.

## Example

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesListenerSet
metadata:
  name: test-listener-set
spec:
  namespace:
    value: test-namespace
  parentRef:
    name:
      value: test-gateway
    namespace: gateway-namespace
  listeners:
    - name: tenant-https
      hostname: tenant.example.com
      port: 443
      protocol: HTTPS
      tls:
        mode: Terminate
        certificateRefs:
          - name:
              value: tenant-tls-cert
        options:
          example.com/minimum-tls-version: "1.3"
      allowedRoutes:
        namespaces:
          from: Selector
          selector:
            matchLabels:
              team: tenant
            matchExpressions:
              - key: environment
                operator: In
                values:
                  - staging
                  - production
        kinds:
          - kind: HTTPRoute
    - name: tenant-tls-passthrough
      hostname: db.tenant.example.com
      port: 8443
      protocol: TLS
      tls:
        mode: Passthrough
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.namespace` | `string \| valueFrom` | yes |  | KubernetesNamespace (`spec.name`) |
| `spec.parentRef` | `KubernetesGatewayApiParentGatewayReference` | yes |  |  |
| `spec.parentRef.group` | `string` |  |  |  |
| `spec.parentRef.kind` | `string` | yes |  |  |
| `spec.parentRef.namespace` | `string` | yes |  |  |
| `spec.parentRef.name` | `string \| valueFrom` | yes |  | KubernetesGateway (`status.outputs.gateway_name`) |
| `spec.listeners` | `[]KubernetesListenerSetListener` | yes |  |  |
| `spec.listeners[].name` | `string` | yes |  |  |
| `spec.listeners[].hostname` | `string` | yes |  |  |
| `spec.listeners[].port` | `int32` | yes |  |  |
| `spec.listeners[].protocol` | `string` | yes |  |  |
| `spec.listeners[].tls` | `KubernetesGatewayApiListenerTlsConfig` |  |  |  |
| `spec.listeners[].tls.mode` | `string` |  |  |  |
| `spec.listeners[].tls.certificateRefs` | `[]KubernetesGatewayApiSecretObjectReference` |  |  |  |
| `spec.listeners[].tls.certificateRefs[].group` | `string` |  |  |  |
| `spec.listeners[].tls.certificateRefs[].kind` | `string` | yes |  |  |
| `spec.listeners[].tls.certificateRefs[].name` | `string \| valueFrom` | yes |  | KubernetesSecret (`status.outputs.secret_name`) |
| `spec.listeners[].tls.certificateRefs[].namespace` | `string` |  |  |  |
| `spec.listeners[].tls.options` | `map<string, string>` |  |  |  |
| `spec.listeners[].allowedRoutes` | `KubernetesGatewayApiAllowedRoutes` |  |  |  |
| `spec.listeners[].allowedRoutes.namespaces` | `KubernetesGatewayApiRouteNamespaces` |  |  |  |
| `spec.listeners[].allowedRoutes.namespaces.from` | `string` |  |  |  |
| `spec.listeners[].allowedRoutes.namespaces.selector` | `KubernetesGatewayApiLabelSelector` |  |  |  |
| `spec.listeners[].allowedRoutes.namespaces.selector.matchLabels` | `map<string, string>` |  |  |  |
| `spec.listeners[].allowedRoutes.namespaces.selector.matchExpressions` | `[]KubernetesGatewayApiLabelSelectorRequirement` |  |  |  |
| `spec.listeners[].allowedRoutes.namespaces.selector.matchExpressions[].key` | `string` | yes |  |  |
| `spec.listeners[].allowedRoutes.namespaces.selector.matchExpressions[].operator` | `string` | yes |  |  |
| `spec.listeners[].allowedRoutes.namespaces.selector.matchExpressions[].values` | `[]string` |  |  |  |
| `spec.listeners[].allowedRoutes.kinds` | `[]KubernetesGatewayApiRouteGroupKind` |  |  |  |
| `spec.listeners[].allowedRoutes.kinds[].group` | `string` |  |  |  |
| `spec.listeners[].allowedRoutes.kinds[].kind` | `string` | yes |  |  |

## Field Details

### spec.namespace

`string | valueFrom` · required

Namespace in which the ListenerSet is created. The parent Gateway must
allow attachment from this namespace (allowed_listeners.namespaces); a
ListenerSet can reference TLS Secrets in its OWN namespace without a
ReferenceGrant, but not the parent Gateway's.

- references: KubernetesNamespace (`spec.name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesNamespace, name: <that resource's name>, fieldPath: spec.name}} -- a bare string does not parse

### spec.parentRef

`KubernetesGatewayApiParentGatewayReference` · required

Reference to the Gateway these listeners attach to. The Gateway merges
listeners from itself and all attached ListenerSets (its own listeners
take precedence, then ListenerSets by creation time, then alphabetical
namespace/name order). The reference's name defaults to a
KubernetesGateway foreign key — wire it with valueFrom in an infra chart
and the ListenerSet deploys after its Gateway.

- rule: {"required":true}

### spec.parentRef.group

`string` · optional (explicit presence)

Group of the referent.

Upstream default: "gateway.networking.k8s.io"
Group pattern: empty or an RFC 1123 subdomain (max 253).

- rule: {"string":{"maxLen":"253","pattern":"^$|^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$"}}

### spec.parentRef.kind

`string` · required · optional (explicit presence)

Kind of the referent.

Upstream default: "Gateway"
Kind pattern: 1-63 chars, ^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$

- rule: {"string":{"minLen":"1","maxLen":"63","pattern":"^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$"}}

### spec.parentRef.namespace

`string` · required · optional (explicit presence)

Namespace of the referent. When unspecified, the ListenerSet's own
namespace is assumed. The parent Gateway must allow attachment from the
ListenerSet's namespace via its allowed_listeners configuration.

- rule: {"string":{"minLen":"1","maxLen":"63","pattern":"^[a-z0-9]([-a-z0-9]*[a-z0-9])?$"}}

### spec.parentRef.name

`string | valueFrom` · required

Name of the parent Gateway. Defaults to a KubernetesGateway foreign key:
in an infra chart, wire it with valueFrom against the Gateway resource so
the ListenerSet deploys after its Gateway. Pass the literal name with
`value:` when the Gateway is not Planton-managed.

- references: KubernetesGateway (`status.outputs.gateway_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesGateway, name: <that resource's name>, fieldPath: status.outputs.gateway_name}} -- a bare string does not parse

### spec.listeners

`[]KubernetesListenerSetListener` · required

Listeners merged into the parent Gateway. Each listener defines a port,
protocol, and optional TLS and hostname — the same shape as a Gateway's
own listeners. At least one listener is required; up to 64.

- rule: {"repeated":{"minItems":"1","maxItems":"64"}}
- rule: tls must not be set when protocol is HTTP, TCP, or UDP
- rule: tls mode must be Terminate (or left unset) when protocol is HTTPS
- rule: tls and tls.mode must be set when protocol is TLS
- rule: hostname must not be set when protocol is TCP or UDP

### spec.listeners[].name

`string` · required

Name of the listener; must be unique within this ListenerSet. Routes
attach to a specific listener by this name (parentRef.sectionName against
the ListenerSet). It does NOT need to be unique across the parent Gateway
and its other ListenerSets.

Upstream SectionName constraints: 1-253 chars; lowercase RFC 1123
subdomain. Pattern: ^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$

- rule: {"required":true,"string":{"minLen":"1","maxLen":"253","pattern":"^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$"}}

### spec.listeners[].hostname

`string` · required · optional (explicit presence)

Virtual hostname to match for protocols that support it (HTTP, HTTPS, TLS).
When unset, all hostnames match. A leading wildcard label (for example
"*.example.com") is a suffix match. Ignored for TCP and UDP. For TLS
listeners, the SNI-to-certificate association is derived from this value.

Upstream Hostname constraints: 1-253 chars. Pattern (wildcard prefix
allowed): ^(\*\.)?[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$

- rule: {"string":{"minLen":"1","maxLen":"253","pattern":"^(\\*\\.)?[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$"}}

### spec.listeners[].port

`int32` · required

Network port this listener binds. Multiple listeners may share a port
subject to the distinctness rules above.

- rule: {"required":true,"int32":{"lte":65535,"gte":1}}

### spec.listeners[].protocol

`string` · required

Application protocol the listener accepts. Core values are HTTP, HTTPS,
TLS, TCP, and UDP. Implementation-specific protocols use a domain-prefixed
path such as "example.com/my-protocol".

Upstream ProtocolType is an open set validated by pattern (not a closed
enum), so custom domain-prefixed protocols remain valid.
Pattern: ^[a-zA-Z0-9]([-a-zA-Z0-9]*[a-zA-Z0-9])?$|<domain>/<path>

- rule: {"required":true,"string":{"minLen":"1","maxLen":"255","pattern":"^[a-zA-Z0-9]([-a-zA-Z0-9]*[a-zA-Z0-9])?$|[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*\\/[A-Za-z0-9]+$"}}

### spec.listeners[].tls

`KubernetesGatewayApiListenerTlsConfig`

TLS configuration for this listener. Required for HTTPS and TLS protocols;
must be absent for HTTP, TCP, and UDP (enforced above). Certificate
references default to KubernetesSecret foreign keys — wire valueFrom
against a KubernetesCertificate's status.outputs.secret_name to terminate
with cert-manager-issued material. Shared Gateway API type (also used by
KubernetesGateway listeners).

- rule: certificate_refs or options must be provided when tls mode is Terminate (the default)

### spec.listeners[].tls.mode

`string` · optional (explicit presence)

How the listener handles the client's TLS session. "Terminate" (default)
decrypts at the Gateway and requires a certificate. "Passthrough" forwards
the encrypted stream untouched (TLSRoute only) and ignores certificate_refs.

Upstream default: "Terminate". Closed enum: Terminate | Passthrough.

- rule: tls mode must be either 'Terminate' or 'Passthrough'

### spec.listeners[].tls.certificateRefs

`[]KubernetesGatewayApiSecretObjectReference`

References to Kubernetes Secrets holding the TLS certificate/key used to
terminate the listener. A single reference to a kubernetes.io/tls Secret
has Core support; multiple references are implementation-specific.

Each reference's name is a KubernetesSecret foreign key; the Secret is
typically produced by a KubernetesCertificate (cert-manager) — wire
valueFrom against the Certificate's status.outputs.secret_name so the
listener terminates with the issued certificate. See
KubernetesGatewayApiSecretObjectReference.

- rule: {"repeated":{"maxItems":"64"}}

### spec.listeners[].tls.certificateRefs[].group

`string` · optional (explicit presence)

Group of the referent. Empty string infers the core API group.

Upstream default: "" (core API group)
Group pattern: empty or an RFC 1123 subdomain (max 253).

- rule: {"string":{"maxLen":"253","pattern":"^$|^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$"}}

### spec.listeners[].tls.certificateRefs[].kind

`string` · required · optional (explicit presence)

Kind of the referent.

Upstream default: "Secret"
Kind pattern: 1-63 chars, ^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$

- rule: {"string":{"minLen":"1","maxLen":"63","pattern":"^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$"}}

### spec.listeners[].tls.certificateRefs[].name

`string | valueFrom` · required

Name of the referent. Defaults to a KubernetesSecret foreign key. The
Secret is typically produced by cert-manager: reference a
KubernetesCertificate's status.outputs.secret_name with valueFrom to wire
a certificate the moment it is issued, or a KubernetesSecret directly for
externally provisioned material. Pass the literal name with `value:` when
the Secret is not Planton-managed.

- references: KubernetesSecret (`status.outputs.secret_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesSecret, name: <that resource's name>, fieldPath: status.outputs.secret_name}} -- a bare string does not parse

### spec.listeners[].tls.certificateRefs[].namespace

`string` · optional (explicit presence)

Namespace of the referenced object. When unspecified, the local namespace
is inferred. Cross-namespace references require a ReferenceGrant.

### spec.listeners[].tls.options

`map<string, string>`

Implementation-specific TLS options (for example minimum TLS version or
cipher suites). Keys should be domain-prefixed to avoid ambiguity.

Upstream key/value: Gateway API AnnotationKey (1-253) / AnnotationValue
(0-4096).

- rule: {"map":{"maxPairs":"16","keys":{"string":{"minLen":"1","maxLen":"253","pattern":"^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?([A-Za-z0-9][-A-Za-z0-9_.]{0,61})?[A-Za-z0-9]$"}},"values":{"string":{"maxLen":"4096"}}}}

### spec.listeners[].allowedRoutes

`KubernetesGatewayApiAllowedRoutes`

Restricts which Routes (by kind and namespace) may attach to this
listener. When unset, the controller selects routes by listener protocol
and — for ListenerSets — the ListenerSet's OWN namespace by default.
Shared Gateway API type (also used by KubernetesGateway listeners).

### spec.listeners[].allowedRoutes.namespaces

`KubernetesGatewayApiRouteNamespaces`

Namespaces from which Routes may attach. Defaults to the Gateway's own
namespace.

### spec.listeners[].allowedRoutes.namespaces.from

`string` · optional (explicit presence)

Where Routes are selected from: "All" (any namespace), "Selector"
(namespaces matching the selector), or "Same" (the Gateway's namespace).

Upstream default: "Same". Closed enum: All | Selector | Same.

- rule: from must be one of 'All', 'Selector', or 'Same'

### spec.listeners[].allowedRoutes.namespaces.selector

`KubernetesGatewayApiLabelSelector`

Namespace label selector. Required (and only honored) when `from` is
"Selector".

### spec.listeners[].allowedRoutes.namespaces.selector.matchLabels

`map<string, string>`

Map of {key,value} pairs. A namespace matches if it carries all of these
labels.

### spec.listeners[].allowedRoutes.namespaces.selector.matchExpressions

`[]KubernetesGatewayApiLabelSelectorRequirement`

List of label selector requirements, ANDed together with match_labels.

### spec.listeners[].allowedRoutes.namespaces.selector.matchExpressions[].key

`string` · required

Label key the requirement applies to.

- rule: {"required":true}

### spec.listeners[].allowedRoutes.namespaces.selector.matchExpressions[].operator

`string` · required

Relationship between the key and values: "In", "NotIn", "Exists", or
"DoesNotExist". For "Exists"/"DoesNotExist", values must be empty.

- rule: operator must be one of 'In', 'NotIn', 'Exists', or 'DoesNotExist'
- rule: {"required":true}

### spec.listeners[].allowedRoutes.namespaces.selector.matchExpressions[].values

`[]string`

Values for the requirement. Must be non-empty for "In"/"NotIn" and empty
for "Exists"/"DoesNotExist".

### spec.listeners[].allowedRoutes.kinds

`[]KubernetesGatewayApiRouteGroupKind`

Route kinds allowed to bind to this listener. When empty, allowed kinds
are inferred from the listener protocol.

- rule: {"repeated":{"maxItems":"8"}}

### spec.listeners[].allowedRoutes.kinds[].group

`string` · optional (explicit presence)

API group of the Route. Empty string selects the core API group.

Upstream default: "gateway.networking.k8s.io".
Group pattern: empty or an RFC 1123 subdomain (max 253).

- rule: {"string":{"maxLen":"253","pattern":"^$|^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$"}}

### spec.listeners[].allowedRoutes.kinds[].kind

`string` · required

Kind of the Route (for example "HTTPRoute", "GRPCRoute", "TLSRoute").

Kind pattern: 1-63 chars, ^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$

- rule: {"required":true,"string":{"minLen":"1","maxLen":"63","pattern":"^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$"}}

## Validation Rules

- `listener_set.listener_names_unique`: each listener name must be unique within the ListenerSet
- `listener_set.listener_port_protocol_hostname_unique`: each listener must have a unique combination of port, protocol, and hostname

## Outputs

Reference an output from another manifest as `valueFrom: {kind: KubernetesListenerSet, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.listener_set_name` | `string` | Name of the created ListenerSet (equals metadata.name). Routes reference this name in their parent_refs (kind: ListenerSet) to attach to the merged listeners. |
| `status.outputs.namespace` | `string` | Namespace the ListenerSet was created in (the resolved spec.namespace). |
| `status.outputs.gateway_name` | `string` | Name of the parent Gateway the listeners attach to (the resolved spec.parent_ref.name). Exposed so downstream resources can follow the chain to the Gateway without re-resolving the reference. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.namespace` | KubernetesNamespace | `spec.name` |
| `spec.parentRef.name` | KubernetesGateway | `status.outputs.gateway_name` |
| `spec.listeners[].tls.certificateRefs[].name` | KubernetesSecret | `status.outputs.secret_name` |

## See Also

- [Overview](../README.md)
