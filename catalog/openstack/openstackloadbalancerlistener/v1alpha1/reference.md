# OpenStackLoadBalancerListener

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `openstack.planton.dev/v1alpha1`

OpenStackLoadBalancerListenerSpec defines the configuration for an Octavia listener
in OpenStack. A listener binds a protocol and port to a load balancer, accepting
incoming traffic and forwarding it to a backend pool.

Terraform resource: openstack_lb_listener_v2
Pulumi resource:    loadbalancer.Listener

Validations:
- default_tls_container_ref is required when protocol is TERMINATED_HTTPS.
- protocol_port must be between 1 and 65535.

## Example

```yaml
apiVersion: openstack.planton.dev/v1alpha1
kind: OpenStackLoadBalancerListener
metadata:
  name: test-http-listener
spec:
  loadbalancer_id:
    value: "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
  protocol: "HTTP"
  protocol_port: 80
  description: "Test HTTP listener for local development"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.loadbalancerId` | `string \| valueFrom` | yes |  | OpenStackLoadBalancer (`status.outputs.loadbalancer_id`) |
| `spec.protocol` | `string` | yes |  |  |
| `spec.protocolPort` | `int32` | yes |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.connectionLimit` | `int32` |  |  |  |
| `spec.defaultTlsContainerRef` | `string` |  |  |  |
| `spec.insertHeaders` | `map<string, string>` |  |  |  |
| `spec.allowedCidrs` | `[]string` |  |  |  |
| `spec.adminStateUp` | `bool` |  | `true` |  |
| `spec.tags` | `[]string` |  |  |  |
| `spec.region` | `string` |  |  |  |

## Field Details

### spec.loadbalancerId

`string | valueFrom` · required

(Required) The load balancer to attach this listener to.
ForceNew: changing this requires recreating the listener.

FK: OpenStackLoadBalancer.status.outputs.loadbalancer_id

- references: OpenStackLoadBalancer (`status.outputs.loadbalancer_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OpenStackLoadBalancer, name: <that resource's name>, fieldPath: status.outputs.loadbalancer_id}} -- a bare string does not parse

### spec.protocol

`string` · required

(Required) The protocol the listener accepts.
ForceNew: changing this requires recreating the listener.

- HTTP: Unencrypted HTTP traffic (Layer 7)
- HTTPS: Pass-through encrypted traffic (Layer 4, no TLS termination)
- TCP: Raw TCP traffic (Layer 4)
- UDP: Raw UDP traffic (Layer 4)
- TERMINATED_HTTPS: TLS termination at the load balancer (requires default_tls_container_ref)

- rule: {"required":true,"string":{"in":["HTTP","HTTPS","TCP","UDP","TERMINATED_HTTPS"]}}

### spec.protocolPort

`int32` · required

(Required) The port on which the listener accepts traffic.
ForceNew: changing this requires recreating the listener.
Must be between 1 and 65535.

- rule: protocol_port must be between 1 and 65535
- rule: {"required":true}

### spec.description

`string`

(Optional) Human-readable description of the listener.

### spec.connectionLimit

`int32` · optional (explicit presence)

(Optional) Maximum number of connections the listener allows.
-1 means unlimited (Octavia default). Leave unset to use the Octavia default.

### spec.defaultTlsContainerRef

`string`

(Optional) URI of the Barbican TLS secret container for TLS termination.
Required when protocol is TERMINATED_HTTPS. The container must hold the
certificate, private key, and optional intermediates.

### spec.insertHeaders

`map<string, string>`

(Optional) Headers to insert into HTTP requests before forwarding to backends.
Common use: {"X-Forwarded-For": "true", "X-Forwarded-Proto": "true"}
Only applicable to HTTP and TERMINATED_HTTPS protocols.

### spec.allowedCidrs

`[]string`

(Optional) List of CIDRs allowed to access this listener.
When set, only traffic from these CIDRs reaches the listener; all other
traffic is dropped. When empty, all traffic is allowed.

### spec.adminStateUp

`bool` · optional (explicit presence)

(Optional) Administrative state of the listener.
When false, the listener stops accepting traffic. Default: true.

- default: `true`

### spec.tags

`[]string`

(Optional) Tags applied to the listener in OpenStack.
Must be unique within this resource.

- rule: {"repeated":{"unique":true}}

### spec.region

`string`

(Optional) Override the region from the provider configuration.

## Validation Rules

- `tls_ref_required_for_terminated_https`: default_tls_container_ref is required when protocol is TERMINATED_HTTPS

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OpenStackLoadBalancerListener, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.listener_id` | `string` | listener_id is the unique identifier (UUID) of the listener. This is the primary output used as a foreign key by pools. |
| `status.outputs.name` | `string` | name is the name of the listener (derived from metadata.name). |
| `status.outputs.protocol` | `string` | protocol is the protocol the listener accepts (HTTP, HTTPS, TCP, UDP, TERMINATED_HTTPS). |
| `status.outputs.protocol_port` | `int32` | protocol_port is the port on which the listener accepts traffic. |
| `status.outputs.region` | `string` | region is the OpenStack region where the listener was created. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.loadbalancerId` | OpenStackLoadBalancer | `status.outputs.loadbalancer_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| OpenStackLoadBalancerPool | `spec.listenerId` | `status.outputs.listener_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
