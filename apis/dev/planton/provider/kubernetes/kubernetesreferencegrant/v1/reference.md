# KubernetesReferenceGrant

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `kubernetes.planton.dev/v1`

KubernetesReferenceGrantSpec defines a Kubernetes Gateway API ReferenceGrant:
a namespaced authorization that permits resources in OTHER namespaces to
reference specified kinds of resources in THIS grant's namespace. It is a form
of runtime verification -- all cross-namespace references in the Gateway API
(except cross-namespace Gateway-route attachment) require a ReferenceGrant in
the referenced (the "to") namespace.

100% fidelity with the upstream Gateway API v1.6.1 ReferenceGrantSpec
(kubernetes-sigs/gateway-api apis/v1/referencegrant_types.go), standard channel,
served as gateway.networking.k8s.io/v1. Upstream spec fields are flattened after
the Planton namespaced envelope (namespace).

Unlike the route kinds, ReferenceGrant has NO parent_refs/backend_refs and does
not reuse the shared gateway_api.proto reference types: its from/to entries are
trust assertions about KINDS of resources, not pointers to specific objects.
Upstream deliberately omits a status subresource on ReferenceGrant, so there is
no controller-managed status to surface.

## Example

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesReferenceGrant
metadata:
  name: test-reference-grant
spec:
  namespace:
    value: backend-namespace
  from:
    - group: gateway.networking.k8s.io
      kind: HTTPRoute
      namespace: frontend-namespace
    - group: gateway.networking.k8s.io
      kind: Gateway
      namespace: gateway-namespace
  to:
    - group: ""
      kind: Service
    - group: ""
      kind: Secret
      name: shared-tls-cert
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.namespace` | `string \| valueFrom` | yes |  | KubernetesNamespace (`spec.name`) |
| `spec.from` | `[]KubernetesReferenceGrantFrom` | yes |  |  |
| `spec.from[].group` | `string` | yes |  |  |
| `spec.from[].kind` | `string` | yes |  |  |
| `spec.from[].namespace` | `string` | yes |  |  |
| `spec.to` | `[]KubernetesReferenceGrantTo` | yes |  |  |
| `spec.to[].group` | `string` | yes |  |  |
| `spec.to[].kind` | `string` | yes |  |  |
| `spec.to[].name` | `string` | yes |  |  |

## Field Details

### spec.namespace

`string | valueFrom` · required

Namespace the ReferenceGrant is created in. This is the namespace being
referenced INTO (the "to" side): the grant lives alongside the resources it
authorizes inbound references to, and revoking the grant revokes that access.

- references: KubernetesNamespace (`spec.name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesNamespace, name: <that resource's name>, fieldPath: spec.name}} -- a bare string does not parse

### spec.from

`[]KubernetesReferenceGrantFrom` · required

The trusted namespaces and kinds that may reference the resources described in
`to`. Entries are combined with OR -- each is an additional place inbound
references are valid from. At least one entry is required (upstream MinItems=1),
at most 16. Upstream support level: Core.

- rule: {"repeated":{"minItems":"1","maxItems":"16"}}

### spec.from[].group

`string` · required · optional (explicit presence)

Group of the referent. The empty string ("") infers the Kubernetes core API
group. Upstream requires the KEY to be present, but its Group type explicitly
allows the empty value -- so this is a proto3 `optional` string with a presence
`required` rule: it must be SET (and is therefore emitted to the CRD, which
rejects a missing key) but may be empty. The `optional` is what lets
ReferenceGrant be a faithful kubernetes_manifest projection: protojson omits
unset proto3 scalars, so a non-optional empty-string group would be dropped
and rejected by the API server. Pattern: empty or an RFC 1123 subdomain (253).

INFRA-CHART COMPOSABILITY: group is a trust assertion about a KIND of
resource, not a pointer to a specific Planton resource instance -- do NOT add a
metadata.relationships hint for it.

- rule: {"required":true,"string":{"maxLen":"253","pattern":"^$|^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$"}}

### spec.from[].kind

`string` · required

Kind of the referent. Required. Core-supported source kinds: Gateway (when
permitting a SecretObjectReference) and GRPCRoute / HTTPRoute / TCPRoute /
TLSRoute / UDPRoute (when permitting a BackendObjectReference). Pattern: 1-63
chars, ^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$.

INFRA-CHART COMPOSABILITY: kind is a trust assertion about a KIND, not
an instance pointer -- no metadata.relationships hint applies.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"63","pattern":"^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$"}}

### spec.from[].namespace

`string` · required

Namespace of the referent: the source namespace trusted to reference into this
grant's namespace. Required. Pattern: an RFC 1123 label (1-63 chars).

INFRA-CHART COMPOSABILITY: this is the one genuine cross-resource
reference in from/to -- it names another KubernetesNamespace. It stays a PLAIN
string (not a StringValueOrRef foreign key): the grant trusts a namespace by
its literal name (the trust boundary), it does not consume another resource's
output. Because a plain string creates NO automatic DAG edge, when the source
namespace is Planton-managed an infra-chart author should express the edge
via metadata.relationships, e.g.:
  metadata:
    relationships:
      - kind: KubernetesNamespace
        name: "{{ values.source_ns }}"
        type: uses

- rule: {"required":true,"string":{"minLen":"1","maxLen":"63","pattern":"^[a-z0-9]([-a-z0-9]*[a-z0-9])?$"}}

### spec.to

`[]KubernetesReferenceGrantTo` · required

The resources in this grant's namespace that may be referenced by the
resources described in `from`. Entries are combined with OR. At least one
entry is required (upstream MinItems=1), at most 16. Upstream support
level: Core.

- rule: {"repeated":{"minItems":"1","maxItems":"16"}}

### spec.to[].group

`string` · required · optional (explicit presence)

Group of the referent. The empty string ("") infers the Kubernetes core API
group (e.g. for the Secret and Service kinds). Modeled like from.group: a
proto3 `optional` string with a presence `required` rule -- the key must be
SET (so protojson emits it and the CRD accepts the resource) but may be empty.
Pattern: empty or an RFC 1123 subdomain (max 253).

INFRA-CHART COMPOSABILITY: a trust assertion about a KIND, not an
instance pointer -- no metadata.relationships hint applies.

- rule: {"required":true,"string":{"maxLen":"253","pattern":"^$|^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$"}}

### spec.to[].kind

`string` · required

Kind of the referent. Required. Core-supported target kinds: Secret (when
permitting a SecretObjectReference) and Service (when permitting a
BackendObjectReference). Pattern: 1-63 chars.

INFRA-CHART COMPOSABILITY: a trust assertion about a KIND, not an
instance pointer -- no metadata.relationships hint applies.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"63","pattern":"^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$"}}

### spec.to[].name

`string` · required · optional (explicit presence)

Name of the referent. When unspecified, the grant refers to ALL resources of
the specified group and kind in this grant's namespace; when set, it narrows
the grant to a single named resource. Optional (so absence -- meaning
"all resources" -- is distinct from any concrete name). Upstream ObjectName
has a length bound only (1-253) and no character pattern, so none is invented.

INFRA-CHART COMPOSABILITY: even when set, this is an in-namespace
filter on the trust grant, not a deploy-ordered Planton foreign key -- no
metadata.relationships hint applies.

- rule: {"string":{"minLen":"1","maxLen":"253"}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: KubernetesReferenceGrant, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.reference_grant_name` | `string` | Name of the created ReferenceGrant (equals metadata.name). In InfraCharts this lets the grant be referenced as a low-dependency leaf that consumers (Gateways/Routes making the cross-namespace reference) order themselves after. |
| `status.outputs.namespace` | `string` | Namespace the ReferenceGrant was created in (the resolved spec.namespace). This is the "to" namespace -- the one whose resources the grant authorizes inbound references to. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.namespace` | KubernetesNamespace | `spec.name` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
