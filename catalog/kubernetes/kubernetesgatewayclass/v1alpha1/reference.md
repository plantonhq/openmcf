# KubernetesGatewayClass

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `kubernetes.planton.dev/v1alpha1`

KubernetesGatewayClassSpec defines a Kubernetes Gateway API GatewayClass:
a cluster-scoped template that identifies the controller responsible for
managing Gateways of this class (e.g. Istio, Envoy Gateway, NGINX).

100% fidelity with the upstream Gateway API v1.6.1 GatewayClassSpec
(kubernetes-sigs/gateway-api apis/v1/gatewayclass_types.go). The three
upstream fields (controller_name, parameters_ref, description) are flattened
after the Planton cluster-scoped envelope.

GatewayClass is a cluster-scoped resource (+kubebuilder:resource:scope=Cluster
upstream), so this spec intentionally has NO namespace field.

## Example

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesGatewayClass
metadata:
  name: test-gateway-class
spec:
  controllerName: istio.io/gateway-controller
  parametersRef:
    group: ""
    kind: ConfigMap
    name: gateway-class-params
    namespace: gateway-system
  description: "Full-surface gateway class for offline plan proofs"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.controllerName` | `string` | yes |  |  |
| `spec.parametersRef` | `KubernetesGatewayApiParametersReference` |  |  |  |
| `spec.parametersRef.group` | `string` | yes |  |  |
| `spec.parametersRef.kind` | `string` | yes |  |  |
| `spec.parametersRef.name` | `string` | yes |  |  |
| `spec.parametersRef.namespace` | `string` |  |  |  |
| `spec.description` | `string` |  |  |  |

## Field Details

### spec.controllerName

`string` · required

Name of the controller that manages Gateways of this class, expressed as a
domain-prefixed path (for example, "istio.io/gateway-controller" or
"gateway.envoyproxy.io/gatewayclass-controller").

This value is immutable once the GatewayClass exists in the cluster: the
Gateway API admission webhook rejects changes to controllerName. Choose it
carefully, because changing it requires recreating the GatewayClass.

Upstream support level: Core.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"253","pattern":"^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*\\/[A-Za-z0-9\\/\\-._~%!$&'()*+,;=:]+$"}}

### spec.parametersRef

`KubernetesGatewayApiParametersReference`

Optional reference to a controller-specific resource (such as a ConfigMap
or an implementation-defined custom resource) that holds configuration
parameters for this GatewayClass. Leave unset if the controller needs no
additional configuration.

Shared Gateway API type (see gateway_api.proto). This is a structured
reference (group/kind/name/namespace) to an arbitrary Kubernetes object,
not a single Planton resource output, so it is not a foreign-key field.

### spec.parametersRef.group

`string` · required · optional (explicit presence)

Group of the referent. Empty string infers the core API group
(e.g. for ConfigMap parameters).

Upstream requires the KEY to be present, but its Group type explicitly
allows the empty value -- so this is a proto3 `optional` string with a
presence `required` rule: it must be SET (and is therefore emitted to
the rendered CR, whose CRD rejects a missing key) but may be empty.
The `optional` is what keeps the projection faithful: protojson omits
unset proto3 scalars, so a non-optional empty-string group would be
dropped from the manifest and rejected by the API server.

Group pattern: empty or an RFC 1123 subdomain (max 253).

- rule: {"required":true,"string":{"maxLen":"253","pattern":"^$|^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$"}}

### spec.parametersRef.kind

`string` · required

Kind of the referent.

Upstream models Kind as a required value type.
Kind pattern: 1-63 chars, ^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$

- rule: {"required":true,"string":{"minLen":"1","maxLen":"63","pattern":"^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$"}}

### spec.parametersRef.name

`string` · required

Name of the referent.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"253"}}

### spec.parametersRef.namespace

`string` · optional (explicit presence)

Namespace of the referent. Required for namespace-scoped resources,
must be unset for cluster-scoped resources.

### spec.description

`string` · optional (explicit presence)

Optional human-friendly description of this GatewayClass. Upstream limits
this to 64 characters.

- rule: {"string":{"maxLen":"64"}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: KubernetesGatewayClass, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.gateway_class_name` | `string` | Name of the created GatewayClass (equals metadata.name). Use this value in KubernetesGateway.spec.gateway_class_name to attach a Gateway to this class; in InfraCharts it is the foreign-key target (status.outputs.gateway_class_name) that orders GatewayClass before Gateway. |
| `status.outputs.controller_name` | `string` | The controller managing this class (echoes spec.controller_name), exposed for observability and for confirming which implementation owns the class. |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| KubernetesGateway | `spec.gatewayClassName` | `status.outputs.gateway_class_name` |

## See Also

- [Overview](../README.md)
