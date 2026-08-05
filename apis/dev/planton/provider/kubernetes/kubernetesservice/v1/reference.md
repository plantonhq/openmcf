# KubernetesService

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `kubernetes.planton.dev/v1`

**KubernetesServiceSpec** defines a standalone Kubernetes Service — the stable
network identity in front of a set of pods. A Service gives its backends one
durable virtual IP and DNS name while the pods behind it come and go, and it is
the unit every other networking construct composes on: Ingress backends point at
a Service, NetworkPolicies allow traffic to the pods a Service selects, and
sibling workloads connect to its in-cluster DNS name.

Workload kinds (KubernetesDeployment, KubernetesStatefulSet) already create a
Service for their own pods; this standalone kind covers everything else:
exposing pods managed outside Planton, LoadBalancer/NodePort exposure with
cloud-provider annotations, ExternalName aliases to services outside the
cluster, headless services for custom discovery, dual-stack addressing, and
selectorless services fronting manually-managed endpoints.

The spec covers the complete core/v1 ServiceSpec surface. The single deliberate
omission is the deprecated `loadBalancerIP` field — upstream deprecated it as
under-specified and non-portable; every cloud expresses a pinned LB address
through provider-specific annotations instead (set them in `annotations`).

## Example

```yaml
# Full-surface test manifest for the offline plan proofs. Exercises the
# LoadBalancer arm (source ranges, class, node-port allocation, health-check
# port), dual-stack fields, session affinity with timeout, and multi-port
# shapes — the arms the kind-cluster E2E lanes cannot exercise live.
# traffic_distribution is deliberately absent: the proof runs through BOTH
# engines from this one manifest, and the Terraform module fails the plan
# loudly when that (Pulumi-only) field is set.
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesService
metadata:
  name: test-service
spec:
  namespace:
    value: default
  name: test-service
  labels:
    team: platform-engineering
  annotations:
    service.beta.kubernetes.io/aws-load-balancer-type: nlb
  type: load_balancer
  selector:
    app: test-app
  ports:
    - name: http
      port: 80
      target_port: "8080"
      protocol: TCP
      app_protocol: http
    - name: metrics
      port: 9090
      target_port: metrics
      protocol: TCP
      node_port: 30990
  external_traffic_policy: local
  health_check_node_port: 30991
  session_affinity: client_ip
  session_affinity_timeout_seconds: 3600
  load_balancer_source_ranges:
    - 203.0.113.0/24
  load_balancer_class: example.com/internal-vip
  allocate_load_balancer_node_ports: false
  publish_not_ready_addresses: false
  ip_families:
    - ipv4
  ip_family_policy: single_stack
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.namespace` | `string \| valueFrom` |  |  | KubernetesNamespace (`spec.name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.labels` | `map<string, string>` |  |  |  |
| `spec.annotations` | `map<string, string>` |  |  |  |
| `spec.type` | `enum` |  | `cluster_ip` |  |
| `spec.selector` | `map<string, string>` |  |  |  |
| `spec.ports` | `[]KubernetesServicePort` |  |  |  |
| `spec.ports[].name` | `string` |  |  |  |
| `spec.ports[].protocol` | `enum` |  | `TCP` |  |
| `spec.ports[].appProtocol` | `string` |  |  |  |
| `spec.ports[].port` | `int32` |  |  |  |
| `spec.ports[].targetPort` | `string` |  |  |  |
| `spec.ports[].nodePort` | `int32` |  |  |  |
| `spec.headless` | `bool` |  |  |  |
| `spec.clusterIpAddress` | `string` |  |  |  |
| `spec.externalDnsName` | `string` |  |  |  |
| `spec.externalIps` | `[]string` |  |  |  |
| `spec.externalTrafficPolicy` | `enum` |  | `cluster` |  |
| `spec.internalTrafficPolicy` | `enum` |  | `internal_cluster` |  |
| `spec.trafficDistribution` | `enum` |  |  |  |
| `spec.sessionAffinity` | `enum` |  | `none` |  |
| `spec.sessionAffinityTimeoutSeconds` | `int32` |  |  |  |
| `spec.loadBalancerSourceRanges` | `[]string` |  |  |  |
| `spec.loadBalancerClass` | `string` |  |  |  |
| `spec.allocateLoadBalancerNodePorts` | `bool` |  |  |  |
| `spec.healthCheckNodePort` | `int32` |  |  |  |
| `spec.publishNotReadyAddresses` | `bool` |  |  |  |
| `spec.ipFamilies` | `[]enum` |  |  |  |
| `spec.ipFamilyPolicy` | `enum` |  |  |  |

## Field Details

### spec.namespace

`string | valueFrom`

The namespace to create the Service in. Accepts a literal namespace name or a
reference to a KubernetesNamespace resource, so an infra chart creates the
namespace and the Service in one run. When omitted, the Service lands in the
cluster's `default` namespace — the same behavior as kubectl without a
namespace flag.

- references: KubernetesNamespace (`spec.name`)
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesNamespace, name: <that resource's name>, fieldPath: spec.name}} -- a bare string does not parse

### spec.name

`string` · required

The name of the Service (its `metadata.name` in the cluster). This becomes the
service's DNS name, so the Kubernetes API enforces the stricter RFC 1035 label
form: lowercase alphanumeric and hyphens, at most 63 characters, and it MUST
start with a letter (a leading digit is valid for most Kubernetes names but is
rejected live for Services because DNS labels used in SRV records cannot start
with a digit).

- rule: Service name must be a valid DNS-1035 label: lowercase alphanumeric and hyphens, starting with a letter and ending with an alphanumeric (e.g. "my-service", "web")
- rule: {"string":{"minLen":"1","maxLen":"63"}}

### spec.labels

`map<string, string>`

Additional labels to apply to the Service object. Merged with the standard
Planton governance labels. These label the Service itself — pod selection is
controlled by `selector`, not by these labels.

### spec.annotations

`map<string, string>`

Annotations to apply to the Service object. For LoadBalancer services this is
where cloud-provider behavior is configured — the annotation set IS the
portable way to tune the provisioned load balancer:
- `service.beta.kubernetes.io/aws-load-balancer-type: "nlb"` — AWS NLB
- `service.beta.kubernetes.io/aws-load-balancer-internal: "true"` — AWS internal LB
- `cloud.google.com/load-balancer-type: "Internal"` — GCP internal LB
- `networking.gke.io/load-balancer-ip-addresses: "<name>"` — GKE pinned address
- `service.beta.kubernetes.io/azure-load-balancer-internal: "true"` — Azure internal LB
- `external-dns.alpha.kubernetes.io/hostname: "app.example.com"` — external-dns record

Which controller answers matters (verified live on EKS): the
`aws-load-balancer-type: "nlb"` value is handled by EKS's built-in
cloud controller — no extra install needed — while the
`aws-load-balancer-type: "external"` family of annotations is handled
ONLY by the AWS Load Balancer Controller; with those set and the
controller absent, the Service simply never receives an address (no
error anywhere — the annotation just has no reader).

### spec.type

`enum` · optional (explicit presence)

How the Service is exposed.
Default: cluster_ip

- default: `cluster_ip`
- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `kubernetes_service_type_unspecified` -- Unspecified. Defaults to cluster_ip.
- `cluster_ip` -- ClusterIP: a cluster-internal virtual IP. Reachable only from inside the cluster. The default and most common type.
- `node_port` -- NodePort: builds on ClusterIP and additionally opens a static port on every node. Reachable externally at <NodeIP>:<NodePort>.
- `load_balancer` -- LoadBalancer: builds on NodePort and asks the cloud provider to provision an external load balancer routing to the same endpoints.
- `external_name` -- ExternalName: a DNS alias — cluster DNS returns a CNAME to `external_name`. No proxying, no selectors, no ports are involved.

### spec.selector

`map<string, string>`

Label selector identifying the pods this Service routes to. Traffic goes to
pods whose labels match ALL key-value pairs listed here. Planton workloads
stamp a stable selector identity on their pods — the `app` label set to the
workload's `metadata.name` — and export the full set as their
`selector_labels` output, so `app: <workload-name>` is sufficient to select
a Planton-managed workload's pods.

Leave empty for ExternalName services, or to create a selectorless service
whose endpoints are managed manually (via EndpointSlice objects) or by an
external controller.

### spec.ports

`[]KubernetesServicePort`

The ports this Service exposes. At least one port is required for every type
except ExternalName (which forwards by DNS, not by port).

- rule: a numeric target_port must be in the range 1-65535

### spec.ports[].name

`string`

Name of this port. Required (and must be unique) when the Service exposes
more than one port; optional for a single port. Must be a valid IANA service
name: lowercase alphanumeric and hyphens, at most 15 characters, at least
one letter. Named ports let consumers reference "http" instead of a number.

- rule: Port name must be a valid IANA service name: lowercase alphanumeric and hyphens, at most 15 characters, containing at least one letter (e.g. "http", "grpc-web")

### spec.ports[].protocol

`enum` · optional (explicit presence)

The IP protocol for this port.
Default: TCP

- default: `TCP`
- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `kubernetes_service_protocol_unspecified` -- Unspecified. Defaults to TCP.
- `TCP` -- TCP — the default and most common.
- `UDP` -- UDP — DNS, QUIC, media streaming.
- `SCTP` -- SCTP — telecom workloads; requires cluster SCTP support.

### spec.ports[].appProtocol

`string`

The application protocol for this port — a hint richer dataplanes (meshes,
L7 load balancers) use to pick the right proxying behavior. Either an IANA
service name ("http", "https") or a prefixed name: "kubernetes.io/h2c"
(HTTP/2 cleartext), "kubernetes.io/ws" (WebSocket), "kubernetes.io/wss"
(WebSocket over TLS).

- rule: {"string":{"maxLen":"316"}}

### spec.ports[].port

`int32`

The port number the Service exposes — what clients connect to.
Range: 1–65535.

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.ports[].targetPort

`string`

The port on the backend pods to route to: a number ("8080") or a named
container port ("http") declared in the pod's container ports. Defaults to
the same value as `port`. Ignored for headless services (which have no
proxying — set it equal to `port` or omit it).

- rule: target_port must be a port number ("8080") or a named container port ("http" — a valid IANA service name)

### spec.ports[].nodePort

`int32`

The static port opened on every node for NodePort and LoadBalancer services.
Omit (0) to let Kubernetes allocate from the node-port range. Setting it on
a ClusterIP service is rejected by the API.
Range when set: 30000–32767 (the default node-port range).

- rule: node_port must be 0 (auto-allocate) or in the node-port range 30000-32767

### spec.headless

`bool`

When true, creates a headless service (`clusterIP: None`): no virtual IP is
allocated and DNS returns the pod IPs directly. The tool of choice for
StatefulSet peer discovery and any client that wants to talk to each pod
individually. Incompatible with NodePort/LoadBalancer types and with a static
`cluster_ip`.

### spec.clusterIpAddress

`string`

A specific cluster-internal IP to assign instead of an allocated one. Must be
a valid IP inside the cluster's service CIDR and not already in use — the API
server rejects the Service otherwise. Rarely needed; omit to let Kubernetes
allocate. For a headless service use `headless: true`, never "None" here.
Immutable after creation.

- rule: cluster_ip_address must be a valid IPv4 or IPv6 address (to create a headless service set headless: true instead of clusterIP "None")

### spec.externalDnsName

`string`

The DNS name an ExternalName service aliases to (e.g.
"db.prod.example.com") — cluster DNS returns it as a CNAME. Required when
(and only meaningful when) `type` is external_name. Must be a lowercase
RFC-1123 hostname.

- rule: external_dns_name must be a lowercase RFC-1123 hostname (e.g. "db.prod.example.com")

### spec.externalIps

`[]string`

IP addresses outside Kubernetes' management for which nodes also accept
traffic for this Service — typically VIPs owned by an external load balancer
or router that fronts the cluster. The user is responsible for routing
traffic to these IPs; Kubernetes only accepts it once it arrives at a node.

- rule: {"repeated":{"items":{"cel":[{"id":"external_ips.format","message":"each external IP must be a valid IPv4 or IPv6 address","expression":"this.isIp()"}]}}}

### spec.externalTrafficPolicy

`enum` · optional (explicit presence)

External traffic policy for NodePort and LoadBalancer services. Choose
`local` when the application needs the real client source IP.
Default: cluster

- default: `cluster`
- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `external_traffic_policy_unspecified` -- Unspecified. Defaults to cluster.
- `cluster` -- Cluster: routes to all ready endpoints across all nodes. Even load distribution, but adds a possible second hop and masquerades the client source IP.
- `local` -- Local: routes only to endpoints on the node that received the traffic. Preserves the client source IP and removes the extra hop, but nodes without endpoints drop the traffic — pair with a load balancer health check (see health_check_node_port).

### spec.internalTrafficPolicy

`enum` · optional (explicit presence)

Internal traffic policy for ClusterIP traffic.
Default: internal_cluster

- default: `internal_cluster`
- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `internal_traffic_policy_unspecified` -- Unspecified. Defaults to cluster.
- `internal_cluster` -- Cluster: routes to all ready endpoints evenly (the standard behavior).
- `internal_local` -- Local: routes only to endpoints on the same node as the client pod, dropping traffic when the node has none. Useful for node-local agents (DNS caches, log collectors) exposed behind a Service.

### spec.trafficDistribution

`enum` · optional (explicit presence)

Topology-aware routing preference. Omit for the cluster's default routing.

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `traffic_distribution_unspecified` -- Unspecified: the implementation applies its default routing strategy.
- `prefer_same_zone` -- PreferSameZone: prefer endpoints in the client's zone — cuts cross-zone data transfer cost and latency. Only set when endpoints are spread evenly enough that same-zone preference cannot overload a zone's endpoints.
- `prefer_same_node` -- PreferSameNode: prefer endpoints on the client's node. For node-local patterns where every node runs an endpoint (e.g. DaemonSet-backed services).

### spec.sessionAffinity

`enum` · optional (explicit presence)

Session affinity mode.
Default: none

- default: `none`
- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `session_affinity_unspecified` -- Unspecified. Defaults to none.
- `none` -- None: every request may land on any backend pod.
- `client_ip` -- ClientIP: requests from one client IP stick to one backend pod for the affinity timeout window. For stateful applications needing sticky sessions without a smarter L7 layer.

### spec.sessionAffinityTimeoutSeconds

`int32` · optional (explicit presence)

How long (in seconds) a ClientIP affinity pin lasts. Only meaningful when
`session_affinity` is client_ip. Kubernetes default: 10800 (3 hours).
Range: 1–86400 (1 day).

- rule: {"int32":{"lte":86400,"gte":1}}

### spec.loadBalancerSourceRanges

`[]string`

Restrict which client CIDRs may reach the provisioned load balancer. Only
meaningful for LoadBalancer services, and only enforced by cloud providers
that support it. An empty list means open to all sources.

- rule: {"repeated":{"items":{"cel":[{"id":"load_balancer_source_ranges.cidr","message":"each source range must be a valid CIDR (e.g. \"203.0.113.0/24\", \"2001:db8::/64\")","expression":"this.isIpPrefix()"}]}}}

### spec.loadBalancerClass

`string`

Selects which load-balancer implementation provisions this Service when the
cluster runs more than one (e.g. the cloud default plus MetalLB). A
label-style identifier, optionally domain-prefixed (e.g.
"example.com/internal-vip"). Omit for the cluster's default implementation.
Only settable for LoadBalancer services; immutable once set.

- rule: {"string":{"maxLen":"317"}}

### spec.allocateLoadBalancerNodePorts

`bool` · optional (explicit presence)

Whether NodePorts are auto-allocated for this LoadBalancer service.
Unset defers to the Kubernetes default (true) — deliberately NOT a
platform default, because the field may only exist on LoadBalancer
services and a platform-applied value would be rejected on every other
type. Set false only when the load-balancer implementation routes to pods
directly (e.g. VIP-mode MetalLB, some NLB IP-target setups) and the
NodePort hop is dead weight.

### spec.healthCheckNodePort

`int32`

A specific health-check NodePort for `external_traffic_policy: local`
LoadBalancer services — the port external load balancers probe to learn
which nodes hold endpoints. Omit to let Kubernetes allocate one. Only
settable when type is load_balancer AND external_traffic_policy is local
(the API rejects it otherwise). Immutable once set.

- rule: health_check_node_port must be 0 (auto-allocate) or in the node-port range 30000-32767

### spec.publishNotReadyAddresses

`bool`

When true, endpoint controllers publish pod addresses even before the pods
report Ready. The canonical use is a StatefulSet's headless governing
Service, where peers must discover each other DURING startup (databases
bootstrapping a quorum). Leave false for ordinary traffic-serving services —
publishing not-ready backends sends real traffic to pods that cannot
handle it.

### spec.ipFamilies

`[]enum`

The IP families assigned to this Service, in preference order. Usually left
empty (the cluster assigns based on its configuration and
`ip_family_policy`). Set explicitly to pin the primary family or the
dual-stack order, e.g. [ipv6, ipv4] for IPv6-primary. At most two entries,
and they must differ. The requested families must be available in the
cluster or creation fails.

- rule: {"repeated":{"maxItems":"2","items":{"enum":{"definedOnly":true,"notIn":[0]}}}}

Allowed values (use exactly as shown):

- `ip_family_unspecified` -- Unspecified.
- `ipv4` -- IPv4.
- `ipv6` -- IPv6.

### spec.ipFamilyPolicy

`enum` · optional (explicit presence)

The dual-stack policy for this Service. Omit for SingleStack.

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `ip_family_policy_unspecified` -- Unspecified. Kubernetes defaults to SingleStack.
- `single_stack` -- SingleStack: one IP family (the cluster default or ip_families[0]).
- `prefer_dual_stack` -- PreferDualStack: two families on dual-stack clusters, one on single-stack clusters — the portable way to opt into dual-stack.
- `require_dual_stack` -- RequireDualStack: two families or fail — creation errors on single-stack clusters.

## Validation Rules

- `external_name_requires_dns_name`: external_dns_name must be set when type is external_name
- `external_dns_name_only_for_external_name_type`: external_dns_name can only be set when type is external_name
- `headless_incompatible_with_nodeport_loadbalancer`: headless cannot be true when type is node_port or load_balancer
- `headless_excludes_static_cluster_ip`: cluster_ip_address cannot be set on a headless service (headless means clusterIP "None")
- `non_external_name_requires_ports`: at least one port must be specified for every service type except external_name
- `external_name_excludes_proxy_fields`: selector, cluster_ip_address, and external_ips cannot be set when type is external_name (an ExternalName service is a DNS alias, not a proxy)
- `health_check_node_port_requires_local_lb`: health_check_node_port can only be set when type is load_balancer and external_traffic_policy is local
- `session_affinity_timeout_requires_client_ip`: session_affinity_timeout_seconds can only be set when session_affinity is client_ip
- `load_balancer_fields_require_load_balancer_type`: load_balancer_source_ranges, load_balancer_class, and allocate_load_balancer_node_ports can only be set when type is load_balancer
- `ip_families_distinct`: ip_families entries must be distinct (at most one ipv4 and one ipv6)
- `single_stack_allows_one_family`: ip_families may list at most one family when ip_family_policy is single_stack

## Outputs

Reference an output from another manifest as `valueFrom: {kind: KubernetesService, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.service_name` | `string` | The name of the Service object as created in the cluster. |
| `status.outputs.namespace` | `string` | The namespace the Service was created in. |
| `status.outputs.type` | `string` | The service type as deployed (ClusterIP, NodePort, LoadBalancer, ExternalName). |
| `status.outputs.cluster_ip` | `string` | The cluster-internal virtual IP assigned to the Service. Empty for headless services (which have no virtual IP) and ExternalName services (which are DNS aliases). |
| `status.outputs.load_balancer_ip` | `string` | The IP address of the provisioned load balancer. Populated only for LoadBalancer services on providers that expose an IP (GCP, Azure, MetalLB); AWS-style providers populate `load_balancer_hostname` instead. |
| `status.outputs.load_balancer_hostname` | `string` | The DNS hostname of the provisioned load balancer. Populated only for LoadBalancer services on providers that expose a hostname (AWS ELB/NLB); IP-based providers populate `load_balancer_ip` instead. |
| `status.outputs.kube_endpoint` | `string` | In-cluster DNS endpoint of the Service — the handle exposure kinds and sibling workloads connect to. ex: my-app.my-ns.svc.cluster.local |
| `status.outputs.port_forward_command` | `string` | Ready-to-run port-forward command for reaching the Service from a developer machine without any external exposure. Empty for ExternalName services (there is nothing to forward to). ex: kubectl port-forward -n my-ns service/my-app 8080:80 |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.namespace` | KubernetesNamespace | `spec.name` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| KubernetesBackendTlsPolicy | `spec.targetRefs[].name` | `status.outputs.service_name` |
| KubernetesDestinationRule | `spec.host` | `status.outputs.kube_endpoint` |
| KubernetesGrpcRoute | `spec.rules[].filters[].requestMirror.backendRef.name` | `status.outputs.service_name` |
| KubernetesGrpcRoute | `spec.rules[].backendRefs[].name` | `status.outputs.service_name` |
| KubernetesGrpcRoute | `spec.rules[].backendRefs[].filters[].requestMirror.backendRef.name` | `status.outputs.service_name` |
| KubernetesHttpRoute | `spec.rules[].filters[].requestMirror.backendRef.name` | `status.outputs.service_name` |
| KubernetesHttpRoute | `spec.rules[].backendRefs[].name` | `status.outputs.service_name` |
| KubernetesHttpRoute | `spec.rules[].backendRefs[].filters[].requestMirror.backendRef.name` | `status.outputs.service_name` |
| KubernetesIngress | `spec.defaultBackend.serviceName` | `status.outputs.service_name` |
| KubernetesIngress | `spec.rules[].paths[].backend.serviceName` | `status.outputs.service_name` |
| KubernetesTcpRoute | `spec.rules[].backendRefs[].name` | `status.outputs.service_name` |
| KubernetesTlsRoute | `spec.rules[].backendRefs[].name` | `status.outputs.service_name` |
| KubernetesUdpRoute | `spec.rules[].backendRefs[].name` | `status.outputs.service_name` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
