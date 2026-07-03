# GCP Global Forwarding Rules: The VIP Node (and PSC's Front Door)

## Where Traffic Enters

Everything else in the global load balancing family is wiring; the forwarding rule is the doorway. It owns exactly three entry facts — which IP, which port(s), which protocol — and one pointer: the `target` proxy that receives matched connections. DNS records point at its IP; the whole proxy → URL map → backend service → backend chain hangs off its target.

Two design consequences follow:

1. **The VIP's stability is the frontend's stability.** Every field except `target` and `labels` is ForceNew. If the rule was created with an ephemeral IP, recreating it means a NEW IP and broken DNS. Production frontends therefore reference a reserved `GcpGlobalAddress` — the address outlives any frontend rebuild, and re-pointing DNS is never needed.
2. **`target` being mutable is the zero-downtime swap.** GCP's dedicated `setTarget` call repoints a live VIP at a new proxy — new certificates, new routing table, even a blue/green frontend — with zero connection-level disruption and zero DNS churn.

## One IP, Two Rules: The Port-Range Contract

Two external forwarding rules cannot share the same `[IPAddress, IPProtocol]` with overlapping port ranges — which is precisely the mechanism that lets the standard pair share one IP: a port-80 rule pointing at a target HTTP proxy (serving the http→https redirect map) and a port-443 rule pointing at the target HTTPS proxy. Proxy-based load balancers accept only specific ports (80/8080/443); the spec validates the `"n"` / `"n-m"` shape and leaves port-product rules to the API.

## The Scheme Decides Everything

`load_balancing_scheme` is the field the rest of the spec pivots on, and the spec's CEL rules enforce the pivots pre-deploy:

- **EXTERNAL / EXTERNAL_MANAGED** — internet-facing ALBs living on Google's edge. No VPC `network` (CEL rejects it — a common misconfiguration GCP would also reject, later and more cryptically).
- **INTERNAL_MANAGED** — the cross-region internal ALB; lives in a VPC.
- **INTERNAL_SELF_MANAGED** — Traffic Director. The only scheme where `metadataFilters` (xDS client scoping) applies.
- **NONE** — Private Service Connect. The API expresses this as an EMPTY scheme string; the spec models it as the explicit sentinel `NONE` so "unset (default EXTERNAL)" and "PSC" cannot be confused. Both engines translate `NONE` → `""` identically.

## The PSC Face

With scheme `NONE`, the same resource stops being a load balancer frontend and becomes a Private Service Connect endpoint:

- `target: all-apis` or `vpc-sc` forwards a VPC's Google-API traffic to a private internal IP — no public internet, VPC-SC-compatible.
- `target: <service attachment URI>` connects to a producer's published service; the `psc_connection_status` output reports whether the producer `ACCEPTED` the connection.
- PSC extras are CEL-fenced to this scheme: `serviceDirectoryRegistration` (discovery by name) and `noAutomateDnsZone` (skip the auto-created `googleapis.com` private zone).
- PSC-for-Google-APIs names double as discovery handles and are capped at 20 letters/digits — documented on the name field; the API enforces the shorter cap.

## The Migration Canary

`externalManagedBackendBucketMigrationState` + `...TestingPercentage` drive the EXTERNAL → EXTERNAL_MANAGED migration for backend buckets behind this rule without recreating the VIP: PREPARE → TEST_BY_PERCENTAGE (with a 0-100 traffic fraction) → TEST_ALL_TRAFFIC → flip the scheme. This is the forwarding-rule-side counterpart of the backend service's own migration dials.

## The 90/10 Coverage Decision

| Provider field | Modeled | Notes |
|---|---|---|
| `name` | ✅ `forwardingRuleName` | Defaults to `metadata.name`; RFC1035 validated |
| `project` | ✅ `projectId` | `StringValueOrRef` → GcpProject; empty → provider default |
| `description` | ✅ | |
| `target` | ✅ `target` | Required ref, default_kind → GcpTargetHttpsProxy self_link; in-place |
| `ip_address` | ✅ `ipAddress` | Ref → GcpGlobalAddress `address` output (literal IP keeps state drift-free) |
| `ip_protocol` | ✅ `ipProtocol` | TCP default via middleware |
| `ip_version` | ✅ `ipVersion` | |
| `load_balancing_scheme` | ✅ | NONE sentinel ↔ API empty scheme |
| `port_range` | ✅ `portRange` | `n` / `n-m` pattern CEL |
| `network` / `subnetwork` | ✅ | Refs → GcpVpc / GcpSubnetwork; scheme coherence CEL |
| `network_tier` | ✅ `networkTier` | PREMIUM-only CEL (STANDARD is regional-only) |
| `metadata_filters` | ✅ | Traffic Director CEL fence |
| `service_directory_registrations` | ✅ `serviceDirectoryRegistration` | Provider caps the list at 1 → modeled singular |
| `no_automate_dns_zone` | ✅ | PSC CEL fence |
| `labels` | ✅ | Mutable |
| `external_managed_backend_bucket_migration_state` / `..._testing_percentage` | ✅ | Percentage-requires-state CEL |
| `source_ip_ranges` | ❌ | Present in the provider's global schema but documented by the API as valid only on REGIONAL EXTERNAL rules — modeling it here would offer a knob that always fails |
| `allow_psc_global_access` | ❌ | google-beta-only on the released 6.x line; add when it reaches GA |
| `base_forwarding_rule` (computed) | ❌ | Only meaningful with source_ip_ranges |
| `deletion_policy` | ❌ | Absent from the released 6.x line |
| `timeouts` | ❌ | Operation plumbing |

## Composition

The complete global external HTTPS serving path, VIP inward:

1. **GcpGlobalAddress** — the reserved static IP this rule binds (`ipAddress`).
2. **GcpGlobalForwardingRule** (this component) — IP + port → target.
3. **GcpTargetHttpsProxy** — TLS termination (+ `GcpTargetHttpProxy` for the port-80 redirect twin).
4. **GcpUrlMap → GcpBackendService / GcpBackendBucket → GcpHealthCheck / GcpRegionNetworkEndpointGroup** — routing and backends.

The kind registry declares `GcpTargetHttpsProxy` and `GcpGlobalAddress` as prerequisites, so the E2E harness (and composed charts) installs the full chain first.

## Operational Notes

- **DNS points at `ip_address`** — the output is always the literal IP, even when the spec referenced an address resource.
- **Frontend swap runbook**: create the replacement proxy, update `target` (in-place), verify, delete the old proxy. The VIP and DNS never move.
- The module enables `compute.googleapis.com` before creating the rule, so a fresh project works on the first deploy.
