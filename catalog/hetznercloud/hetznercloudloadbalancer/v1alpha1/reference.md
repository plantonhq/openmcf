# HetznerCloudLoadBalancer

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `hetzner-cloud.planton.dev/v1alpha1`

HetznerCloudLoadBalancerSpec defines the specification for a Hetzner Cloud
load balancer with its services, targets, and optional network attachment.

A load balancer distributes incoming traffic across one or more backend
targets (servers, label-selected server groups, or external IPs). Each
service defines a listener that accepts traffic on a protocol/port
combination and forwards it to targets on a destination port. Services
support HTTP-level features such as sticky sessions, TLS termination with
Hetzner-managed or uploaded certificates, and HTTP-to-HTTPS redirection.

Targets can be added per-server (via server_id), dynamically via Hetzner
Cloud label selectors, or as raw IP addresses for external backends. The
load balancer can optionally be attached to a private Hetzner Cloud network
so that target traffic flows over the private network instead of the public
internet.

Bundled provider resources:
  - hcloud_load_balancer:         The load balancer itself.
  - hcloud_load_balancer_service: One per entry in the services list.
  - hcloud_load_balancer_target:  One per entry across all three target
                                  lists (server_targets, label_selector_targets,
                                  ip_targets).
  - hcloud_load_balancer_network: Created when the network attachment is set.

Fields not exposed in this spec (hardcoded or derived in IaC modules):
  - name:   Derived from metadata.name.
  - labels: Derived from metadata (CG01 pattern). Standard labels take
            precedence over user-specified metadata.labels.

## Example

```yaml
apiVersion: hetzner-cloud.planton.dev/v1alpha1
kind: HetznerCloudLoadBalancer
metadata:
  name: hetznercloudloadbalancer-demo
  org: demo-org
  env: dev
  labels:
    team: platform
spec:
  loadBalancerType: lb11
  location: fsn1
  algorithm: round_robin
  services:
    - protocol: https
      listenPort: 443
      destinationPort: 8080
      http:
        certificateIds:
          - value: "12345"
        redirectHttp: true
        stickySessions: true
        cookieName: APPSESSION
        cookieLifetime: 3600
      healthCheck:
        protocol: http
        port: 8080
        interval: 10
        timeout: 5
        retries: 3
        http:
          path: /health
    - protocol: http
      listenPort: 80
      destinationPort: 8080
  serverTargets:
    - serverId:
        value: "11111"
      usePrivateIp: true
    - serverId:
        value: "22222"
      usePrivateIp: true
  network:
    networkId:
      value: "99999"
    ip: "10.0.1.100"
    enablePublicInterface: true
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.loadBalancerType` | `string` | yes |  |  |
| `spec.location` | `string` | yes |  |  |
| `spec.algorithm` | `enum` |  |  |  |
| `spec.services` | `[]Service` | yes |  |  |
| `spec.services[].protocol` | `enum` | yes |  |  |
| `spec.services[].listenPort` | `int32` |  |  |  |
| `spec.services[].destinationPort` | `int32` |  |  |  |
| `spec.services[].proxyprotocol` | `bool` |  |  |  |
| `spec.services[].http` | `HttpConfig` |  |  |  |
| `spec.services[].http.stickySessions` | `bool` |  |  |  |
| `spec.services[].http.cookieName` | `string` |  |  |  |
| `spec.services[].http.cookieLifetime` | `int32` |  |  |  |
| `spec.services[].http.certificateIds` | `[]string \| valueFrom` |  |  | HetznerCloudCertificate (`status.outputs.certificate_id`) |
| `spec.services[].http.redirectHttp` | `bool` |  |  |  |
| `spec.services[].healthCheck` | `HealthCheck` |  |  |  |
| `spec.services[].healthCheck.protocol` | `enum` |  |  |  |
| `spec.services[].healthCheck.port` | `int32` |  |  |  |
| `spec.services[].healthCheck.interval` | `int32` |  | `15` |  |
| `spec.services[].healthCheck.timeout` | `int32` |  | `10` |  |
| `spec.services[].healthCheck.retries` | `int32` |  | `3` |  |
| `spec.services[].healthCheck.http` | `HealthCheckHttp` |  |  |  |
| `spec.services[].healthCheck.http.domain` | `string` |  |  |  |
| `spec.services[].healthCheck.http.path` | `string` |  |  |  |
| `spec.services[].healthCheck.http.response` | `string` |  |  |  |
| `spec.services[].healthCheck.http.tls` | `bool` |  |  |  |
| `spec.services[].healthCheck.http.statusCodes` | `[]string` |  |  |  |
| `spec.serverTargets` | `[]ServerTarget` |  |  |  |
| `spec.serverTargets[].serverId` | `string \| valueFrom` | yes |  | HetznerCloudServer (`status.outputs.server_id`) |
| `spec.serverTargets[].usePrivateIp` | `bool` |  |  |  |
| `spec.labelSelectorTargets` | `[]LabelSelectorTarget` |  |  |  |
| `spec.labelSelectorTargets[].selector` | `string` | yes |  |  |
| `spec.labelSelectorTargets[].usePrivateIp` | `bool` |  |  |  |
| `spec.ipTargets` | `[]IpTarget` |  |  |  |
| `spec.ipTargets[].ip` | `string` | yes |  |  |
| `spec.network` | `NetworkAttachment` |  |  |  |
| `spec.network.networkId` | `string \| valueFrom` | yes |  | HetznerCloudNetwork (`status.outputs.network_id`) |
| `spec.network.ip` | `string` |  |  |  |
| `spec.network.enablePublicInterface` | `bool` |  | `true` |  |
| `spec.deleteProtection` | `bool` |  |  |  |

## Field Details

### spec.loadBalancerType

`string` · required

Load balancer size that determines connection and bandwidth limits.

Available types: "lb11" (25 targets, 10k connections/s),
"lb21" (75 targets, 20k connections/s), "lb31" (150 targets,
40k connections/s).

Can be changed after creation (in-place resize).

- rule: {"string":{"minLen":"1"}}

### spec.location

`string` · required

Hetzner Cloud location for the load balancer (e.g., "fsn1", "nbg1",
"hel1", "ash", "hil", "sin"). Determines the physical datacenter.

Server targets must be reachable from this location. When using a
private network, the network must have a subnet in the same network
zone as this location.

Changing this value forces replacement of the load balancer.

- rule: {"string":{"minLen":"1"}}

### spec.algorithm

`enum`

Traffic distribution algorithm. Determines how the load balancer
selects a target for each incoming connection.

Default: round_robin (if unspecified or algorithm_unspecified).

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `algorithm_unspecified`
- `round_robin` -- Distribute connections evenly across all healthy targets in order.
- `least_connections` -- Send each new connection to the target with the fewest active connections. Better for backends with varying response times.

### spec.services

`[]Service` · required

Services (listeners) that the load balancer exposes. Each service
binds to a listen_port and forwards traffic to targets on a
destination_port using the specified protocol.

At least one service is required -- a load balancer without services
cannot accept traffic.

- rule: {"repeated":{"minItems":"1"}}
- rule: listen_port is required when protocol is tcp
- rule: destination_port is required when protocol is tcp

### spec.services[].protocol

`enum` · required

Listener protocol. Determines how the load balancer processes
incoming connections for this service.

- rule: {"required":true,"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `service_protocol_unspecified`
- `http` -- Layer 7 HTTP. Enables sticky sessions, health check HTTP mode, and HTTP-specific features. Default listen_port: 80.
- `https` -- Layer 7 HTTPS with TLS termination at the load balancer. Requires certificates in the http config. Enables HTTP-to-HTTPS redirect. Default listen_port: 443.
- `tcp` -- Layer 4 TCP pass-through. No application-layer inspection. listen_port and destination_port are required.

### spec.services[].listenPort

`int32` · optional (explicit presence)

Port the load balancer listens on for this service.

Must be unique across all services on the same load balancer.
Defaults to 80 for HTTP, 443 for HTTPS. Required for TCP.

Changing this value forces replacement of the service.

### spec.services[].destinationPort

`int32` · optional (explicit presence)

Port on the target servers that receives forwarded traffic.

Defaults to the listen_port value for HTTP and HTTPS. Required
for TCP.

### spec.services[].proxyprotocol

`bool`

Enable PROXY protocol (v1) when forwarding to targets. The target
application must support PROXY protocol to parse the original
client IP from the PROXY header.

### spec.services[].http

`HttpConfig`

HTTP-level configuration. Only applicable when protocol is http or
https. Ignored for tcp services.

### spec.services[].http.stickySessions

`bool`

Enable cookie-based session affinity. When true, the load balancer
sets a cookie on the first response and routes subsequent requests
with that cookie to the same target.

### spec.services[].http.cookieName

`string`

Name of the sticky session cookie. Only used when sticky_sessions
is true.

Default: "HCLBSTICKY" (provider default when not specified).

### spec.services[].http.cookieLifetime

`int32`

Lifetime of the sticky session cookie in seconds. Only used when
sticky_sessions is true.

Default: 300 (provider default when not specified or 0).

### spec.services[].http.certificateIds

`[]string | valueFrom`

TLS certificates for HTTPS termination. Only used when the parent
service protocol is https.

Each entry accepts a literal Hetzner Cloud certificate ID (as a
string) or a reference to a HetznerCloudCertificate resource's
output via valueFrom.

Example (reference):
  certificateIds:
    - valueFrom:
        kind: HetznerCloudCertificate
        name: my-cert
        fieldPath: status.outputs.certificate_id

- references: HetznerCloudCertificate (`status.outputs.certificate_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: HetznerCloudCertificate, name: <that resource's name>, fieldPath: status.outputs.certificate_id}} -- a bare string does not parse

### spec.services[].http.redirectHttp

`bool`

Redirect all HTTP traffic to HTTPS. Only valid for https services
that use the default listen port (443). When enabled, the load
balancer automatically creates an HTTP-to-HTTPS redirect on port
80.

### spec.services[].healthCheck

`HealthCheck`

Health check configuration for this service. Optional -- when not
set, the provider creates a default health check matching the
service protocol and destination port.

Set this to customize the health check path, interval, thresholds,
or to use a different protocol than the service (e.g., TCP health
check on an HTTP service).

### spec.services[].healthCheck.protocol

`enum`

Health check protocol. Determines whether the check performs a TCP
connection test or an HTTP request.

Default: matches the parent service's protocol (http -> http,
https -> http, tcp -> tcp). Handled in IaC modules.

Allowed values (use exactly as shown):

- `service_protocol_unspecified`
- `http` -- Layer 7 HTTP. Enables sticky sessions, health check HTTP mode, and HTTP-specific features. Default listen_port: 80.
- `https` -- Layer 7 HTTPS with TLS termination at the load balancer. Requires certificates in the http config. Enables HTTP-to-HTTPS redirect. Default listen_port: 443.
- `tcp` -- Layer 4 TCP pass-through. No application-layer inspection. listen_port and destination_port are required.

### spec.services[].healthCheck.port

`int32` · optional (explicit presence)

Port to health-check on the target. Required when the health check
block is present.

Default: matches the parent service's destination_port. Handled in
IaC modules.

### spec.services[].healthCheck.interval

`int32` · optional (explicit presence)

Time between health checks in seconds.

Default: 15

- default: `15`

### spec.services[].healthCheck.timeout

`int32` · optional (explicit presence)

Maximum time to wait for a health check response in seconds. Must
be less than interval.

Default: 10

- default: `10`

### spec.services[].healthCheck.retries

`int32` · optional (explicit presence)

Number of consecutive failed checks before a target is marked
unhealthy.

Default: 3

- default: `3`

### spec.services[].healthCheck.http

`HealthCheckHttp`

HTTP-specific health check settings. Only used when the health
check protocol is http or https.

### spec.services[].healthCheck.http.domain

`string`

Domain name to send in the HTTP Host header for the health check
request. If empty, the target's IP address is used.

### spec.services[].healthCheck.http.path

`string`

URL path for the health check request.

Default: "/" (provider default when not specified).

### spec.services[].healthCheck.http.response

`string`

Expected response body substring. If set, the health check only
passes when the response body contains this string.

### spec.services[].healthCheck.http.tls

`bool`

Verify the target's TLS certificate when performing the health
check. Only meaningful when the health check protocol is https.

### spec.services[].healthCheck.http.statusCodes

`[]string`

Expected HTTP status codes. The health check passes when the
response status matches any code in this list.

Default: ["2??", "3??"] (provider default when not specified).
Uses wildcard notation: "2??" matches any 2xx status.

### spec.serverTargets

`[]ServerTarget`

Server targets. Each entry adds a specific server as a backend.
The server must be reachable from the load balancer's location.

### spec.serverTargets[].serverId

`string | valueFrom` · required

Server to add as a target.

Accepts a literal Hetzner Cloud server ID (as a string) or a
reference to a HetznerCloudServer resource's output via valueFrom.

Changing this value forces replacement of the target.

Example (reference):
  serverId:
    valueFrom:
      kind: HetznerCloudServer
      name: web-1
      fieldPath: status.outputs.server_id

- references: HetznerCloudServer (`status.outputs.server_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: HetznerCloudServer, name: <that resource's name>, fieldPath: status.outputs.server_id}} -- a bare string does not parse

### spec.serverTargets[].usePrivateIp

`bool`

Route traffic to the server's private IP instead of its public IP.
Requires the load balancer and server to be attached to the same
private network.

### spec.labelSelectorTargets

`[]LabelSelectorTarget`

Label selector targets. Each entry dynamically adds all servers
matching a Hetzner Cloud label selector as backends. Servers are
auto-discovered and may change as labels are added or removed.

### spec.labelSelectorTargets[].selector

`string` · required

Hetzner Cloud label selector expression. All servers matching this
selector are added as targets.

Example: "env=production,role=web"

Changing this value forces replacement of the target.

- rule: {"string":{"minLen":"1"}}

### spec.labelSelectorTargets[].usePrivateIp

`bool`

Route traffic to servers' private IPs instead of their public IPs.
Requires the load balancer and all matching servers to be attached
to the same private network.

### spec.ipTargets

`[]IpTarget`

IP targets. Each entry adds an external IP address as a backend.
Use this for targets outside of Hetzner Cloud.

### spec.ipTargets[].ip

`string` · required

IP address of the external backend.

Changing this value forces replacement of the target.

- rule: {"string":{"minLen":"1"}}

### spec.network

`NetworkAttachment`

Private network attachment. Optional. When set, the load balancer
is attached to the specified Hetzner Cloud network and receives a
private IP within the network's subnet range. Target traffic can
then flow over the private network instead of the public internet.

A load balancer can be attached to at most one network.

### spec.network.networkId

`string | valueFrom` · required

Network to attach the load balancer to.

Accepts a literal Hetzner Cloud network ID (as a string) or a
reference to a HetznerCloudNetwork resource's output via valueFrom.

Changing this value forces replacement of the network attachment.

Example (reference):
  networkId:
    valueFrom:
      kind: HetznerCloudNetwork
      name: main-vpc
      fieldPath: status.outputs.network_id

- references: HetznerCloudNetwork (`status.outputs.network_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: HetznerCloudNetwork, name: <that resource's name>, fieldPath: status.outputs.network_id}} -- a bare string does not parse

### spec.network.ip

`string`

Specific IP address to assign to the load balancer within the
network's subnet range. If omitted, Hetzner Cloud auto-assigns
an IP.

### spec.network.enablePublicInterface

`bool` · optional (explicit presence)

Enable the load balancer's public interface. When false, the load
balancer is only reachable via its private network IP.

Default: true

IMPORTANT: This uses optional bool because proto3 bool defaults to
false, which would silently disable the public interface. The
optional keyword + default annotation ensures correct behavior.

- default: `true`

### spec.deleteProtection

`bool`

Prevent accidental deletion of the load balancer via the Hetzner
Cloud API. When enabled, the load balancer cannot be deleted until
protection is removed.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: HetznerCloudLoadBalancer, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.load_balancer_id` | `string` | The Hetzner Cloud numeric ID of the created load balancer (as a string). Can be referenced by other components via StringValueOrRef. |
| `status.outputs.ipv4_address` | `string` | The public IPv4 address assigned to the load balancer. Empty if the public interface is disabled via network.enable_public_interface = false. |
| `status.outputs.ipv6_address` | `string` | The public IPv6 address assigned to the load balancer. Empty if the public interface is disabled. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.services[].http.certificateIds` | HetznerCloudCertificate | `status.outputs.certificate_id` |
| `spec.serverTargets[].serverId` | HetznerCloudServer | `status.outputs.server_id` |
| `spec.network.networkId` | HetznerCloudNetwork | `status.outputs.network_id` |

## See Also

- [Overview](../README.md)
