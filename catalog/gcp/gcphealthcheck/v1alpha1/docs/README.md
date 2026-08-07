# GCP Health Checks: The Probe Every Load-Balancer Decision Hangs On

## What a Health Check Actually Decides

Every request a Google Cloud load balancer routes is preceded by a question it answered seconds ago: *which backends are alive?* The health check is the machinery behind that answer. Backend services consult it before including a backend in rotation; managed instance groups consult it before recreating an instance they believe is dead. Get the probe wrong and one of two failure modes follows: traffic flows to dead backends (probe too lenient or probing the wrong thing), or healthy capacity is drained by false alarms (probe too strict, or a "health" endpoint that inherits its dependencies' latency).

Health checks are the leaf of the load balancing family. They reference nothing — only a project — and everything above them (backend services, and through them URL maps, proxies, and forwarding rules) references them by self-link. One check is commonly shared by many backend services, which is exactly why it is a first-class node rather than a block inside a backend service.

## One Kind, Two API Collections

GCP models global and regional health checks as separate resources (`google_compute_health_check` and `google_compute_region_health_check`) with an identical probe surface — the same seven protocol blocks, the same cadence knobs, the same logging toggle. The only genuine deltas: `source_regions` exists only on the global resource, and the regional one carries a `region` and a numeric server ID.

Because the field sets match, this component models both as ONE kind with a `region` scope switch: empty means global, set means regional. This is the same granularity test that keeps regional and global addresses as *separate* kinds (their field sets genuinely diverge — subnetwork attachment, network tiers) applied to the opposite facts. Scope matching matters downstream: a regional backend service can only reference health checks in its own region, and the `region` output exists so composition can verify that.

## Choosing the Probe Protocol

The protocol should match what the backend actually speaks *on the probe port* — which is not necessarily its serving protocol:

- **http / https / http2** probe a `request_path` and pass on any 2xx (optionally requiring the body to start with `response`). This is the right family for anything with an HTTP endpoint, and the ONLY family serverless NEGs behind an external LB accept in practice (with `USE_SERVING_PORT`).
- **tcp / ssl** prove connectability (and optionally a banner exchange via `request`/`response`). Cheap and dependency-free — right for databases and custom protocols, but blind to application-level sickness: a deadlocked process with an open listener passes a TCP probe.
- **grpc / grpc_tls** call the standard `grpc.health.v1.Health/Check` service and pass on `SERVING`. `grpc_service_name` selects one service on a multi-service server; the backend must implement the health service with the same name.

`port_specification` decides where the probe lands: a fixed number, the instance group's named port, or `USE_SERVING_PORT` (whatever the backend serves on — the right default for most modern backends). The spec enforces the coherence rules GCP would otherwise reject at deploy time: named-port probing requires a `port_name` and forbids a numeric `port`; serving-port probing forbids both.

## The Failover Math

Detection time is a product, not a field: a backend is marked unhealthy after `unhealthy_threshold` consecutive failures spaced `check_interval_sec` apart (defaults 2 × 5 s ≈ 10 s worst-case detection), and re-admitted after `healthy_threshold` consecutive successes. Tightening the interval speeds detection but multiplies probe traffic — remember every prober in Google's fleet probes independently, so backends see substantially more probe QPS than the interval suggests. `timeout_sec` must never exceed the interval; GCP rejects it, and the spec enforces it pre-deploy.

`source_regions` (global checks only) pins probing to exactly three named regions so a prober-region outage cannot flip global verdicts. It carries real constraints, all enforced in the spec: exactly 3 regions, HTTP/HTTPS/TCP only, interval ≥ 30 s, no PROXY headers, no TCP request payload — and GCP refuses to use such a check for MIG auto-healing.

## The Firewall Gotcha

Probes originate from Google's prober ranges — `35.191.0.0/16` and `130.211.0.0/22` for most load balancer types. Instance-group backends MUST have an ingress firewall rule admitting those ranges on the probe port, or every backend shows permanently unhealthy while serving traffic fine when reached directly. This is the single most common health-check debugging dead end. Serverless backends (Cloud Run, Cloud Functions) need nothing — Google probes them internally.

## The 90/10 Coverage Decision

| Provider field | Modeled | Notes |
|---|---|---|
| `name` | ✅ `healthCheckName` | Defaults to `metadata.name`; RFC1035 validated |
| `project` | ✅ `projectId` | `StringValueOrRef` → GcpProject; falls back to the provider default project |
| `region` (regional resource) | ✅ `region` | The scope switch: empty = global, set = regional |
| `check_interval_sec` / `timeout_sec` | ✅ | Defaults 5/5; timeout ≤ interval enforced pre-deploy |
| `healthy_threshold` / `unhealthy_threshold` | ✅ | Defaults 2/2 |
| all 7 protocol blocks | ✅ oneof | Exactly-one enforced by the proto oneof; per-block port-spec coherence CEL |
| `log_config.enable` | ✅ `enableLogging` | Flattened — the provider block carries exactly one field |
| `source_regions` | ✅ `sourceRegions` | Global-only; all five API constraints enforced pre-deploy |
| `health_check_id` (regional, computed) | ❌ | The numeric server ID duplicates the self-link handle; nothing composes on it |
| `deletion_policy` | ❌ | A Terraform-provider-level abandon-vs-delete lever, not a property of the health check; Planton's lifecycle management owns this concern |
| `params.resource_manager_tags` | ❌ | Write-only tag bindings, unmodeled across the catalog; adopting them is a catalog-wide decision |
| `timeouts` | ❌ | Operation plumbing, not resource configuration |

The gRPC-with-TLS block is preview-stage on the current released provider line; the Terraform module selects the beta provider for the health check resources so the whole protocol surface is available (the Pulumi provider is beta-bridged by construction — both engines expose identical surfaces).

## Composition

A typical global external serving path composes upward from here:

1. **GcpHealthCheck** (this component) — the probe.
2. **GcpBackendService** — references the check's `self_link`, attaches backends and CDN/IAP/Armor policy.
3. **GcpUrlMap → GcpTargetHttpsProxy → GcpGlobalForwardingRule** — routing, TLS termination, and the VIP.

For instance-group backends, pair the check with a **GcpFirewallRule** admitting `35.191.0.0/16` and `130.211.0.0/22` on the probe port.
