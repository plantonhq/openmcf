# OpenStackLoadBalancerPool

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `openstack.planton.dev/v1`

OpenStackLoadBalancerPoolSpec defines the configuration for an Octavia backend pool
in OpenStack. A pool groups backend members (servers) and defines the protocol and
load-balancing algorithm used to distribute traffic from a listener.

Terraform resource: openstack_lb_pool_v2
Pulumi resource:    loadbalancer.Pool

Design note: The Terraform provider enforces ExactlyOneOf(listener_id, loadbalancer_id).
This component exposes listener_id only (80/20). Shared pools via loadbalancer_id
(used for L7 policy routing) can be added as a v2 feature.

Validations:
- cookie_name is only valid when persistence type is APP_COOKIE.

## Example

```yaml
apiVersion: openstack.planton.dev/v1
kind: OpenStackLoadBalancerPool
metadata:
  name: test-pool
spec:
  listener_id:
    value: "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
  protocol: "HTTP"
  lb_method: "ROUND_ROBIN"
  description: "Test pool for local development"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.listenerId` | `string \| valueFrom` | yes |  | OpenStackLoadBalancerListener (`status.outputs.listener_id`) |
| `spec.protocol` | `string` | yes |  |  |
| `spec.lbMethod` | `string` | yes |  |  |
| `spec.persistence` | `SessionPersistence` |  |  |  |
| `spec.persistence.type` | `string` | yes |  |  |
| `spec.persistence.cookieName` | `string` |  |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.adminStateUp` | `bool` |  | `true` |  |
| `spec.tags` | `[]string` |  |  |  |
| `spec.region` | `string` |  |  |  |

## Field Details

### spec.listenerId

`string | valueFrom` · required

(Required) The listener this pool is the default pool for.
ForceNew: changing this requires recreating the pool.

FK: OpenStackLoadBalancerListener.status.outputs.listener_id

- references: OpenStackLoadBalancerListener (`status.outputs.listener_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OpenStackLoadBalancerListener, name: <that resource's name>, fieldPath: status.outputs.listener_id}} -- a bare string does not parse

### spec.protocol

`string` · required

(Required) The protocol used by pool members to receive traffic.
ForceNew: changing this requires recreating the pool.

- rule: {"required":true,"string":{"in":["HTTP","HTTPS","TCP","UDP","PROXY"]}}

### spec.lbMethod

`string` · required

(Required) The load-balancing algorithm to distribute traffic across members.

- ROUND_ROBIN: Equal distribution across all members
- LEAST_CONNECTIONS: Send to member with fewest active connections
- SOURCE_IP: Hash client IP for sticky routing
- SOURCE_IP_PORT: Hash client IP and port for fine-grained stickiness

- rule: {"required":true,"string":{"in":["ROUND_ROBIN","LEAST_CONNECTIONS","SOURCE_IP","SOURCE_IP_PORT"]}}

### spec.persistence

`SessionPersistence`

(Optional) Session persistence configuration.
Ensures requests from the same client are routed to the same backend member.
Only one persistence config is allowed (Octavia enforces this).

- rule: cookie_name is only valid when persistence type is APP_COOKIE

### spec.persistence.type

`string` · required

(Required) The type of session persistence.

- SOURCE_IP: Hash the client's IP address
- HTTP_COOKIE: Octavia inserts and tracks a cookie
- APP_COOKIE: Application manages the cookie (requires cookie_name)

- rule: {"required":true,"string":{"in":["SOURCE_IP","HTTP_COOKIE","APP_COOKIE"]}}

### spec.persistence.cookieName

`string`

(Optional) The name of the application cookie to use for session affinity.
Only valid when type is APP_COOKIE.

### spec.description

`string`

(Optional) Human-readable description of the pool.

### spec.adminStateUp

`bool` · optional (explicit presence)

(Optional) Administrative state of the pool.
When false, the pool stops receiving traffic. Default: true.

- default: `true`

### spec.tags

`[]string`

(Optional) Tags applied to the pool in OpenStack.
Must be unique within this resource.

- rule: {"repeated":{"unique":true}}

### spec.region

`string`

(Optional) Override the region from the provider configuration.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OpenStackLoadBalancerPool, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.pool_id` | `string` | pool_id is the unique identifier (UUID) of the pool. This is the primary output used as a foreign key by members and monitors. |
| `status.outputs.name` | `string` | name is the name of the pool (derived from metadata.name). |
| `status.outputs.protocol` | `string` | protocol is the backend protocol (HTTP, HTTPS, TCP, UDP, PROXY). |
| `status.outputs.lb_method` | `string` | lb_method is the load-balancing algorithm (ROUND_ROBIN, LEAST_CONNECTIONS, etc.). |
| `status.outputs.region` | `string` | region is the OpenStack region where the pool was created. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.listenerId` | OpenStackLoadBalancerListener | `status.outputs.listener_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| OpenStackLoadBalancerMember | `spec.poolId` | `status.outputs.pool_id` |
| OpenStackLoadBalancerMonitor | `spec.poolId` | `status.outputs.pool_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
