# ScalewayLoadBalancer

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `scaleway.planton.dev/v1alpha1`

ScalewayLoadBalancerSpec defines the specification for a Scaleway Load Balancer.

A Scaleway Load Balancer is a managed Layer 4/7 traffic distribution appliance
that sits in front of backend servers and distributes incoming connections
across them based on configurable forwarding rules, health checks, and
load-balancing algorithms.

This resource bundles several Scaleway resources into a single declarative unit:
  1. A dedicated Flexible IP (public IPv4 address).
  2. The Load Balancer appliance itself.
  3. One or more backends (server pools with health checks).
  4. One or more frontends (listeners that route traffic to backends).
  5. Optional TLS certificates (Let's Encrypt or custom PEM).

Load Balancers are **zonal** resources (e.g., "fr-par-1"), unlike VPCs and
Private Networks which are regional.

**Composition pattern**: The Load Balancer accepts a Private Network reference
via `StringValueOrRef`, enabling it to reach backend servers on private IPs.
Downstream resources like `ScalewayDnsRecord` can reference
`status.outputs.lb_ip_address` to create DNS records pointing to the LB.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.zone` | `string` | yes |  |  |
| `spec.type` | `string` | yes | `LB-S` |  |
| `spec.privateNetworkId` | `string \| valueFrom` |  |  | ScalewayPrivateNetwork (`status.outputs.private_network_id`) |
| `spec.description` | `string` |  |  |  |
| `spec.sslCompatibilityLevel` | `string` |  |  |  |
| `spec.backends` | `[]ScalewayLoadBalancerBackend` | yes |  |  |
| `spec.backends[].name` | `string` | yes |  |  |
| `spec.backends[].serverIps` | `[]string` | yes |  |  |
| `spec.backends[].forwardPort` | `int32` | yes |  |  |
| `spec.backends[].forwardProtocol` | `string` | yes | `http` |  |
| `spec.backends[].forwardPortAlgorithm` | `string` |  | `roundrobin` |  |
| `spec.backends[].stickySessions` | `string` |  |  |  |
| `spec.backends[].stickySessionsCookieName` | `string` |  |  |  |
| `spec.backends[].healthCheck` | `ScalewayLoadBalancerHealthCheck` |  |  |  |
| `spec.backends[].healthCheck.type` | `string` |  | `tcp` |  |
| `spec.backends[].healthCheck.uri` | `string` |  |  |  |
| `spec.backends[].healthCheck.expectedCode` | `int32` |  |  |  |
| `spec.backends[].healthCheck.checkDelay` | `string` |  | `5s` |  |
| `spec.backends[].healthCheck.checkTimeout` | `string` |  | `3s` |  |
| `spec.backends[].healthCheck.checkMaxRetries` | `int32` |  | `3` |  |
| `spec.backends[].healthCheck.port` | `int32` |  |  |  |
| `spec.backends[].timeoutConnect` | `string` |  |  |  |
| `spec.backends[].timeoutServer` | `string` |  |  |  |
| `spec.backends[].onMarkedDownAction` | `string` |  |  |  |
| `spec.backends[].sslBridging` | `bool` |  |  |  |
| `spec.backends[].proxyProtocol` | `string` |  |  |  |
| `spec.frontends` | `[]ScalewayLoadBalancerFrontend` | yes |  |  |
| `spec.frontends[].name` | `string` | yes |  |  |
| `spec.frontends[].inboundPort` | `int32` | yes |  |  |
| `spec.frontends[].backendName` | `string` | yes |  |  |
| `spec.frontends[].certificateNames` | `[]string` |  |  |  |
| `spec.frontends[].timeoutClient` | `string` |  |  |  |
| `spec.frontends[].enableHttp3` | `bool` |  |  |  |
| `spec.certificates` | `[]ScalewayLoadBalancerCertificate` |  |  |  |
| `spec.certificates[].name` | `string` | yes |  |  |
| `spec.certificates[].letsencrypt` | `ScalewayLoadBalancerLetsencrypt` |  |  |  |
| `spec.certificates[].letsencrypt.commonName` | `string` | yes |  |  |
| `spec.certificates[].letsencrypt.subjectAlternativeNames` | `[]string` |  |  |  |
| `spec.certificates[].customCertificate` | `ScalewayLoadBalancerCustomCertificate` |  |  |  |
| `spec.certificates[].customCertificate.certificateChain` | `string` | yes |  |  |

## Field Details

### spec.zone

`string` · required

The Scaleway zone where the Load Balancer will be created.
Examples: "fr-par-1", "nl-ams-1", "pl-waw-1"

Load Balancers are zonal resources. The zone must be within the same
region as any Private Network the LB is attached to. For example, if
the Private Network is in region "fr-par", the LB zone must be
"fr-par-1", "fr-par-2", etc.

This field is required and cannot be changed after creation.

- rule: {"required":true}

### spec.type

`string` · required

Load Balancer type determines bandwidth, throughput, and pricing tier.

Available types (subject to Scaleway's current offering):
  - "LB-S"     -- Small. Up to 400 Mbps. Good for development and small apps.
  - "LB-GP-M"  -- Medium. Up to 4 Gbps. General-purpose production workloads.
  - "LB-GP-L"  -- Large. Up to 8 Gbps. High-traffic applications.
  - "LB-GP-XL" -- Extra Large. Up to 10 Gbps. Maximum throughput.

Choose "LB-S" for development and "LB-GP-M" for most production workloads.
Type can be changed after creation (scale up/down).

- default: `LB-S`
- rule: {"required":true}

### spec.privateNetworkId

`string | valueFrom`

The Private Network to attach the Load Balancer to.

When set, the LB receives a private IP on this network and can reach
backend servers via their private IPs. This is the recommended topology
for production: keep backend servers off the public internet and let
the LB handle ingress.

In infra charts, this is typically wired via valueFrom:

  privateNetworkId:
    valueFrom:
      kind: ScalewayPrivateNetwork
      name: app-network
      fieldPath: status.outputs.private_network_id

Optional. If omitted, the LB operates on the public network only.

- references: ScalewayPrivateNetwork (`status.outputs.private_network_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: ScalewayPrivateNetwork, name: <that resource's name>, fieldPath: status.outputs.private_network_id}} -- a bare string does not parse

### spec.description

`string`

Human-readable description for the Load Balancer.

Optional. Useful for organizational purposes in the Scaleway console.

### spec.sslCompatibilityLevel

`string`

Minimum SSL/TLS compatibility level for HTTPS frontends.

Controls the minimum TLS version clients must support to connect.
Options:
  - "ssl_compatibility_level_intermediate" (default) -- TLS 1.2+. Broad compatibility.
  - "ssl_compatibility_level_modern" -- TLS 1.3 only. Maximum security.

Only relevant when using HTTPS frontends with certificates. If omitted,
Scaleway uses the "intermediate" level.

### spec.backends

`[]ScalewayLoadBalancerBackend` · required

Backend pools. Each backend defines a named set of servers that receive
traffic, along with health check rules and load-balancing configuration.

At least one backend is required. Frontends reference backends by name
to route traffic.

Example: A backend named "web" with two server IPs on port 80.

- rule: {"required":true,"repeated":{"minItems":"1"}}

### spec.backends[].name

`string` · required

Name identifying this backend. Used by frontends to reference this pool.

Must be unique within the Load Balancer spec. Use descriptive names
like "web", "api", "grpc" to make the configuration self-documenting.

- rule: {"required":true}

### spec.backends[].serverIps

`[]string` · required

IP addresses of backend servers.

These are the servers that will receive forwarded traffic. When the LB
is attached to a Private Network, use private IPs (e.g., "10.0.1.5").
Without a Private Network, use public IPs.

At least one server IP is required.

- rule: {"required":true,"repeated":{"minItems":"1"}}

### spec.backends[].forwardPort

`int32` · required

Port on backend servers that receives forwarded traffic.

This is the port your application listens on (e.g., 80 for HTTP,
443 for HTTPS, 8080 for a custom app server).

- rule: {"required":true}

### spec.backends[].forwardProtocol

`string` · required

Protocol for communication between the LB and backend servers.

Options: "http", "https", "tcp"
  - "http"  -- Layer 7. The LB inspects HTTP headers. Use for web apps.
  - "https" -- Layer 7 with TLS to backends (ssl_bridging). Use when
               backends require encrypted connections.
  - "tcp"   -- Layer 4. The LB forwards raw TCP. Use for databases,
               gRPC, or any non-HTTP protocol.

- default: `http`
- rule: {"required":true}

### spec.backends[].forwardPortAlgorithm

`string`

Load-balancing algorithm for distributing connections across servers.

Options:
  - "roundrobin" (default) -- Distribute evenly in rotation.
  - "leastconn"  -- Send to the server with fewest active connections.
  - "first"      -- Send to the first healthy server (active-passive).

- default: `roundrobin`

### spec.backends[].stickySessions

`string`

Sticky session type for maintaining client affinity.

Options:
  - "none" (default)  -- No session affinity.
  - "cookie" -- Insert an HTTP cookie to track client sessions.
  - "table"  -- Use a connection table (Layer 4, for TCP backends).

Use "cookie" for web applications that store session state server-side.
Use "table" for TCP-level session affinity.

### spec.backends[].stickySessionsCookieName

`string`

Cookie name for sticky sessions. Required when sticky_sessions = "cookie".

The LB injects this cookie into HTTP responses to track which backend
server a client was previously routed to.
Example: "SERVERID"

### spec.backends[].healthCheck

`ScalewayLoadBalancerHealthCheck`

Health check configuration for this backend.

Defines how the LB probes backend servers to detect failures. Unhealthy
servers are automatically removed from rotation and re-added when they
recover.

If omitted, a default TCP health check on the forward_port is used
with a 5-second interval, 3-second timeout, and 3 retries.

### spec.backends[].healthCheck.type

`string`

Health check protocol.

Options: "tcp" (default), "http", "https"

Use "tcp" for non-HTTP services (databases, gRPC, custom protocols).
Use "http" for web applications. Use "https" when backends require TLS.

- default: `tcp`

### spec.backends[].healthCheck.uri

`string`

URI path for HTTP/HTTPS health checks.

The LB sends a GET request to this path on each backend server.
Example: "/health", "/ready", "/ping"

Default: "/" -- Only meaningful when type is "http" or "https".

### spec.backends[].healthCheck.expectedCode

`int32`

Expected HTTP status code for a healthy response.

The health check passes if the server returns this status code.
Default: 200 -- Only meaningful when type is "http" or "https".

### spec.backends[].healthCheck.checkDelay

`string`

Interval between health check probes.

Duration string (e.g., "5s", "10s", "30s").
Lower values detect failures faster but generate more probe traffic.
Default: "5s"

- default: `5s`

### spec.backends[].healthCheck.checkTimeout

`string`

Maximum time to wait for a health check response.

Duration string (e.g., "3s", "5s").
Must be less than check_delay. Default: "3s"

- default: `3s`

### spec.backends[].healthCheck.checkMaxRetries

`int32`

Number of consecutive failed checks before marking a server as unhealthy.

Higher values tolerate transient failures but take longer to detect
real outages. Default: 3

- default: `3`

### spec.backends[].healthCheck.port

`int32`

Port to send health check probes to.

If omitted or set to 0, defaults to the backend's `forward_port`.
Set a different port when health checks run on a dedicated monitoring
port (e.g., an application that serves traffic on 8080 but exposes
health at 8081).

### spec.backends[].timeoutConnect

`string`

Maximum time to wait for a connection to a backend server.

Duration string (e.g., "5s", "10s"). If omitted, Scaleway's default applies.
Increase for backends with slow connection establishment (e.g., cold starts).

### spec.backends[].timeoutServer

`string`

Maximum time a backend server connection can be idle before being closed.

Duration string (e.g., "30s", "5m"). If omitted, Scaleway's default applies.
Increase for long-polling, WebSocket, or streaming backends.

### spec.backends[].onMarkedDownAction

`string`

Action when a backend server is marked as down.

Options:
  - "none" (default) -- Keep existing connections open.
  - "shutdown_sessions" -- Immediately close all connections to the
    downed server. Use when fast failover is more important than
    graceful connection draining.

### spec.backends[].sslBridging

`bool`

Enable SSL bridging (re-encrypt traffic between the LB and backends).

When true, the LB establishes a new TLS connection to backend servers.
Use when backends require encrypted connections (e.g., for compliance).

Default: false (traffic between LB and backends is unencrypted).

### spec.backends[].proxyProtocol

`string`

PROXY protocol version for passing client connection metadata to backends.

Options:
  - "none" (default) -- No PROXY protocol.
  - "v1"      -- PROXY protocol v1 (human-readable header).
  - "v2"      -- PROXY protocol v2 (binary header).
  - "v2_ssl"  -- v2 with SSL information.
  - "v2_ssl_cn" -- v2 with SSL and client certificate CN.

Use when backend servers need the original client IP (e.g., Nginx with
`proxy_protocol` directive, HAProxy with `accept-proxy`).

### spec.frontends

`[]ScalewayLoadBalancerFrontend` · required

Frontends. Each frontend defines a named listener on a specific port
that routes incoming traffic to a backend.

At least one frontend is required. Each frontend must reference a
backend by name (the `backend_name` field must match a backend's `name`).

Example: A frontend named "http" listening on port 80 routing to backend "web".

- rule: {"required":true,"repeated":{"minItems":"1"}}

### spec.frontends[].name

`string` · required

Name identifying this frontend.

Must be unique within the Load Balancer spec. Use descriptive names
like "http", "https", "grpc", "api-8080".

- rule: {"required":true}

### spec.frontends[].inboundPort

`int32` · required

TCP port to listen on for incoming connections.

Common ports: 80 (HTTP), 443 (HTTPS), 8080 (alt-HTTP), 8443 (alt-HTTPS).
Each frontend must use a unique port within the Load Balancer.

- rule: {"required":true}

### spec.frontends[].backendName

`string` · required

Name of the backend to route traffic to.

Must match a backend's `name` field defined in `spec.backends`.
All traffic arriving on this frontend's port is forwarded to the
referenced backend's server pool.

- rule: {"required":true}

### spec.frontends[].certificateNames

`[]string`

Names of TLS certificates to attach to this frontend.

Must match certificate names defined in `spec.certificates`.
Required for HTTPS frontends. Not allowed on plaintext HTTP frontends
(Scaleway rejects certificates on port 80).

Multiple certificates can be attached for SNI-based selection (the LB
selects the correct certificate based on the client's requested hostname).

### spec.frontends[].timeoutClient

`string`

Maximum time a client connection can be idle before being closed.

Duration string (e.g., "30s", "5m"). If omitted, Scaleway's default applies.
Increase for long-polling, WebSocket, or Server-Sent Events clients.

### spec.frontends[].enableHttp3

`bool`

Enable HTTP/3 (QUIC) support on this frontend.

When true, the frontend accepts HTTP/3 connections over UDP in addition
to HTTP/1.1 and HTTP/2 over TCP. Requires an HTTPS frontend with a
TLS certificate.

Default: false.

### spec.certificates

`[]ScalewayLoadBalancerCertificate`

TLS certificates for HTTPS frontends.

Each certificate has a name and is either auto-provisioned via Let's Encrypt
or provided as a custom PEM chain. Frontends reference certificates by
name in their `certificate_names` field.

Optional. Only needed when frontends serve HTTPS traffic.

### spec.certificates[].name

`string` · required

Name identifying this certificate.

Must be unique within the Load Balancer spec. Frontends reference
this name in their `certificate_names` field.
Example: "example-com-cert", "wildcard-cert"

- rule: {"required":true}

### spec.certificates[].letsencrypt

`ScalewayLoadBalancerLetsencrypt`

Let's Encrypt auto-provisioned certificate configuration.

When set, Scaleway automatically provisions and renews a TLS certificate
for the specified domain(s) using the ACME protocol. The domain must
resolve to the LB's public IP for validation to succeed.

Exactly one of `letsencrypt` or `custom_certificate` must be set.

### spec.certificates[].letsencrypt.commonName

`string` · required

Primary domain name for the certificate.

Example: "example.com", "app.example.com"

The domain must resolve to the LB's public IP. If using a new domain,
create a DNS A record first, then add the certificate.

- rule: {"required":true}

### spec.certificates[].letsencrypt.subjectAlternativeNames

`[]string`

Subject Alternative Names (additional domains covered by this certificate).

Examples: ["www.example.com", "api.example.com"]

All SANs must also resolve to the LB's public IP. Let's Encrypt
validates each domain independently.

### spec.certificates[].customCertificate

`ScalewayLoadBalancerCustomCertificate`

Custom certificate configuration (user-provided PEM).

When set, the user provides a full certificate chain in PEM format.
Use this for certificates from commercial CAs, internal PKI, or
wildcard certificates that Let's Encrypt doesn't support.

Exactly one of `letsencrypt` or `custom_certificate` must be set.

### spec.certificates[].customCertificate.certificateChain

`string` · required

Full certificate chain in PEM format.

Must include the server certificate followed by any intermediate
certificates, in order. The private key must NOT be included (it is
managed separately by Scaleway).

Example: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----\n"

- rule: {"required":true}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: ScalewayLoadBalancer, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.lb_id` | `string` | The unique identifier of the created Load Balancer. Format: zoned ID (e.g., "fr-par-1/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"). This ID can be used to reference the LB in Scaleway APIs or to manage it outside of Planton. |
| `status.outputs.lb_ip_address` | `string` | The public IPv4 address assigned to the Load Balancer's Flexible IP. This is the address that clients connect to. Use this output for: - DNS A records pointing to the LB (via ScalewayDnsRecord) - Firewall allowlists on backend servers - Monitoring and connectivity diagnostics - External service whitelisting |
| `status.outputs.lb_ip_id` | `string` | The unique identifier of the Flexible IP resource attached to the LB. Format: zoned ID. The Flexible IP is created as a dedicated resource with independent lifecycle, so the IP survives LB replacement. This ID is useful for managing the IP independently. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.privateNetworkId` | ScalewayPrivateNetwork | `status.outputs.private_network_id` |

## See Also

- [Overview](../README.md)
