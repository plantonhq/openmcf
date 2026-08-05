# KubernetesGatewayApiCrds

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `kubernetes.planton.dev/v1`

KubernetesGatewayApiCrdsSpec defines configuration for installing the
Kubernetes Gateway API CRDs on any Kubernetes cluster.

The Gateway API is the next-generation Kubernetes API for managing ingress
and service mesh traffic. Installing these CRDs enables Gateway, HTTPRoute,
GRPCRoute, and other Gateway API resources to be used in the cluster.

## Example

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesGatewayApiCrds
metadata:
  name: gateway-api-crds-test
spec:
  version: v1.6.1
  install_channel:
    channel: standard
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.version` | `string` | yes | `v1.6.1` |  |
| `spec.installChannel` | `InstallChannel` |  |  |  |
| `spec.installChannel.channel` | `enum` |  |  |  |

## Field Details

### spec.version

`string` · required · optional (explicit presence)

Gateway API version to install (e.g., "v1.6.1").
Defaults to v1.6.1 if not specified — the version the catalog's Gateway
API projection kinds (Gateway, routes, ListenerSet, ReferenceGrant) are
designed against. Installing an older version narrows what those kinds
can deploy (for example, TCPRoute and UDPRoute are standard-channel only
from v1.6.0, and ListenerSet from v1.5.0).
Must start with 'v' followed by a valid semver version.

- default: `v1.6.1`
- rule: {"string":{"minLen":"1","pattern":"^v[0-9]+\\.[0-9]+\\.[0-9]+(-[a-zA-Z0-9]+)?$"}}

### spec.installChannel

`InstallChannel`

Installation channel for Gateway API CRDs.
Defaults to standard channel if not specified.

### spec.installChannel.channel

`enum`

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `gateway_api_install_channel_unspecified`
- `standard` -- Standard channel (GA/beta resources). As of v1.6: GatewayClass, Gateway, ListenerSet, HTTPRoute, GRPCRoute, TLSRoute, TCPRoute, UDPRoute, ReferenceGrant, BackendTLSPolicy. This is the channel the catalog's projection kinds target.
- `experimental` -- Experimental channel: all standard resources (with additional experimental fields) plus experimental resources such as XBackend, XBackendTrafficPolicy, and XMesh. Future releases may break or remove experimental resources — prefer standard unless a specific experimental feature is required.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: KubernetesGatewayApiCrds, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.installed_version` | `string` | Gateway API version that was installed (e.g., "v1.2.1", "v1.3.0"). |
| `status.outputs.installed_channel` | `string` | Installation channel that was used (standard or experimental). |
| `status.outputs.installed_manifest_url` | `string` | Full URL of the Gateway API CRD bundle that was applied. This is the exact artifact installed -- it encodes both the version and the channel in a single value (e.g. ".../releases/download/v1.5.1/experimental-install.yaml"), so it is the authoritative record of what landed on the cluster. |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
