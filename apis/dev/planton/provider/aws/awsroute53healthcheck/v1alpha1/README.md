# Overview

The AWS Route 53 Health Check API resource provisions the availability probes that DNS records reference to keep unhealthy endpoints out of DNS answers. Health checks are the mechanism behind DNS-level high availability on AWS: failover routing (primary/secondary), health-aware weighted routing, and multivalue answers all depend on them.

## Why We Created This API Resource

Route 53 health checks are small objects with a deceptively wide surface — eight check types across three fundamentally different monitoring models, each with its own required and forbidden parameters. Getting the combination wrong surfaces only at apply time as an opaque API error. This resource:

- **Validates the per-type contract at authoring time**: endpoint checks require a target, string-match checks require a search string, TCP has no default port, calculated checks require children, CloudWatch checks require an alarm pair — all enforced before anything reaches the cloud.
- **Makes checks first-class graph nodes**: DNS records reference a health check by ID; modeling the check as its own resource lets one check gate many records and lets calculated checks aggregate children.
- **Encodes the three monitoring models honestly**: internet-facing probing, CloudWatch-alarm mirroring (the only way to health-check private resources), and recovery-control mirroring for disaster-recovery runbooks.

## The Three Monitoring Models

### Endpoint Probing (HTTP, HTTPS, HTTP_STR_MATCH, HTTPS_STR_MATCH, TCP)

Route 53's global checker fleet probes an address you specify. The endpoint must be reachable from the public internet. Tune the probe with the interval (10 or 30 seconds), failure threshold (1–10 consecutive results), checker regions (minimum 3), latency measurement, and SNI.

### CloudWatch Alarm Mirroring (CLOUDWATCH_METRIC)

The check mirrors the state of a CloudWatch alarm instead of probing anything. This is the pattern for PRIVATE resources the checker fleet cannot reach, and for gating DNS on application-level metrics (error rates, queue depth). Configure what the check reports while the alarm has insufficient data.

### Aggregation and Recovery (CALCULATED, RECOVERY_CONTROL)

Calculated checks aggregate up to 256 child health checks and report healthy when at least the threshold number of children are healthy — composing per-endpoint checks into service-level health. Recovery-control checks mirror an Application Recovery Controller routing control: a deliberate switch for disaster-recovery runbooks rather than an observation.

## Key Features

- All eight check types with per-type CEL gating (invalid combinations fail at authoring time)
- Child health checks as references to other AwsRoute53HealthCheck resources (calculated aggregation composes in the graph)
- State shaping on any type: inversion and administrative disable (the maintenance-window switch)
- Checker-region selection, latency graphs, SNI for name-based virtual hosts
- The health check ID and ARN exported for DNS records, calculated parents, and CloudWatch metric dimensions

## Composition

- `AwsRoute53DnsRecord.health_check_id` references this resource's `health_check_id` output — failover PRIMARY records, health-aware weighted members, and multivalue answers all compose through it.
- A calculated check's `child_health_checks` reference other health checks' IDs, building service-level health from endpoint-level probes.
- CLOUDWATCH_METRIC checks compose with `AwsCloudwatchAlarm` by alarm name and region.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
