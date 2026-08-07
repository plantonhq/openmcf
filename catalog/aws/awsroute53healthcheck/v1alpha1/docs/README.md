# AWS Route 53 Health Checks: Deployment Research

## What a Health Check Is (and Is Not)

A Route 53 health check is a **global availability signal** consumed by DNS routing. It does not restart anything, scale anything, or page anyone by itself — it flips a healthy/unhealthy bit that:

1. **DNS records consume** via `health_check_id`: an unhealthy record drops out of DNS answers (failover, weighted, multivalue routing).
2. **CloudWatch exposes** as the `HealthCheckStatus` metric, which CAN be alarmed and paged on.
3. **Other health checks aggregate** (calculated checks), composing endpoint probes into service-level health.

This makes health checks the foundation of DNS-level high availability: the cheapest, most globally distributed failover mechanism AWS offers.

## The Three Monitoring Models

### 1. Endpoint probing (HTTP, HTTPS, HTTP_STR_MATCH, HTTPS_STR_MATCH, TCP)

Route 53 operates a fleet of checkers in eight regions (us-east-1, us-west-1, us-west-2, eu-west-1, ap-southeast-1, ap-southeast-2, ap-northeast-1, sa-east-1). Each checker probes independently at the configured interval; the check is healthy when more than ~18% of checkers report success (a quorum, so one region's network weather cannot flip the state).

Implications worth knowing:

- **The endpoint must be internet-reachable.** Checkers have published IP ranges; security groups must admit them (or admit 0.0.0.0/0 on the probed port).
- **The real probe rate is much higher than the interval.** With ~15 checkers at a 30-second interval, the endpoint sees a request every ~2 seconds on average. Fast (10-second) checks roughly triple that.
- **String-match checks read only the first 5,120 bytes** of the body — keep health pages small and put the sentinel early.
- **fqdn vs ip_address**: fqdn alone re-resolves each probe (tracks DNS-based scaling); ip_address pins the target and uses fqdn only as the Host header.

### 2. CloudWatch alarm mirroring (CLOUDWATCH_METRIC)

The check mirrors an alarm's state without probing anything. This is the ONLY health-check path for private resources (internal load balancers, databases, anything the checker fleet cannot reach), and the way to gate DNS on application-level signals (5xx rates, queue depth, saturation). The `insufficient_data_health_status` knob decides what the check reports while the alarm lacks data — `LastKnownStatus` is the operationally sane default.

### 3. Aggregation and recovery (CALCULATED, RECOVERY_CONTROL)

Calculated checks aggregate up to 256 children with a healthy-count threshold — "the service is up if at least 2 of 3 regions are up." Recovery-control checks mirror an Application Recovery Controller routing control: a deliberate, auditable switch for DR runbooks rather than an observation.

## Production Guidance

- **Failover pairs**: give the PRIMARY record a health check; the SECONDARY record usually goes without one (it is the answer of last resort). Set `evaluate_target_health` on alias records instead when the target is an ALB/NLB — that reuses the load balancer's own target health for free.
- **Fast detection**: `request_interval: 10` + `failure_threshold: 2` detects failure in ~20–30 seconds, at higher probe cost. The default 30/3 detects in ~90 seconds.
- **Inverted checks** are for "route AWAY while X is up" arrangements (rare; e.g. an interstitial maintenance page).
- **`disabled` is the maintenance-window switch**: probing stops and the check reports healthy, so planned endpoint downtime does not trigger failover. Pair it with change control.
- **Cost**: basic AWS-endpoint checks are ~$0.50/month (first 50 free); non-AWS endpoints and optional features (fast interval, latency, string match) add per-feature charges.

## The Coverage Line (90/10)

### Covered

Everything on `aws_route53_health_check`: all eight types, probe tuning (interval, threshold, regions, latency, SNI), state shaping (invert, disable), calculated aggregation with child references, CloudWatch mirroring, recovery-control mirroring.

### Deliberately not modeled

- **`reference_name`**: a create-time idempotency token, not configuration — the console display name comes from the `Name` tag, which this component derives from `metadata.name`.
- **`triggers`**: a provider-side recompute map with no AWS-side meaning.
- **The Application Recovery Controller control plane** (clusters, control panels, routing controls, safety rules): a separate service surface; `routing_control_arn` composes by literal ARN.

## Validation Model

The per-type contracts are CEL rules on the spec, so invalid combinations fail at authoring time instead of as an opaque `InvalidInput` at apply: endpoint checks require a target and only endpoint checks may carry one; TCP requires an explicit port; string matching is required for exactly the `*_STR_MATCH` types; children define CALCULATED; the alarm pair defines CLOUDWATCH_METRIC; the routing-control ARN defines RECOVERY_CONTROL; customized checker-region lists need at least three entries.
