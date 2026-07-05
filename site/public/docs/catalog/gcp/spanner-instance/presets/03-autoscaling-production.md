---
title: "Autoscaling Production (Multi-Region)"
description: "Provisions a multi-region Cloud Spanner instance with Spanner's managed autoscaler: capacity follows utilization within explicit bounds, including an asymmetric override that scales one read-heavy..."
type: "preset"
rank: "03"
presetSlug: "03-autoscaling-production"
componentSlug: "spanner-instance"
componentTitle: "Spanner Instance"
provider: "gcp"
icon: "package"
order: 3
---

# Autoscaling Production (Multi-Region)

Provisions a multi-region Cloud Spanner instance with Spanner's managed autoscaler: capacity follows utilization within explicit bounds, including an asymmetric override that scales one read-heavy replica region independently of the rest.

## When to Use

- Variable or spiky workloads where fixed capacity means paying for the peak
- Multi-region deployments where one region serves disproportionate read traffic (analytics, a regional user base)
- Teams that want capacity management owned by GCP with hard cost bounds

## Key Configuration

- **Autoscaling limits 1-5 nodes** — the hard floor and ceiling; the bill can never exceed the ceiling
- **CPU target 45%** — Google's multi-region guidance (lower than the 65% regional recommendation) to preserve failover headroom; storage target 80%
- **Asymmetric option** — the selected replica region gets its own 2-10 node range instead of the instance-wide bounds; requires ENTERPRISE (or ENTERPRISE_PLUS) edition and a multi-region config
- **AUTOMATIC backup schedule** — new databases get a default backup schedule

## Customization Notes

- Replace `config` with your multi-region topology (`nam6`, `eur6`, `nam-eur-asia1`, ...); this choice is immutable
- Replace `replicaLocation` with the actual read-heavy region of your chosen config
- While autoscaling is enabled, `numNodes`/`processingUnits` are read-only reflections of the current allocation
- All autoscaling knobs update in place — tightening bounds never recreates the instance

## Related Presets

- **01-free-instance** — zero-cost instance for development
- **02-regional-production** — fixed-capacity single-region instance
