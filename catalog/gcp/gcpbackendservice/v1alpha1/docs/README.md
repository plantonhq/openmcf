# GCP Backend Services: The Hub of the Load-Balancing Graph

## Where the Backend Service Sits

Every request through a Google external Application Load Balancer flows: forwarding rule → target proxy → URL map → **backend service** → backend group. The pieces above the backend service decide *which* service handles a request; the backend service decides *everything about how* that request reaches compute: which instance groups or NEGs serve it, how their capacity is measured, whether the session sticks, whether the edge cache answers first, whether a Google identity is required, and what gets logged.

That makes this the hub node of the family. Health checks (GcpHealthCheck), Cloud Armor policies (GcpCloudArmorPolicy), and backend groups all attach here by reference; URL maps reference the backend service's self-link. Composition mirrors GCP's own reference graph exactly — each arrow is a `StringValueOrRef`, each node has its own lifecycle.

## Global vs Regional Is a Real Boundary

GCP has two backend-service resources, and they are not the same resource with a location field. The **global** backend service (this kind) carries the external/edge feature set: CDN, compression, custom headers, edge security policy, the EXTERNAL→EXTERNAL_MANAGED migration, and the Traffic Director mesh surface. The **regional** backend service carries passthrough-NLB semantics the global one lacks entirely: failover policy, connection tracking, HA policy, and VPC network scoping. Folding both into one kind would force every field comment into "only if global…/only if regional…" hedges; the honest model is two kinds, and this one is the global hub. (The regional variant arrives with the regional-LB family.)

## One Health Check, Singular

The provider caps `health_checks` at exactly one entry (MinItems 1, MaxItems 1). A repeated field in the spec would advertise capacity the API does not have, so the spec models `healthCheck` as a singular reference. It is optional for exactly one reason: backend services whose backends are internet or serverless NEGs must NOT have a health check — the serverless platform manages its own health. A health-checked service with zero backends is also valid, and is the natural creation order: probe first, then attach capacity.

## Balancing Modes Decide Which Dials Exist

Each backend entry carries a `balancingMode` that decides which capacity fields are meaningful: `UTILIZATION` watches instance CPU (instance groups only — GCP strips `maxUtilization` from NEG backends at the API boundary), `RATE` caps requests per second, `CONNECTION` caps open connections for TCP/SSL, and `CUSTOM_METRICS` balances on utilization values the backends themselves report via ORCA. The spec enforces the mode↔dial coherence GCP would otherwise reject at deploy time: RATE requires a rate target, CONNECTION a connection target, CUSTOM_METRICS at least one metric. `capacityScaler: 0` is the drain lever — the backend stays attached but accepts nothing.

## CDN Is a Policy, Not a Resource

As with backend buckets, there is no standalone Cloud CDN object in GCP — `enableCdn` plus a `cdnPolicy` block on the backend service is the entire story. The backend-service cache key is richer than the bucket's: host, protocol, and query string membership are toggleable, query parameters can be white- or black-listed (mutually exclusive — enforced pre-deploy), and named cookies and headers can join the key for per-variant caching. Cache-mode/TTL coherence (USE_ORIGIN_HEADERS forbids explicit TTLs; FORCE_CACHE_ALL forbids maxTtl) is enforced by the spec before GCP would silently strip the values. CDN only fronts external schemes — the spec rejects it on INTERNAL_* before deploy.

## IAP: Zero-Trust in One Block

Identity-Aware Proxy turns the load balancer into an authentication boundary: every request must carry a valid Google identity, unauthenticated ones get a login redirect, and the backend receives signed assertion headers it can trust. The common configuration is a single flag (`iap.enabled: true` with the Google-managed OAuth client); a custom OAuth client is only needed for a branded consent screen, and its `oauth2ClientSecret` is handled as a secret end to end — reference-only in the control plane, secret in Pulumi state, never in outputs. GCP itself never returns the secret after creation, only its SHA-256.

## The Service-Mesh Arm

Four blocks exist for Traffic Director (`INTERNAL_SELF_MANAGED`) deployments: `circuitBreakers` (connection-volume limits), `outlierDetection` (passive ejection of erroring hosts — also valid on EXTERNAL_MANAGED), `consistentHash` (soft session affinity on the MAGLEV/RING_HASH ring), and `maxStreamDuration`. Their scheme applicability is enforced pre-deploy so a plain external ALB manifest cannot silently carry dead mesh configuration.

## The 90/10 Coverage Decision

| Provider field | Modeled | Notes |
|---|---|---|
| `name`, `project`, `description` | ✅ | Name defaults to `metadata.name`; RFC1035 validated |
| `protocol` (incl. H2C) | ✅ | Full released enum |
| `load_balancing_scheme` | ✅ | All four schemes; canary-migration constraint enforced |
| `port_name`, `timeout_sec`, `connection_draining_timeout_sec` | ✅ | |
| `health_checks` | ✅ `healthCheck` (singular) | GCP caps at 1; `StringValueOrRef` → GcpHealthCheck |
| `backend` (full block incl. custom_metrics) | ✅ `backends` | Mode↔dial coherence CEL; group is a ref (instance-group/NEG kinds attach as they land in the catalog) |
| `session_affinity` + `affinity_cookie_ttl_sec` + `strong_session_affinity_cookie` | ✅ | Mode↔block coherence CEL |
| `locality_lb_policy` + `locality_lb_policies` + `consistent_hash` | ✅ | Built-in + custom xDS policies; hash-policy coherence CEL |
| `enable_cdn` + `cdn_policy` (full block) | ✅ | Richer cache key than the bucket; coherence CELs |
| `security_policy` + `edge_security_policy` | ✅ | `StringValueOrRef` → GcpCloudArmorPolicy; applied via dedicated set-policy API calls |
| `iap` | ✅ | First IAP surface in the catalog; client secret `(sensitive)` |
| `log_config` | ✅ | Sampling + optional-field control |
| `custom_request_headers` / `custom_response_headers` | ✅ | Header-form validated |
| `compression_mode` | ✅ | |
| `circuit_breakers`, `outlier_detection`, `max_stream_duration` | ✅ | Scheme applicability CELs |
| `security_settings` (incl. `aws_v4_authentication`) | ✅ | `access_key` `(sensitive)` |
| `tls_settings` | ✅ | Protocol applicability CEL; SAN oneof |
| `ip_address_selection_policy` | ✅ | |
| `external_managed_migration_state` / `_testing_percentage` | ✅ | Scheme + state coherence CELs |
| `custom_metrics` (service-level) | ✅ | WEIGHTED_ROUND_ROBIN coherence CEL |
| `service_lb_policy` | ✅ | Plain URL (Network Services resource, outside the compute family) |
| signed-URL keys (separate resource) | ✅ folded as `signedUrlKeys` | Max 3; `key_value` sensitive; never FK-referenced → folded, not a kind |
| `dynamic_forwarding` | ❌ | Beta-only (Service Extensions dynamic forwarding); unreleased on the GA 6.x line |
| `network_pass_through_lb_traffic_policy` | ❌ | Beta-only on global; passthrough-NLB zonal affinity belongs to the regional family |
| backend `max_in_flight_requests*` / `traffic_duration` | ❌ | Beta-only backend dials |
| backend-service IAM (policy/binding/member) | ❌ | Beta-only resources; resource-scoped IAM trios are deliberately not modeled as kinds |
| `deletion_policy`, `params.resource_manager_tags` | ❌ | Not present in the released 6.x provider |
| `timeouts` | ❌ | Operation plumbing, not resource configuration |

## Composition

The dynamic half of a global serving path:

1. **GcpHealthCheck** — the probe, referenced by self-link.
2. **GcpBackendService** (this component) — owns backends, affinity, CDN, IAP, logging; references the health check and Cloud Armor policies.
3. **GcpUrlMap → GcpTargetHttpsProxy → GcpGlobalForwardingRule** — route dynamic paths here by self-link; static paths go to backend buckets.

Serverless NEGs bridge Cloud Run and Cloud Functions into `backends[].group` as those kinds land; instance groups attach the VM path the same way.
