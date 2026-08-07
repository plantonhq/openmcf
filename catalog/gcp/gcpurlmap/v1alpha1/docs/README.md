# GCP URL Maps: The Routing Brain Behind Every Global External ALB

## What a URL Map Actually Decides

Every request hitting a Google Cloud global external Application Load Balancer passes through a URL map before it reaches a backend. The map reads the request's Host header and path, walks a routing table, and picks one outcome: send it to a backend service or bucket, split it across weighted backends, rewrite or redirect it, mutate headers, or serve a custom error page from a bucket.

URL maps sit in the middle of the load balancing family. They reference backend services and buckets by self-link; target HTTP(S) proxies reference the map's self-link. Change a path rule incorrectly and traffic silently misroutes — which is why the spec supports `tests[]` that GCP evaluates at update time and blocks the change if a path no longer resolves as expected.

## Evaluation Order (Why the Table Shape Matters)

Routing is not "first matching rule wins" in one flat list. GCP evaluates in layers:

1. **Host rules** match the Host header (with optional wildcards like `*.example.com`) and select a named **path matcher**.
2. Inside that path matcher, **route rules** (priority-ordered, with header/query/path matching) are tried first.
3. Then **path rules** (longest prefix match on path patterns like `/api/*`).
4. Then the path matcher's own default (service, redirect, or weighted route action).
5. If no host rule matched, the URL map's **top-level default** applies.

The spec models this structure directly — not as an opaque blob — so operators can read a routing incident from the declared table rather than reverse-engineering GCP console JSON.

## Default Targets: Service, Redirect, or Weighted Split

At every level (URL map default, path matcher default, path rule, route rule) the "what to do" is exactly one of:

- **service** — forward to a backend service or bucket self-link
- **urlRedirect** — return a redirect response (apex→www, http→https, path rewrite in Location)
- **routeAction** — weight traffic across several backend services and optionally rewrite host/path before forwarding

The top-level default is special: `defaultRouteAction` only counts when it carries `weightedBackendServices` (CEL-enforced). A rewrite-only route action belongs inside a path or route rule, not as the URL map catch-all.

## Path Matchers: pathRules vs routeRules

GCP allows two path-level rule systems inside a matcher, but **never both on the same matcher**:

- **pathRules** — simple longest-prefix patterns (`/api/*`, `/static/*`). The workhorse for host+path fan-out.
- **routeRules** — priority-ordered rules with rich matching (prefix, exact, regex, path templates with captured variables, header/query/metadata filters). Required for path template rewrites and Traffic Director-style routing.

Choosing the wrong one is a common design mistake: path rules are simpler and cheaper to reason about; route rules buy expressiveness at the cost of priority management and uniqueness constraints.

## The routeAction Coverage Boundary

The provider's `route_action` block also carries mesh-advanced sub-policies — timeout, retry, CORS, fault injection, request mirroring, max stream duration. These overlap backend service resilience settings and are rarely used on a standard external Application Load Balancer.

This kind models the routing-defining parts only: **weighted_backend_services** and **url_rewrite** (host, path prefix, and path template on route rules). Adding the advanced sub-policies later is purely additive — no consumer breaks when they land.

## The 90/10 Coverage Decision

| Provider field | Modeled | Notes |
|---|---|---|
| `name` | ✅ `urlMapName` | Defaults to `metadata.name`; RFC1035 validated |
| `project` | ✅ `projectId` | `StringValueOrRef` → GcpProject; empty → provider default |
| `description` | ✅ | |
| `default_service` | ✅ `defaultService` | `StringValueOrRef`; exactly-one default CEL |
| `default_url_redirect` | ✅ `defaultUrlRedirect` | path/prefix exclusivity CEL |
| `default_route_action.weighted_backend_services` | ✅ | Top-level requires non-empty weighted list |
| `default_route_action.url_rewrite` | ✅ | prefix/template exclusivity CEL; template route-rule only |
| `default_custom_error_response_policy` | ✅ | Global external ALBs |
| `header_action` | ✅ | Request/response add/remove at URL-map level |
| `host_rules` | ✅ | |
| `path_matchers` (full tree) | ✅ | path_rules XOR route_rules CEL per matcher |
| `path_rules` / `route_rules` targets | ✅ | service / redirect / route_action exactly-one CEL |
| `route_rules.match_rules` | ✅ | Single path matcher dimension CEL |
| `route_rules.header/query/metadata matches` | ✅ | |
| `test` | ✅ `tests` | Self-tests block bad updates |
| `route_action.timeout/retry/cors/fault/mirror/cache` | ❌ | Documented boundary; additive follow-up |
| `map_id` / `fingerprint` (computed) | outputs only | `map_id` exported for debugging; not a composition handle |
| `deletion_policy` | ❌ | Terraform-provider lifecycle; Planton owns this |
| `params.resource_manager_tags` | ❌ | Catalog-wide tag binding decision |
| `timeouts` | ❌ | Operation plumbing |

## Composition

A typical global external serving path composes upward from backends:

1. **GcpHealthCheck** — the probe (when backends need one).
2. **GcpBackendService** — references the check, attaches backends/NEG/instance groups, CDN, IAP, Armor.
3. **GcpUrlMap** (this component) — routes host/path to backend self-links; supports weighted canary and redirects.
4. **GcpTargetHttpsProxy → GcpGlobalForwardingRule** — TLS termination and the VIP.

Reference backend self-links via `valueFrom` on `GcpBackendService` (the kind registry declares it as a prerequisite). Custom error pages reference `GcpBackendBucket`.

## Presets

| Preset | Pattern |
|--------|---------|
| [01-host-path-fanout](../presets/01-host-path-fanout.yaml) | Host rule + path prefix fan-out to API and static backends |
| [02-weighted-canary](../presets/02-weighted-canary.yaml) | 90/10 weighted split for canary rollouts |
| [03-apex-redirect](../presets/03-apex-redirect.yaml) | Apex→www HTTPS redirect catch-all |
