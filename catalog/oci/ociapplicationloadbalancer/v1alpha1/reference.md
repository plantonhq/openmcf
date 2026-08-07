# OciApplicationLoadBalancer

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `oci.planton.dev/v1alpha1`

OciApplicationLoadBalancerSpec defines the specification for an Oracle Cloud
Infrastructure Application Load Balancer (Layer 7).

An OCI load balancer provides automated traffic distribution across multiple
backend servers. It supports HTTP, HTTP/2, TCP, and gRPC protocols with
features including SSL termination, cookie-based session persistence,
health checking, virtual hostname routing, and rule-based request
manipulation (header injection, HTTP redirects, access control).

This component bundles the load balancer with its backend sets, backends,
listeners, certificates, hostnames, and rule sets into a single deployment
unit. All sub-resources are created atomically -- a load balancer without
at least one backend set and listener is not functional.

For Layer 4 (TCP/UDP) load balancing with source IP preservation, use
OciNetworkLoadBalancer instead.

Deprecated/excluded resources:
  - path_route_set: deprecated by Oracle in favor of routing policies

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.compartmentId` | `string \| valueFrom` | yes |  | OciCompartment (`status.outputs.compartment_id`) |
| `spec.displayName` | `string` |  |  |  |
| `spec.shape` | `string` | yes |  |  |
| `spec.shapeDetails` | `ShapeDetails` |  |  |  |
| `spec.shapeDetails.minimumBandwidthInMbps` | `int32` |  |  |  |
| `spec.shapeDetails.maximumBandwidthInMbps` | `int32` |  |  |  |
| `spec.subnetIds` | `[]string \| valueFrom` | yes |  | OciSubnet (`status.outputs.subnet_id`) |
| `spec.isPrivate` | `bool` |  |  |  |
| `spec.networkSecurityGroupIds` | `[]string \| valueFrom` |  |  | OciSecurityGroup (`status.outputs.network_security_group_id`) |
| `spec.isDeleteProtectionEnabled` | `bool` |  |  |  |
| `spec.ipMode` | `string` |  |  |  |
| `spec.reservedIps` | `[]ReservedIp` |  |  |  |
| `spec.reservedIps[].id` | `string` | yes |  |  |
| `spec.isRequestIdEnabled` | `bool` |  |  |  |
| `spec.requestIdHeader` | `string` |  |  |  |
| `spec.backendSets` | `[]BackendSet` | yes |  |  |
| `spec.backendSets[].name` | `string` | yes |  |  |
| `spec.backendSets[].policy` | `enum` |  |  |  |
| `spec.backendSets[].healthChecker` | `HealthChecker` | yes |  |  |
| `spec.backendSets[].healthChecker.protocol` | `enum` |  |  |  |
| `spec.backendSets[].healthChecker.port` | `int32` |  |  |  |
| `spec.backendSets[].healthChecker.urlPath` | `string` |  |  |  |
| `spec.backendSets[].healthChecker.returnCode` | `int32` |  |  |  |
| `spec.backendSets[].healthChecker.responseBodyRegex` | `string` |  |  |  |
| `spec.backendSets[].healthChecker.intervalMs` | `int32` |  |  |  |
| `spec.backendSets[].healthChecker.timeoutInMillis` | `int32` |  |  |  |
| `spec.backendSets[].healthChecker.retries` | `int32` |  |  |  |
| `spec.backendSets[].healthChecker.isForcePlainText` | `bool` |  |  |  |
| `spec.backendSets[].backends` | `[]Backend` |  |  |  |
| `spec.backendSets[].backends[].ipAddress` | `string` | yes |  |  |
| `spec.backendSets[].backends[].port` | `int32` |  |  |  |
| `spec.backendSets[].backends[].weight` | `int32` |  |  |  |
| `spec.backendSets[].backends[].backup` | `bool` |  |  |  |
| `spec.backendSets[].backends[].drain` | `bool` |  |  |  |
| `spec.backendSets[].backends[].offline` | `bool` |  |  |  |
| `spec.backendSets[].backends[].maxConnections` | `int32` |  |  |  |
| `spec.backendSets[].sslConfiguration` | `SslConfiguration` |  |  |  |
| `spec.backendSets[].sslConfiguration.certificateIds` | `[]string` |  |  |  |
| `spec.backendSets[].sslConfiguration.certificateName` | `string` |  |  |  |
| `spec.backendSets[].sslConfiguration.cipherSuiteName` | `string` |  |  |  |
| `spec.backendSets[].sslConfiguration.protocols` | `[]string` |  |  |  |
| `spec.backendSets[].sslConfiguration.serverOrderPreference` | `string` |  |  |  |
| `spec.backendSets[].sslConfiguration.trustedCertificateAuthorityIds` | `[]string` |  |  |  |
| `spec.backendSets[].sslConfiguration.verifyDepth` | `int32` |  |  |  |
| `spec.backendSets[].sslConfiguration.verifyPeerCertificate` | `bool` |  |  |  |
| `spec.backendSets[].sslConfiguration.hasSessionResumption` | `bool` |  |  |  |
| `spec.backendSets[].backendMaxConnections` | `int32` |  |  |  |
| `spec.backendSets[].lbCookieSessionPersistence` | `LbCookieSessionPersistenceConfig` |  |  |  |
| `spec.backendSets[].lbCookieSessionPersistence.cookieName` | `string` |  |  |  |
| `spec.backendSets[].lbCookieSessionPersistence.disableFallback` | `bool` |  |  |  |
| `spec.backendSets[].lbCookieSessionPersistence.domain` | `string` |  |  |  |
| `spec.backendSets[].lbCookieSessionPersistence.isHttpOnly` | `bool` |  |  |  |
| `spec.backendSets[].lbCookieSessionPersistence.isSecure` | `bool` |  |  |  |
| `spec.backendSets[].lbCookieSessionPersistence.maxAgeInSeconds` | `int32` |  |  |  |
| `spec.backendSets[].lbCookieSessionPersistence.path` | `string` |  |  |  |
| `spec.backendSets[].appCookieSessionPersistence` | `SessionPersistenceConfig` |  |  |  |
| `spec.backendSets[].appCookieSessionPersistence.cookieName` | `string` | yes |  |  |
| `spec.backendSets[].appCookieSessionPersistence.disableFallback` | `bool` |  |  |  |
| `spec.listeners` | `[]Listener` | yes |  |  |
| `spec.listeners[].name` | `string` | yes |  |  |
| `spec.listeners[].port` | `int32` |  |  |  |
| `spec.listeners[].protocol` | `enum` |  |  |  |
| `spec.listeners[].defaultBackendSetName` | `string` | yes |  |  |
| `spec.listeners[].sslConfiguration` | `SslConfiguration` |  |  |  |
| `spec.listeners[].sslConfiguration.certificateIds` | `[]string` |  |  |  |
| `spec.listeners[].sslConfiguration.certificateName` | `string` |  |  |  |
| `spec.listeners[].sslConfiguration.cipherSuiteName` | `string` |  |  |  |
| `spec.listeners[].sslConfiguration.protocols` | `[]string` |  |  |  |
| `spec.listeners[].sslConfiguration.serverOrderPreference` | `string` |  |  |  |
| `spec.listeners[].sslConfiguration.trustedCertificateAuthorityIds` | `[]string` |  |  |  |
| `spec.listeners[].sslConfiguration.verifyDepth` | `int32` |  |  |  |
| `spec.listeners[].sslConfiguration.verifyPeerCertificate` | `bool` |  |  |  |
| `spec.listeners[].sslConfiguration.hasSessionResumption` | `bool` |  |  |  |
| `spec.listeners[].connectionConfiguration` | `ConnectionConfiguration` |  |  |  |
| `spec.listeners[].connectionConfiguration.idleTimeoutInSeconds` | `int64` |  |  |  |
| `spec.listeners[].connectionConfiguration.backendTcpProxyProtocolVersion` | `int32` |  |  |  |
| `spec.listeners[].hostnameNames` | `[]string` |  |  |  |
| `spec.listeners[].ruleSetNames` | `[]string` |  |  |  |
| `spec.listeners[].routingPolicyName` | `string` |  |  |  |
| `spec.certificates` | `[]Certificate` |  |  |  |
| `spec.certificates[].certificateName` | `string` | yes |  |  |
| `spec.certificates[].caCertificate` | `string` |  |  |  |
| `spec.certificates[].publicCertificate` | `string` |  |  |  |
| `spec.certificates[].privateKey` | `string` (sensitive) |  |  |  |
| `spec.certificates[].passphrase` | `string` (sensitive) |  |  |  |
| `spec.hostnames` | `[]Hostname` |  |  |  |
| `spec.hostnames[].name` | `string` | yes |  |  |
| `spec.hostnames[].hostname` | `string` | yes |  |  |
| `spec.ruleSets` | `[]RuleSet` |  |  |  |
| `spec.ruleSets[].name` | `string` | yes |  |  |
| `spec.ruleSets[].items` | `[]RuleSetItem` | yes |  |  |
| `spec.ruleSets[].items[].action` | `enum` |  |  |  |
| `spec.ruleSets[].items[].header` | `string` |  |  |  |
| `spec.ruleSets[].items[].value` | `string` |  |  |  |
| `spec.ruleSets[].items[].prefix` | `string` |  |  |  |
| `spec.ruleSets[].items[].suffix` | `string` |  |  |  |
| `spec.ruleSets[].items[].redirectUri` | `RedirectUri` |  |  |  |
| `spec.ruleSets[].items[].redirectUri.protocol` | `string` |  |  |  |
| `spec.ruleSets[].items[].redirectUri.host` | `string` |  |  |  |
| `spec.ruleSets[].items[].redirectUri.port` | `int32` |  |  |  |
| `spec.ruleSets[].items[].redirectUri.path` | `string` |  |  |  |
| `spec.ruleSets[].items[].redirectUri.query` | `string` |  |  |  |
| `spec.ruleSets[].items[].responseCode` | `int32` |  |  |  |
| `spec.ruleSets[].items[].conditions` | `[]RuleSetItemCondition` |  |  |  |
| `spec.ruleSets[].items[].conditions[].attributeName` | `string` | yes |  |  |
| `spec.ruleSets[].items[].conditions[].attributeValue` | `string` | yes |  |  |
| `spec.ruleSets[].items[].conditions[].operator` | `string` |  |  |  |
| `spec.ruleSets[].items[].allowedMethods` | `[]string` |  |  |  |
| `spec.ruleSets[].items[].statusCode` | `int32` |  |  |  |
| `spec.ruleSets[].items[].areInvalidCharactersAllowed` | `bool` |  |  |  |
| `spec.ruleSets[].items[].httpLargeHeaderSizeInKb` | `int32` |  |  |  |
| `spec.ruleSets[].items[].defaultMaxConnections` | `int32` |  |  |  |
| `spec.ruleSets[].items[].ipMaxConnections` | `[]IpMaxConnection` |  |  |  |
| `spec.ruleSets[].items[].ipMaxConnections[].ipAddresses` | `[]string` |  |  |  |
| `spec.ruleSets[].items[].ipMaxConnections[].maxConnections` | `int32` |  |  |  |
| `spec.ruleSets[].items[].description` | `string` |  |  |  |

## Field Details

### spec.compartmentId

`string | valueFrom` · required

OCID of the compartment where the load balancer will be created.

- references: OciCompartment (`status.outputs.compartment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciCompartment, name: <that resource's name>, fieldPath: status.outputs.compartment_id}} -- a bare string does not parse

### spec.displayName

`string`

Human-readable name for the load balancer shown in the OCI Console.
Falls back to metadata.name if not provided.

### spec.shape

`string` · required

Load balancer shape. Use "flexible" for configurable bandwidth (recommended).
Deprecated fixed shapes ("100Mbps", "400Mbps", "8000Mbps") are accepted
for backward compatibility.

- rule: {"string":{"minLen":"1"}}

### spec.shapeDetails

`ShapeDetails`

Bandwidth configuration for flexible-shape load balancers.
Required when shape is "flexible". Ignored for fixed shapes.

### spec.shapeDetails.minimumBandwidthInMbps

`int32`

Minimum bandwidth in Mbps. The load balancer always provides at least
this bandwidth. Valid range: 10-8000.

- rule: {"int32":{"lte":8000,"gte":10}}

### spec.shapeDetails.maximumBandwidthInMbps

`int32`

Maximum bandwidth in Mbps. The load balancer can burst up to this
bandwidth. Must be >= minimum_bandwidth_in_mbps. Valid range: 10-8000.

- rule: {"int32":{"lte":8000,"gte":10}}

### spec.subnetIds

`[]string | valueFrom` · required

OCIDs of subnets where the load balancer will be provisioned.
At least one subnet is required. For regional (recommended) load balancers,
provide subnets in two different availability domains for high availability.
Changing subnets after creation forces load balancer recreation.

- references: OciSubnet (`status.outputs.subnet_id`)
- rule: {"repeated":{"minItems":"1"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.isPrivate

`bool`

When true, creates a private load balancer that is not accessible from
the public internet. Private load balancers receive only private IP
addresses from the assigned subnets.
Changing this after creation forces load balancer recreation.

### spec.networkSecurityGroupIds

`[]string | valueFrom`

OCIDs of network security groups applied to the load balancer.

- references: OciSecurityGroup (`status.outputs.network_security_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.network_security_group_id}} -- a bare string does not parse

### spec.isDeleteProtectionEnabled

`bool`

When true, prevents accidental deletion of the load balancer.
The protection must be explicitly disabled before the load balancer
can be deleted.

### spec.ipMode

`string`

IP version mode for the load balancer. Accepted values: "IPV4", "IPV6".
When omitted, defaults to "IPV4".

### spec.reservedIps

`[]ReservedIp`

Pre-created reserved public IPs to assign to the load balancer.
When omitted, OCI assigns ephemeral public IPs.

### spec.reservedIps[].id

`string` · required

OCID of the reserved public IP to assign to the load balancer.

- rule: {"string":{"minLen":"1"}}

### spec.isRequestIdEnabled

`bool`

When true, the load balancer adds a request ID header to each request
for tracing and debugging purposes.

### spec.requestIdHeader

`string`

Custom header name for the request ID. Only effective when
is_request_id_enabled is true. When omitted, OCI uses a default header.

### spec.backendSets

`[]BackendSet` · required

Backend sets define groups of backend servers with load balancing
policies and health checking. At least one backend set is required.
Each listener routes traffic to exactly one default backend set.

- rule: {"repeated":{"minItems":"1"}}

### spec.backendSets[].name

`string` · required

Unique name for this backend set within the load balancer.
Listeners reference backend sets by this name.

- rule: {"string":{"minLen":"1"}}

### spec.backendSets[].policy

`enum`

Load balancing policy that determines how traffic is distributed
across backends.

- rule: {"enum":{"notIn":[0]}}

Allowed values (use exactly as shown):

- `policy_unspecified`
- `round_robin`
- `least_connections`
- `ip_hash`

### spec.backendSets[].healthChecker

`HealthChecker` · required

Health checker configuration that monitors backend availability.
The load balancer removes unhealthy backends from the rotation
until they pass health checks again.

- rule: {"required":true}

### spec.backendSets[].healthChecker.protocol

`enum`

Protocol used for health checks.

- rule: {"enum":{"notIn":[0]}}

Allowed values (use exactly as shown):

- `protocol_unspecified`
- `http`
- `tcp`

### spec.backendSets[].healthChecker.port

`int32`

Port on the backend server to probe. When omitted or set to 0,
the health checker uses the backend's traffic port.

### spec.backendSets[].healthChecker.urlPath

`string`

URL path for HTTP health checks (e.g., "/health", "/ready").
Required when protocol is http. Ignored for tcp.

### spec.backendSets[].healthChecker.returnCode

`int32`

Expected HTTP status code from healthy backends (e.g., 200).
When omitted, any 2xx status is considered healthy.

### spec.backendSets[].healthChecker.responseBodyRegex

`string`

Regex pattern to match against the response body. The backend is
considered healthy only if the response body matches this pattern.
When omitted, body content is not checked.

### spec.backendSets[].healthChecker.intervalMs

`int32`

Interval between consecutive health checks in milliseconds.
When omitted, defaults to 30000 (30 seconds).

### spec.backendSets[].healthChecker.timeoutInMillis

`int32`

Maximum time to wait for a health check response in milliseconds.
When omitted, defaults to 3000 (3 seconds).

### spec.backendSets[].healthChecker.retries

`int32`

Number of consecutive failed health checks before marking a backend
as unhealthy. When omitted, defaults to 3.

### spec.backendSets[].healthChecker.isForcePlainText

`bool`

When true, forces health checks over plain text even when the
backend set has SSL configured.

### spec.backendSets[].backends

`[]Backend`

Backend servers in this set. When omitted, the backend set is
created without backends (useful when backends are added dynamically).

### spec.backendSets[].backends[].ipAddress

`string` · required

IP address of the backend server.

- rule: {"string":{"minLen":"1"}}

### spec.backendSets[].backends[].port

`int32`

Port on which the backend server listens for traffic.

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.backendSets[].backends[].weight

`int32`

Relative weight for traffic distribution. Higher weights receive
proportionally more traffic. When omitted, defaults to 1.

### spec.backendSets[].backends[].backup

`bool`

When true, this backend only receives traffic when all non-backup
backends are unhealthy.

### spec.backendSets[].backends[].drain

`bool`

When true, the backend is in drain mode -- existing connections
complete but no new connections are sent to this backend.

### spec.backendSets[].backends[].offline

`bool`

When true, the backend is temporarily taken offline. No traffic
is sent to offline backends.

### spec.backendSets[].backends[].maxConnections

`int32`

Maximum number of simultaneous connections to this backend.
When omitted, the backend set's backend_max_connections applies.

### spec.backendSets[].sslConfiguration

`SslConfiguration`

SSL configuration for encrypting traffic between the load balancer
and the backend servers (backend SSL / re-encryption).

### spec.backendSets[].sslConfiguration.certificateIds

`[]string`

OCIDs of OCI Certificate Service certificates.
Preferred over certificate_name for certificate lifecycle management.

### spec.backendSets[].sslConfiguration.certificateName

`string`

Name of a certificate resource defined in this load balancer's
certificates list. Use certificate_ids instead for managed certificates.

### spec.backendSets[].sslConfiguration.cipherSuiteName

`string`

Name of the cipher suite for SSL negotiation.
Example: "oci-default-ssl-cipher-suite-v1".

### spec.backendSets[].sslConfiguration.protocols

`[]string`

TLS protocol versions to accept.
Example: ["TLSv1.2", "TLSv1.3"].

### spec.backendSets[].sslConfiguration.serverOrderPreference

`string`

Server cipher order preference. Accepted values:
"ENABLED" (server preference) or "DISABLED" (client preference).

### spec.backendSets[].sslConfiguration.trustedCertificateAuthorityIds

`[]string`

OCIDs of trusted CA certificates for verifying backend server
certificates (backend SSL) or client certificates (mutual TLS).

### spec.backendSets[].sslConfiguration.verifyDepth

`int32`

Maximum depth for certificate chain verification.
When omitted, defaults to 5.

### spec.backendSets[].sslConfiguration.verifyPeerCertificate

`bool`

When true, the load balancer verifies the peer's certificate.
For backend SSL: verifies the backend server certificate.
For listener SSL with mutual TLS: verifies the client certificate.

### spec.backendSets[].sslConfiguration.hasSessionResumption

`bool`

When true, enables TLS session resumption for improved performance.
Only applicable in listener SSL context. Ignored for backend set SSL.

### spec.backendSets[].backendMaxConnections

`int32`

Maximum number of simultaneous connections to allow per backend.
When omitted, connections are unlimited.

### spec.backendSets[].lbCookieSessionPersistence

`LbCookieSessionPersistenceConfig`

Load-balancer-managed cookie persistence. The LB injects and tracks
a cookie to pin clients to specific backends.

### spec.backendSets[].lbCookieSessionPersistence.cookieName

`string`

Name of the cookie. When omitted, OCI generates a default name.

### spec.backendSets[].lbCookieSessionPersistence.disableFallback

`bool`

When true, connections from clients without the cookie are rejected
rather than being assigned to a new backend. Useful for preventing
session migration.

### spec.backendSets[].lbCookieSessionPersistence.domain

`string`

Domain attribute for the Set-Cookie header.
When omitted, the cookie applies to the request domain.

### spec.backendSets[].lbCookieSessionPersistence.isHttpOnly

`bool`

When true, the cookie is marked HttpOnly (not accessible to JavaScript).

### spec.backendSets[].lbCookieSessionPersistence.isSecure

`bool`

When true, the cookie is marked Secure (only sent over HTTPS).

### spec.backendSets[].lbCookieSessionPersistence.maxAgeInSeconds

`int32`

Cookie lifetime in seconds. When omitted or set to 0, the cookie
is a session cookie (expires when the browser closes).

### spec.backendSets[].lbCookieSessionPersistence.path

`string`

Path attribute for the Set-Cookie header.
When omitted, defaults to "/".

### spec.backendSets[].appCookieSessionPersistence

`SessionPersistenceConfig`

Application-managed cookie persistence. The LB reads an existing
application cookie to determine backend affinity.

### spec.backendSets[].appCookieSessionPersistence.cookieName

`string` · required

Name of the application cookie used for session affinity.

- rule: {"string":{"minLen":"1"}}

### spec.backendSets[].appCookieSessionPersistence.disableFallback

`bool`

When true, connections from clients without the cookie are rejected
rather than being assigned to a new backend.

### spec.listeners

`[]Listener` · required

Listeners define the ports and protocols on which the load balancer
accepts client connections. At least one listener is required.

- rule: {"repeated":{"minItems":"1"}}

### spec.listeners[].name

`string` · required

Unique name for this listener within the load balancer.

- rule: {"string":{"minLen":"1"}}

### spec.listeners[].port

`int32`

Port on which the listener accepts connections.

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.listeners[].protocol

`enum`

Protocol for the listener.

- rule: {"enum":{"notIn":[0]}}

Allowed values (use exactly as shown):

- `protocol_unspecified`
- `http`
- `http2`
- `tcp`
- `grpc`

### spec.listeners[].defaultBackendSetName

`string` · required

Name of the backend set that receives traffic from this listener.
Must match a backend set defined in backend_sets.

- rule: {"string":{"minLen":"1"}}

### spec.listeners[].sslConfiguration

`SslConfiguration`

SSL configuration for encrypting traffic between clients and
the load balancer (SSL termination). Required for HTTPS listeners.

### spec.listeners[].sslConfiguration.certificateIds

`[]string`

OCIDs of OCI Certificate Service certificates.
Preferred over certificate_name for certificate lifecycle management.

### spec.listeners[].sslConfiguration.certificateName

`string`

Name of a certificate resource defined in this load balancer's
certificates list. Use certificate_ids instead for managed certificates.

### spec.listeners[].sslConfiguration.cipherSuiteName

`string`

Name of the cipher suite for SSL negotiation.
Example: "oci-default-ssl-cipher-suite-v1".

### spec.listeners[].sslConfiguration.protocols

`[]string`

TLS protocol versions to accept.
Example: ["TLSv1.2", "TLSv1.3"].

### spec.listeners[].sslConfiguration.serverOrderPreference

`string`

Server cipher order preference. Accepted values:
"ENABLED" (server preference) or "DISABLED" (client preference).

### spec.listeners[].sslConfiguration.trustedCertificateAuthorityIds

`[]string`

OCIDs of trusted CA certificates for verifying backend server
certificates (backend SSL) or client certificates (mutual TLS).

### spec.listeners[].sslConfiguration.verifyDepth

`int32`

Maximum depth for certificate chain verification.
When omitted, defaults to 5.

### spec.listeners[].sslConfiguration.verifyPeerCertificate

`bool`

When true, the load balancer verifies the peer's certificate.
For backend SSL: verifies the backend server certificate.
For listener SSL with mutual TLS: verifies the client certificate.

### spec.listeners[].sslConfiguration.hasSessionResumption

`bool`

When true, enables TLS session resumption for improved performance.
Only applicable in listener SSL context. Ignored for backend set SSL.

### spec.listeners[].connectionConfiguration

`ConnectionConfiguration`

Connection configuration for idle timeout and proxy protocol settings.

### spec.listeners[].connectionConfiguration.idleTimeoutInSeconds

`int64`

Maximum idle time in seconds before the load balancer closes the
connection. Applies to both client-side and backend-side connections.

### spec.listeners[].connectionConfiguration.backendTcpProxyProtocolVersion

`int32`

Backend TCP proxy protocol version. When set, the load balancer
prepends proxy protocol headers to backend connections, allowing
backends to see the original client IP.
Accepted values: 1 or 2.

### spec.listeners[].hostnameNames

`[]string`

Names of hostname resources defined in this load balancer's hostnames
list. When set, the listener only handles requests matching these
hostnames (virtual host routing).

### spec.listeners[].ruleSetNames

`[]string`

Names of rule set resources defined in this load balancer's rule_sets
list. Rule sets are applied in the order specified.

### spec.listeners[].routingPolicyName

`string`

Name of a routing policy for content-based routing.
Routing policies are an advanced OCI feature managed outside
this component.

### spec.certificates

`[]Certificate`

TLS/SSL certificates for HTTPS termination. Certificates are referenced
by name in listener and backend set SSL configurations.

### spec.certificates[].certificateName

`string` · required

Unique name for this certificate within the load balancer.
SSL configurations reference certificates by this name.

- rule: {"string":{"minLen":"1"}}

### spec.certificates[].caCertificate

`string`

PEM-encoded CA certificate chain.

### spec.certificates[].publicCertificate

`string`

PEM-encoded public certificate.

### spec.certificates[].privateKey

`string` · sensitive

PEM-encoded private key for the certificate. Sensitive.

### spec.certificates[].passphrase

`string` · sensitive

Passphrase for an encrypted private key. Sensitive.

### spec.hostnames

`[]Hostname`

Virtual hostnames for host-based routing. Hostnames are referenced
by name in listener configurations to route requests based on the
HTTP Host header.

### spec.hostnames[].name

`string` · required

Unique name for this hostname resource within the load balancer.

- rule: {"string":{"minLen":"1"}}

### spec.hostnames[].hostname

`string` · required

The fully qualified domain name (FQDN) to match against the
HTTP Host header. Example: "app.example.com".

- rule: {"string":{"minLen":"1"}}

### spec.ruleSets

`[]RuleSet`

Rule sets for advanced request/response manipulation. Rule sets are
referenced by name in listener configurations.

### spec.ruleSets[].name

`string` · required

Unique name for this rule set within the load balancer.

- rule: {"string":{"minLen":"1"}}

### spec.ruleSets[].items

`[]RuleSetItem` · required

Rules in this set. At least one rule is required.

- rule: {"repeated":{"minItems":"1"}}

### spec.ruleSets[].items[].action

`enum`

The type of rule. Determines which fields are applicable.

- rule: {"enum":{"notIn":[0]}}

Allowed values (use exactly as shown):

- `action_unspecified`
- `add_http_request_header`
- `add_http_response_header`
- `extend_http_request_header_value`
- `extend_http_response_header_value`
- `remove_http_request_header`
- `remove_http_response_header`
- `redirect`
- `allow`
- `control_access_using_http_methods`
- `http_header`
- `ip_based_max_connections`

### spec.ruleSets[].items[].header

`string`

HTTP header name. Used by header add/remove/extend actions.

### spec.ruleSets[].items[].value

`string`

Header value. Used by add_http_request_header and
add_http_response_header actions.

### spec.ruleSets[].items[].prefix

`string`

Prefix to prepend to an existing header value.
Used by extend_http_request_header_value and
extend_http_response_header_value actions.

### spec.ruleSets[].items[].suffix

`string`

Suffix to append to an existing header value.
Used by extend_http_request_header_value and
extend_http_response_header_value actions.

### spec.ruleSets[].items[].redirectUri

`RedirectUri`

Redirect URI template. Used by the redirect action.

### spec.ruleSets[].items[].redirectUri.protocol

`string`

Target protocol. Use "{protocol}" to preserve the original.

### spec.ruleSets[].items[].redirectUri.host

`string`

Target hostname. Use "{host}" to preserve the original.

### spec.ruleSets[].items[].redirectUri.port

`int32`

Target port. Use 0 to preserve the original.

### spec.ruleSets[].items[].redirectUri.path

`string`

Target path. Use "{path}" to preserve the original.

### spec.ruleSets[].items[].redirectUri.query

`string`

Target query string. Use "{query}" to preserve the original.

### spec.ruleSets[].items[].responseCode

`int32`

HTTP response code for the redirect (e.g., 301, 302, 307, 308).
Used by the redirect action.

### spec.ruleSets[].items[].conditions

`[]RuleSetItemCondition`

Conditions that must be met for this rule to apply.
Used by redirect and allow actions.

### spec.ruleSets[].items[].conditions[].attributeName

`string` · required

Attribute to evaluate. Accepted values:
  "PATH", "SOURCE_IP_ADDRESS", "SOURCE_VCN_ID", "SOURCE_VCN_IP_ADDRESS".

- rule: {"string":{"minLen":"1"}}

### spec.ruleSets[].items[].conditions[].attributeValue

`string` · required

Value to match against the attribute.

- rule: {"string":{"minLen":"1"}}

### spec.ruleSets[].items[].conditions[].operator

`string`

Matching operator. Accepted values:
  "EXACT_MATCH", "FORCE_LONGEST_PREFIX_MATCH", "PREFIX_MATCH", "SUFFIX_MATCH".
When omitted, defaults to "EXACT_MATCH".

### spec.ruleSets[].items[].allowedMethods

`[]string`

Allowed HTTP methods. Requests using other methods receive the
status_code response. Used by control_access_using_http_methods.

### spec.ruleSets[].items[].statusCode

`int32`

HTTP status code returned when access is denied (e.g., 403, 405).
Used by control_access_using_http_methods.

### spec.ruleSets[].items[].areInvalidCharactersAllowed

`bool`

When true, allows invalid characters in HTTP headers.
Used by the http_header action.

### spec.ruleSets[].items[].httpLargeHeaderSizeInKb

`int32`

Maximum HTTP header size in KB.
Used by the http_header action.

### spec.ruleSets[].items[].defaultMaxConnections

`int32`

Default maximum connections per IP when no specific IP rule matches.
Used by ip_based_max_connections.

### spec.ruleSets[].items[].ipMaxConnections

`[]IpMaxConnection`

Per-IP connection limits.
Used by ip_based_max_connections.

### spec.ruleSets[].items[].ipMaxConnections[].ipAddresses

`[]string`

IP addresses to apply this connection limit to.

### spec.ruleSets[].items[].ipMaxConnections[].maxConnections

`int32`

Maximum simultaneous connections allowed from these IPs.

### spec.ruleSets[].items[].description

`string`

Description of the rule.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OciApplicationLoadBalancer, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.load_balancer_id` | `string` | OCID of the load balancer. |
| `status.outputs.ip_addresses` | `string` | Comma-separated IP addresses assigned to the load balancer. For public load balancers, this includes the public IP(s). For private load balancers, this includes the private IP(s). |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.compartmentId` | OciCompartment | `status.outputs.compartment_id` |
| `spec.subnetIds` | OciSubnet | `status.outputs.subnet_id` |
| `spec.networkSecurityGroupIds` | OciSecurityGroup | `status.outputs.network_security_group_id` |

## See Also

- [Overview](../README.md)
